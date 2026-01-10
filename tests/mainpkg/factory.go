package mainpkg

import "fmt"

var _ IFactory = (*Factory)(nil)

// IFactory .
type IFactory interface {
	// NewObject1 comment
	NewObject1() IObject1
	// NewObject2 comment
	NewObject2() IObject2
}

// IObject1 comment .
type IObject1 interface {
	String() string
}

// IObject2 comment .
type IObject2 interface {
	String() string
}

// Factory comment
type Factory struct{}

// NewObject1 comment
func (f *Factory) NewObject1() IObject1 { return NewObject1() }

// NewObject2 comment
func (f *Factory) NewObject2() IObject2 { return NewObject2() }

// privateMethod comment
//
//nolint:unused
func (f *Factory) privateMethod() { fmt.Println("privateMethod call") }

type Object1 struct{}

func NewObject1() Object1        { return Object1{} }
func (o Object1) String() string { return "object1" }

type Object2 struct{}

func NewObject2() Object2        { return Object2{} }
func (o Object2) String() string { return "object2" }
