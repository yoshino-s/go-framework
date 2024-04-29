package common

import (
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

func autoConvertHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) {
			switch f.Kind() {
			case reflect.String:
				return time.Parse(time.RFC3339, data.(string))
			case reflect.Float64:
				return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
			case reflect.Int64:
				return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
			default:
				return data, nil
			}
		} else if t == reflect.TypeOf(time.Duration(0)) {
			switch f.Kind() {
			case reflect.String:
				return time.ParseDuration(data.(string))
			case reflect.Float64:
				return time.Duration(data.(float64)), nil
			case reflect.Int64:
				return time.Duration(data.(int64)), nil
			default:
				return data, nil
			}
		} else if f.Kind() == reflect.String {
			switch t.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(data.(string), 10, 64)
				if err != nil {
					return data, err
				}
				return reflect.ValueOf(i).Convert(t).Interface(), nil
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strconv.ParseUint(data.(string), 10, 64)
				if err != nil {
					return data, err
				}
				return reflect.ValueOf(i).Convert(t).Interface(), nil
			}
		}
		return data, nil
	}
}

func DecodeFromMapstructure(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			autoConvertHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err = decoder.Decode(input); err != nil {
		return err
	}
	return err
}

func MustDecodeFromMapstructure(input interface{}, result interface{}) {
	MustNoError(DecodeFromMapstructure(input, result))
}
