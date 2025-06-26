package magic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-errors/errors"
	"github.com/swaggest/jsonschema-go"
)

func mapToStruct(v any) (*jsonschema.Struct, error) {
	s := &jsonschema.Struct{Nullable: false}
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map {
		return nil, errors.Errorf("expected map, got %T", v)
	}

	for _, key := range reflectValue.MapKeys() {
		prop := jsonschema.Field{
			Name: key.String(),
		}

		switch reflect.TypeOf(reflectValue.MapIndex(key).Interface()).Kind() {
		case reflect.Map:
			value, err := mapToStruct(reflectValue.MapIndex(key).Interface())
			if err != nil {
				return nil, errors.Errorf("failed to convert map to struct: %w", err)
			}
			prop.Value = value
		default:
			prop.Value = reflectValue.MapIndex(key).Interface()
		}
		prop.Tag = reflect.StructTag(fmt.Sprintf(`json:"%s"`, key.String()))

		s.Fields = append(s.Fields, prop)
	}

	s.Fields = append(s.Fields, jsonschema.Field{
		Name:  "_",
		Value: struct{}{},
		Tag:   reflect.StructTag(`additionalProperties:"false"`),
	})

	return s, nil
}

func MarshalJsonSchemaWithComments(
	v map[string]any,
	comments map[string]string,
) (string, error) {
	fmt.Println(comments)

	var reflector = jsonschema.Reflector{}

	s, err := mapToStruct(v)
	if err != nil {
		return "", err
	}

	schema, err := reflector.Reflect(s, jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
		if !params.Processed {
			return nil
		}
		comment := comments["$."+strings.Join(params.Path[1:], ".")+"."+params.Name]
		if comment != "" {
			params.PropertySchema.Description = &comment
		}
		return nil
	}), jsonschema.InlineRefs)
	if err != nil {
		return "", err
	}

	r, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(r), nil
}
