# go2flow

`go2flow` is a tool written in Go that generates Flow type definitions given a file with Go types.

It relies on the [JSON encoding](https://golang.org/pkg/encoding/json/) of a Go struct (defined in struct field tags), generating Flow type definitions to be used in the consumer's JS code base.

# Getting go2flow

On Linux and macOS, this command will copy the go2flow executable to your current working directory:

```bash
$ docker pull kristiehoward/go2flow:latest &&
    id=$(docker create kristiehoward/go2flow:latest) &&
    docker cp $id:/go2flow-$(uname -s)-$(uname -m) go2flow &&
    (docker rm $id >/dev/null)
```

On windows

```bash
$ docker pull kristiehoward/go2flow:latest &&
    id=$(docker create kristiehoward/go2flow:latest) &&
    docker cp $id:/go2flow-Windows-x86_64 go2flow.exe &&
    (docker rm $id >/dev/null)
```

Run `./go2flow -h` or `./go2flow --help` to print usage.

# Development

Run the Proof of Concept on the sample file
```
go run main.go -f samples/kube_types_sample.go
```

Print usage
```
go run main.go -h
```

Run the tests
```
go test
```

# TODO
- [ ] Examples of use
- [ ] More sample files
- [ ] Test output (return a string instead of printing)
- [ ] Document the decisions made for translation from Go type --> JSON output --> Flow definition
- [x] Accept CLI args


# Rules

We handle the following `TypeSpec` definitions:

**`ast.StructType`**
A struct type definition.

Example Go Code:
```go
type MyStruct struct {
    field1: bool,
    field2: string,
}
```

Rule: For each of the fields defined in this struct, handle them according to the rules of struct fields (below).

Generated Flow Code:
```js
type MyStruct = {
    field1: boolean,
    field2: string,
}
```

**`ast.Ident`**
A type definition, typically an alias.

Example Go Code:
```go
type MyStruct string
type MyStruct2 AnotherType
```

Rule: Create a flow alias to whatever the type is. If it exists in the map of go types to flow types, use that mapping. Else, use the name of the custom type.

Generated Flow Code:
```js
type MyStruct = string;
type MyStruct2 = AnotherType;
```

**`ast.ArrayType`**
A type definition, typically an alias, that is an array of another type.

Example Go Code:
```go
type MyStruct []string
type MyStruct2 []*AnotherType
```

Rule: Create a flow alias to an array of whatever the included type is. If it exists in the map of go types to flow types, use that mapping. Else, use the name of the custom type. If the type is a pointer, the pointer will either resolve to JSON or nil (it won't exist), so ignore the pointer value and use the type.

Generated Flow Code:
```js
type MyStruct = Array<string>;
type MyStruct2 = Array<AnotherType>;
```

**`ast.MapType`**
A type definition, typically an alias, that is a map.

Example Go Code:
```go
type MyStruct map[bool]AnotherType
```

Rule: Create a flow alias to a map of the appropriate types, using the following rules: If it exists in the map of go types to flow types, use that mapping. Else, use the name of the custom type. If the type is a pointer, the pointer will either resolve to JSON or nil (it won't exist), so ignore the pointer value and use the type.

Generated Flow Code:
```js
type MyStruct = {
    [boolean]: AnotherType,
}
```
