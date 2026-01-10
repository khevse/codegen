package astpkg

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFieldList(t *testing.T) {
	expr, err := parser.ParseExpr(`func(val string) any {return val}`)
	require.NoError(t, err)

	castedExpr := expr.(*ast.FuncLit)
	paramsFields := NewFieldList(castedExpr.Type.Params)
	resultsFields := NewFieldList(castedExpr.Type.Results)
	require.Equal(
		t,
		[]*Field{{Name: "val", Type: &Ident{Name: "string"}}},
		paramsFields,
	)
	require.Equal(
		t,
		[]*Field{{Name: "", Type: &Ident{Name: "any"}}},
		resultsFields,
	)
}
