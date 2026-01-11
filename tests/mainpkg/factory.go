package mainpkg

import "fmt"

var _ IFactory = (*Factory)(nil)

// IFactory .
type IFactory interface {
	// NewObject1 comment
	NewObject1(string) IObject1
	// NewObject2 comment
	NewObject2(val string) IObject2
}

// Factory comment
type Factory struct{}

// NewObject1 comment
func (f *Factory) NewObject1(val string) IObject1 { return NewObject1(val) }

// NewObject2 comment
func (f *Factory) NewObject2(val string) IObject2 { return NewObject2(val) }

// privateMethod comment
//
//nolint:unused
func (f *Factory) privateMethod() { fmt.Println("privateMethod call") }

type Object1 struct{ value string }

func NewObject1(value string) Object1 { return Object1{value: value} }
func (o Object1) String() string      { return "object1:" + o.value }

type Object2 struct{ value string }

func NewObject2(value string) Object2 { return Object2{value: value} }
func (o Object2) String() string      { return "object2:" + o.value }
