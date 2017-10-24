package main

/*
import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Map the string representation of each reflect.Type to the Flow type for that
// primative
var goTypeToFlowType = map[string]string{
	"bool":      "boolean",
	"int":       "number",
	"int64":     "number",
	"string":    "string",
	"time.Time": "string",
}

type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Sub struct {
	Name            string      `json:"name"`
	ID              int         `json:"subscription_id,omitempty"`
	CreatedAt       []time.Time `json:"created_at"`
	Publisher       Publisher   `json:"publisher"`
	AllPublishers   []Publisher `json:"all_publishers"`
	PopularityScore int64       `json:"popularity"`
	IsOffline       bool        `json:"is_offline"`
}

// TODO Kristie
// - Handle imported types from other schemas
// - Handle AST - walk a tree
// - Handle ENUMs
// - Handle all types in the file
// - Handle embedded types
// - Add tests
// - Handle pointers (nullable types)
// - Handle map
// - Keep comments?
// - Handle unexported fields
// Primitive version of the translator that requires a copy/paste and calling of
// each struct you want translated
// `go run inline_defs.go`
func main() {
	printf := fmt.Printf

	structType := reflect.TypeOf(Sub{})
	printf("type %s {\n", structType.Name())

	// Loop through and do something with the tag of each struct field
	for i := 0; i < structType.NumField(); i++ {
		fieldStructField := structType.Field(i)

		jsonTag := fieldStructField.Tag.Get("json")
		s := strings.Split(jsonTag, ",")

		tagName := s[0]
		// https://flow.org/en/docs/types/primitives/#toc-optional-object-properties
		isOptional := len(s) > 1 && s[1] == "omitempty"

		printf("  %s", tagName)
		if isOptional {
			printf("?: ")
		} else {
			printf(": ")
		}

		typeStr := fieldStructField.Type.String()
		lookupTypeStr := typeStr
		isArray := strings.HasPrefix(typeStr, "[]")
		if isArray {
			lookupTypeStr = strings.TrimPrefix(typeStr, "[]")
		}
		flowType, ok := goTypeToFlowType[lookupTypeStr]

		// Open the flow type definition for an Array of type T: Array<T>
		if isArray {
			printf("Array<")
		}

		if !ok {
			// This is not a primitive type, so print the custom type
			typeName := strings.Split(typeStr, ".")
			// grab the type without the package definition: i.e. get `Publisher`
			// from `main.Publisher`
			customType := typeName[1]
			printf("%s", customType)
		} else {
			printf("%s", flowType)
		}

		// Close the flow type definition for an Array of type T: Array<T>
		if isArray {
			printf(">")
		}
		printf(",\n")
	}
	printf("}\n")
}
*/
