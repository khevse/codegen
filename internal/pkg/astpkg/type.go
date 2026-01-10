package astpkg

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

var (
	_ Type = (*Ident)(nil)
	_ Type = (*StarExpr)(nil)
	_ Type = (*ArrayType)(nil)
	_ Type = (*MapType)(nil)
	_ Type = (*SelectorExpr)(nil)
	_ Type = (*EllipsisType)(nil)
	_ Type = (*FuncType)(nil)
	_ Type = (*StructType)(nil)
	_ Type = (*InterfaceType)(nil)
	_ Type = (*ChanType)(nil)

	_ PackageGetterType = (*Ident)(nil)
	_ PackageGetterType = (*SelectorExpr)(nil)

	_ PackageSetterType = (*Ident)(nil)
	_ PackageSetterType = (*SelectorExpr)(nil)
)

type Type interface {
	String() string
	ExprString() string
}

type PackageGetterType interface {
	GetPackage() string
	GetPackagePath() string
}

type PackageSetterType interface {
	SetPackage(Import)
}

type PackageCarrierType interface {
	PackageGetterType
	PackageSetterType
}

type Ident struct {
	Package     string
	PackagePath string
	Name        string
	Type        Type
}

func (t Ident) String() string { return t.ExprString() }
func (t Ident) ExprString() string {
	if t.Package == "" {
		return t.Name
	}
	return fmt.Sprintf("%s.%s", t.Package, t.Name)
}

func (t Ident) GetPackage() string     { return t.Package }
func (t Ident) GetPackagePath() string { return t.PackagePath }
func (t *Ident) SetPackage(i Import)   { t.Package = i.Alias; t.PackagePath = i.Path }

type StarExpr struct {
	Type Type
}

func (t StarExpr) String() string     { return t.ExprString() }
func (t StarExpr) ExprString() string { return fmt.Sprintf("*%s", t.Type.ExprString()) }

type ArrayType struct {
	Type Type
}

func (t ArrayType) String() string     { return t.ExprString() }
func (t ArrayType) ExprString() string { return fmt.Sprintf("[]%s", t.Type.ExprString()) }

type MapType struct {
	Key   Type
	Value Type
}

func (t MapType) String() string { return t.ExprString() }
func (t MapType) ExprString() string {
	return fmt.Sprintf("map[%s]%s", t.Key.ExprString(), t.Value.ExprString())
}

type SelectorExpr struct {
	Package     string
	PackagePath string
	Name        string
	Type        Type
}

func (t SelectorExpr) String() string {
	if t.Package == "" {
		return t.Name
	}

	return fmt.Sprintf("%s(%v)", t.ExprString(), t.Type)
}

func (t SelectorExpr) ExprString() string {
	if t.Package == "" {
		return t.Name
	}

	return fmt.Sprintf("%s.%s", t.Package, t.Name)
}
func (t SelectorExpr) GetPackage() string     { return t.Package }
func (t SelectorExpr) GetPackagePath() string { return t.PackagePath }
func (t *SelectorExpr) SetPackage(i Import)   { t.Package = i.Alias; t.PackagePath = i.Path }

type EllipsisType struct {
	Type Type
}

func (t EllipsisType) String() string     { return t.ExprString() }
func (t EllipsisType) ExprString() string { return fmt.Sprintf("...%s", t.Type.ExprString()) }

type FuncType struct {
	Params  []*Field
	Results []*Field
}

func (t FuncType) String() string { return t.ExprString() }
func (t FuncType) ExprString() string {
	writer := strings.Builder{}

	writeFields := func(fieldList []*Field) {
		for i, f := range fieldList {
			if i > 0 {
				writer.WriteByte(',')
			}
			name := lo.Ternary(f.Name == "", "_", f.Name)
			writer.WriteString(fmt.Sprintf("%s %s", name, f.Type.ExprString()))
		}
	}

	writer.WriteString("func(")
	writeFields(t.Params)
	writer.WriteByte(')')

	if len(t.Results) > 0 {
		writer.WriteString(" (")
		writeFields(t.Results)
		writer.WriteByte(')')
	}

	return writer.String()
}

type InterfaceType struct {
	Methods []*Field
}

func (t InterfaceType) String() string { return t.ExprString() }
func (t InterfaceType) ExprString() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("interface{")
	for i, m := range t.Methods {
		if i > 0 {
			strBuilder.WriteByte(';')
		}
		strBuilder.WriteString(m.String())
	}
	strBuilder.WriteByte('}')

	return strBuilder.String()
}

type StructType struct {
	Fields []*Field
}

func (t StructType) String() string { return t.ExprString() }
func (t StructType) ExprString() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("struct{")
	for i, f := range t.Fields {
		if i > 0 {
			strBuilder.WriteByte(';')
		}
		strBuilder.WriteString(f.String())
	}
	strBuilder.WriteByte('}')

	return strBuilder.String()
}

type ChanType struct {
	Type      Type
	Direction ast.ChanDir
}

func (t ChanType) String() string { return t.ExprString() }
func (t ChanType) ExprString() string {
	switch t.Direction {
	case ast.SEND | ast.RECV:
		return fmt.Sprintf("chan %s", t.Type.ExprString())
	case ast.SEND:
		return fmt.Sprintf("<-chan %s", t.Type.ExprString())
	case ast.RECV:
		return fmt.Sprintf("chan<- %s", t.Type.ExprString())
	default:
		return fmt.Sprintf("chan %s(direction: %v)", t.Type, t.Direction)
	}
}

func NewType(expr ast.Expr) Type {
	switch casted := expr.(type) {
	case *ast.Ident:
		var typeSpec Type

		if obj := casted.Obj; obj != nil && obj.Decl != nil {
			if spec, ok := obj.Decl.(*ast.TypeSpec); ok {
				typeSpec = NewType(spec.Type)
			}
		}

		return &Ident{
			Package:     "",
			PackagePath: "",
			Name:        casted.Name,
			Type:        typeSpec,
		}
	case *ast.StarExpr:
		return &StarExpr{
			Type: NewType(casted.X),
		}
	case *ast.ArrayType:
		return &ArrayType{
			Type: NewType(casted.Elt),
		}
	case *ast.MapType:
		return &MapType{
			Key:   NewType(casted.Key),
			Value: NewType(casted.Value),
		}
	case *ast.SelectorExpr:
		var pkg string
		if x := casted.X; x != nil {
			if castedX, ok := x.(*ast.Ident); ok {
				pkg = castedX.Name
			}
		}

		if obj := casted.Sel.Obj; obj != nil && obj.Decl != nil {
			if typeSpec, ok := obj.Decl.(*ast.TypeSpec); ok {
				return &SelectorExpr{
					Package:     pkg,
					PackagePath: "",
					Name:        typeSpec.Name.Name,
					Type:        NewType(typeSpec.Type),
				}
			}
		}

		return &SelectorExpr{
			Package:     pkg,
			PackagePath: "",
			Name:        casted.Sel.Name,
			Type:        nil,
		}
	case *ast.Ellipsis:
		return &EllipsisType{
			Type: NewType(casted.Elt),
		}
	case *ast.FuncType:
		return &FuncType{
			Params:  NewFieldList(casted.Params),
			Results: NewFieldList(casted.Results),
		}
	case *ast.StructType:
		return &StructType{
			Fields: NewFieldList(casted.Fields),
		}
	case *ast.InterfaceType:
		return &InterfaceType{
			Methods: NewFieldList(casted.Methods),
		}
	case *ast.ChanType:
		return &ChanType{
			Type:      NewType(casted.Value),
			Direction: casted.Dir,
		}
	default:
		panic(fmt.Sprintf("unknown type: %+[1]v(%[1]T)", casted))
	}
}

func InspectType(t Type, fn func(Type) error) error {
	if t == nil {
		return nil
	}

	inspectSelf := func(result Type) error {
		if err := fn(result); err != nil {
			return fmt.Errorf("%[1]T(%[1]s): %[2]w", result, err)
		}
		return nil
	}

	switch casted := t.(type) {
	case *Ident:
		if casted.Type != nil {
			if err := InspectType(casted.Type, fn); err != nil {
				return fmt.Errorf("Ident(%s): %w", casted, err)
			}
		}
		return inspectSelf(casted)
	case *StarExpr:
		if err := InspectType(casted.Type, fn); err != nil {
			return fmt.Errorf("StarExpr(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *ArrayType:
		if err := InspectType(casted.Type, fn); err != nil {
			return fmt.Errorf("ArrayType(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *MapType:
		if err := InspectType(casted.Key, fn); err != nil {
			return fmt.Errorf("MapType key(%s): %w", casted, err)
		}
		if err := InspectType(casted.Value, fn); err != nil {
			return fmt.Errorf("MapType value(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *SelectorExpr:
		return inspectSelf(casted)
	case *EllipsisType:
		err := InspectType(casted.Type, fn)
		if err != nil {
			return fmt.Errorf("EllipsisType(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *FuncType:
		if err := inspectFieldsTypes(casted.Params, fn); err != nil {
			return fmt.Errorf("FuncType.Params(%s): %w", casted, err)
		}
		if err := inspectFieldsTypes(casted.Results, fn); err != nil {
			return fmt.Errorf("FuncType.Results(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *StructType:
		if err := inspectFieldsTypes(casted.Fields, fn); err != nil {
			return fmt.Errorf("StructType.Fields(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *InterfaceType:
		if err := inspectFieldsTypes(casted.Methods, fn); err != nil {
			return fmt.Errorf("InterfaceType.Methods(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	case *ChanType:
		if err := InspectType(casted.Type, fn); err != nil {
			return fmt.Errorf("ChanType(%s): %w", casted, err)
		}
		return inspectSelf(casted)
	default:
		return fmt.Errorf("unknown type: %T", casted)
	}
}

func SetPackageInformation(t Type, imports ImportList) error {
	if t == nil {
		return nil
	}

	casted, ok := t.(PackageCarrierType)
	if !ok {
		return nil
	}

	if pkgPath := casted.GetPackagePath(); pkgPath != "" {
		pkg, ok := imports.GetByPath(pkgPath)
		if !ok {
			return fmt.Errorf("get import for: %s(package path: %s)", casted, pkgPath)
		}
		casted.SetPackage(pkg)
	} else if pkgAlias := casted.GetPackage(); pkgAlias != "" {
		pkg, err := imports.Get(pkgAlias)
		if err != nil {
			return fmt.Errorf("get import for %s(package: %s): %w", casted, pkgAlias, err)
		}
		casted.SetPackage(pkg)
	}

	return nil
}

func ReplaceImportAliasByImportPath(t Type, importList ImportList) error {
	return InspectType(t, func(t Type) error {
		if casted, ok := t.(PackageCarrierType); ok && !isBaseType(t) {
			importByPath, ok := importList.GetByPath(casted.GetPackagePath())
			if !ok {
				return fmt.Errorf("get import by path: %s", casted.GetPackagePath())
			}

			casted.SetPackage(importByPath)
		}

		return nil
	})
}

func IsExported(name string) bool {
	var firstChar rune
	if len(name) > 0 {
		firstChar = ([]rune(name))[0]
	}

	return unicode.IsUpper(firstChar)
}

func CastToType[T Type](val any) (*T, bool) {
	if casted, ok := val.(*T); ok {
		return casted, ok
	}
	var empty T
	return &empty, false
}

func inspectFieldsTypes(fieldList []*Field, fn func(Type) error) error {
	return InspectFields(
		fieldList,
		func(f *Field) error {
			return InspectType(f.Type, fn)
		},
	)
}

func isBaseType(t Type) bool {
	casted, ok := t.(*Ident)
	return ok && casted.Type == nil
}
