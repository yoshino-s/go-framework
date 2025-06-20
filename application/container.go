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

func (c *Container) doSet(app Application) error {
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
			return errors.Errorf("No instance found for type %d arguments of %T.Set, type is %v", i, appType, paramType)
		}
		args = append(args, reflect.ValueOf(instance))
	}
	resp := method.Func.Call(append([]reflect.Value{reflect.ValueOf(app)}, args...))
	if len(resp) > 0 && !resp[0].IsNil() {
		return resp[0].Interface().(error)
	}
	return nil
}
