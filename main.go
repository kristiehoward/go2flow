package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
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

// GetTagInfo returns a Struct Field's tag name and whether or not it is optional
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

// TODO Test the Recursion
// Better Map --> Object handling
func GetTypeInfo(fieldType ast.Expr) string {
	// Handle Arrays
	switch t := fieldType.(type) {
	// []T
	case *ast.ArrayType:
		elementType := GetTypeInfo(t.Elt)
		return fmt.Sprintf("Array<%s>", elementType)
	// map[T1]T2
	case *ast.MapType:
		keyType := GetTypeInfo(t.Key)
		valueType := GetTypeInfo(t.Value)
		return fmt.Sprintf("{[%s]: %s}", keyType, valueType)
	// T
	case *ast.Ident:
		flowType, ok := goTypeToFlowType[t.Name]
		if !ok {
			return "MISSING_TYPE_DEF_IN_MAP"
		} else {
			return flowType
		}
	}
	return "UNKNOWN_EXPR_TYPE"
}

func handleField(f ast.Field) {
	printf := fmt.Printf
	tag := f.Tag.Value
	name, isOptional := GetTagInfo(tag)
	if name == "" {
		return
	}

	printf("  %s", name)
	// https://flow.org/en/docs/types/primitives/#toc-optional-object-properties
	if isOptional {
		printf("?: ")
	} else {
		printf(": ")
	}

	fieldType := GetTypeInfo(f.Type)
	fmt.Print(fieldType)
	printf(",\n")
}

func handleTypeDef(t ast.TypeSpec) {
	if !t.Name.IsExported() {
		// Do not handle unexported structs
		return
	}
	structDecl, ok := t.Type.(*ast.StructType)
	if !ok {
		// Do not handle non-struct types
		return
	}
	printf := fmt.Printf
	printf("type %s {\n", t.Name)
	fields := structDecl.Fields.List
	for _, field := range fields {
		handleField(*field)
	}
	printf("}\n")
	printf("\n")
	return
}

// Given an ast node, handle a node if it is a type definition
func TypeDefInspector(node ast.Node) bool {
	// Check if this node is a type definition
	ts, ok := node.(*ast.TypeSpec)
	if ok {
		handleTypeDef(*ts)
	}
	return true
}

func main() {
	// Create a new set of source files
	fset := token.NewFileSet()
	// Parse the src file's information into the astNode, including the comments
	astNode, err := parser.ParseFile(fset, "fixtures/test_program.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(astNode, TypeDefInspector)
	// fmt.Println(astNode.Scope)

}
