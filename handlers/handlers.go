package handlers

import (
	"fmt"
	"go/ast"

	"github.com/kristiehoward/go2flow/typeutils"
)

func handleField(f ast.Field) {
	tag := f.Tag.Value
	// A field is optional if the json tag includes `omitempty`
	name, isOptional := typeutils.GetTagInfo(tag)
	// A field is nullable if the identifier is a pointer (nil pointer --> null JSON)
	isNullable := typeutils.IsNullable(f)
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

	fieldType := typeutils.GetTypeInfo(f.Type)
	fmt.Print(fieldType)
	fmt.Printf(",\n")
}

// HandleTypeDef
func HandleTypeDef(ts ast.TypeSpec) {
	if !ts.Name.IsExported() {
		// Do not handle unexported structs
		return
	}

	switch t := ts.Type.(type) {
	// type MyAlias string
	// type MyAlias2 AnotherType
	case *ast.Ident:
		fmt.Printf("export type %s = %s;\n\n", ts.Name, typeutils.GetTypeInfo(t))
		return
	// type MyAlias []AnotherType
	case *ast.ArrayType:
		elementType := typeutils.GetTypeInfo(t.Elt)
		fmt.Printf("export type %s = Array<%s>;\n\n", ts.Name, elementType)
		return
	// type MyAlias map[boolean]AnotherType
	case *ast.MapType:
		keyType := typeutils.GetTypeInfo(t.Key)
		valueType := typeutils.GetTypeInfo(t.Value)
		fmt.Printf("export type %s = {[%s]: %s};\n\n", ts.Name, keyType, valueType)
		return
	case *ast.StructType:
		fmt.Printf("export type %s {\n", ts.Name)
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
