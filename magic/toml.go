package magic

import (
	"bytes"

	"github.com/pelletier/go-toml"
)

func MarshalTomlWithComments(
	v map[string]any,
	comments map[string]string,
) (string, error) {
	tree, err := toml.TreeFromMap(v)
	if err != nil {
		return "", err
	}

	for path, c := range comments {
		// 不存在的键可以选择忽略或报错
		if !tree.Has(path[2:]) {
			continue
		}
		v := tree.Get(path[2:])
		tree.SetWithOptions(path[2:], toml.SetOptions{
			Comment: c,
		}, v)
	}

	var buf bytes.Buffer
	toml.NewEncoder(&buf).
		CompactComments(true).
		Indentation("  ").
		Encode(tree)

	out := buf.String()
	return out, nil
}
