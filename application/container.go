package application

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-errors/errors"
)

type Container struct {
	mapper sync.Map
}

func (c *Container) Register(app any) {
	if app, ok := c.mapper.Load(reflect.TypeOf(app)); ok {
		panic(errors.Errorf("Instance %T already registered", app))
	}

	c.mapper.Store(reflect.TypeOf(app), app)
}

// Get retrieves an instance of the specified type from the container.
func (c *Container) GetInstanceByType(v reflect.Type) any {
	if app, ok := c.mapper.Load(v); ok {
		return app
	}
	return nil
}

func (c *Container) GetInstanceByName(name string) any {
	var found any
	c.mapper.Range(func(key, value any) bool {
		k := key.(reflect.Type)
		if k.Name() == name {
			found = value
			return false // Stop iteration
		}
		return true // Continue iteration
	})
	if found == nil {
		return nil
	}
	return found
}

func (c *Container) GetInstanceByInterface(ifaceType reflect.Type) any {
	var found any
	c.mapper.Range(func(key, value any) bool {
		k := key.(reflect.Type)
		if k.Implements(ifaceType) {
			found = value
			return false // Stop iteration
		}
		return true // Continue iteration
	})
	if found == nil {
		return nil
	}
	return found
}

func (c *Container) GetInstance(typ reflect.Type, name string) any {
	if name != "" {
		return c.GetInstanceByName(name)
	} else if typ.Kind() == reflect.Interface {
		return c.GetInstanceByInterface(typ)
	} else {
		return c.GetInstanceByType(typ)
	}
}

func (c *Container) doSet(app Application) error {
	if err := c.doInject(app); err != nil {
		return err
	}
	appType := reflect.TypeOf(app)
	method, ok := appType.MethodByName("Set")
	if !ok {
		return nil
	}
	numIn := method.Type.NumIn()
	args := make([]reflect.Value, 0, numIn-1)
	for i := 1; i < numIn; i++ {
		paramType := method.Type.In(i)
		instance := c.GetInstance(paramType, "")
		if instance == nil {
			return errors.Errorf("No instance found for type %d arguments of (%T).Set, type is %v", i, appType, paramType)
		}
		args = append(args, reflect.ValueOf(instance))
	}
	resp := method.Func.Call(append([]reflect.Value{reflect.ValueOf(app)}, args...))
	if len(resp) > 0 && !resp[0].IsNil() {
		return resp[0].Interface().(error)
	}
	return nil
}

func (c *Container) doInject(app any) error {
	appValue := reflect.ValueOf(app)
	if appValue.Kind() != reflect.Ptr || appValue.IsNil() {
		return errors.Errorf("Application must be a non-nil pointer, got %T", app)
	}
	appValue = appValue.Elem()

	// 查找所有的fields
	for i := 0; i < appValue.Type().NumField(); i++ {
		field := appValue.Type().Field(i)

		if _, ok := field.Tag.Lookup("nested-inject"); ok {
			fieldValue := appValue.Field(i)
			if err := c.doInject(fieldValue.Interface()); err != nil {
				return err
			}
			continue
		}

		// 如果有 `inject` 标签，则尝试从容器中获取实例
		if tag, ok := field.Tag.Lookup("inject"); ok {
			name := ""
			isOptional := false
			if len(tag) > 0 {
				parts := strings.Split(tag, ",")
				for _, part := range parts[1:] {
					if strings.TrimSpace(part) == "optional" {
						isOptional = true
					}
				}
				name = strings.TrimSpace(parts[0])
			}

			injectApp := c.GetInstance(field.Type, name)

			if injectApp == nil && isOptional {
				continue
			}

			if injectApp == nil {
				return errors.Errorf("No instance found for field (%T).%s of type %s", app, field.Name, field.Type)
			}
			// 将实例设置到字段中
			fieldValue := appValue.Field(i)
			if !fieldValue.CanSet() {
				return errors.Errorf("Cannot set field (%T).%s, field is not settable", app, field.Name)
			}
			fieldValue.Set(reflect.ValueOf(injectApp))
		}
	}

	return nil
}
