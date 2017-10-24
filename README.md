# go2flow

`go2flow` is a tool written in Go that generates Flow type definitions given a file with Go structs.

It relies on the [JSON encoding](https://golang.org/pkg/encoding/json/) of a Go struct (defined in struct field tags), generating Flow type definitions to be used in the consumer's JS code base.

# Development

Run the Proof of Concept on the sample file
```
go run main.go
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
- [ ] Accept CLI args