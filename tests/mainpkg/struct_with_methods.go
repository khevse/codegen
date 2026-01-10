package mainpkg

import (
	"github.com/khevse/codegen/tests/mainpkg/childpkg"               //nolint:staticcheck,gci
	childpkgalias "github.com/khevse/codegen/tests/mainpkg/childpkg" //nolint:staticcheck
)

// StructWithMethods comment
type StructWithMethods struct {
	FieldStruct childpkg.Struct
	FieldString string
}

// GetFieldStruct comment
func (s StructWithMethods) GetFieldStruct() childpkg.Struct {
	return s.FieldStruct
}

func (s StructWithMethods) GetFieldString() string {
	return s.FieldString
}

func (s *StructWithMethods) SetFieldStruct(val childpkgalias.Struct) {
	s.FieldStruct = val
}

func (s *StructWithMethods) SetAllFields(val StructWithMethods) {
	*s = val
}

func (s *StructWithMethods) SetFieldStringFromInterface(val childpkg.Interface) {
	s.FieldString = val.String()
}
