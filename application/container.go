package application

import (
	"reflect"
	"sync"

	"github.com/go-errors/errors"
)

type Container struct {
	mapper sync.Map
}

// Register the application instance to the container.
func (c *Container) Register(app Application) {
	if app, ok := c.mapper.Load(reflect.TypeOf(app)); ok {
		panic(errors.Errorf("Application %T already registered", app))
	}

	c.mapper.Store(reflect.TypeOf(app), app)
}

// Get retrieves an instance of the specified type from the container.
func (c *Container) GetInstance(v reflect.Type) Application {
	if app, ok := c.mapper.Load(v); ok {
		return app.(Application)
	}
	return nil
}

func (c *Container) GetInstanceByName(name string) Application {
	var found Application
	c.mapper.Range(func(key, value any) bool {
		k := key.(reflect.Type)
		if k.Name() == name {
			found = value.(Application)
			return false // Stop iteration
		}
		return true // Continue iteration
	})
	if found == nil {
		return nil
	}
	return found
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
		instance := c.GetInstance(paramType)
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

func (c *Container) doInject(app Application) error {
	appValue := reflect.ValueOf(app)
	if appValue.Kind() != reflect.Ptr || appValue.IsNil() {
		return errors.Errorf("Application must be a non-nil pointer, got %T", app)
	}
	appValue = appValue.Elem()

	// 查找所有的fields
	for i := 0; i < appValue.Type().NumField(); i++ {
		field := appValue.Type().Field(i)

		// 如果有 `inject` 标签，则尝试从容器中获取实例
		if tag, ok := field.Tag.Lookup("inject"); ok {
			var app Application
			if tag == "" {
				// 如果没有指定名称，则使用字段类型作为键
				app = c.GetInstance(field.Type)
			} else {
				// 如果指定了名称，则使用名称查找
				app = c.GetInstanceByName(tag)
			}
			if app == nil {
				return errors.Errorf("No instance found for field (%T).%s of type %s", app, field.Name, field.Type)
			}
			// 将实例设置到字段中
			fieldValue := appValue.Field(i)
			if !fieldValue.CanSet() {
				return errors.Errorf("Cannot set field (%T).%s, field is not settable", app, field.Name)
			}
			if fieldValue.Type() != field.Type {
				return errors.Errorf("Field (%T).%s type mismatch: expected %s, got %s", app, field.Name, field.Type, fieldValue.Type())
			}
			fieldValue.Set(reflect.ValueOf(app))
		}
	}

	return nil
}
