package magic

import (
	"github.com/go-errors/errors"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

func MarshalYamlWithComments(
	v any,
	comments map[string]string,
) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", errors.New(err)
	}

	f, err := parser.ParseBytes(b, 0)
	if err != nil {
		return "", errors.New(err)
	}

	for path, comment := range comments {
		p, err := yaml.PathString(path)
		if err != nil {
			return "", errors.New(err)
		}
		node, err := p.ReadNode(f)
		if err != nil {
			return "", errors.New(err)
		}
		err = node.SetComment(ast.CommentGroup([]*token.Token{
			token.New(" "+comment, "", nil),
		}))
		if err != nil {
			return "", errors.New(err)
		}
		p.ReplaceWithNode(f, node)
	}

	return f.String(), nil
}
