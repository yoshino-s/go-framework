package application

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type demoInterface interface {
	DemoMethod() string
}

type demoStruct struct {
	*EmptyApplication
}

type demo1Struct struct {
	*EmptyApplication
	Int demoInterface `inject:""`
}

func (d *demoStruct) DemoMethod() string {
	return "demo"
}

func TestGetInstance(t *testing.T) {
	container := &Container{}

	demoApp := &demoStruct{
		EmptyApplication: NewEmptyApplication("DemoStruct"),
	}
	container.Register(demoApp)

	// Test GetInstanceByType
	instanceByType := container.GetInstanceByType(reflect.TypeOf(demoApp))
	assert.Equal(t, demoApp, instanceByType, "GetInstanceByType failed")
	// Test GetInstanceByInterface
	instanceByInterface := container.GetInstanceByInterface(reflect.TypeOf((*demoInterface)(nil)).Elem())
	assert.Equal(t, demoApp, instanceByInterface, "GetInstanceByInterface failed")

	demo1App := &demo1Struct{
		EmptyApplication: NewEmptyApplication("Demo1Struct"),
	}
	container.Register(demo1App)

	err := container.doInject(demo1App)
	assert.NoError(t, err, "doInject failed")
	assert.Equal(t, demoApp, demo1App.Int, "Dependency injection failed")
}
