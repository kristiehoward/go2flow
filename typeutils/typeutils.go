package typeutils

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

// Map the string representation of each reflect.Type to the Flow type for that
// primitive once it is sent as a JSON object.
var goTypeToFlowType = map[string]string{
	"bool":      "boolean",
	"int":       "number",
	"int64":     "number",
	"string":    "string",
	"time.Time": "string",
}

// IsNullable Given a field, return if it is nullable. A field is nullable if it is a pointer.
// A nil pointer generates `null` in the JSON output
func IsNullable(f ast.Field) bool {
	_, ok := f.Type.(*ast.StarExpr)
	// This is a StarExpr, which is a pointer
	return ok
}

// GetTagInfo Returns the name of the JSON field and whether or not the field is
// optional based on a struct field's tag
// TODO Kristie 10/24/17 - Update to include the edge cases in
// https://golang.org/pkg/encoding/json/#Marshal
func GetTagInfo(tag string) (name string, isOptional bool) {
	name = ""
	isOptional = false

	// Only parse tags in the form `json:"FIELDNAME,omitempty"` where `,omitempty`
	// is optional
	if !strings.Contains(tag, `json:"`) {
		return
	}
	// This regex matches a json tag only, omitting the `,omitempty` if present
	re := regexp.MustCompile(`json:"([^,"]+)"?`)
	matches := re.FindStringSubmatch(tag)

	// Capture group match will be the 2nd element
	if len(matches) < 2 {
		return
	}

	name = matches[1]
	isOptional = strings.Contains(tag, ",omitempty")
	return
}

// GetTypeInfo returns a string representing the Flow type for a given fieldType
// that's an ast.Expr
// TODO Kristie 10/24/17
// - Add tests
// - Specifically test the recursion, nullable, and optional types
// - Better Map --> Object handling
// - Figure out how to handle imported packages and their definitions in Flow
// - Handle embedded types
// - Option to keep comments?
// - Handle unexported fields
func GetTypeInfo(fieldType ast.Expr) string {
	switch t := fieldType.(type) {
	// *T
	case *ast.StarExpr:
		// Return the type of T, assume that the meaning of the pointer was
		// handled in the calling function
		return GetTypeInfo(t.X)
	// []T
	case *ast.ArrayType:
		elementType := GetTypeInfo(t.Elt)
		return fmt.Sprintf("Array<%s>", elementType)
	// map[T1]T2
	case *ast.MapType:
		keyType := GetTypeInfo(t.Key)
		valueType := GetTypeInfo(t.Value)
		return fmt.Sprintf("{[%s]: %s}", keyType, valueType)
	// Imported type package.T
	case *ast.SelectorExpr:
		typeStr := fmt.Sprintf("%s.%s", t.X, t.Sel)
		flowType, ok := goTypeToFlowType[typeStr]
		if !ok {
			// TODO What to do here when we don't recognize this package?
			return typeStr
		}
		return flowType
	// T
	case *ast.Ident:
		// Primitives will exist in the map
		flowType, ok := goTypeToFlowType[t.Name]
		// Custom type definitions in this package will have a non-nil t.Obj
		isCustomType := t.Obj != nil
		if ok {
			return flowType
		} else if isCustomType {
			return t.Name
		} else {
			return "MISSING_TYPE_DEF_IN_MAP"
		}
	}
	return "UNKNOWN_EXPR_TYPE"
}
