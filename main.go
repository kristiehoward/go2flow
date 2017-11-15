package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
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

const (
	appName  = "Go2Flow"
	appUsage = `Convert Golang types to Flow types`
)

var (
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: ".go file to consume",
		},
	}
)

// GetTagInfo returns the name of the JSON field and whether or not the field is
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

// Given a field, return if it is nullable. A field is nullable if it is a pointer.
// A nil pointer generates `null` in the JSON output
func isNullable(f ast.Field) bool {
	_, ok := f.Type.(*ast.StarExpr)
	// This is a StarExpr, which is a pointer
	return ok
}

func handleField(f ast.Field) {
	tag := f.Tag.Value
	// A field is optional if the json tag includes `omitempty`
	name, isOptional := GetTagInfo(tag)
	// A field is nullable if the identifier is a pointer (nil pointer --> null JSON)
	isNullable := isNullable(f)
	if name == "" {
		return
	}

	fmt.Printf("  %s", name)
	if isOptional {
		// https://flow.org/en/docs/types/primitives/#toc-optional-object-properties
		fmt.Printf("?: ")
	} else if isNullable {
		// If a type is optional AND nullable, it will not show up in the json
		// response, so we can assume the types here are required
		// https://flow.org/en/docs/types/primitives/#toc-maybe-types
		fmt.Printf(": ?")
	} else {
		fmt.Printf(": ")
	}

	fieldType := GetTypeInfo(f.Type)
	fmt.Print(fieldType)
	fmt.Printf(",\n")
}

func handleTypeDef(ts ast.TypeSpec) {
	if !ts.Name.IsExported() {
		// Do not handle unexported structs
		return
	}

	switch t := ts.Type.(type) {
	// type MyAlias string
	// type MyAlias2 AnotherType
	case *ast.Ident:
		fmt.Printf("type %s = %s;\n\n", ts.Name, GetTypeInfo(t))
		return
	// type MyAlias []AnotherType
	case *ast.ArrayType:
		elementType := GetTypeInfo(t.Elt)
		fmt.Printf("type %s = Array<%s>;\n\n", ts.Name, elementType)
		return
	// type MyAlias map[boolean]AnotherType
	case *ast.MapType:
		keyType := GetTypeInfo(t.Key)
		valueType := GetTypeInfo(t.Value)
		fmt.Printf("type %s = {[%s]: %s};\n\n", ts.Name, keyType, valueType)
		return
	case *ast.StructType:
		fmt.Printf("type %s {\n", ts.Name)
		fields := t.Fields.List
		for _, field := range fields {
			handleField(*field)
		}
		fmt.Printf("}\n\n")
		return
		// Don't handle anything else
	}
	return
}

// TypeDefInspector handles ast nodes if they are type definitions
func TypeDefInspector(node ast.Node) bool {
	// Check if this node is a type definition
	ts, ok := node.(*ast.TypeSpec)
	if ok {
		handleTypeDef(*ts)
	}
	return true
}

func run(c *cli.Context) error {
	file := c.String("file")

	// TODO Maxime 11/5/2017
	// Check if the file passed in the CLI has the .go extension
	if file == "" {
		fmt.Println("Please specify a .go file to consume")
		return nil
	}
	// Create a new set of source files
	fset := token.NewFileSet()
	// Parse the src file's information into the astNode, including the comments
	astNode, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Inspect the AST using the inspector that handles only type definitions
	ast.Inspect(astNode, TypeDefInspector)
	return nil
}

// TODO Kristie 10/24/17
// - Accept a file from the CLI arg
// - Accept a folder from the CLI arg
// - Dockerize development
// - Put the output through Prettier (use a container)
// - Optionally keep the comments by the struct defs?
// - Handle definitions not in the struct tags (talk to Maxime)
func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Version = "0.0.1"
	app.Flags = flags
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
