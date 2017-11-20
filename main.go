package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/kristiehoward/go2flow/handlers"
	"github.com/urfave/cli"
)

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
		cli.StringFlag{
			Name:  "dir, d",
			Usage: "directory containing .go file to consume",
		},
	}
)

// TypeDefInspector handles ast nodes if they are type definitions
func TypeDefInspector(node ast.Node) bool {
	// Check if this node is a type definition
	ts, ok := node.(*ast.TypeSpec)
	if ok {
		handlers.HandleTypeDef(*ts)
	}
	return true
}

func handleFile(file string) error {
	fset := token.NewFileSet()
	// Parse the src file's information into the astNode, including the comments
	astNode, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// Inspect the AST using the inspector that handles only type definitions
	ast.Inspect(astNode, TypeDefInspector)
	return nil
}

func run(c *cli.Context) error {
	file := c.String("file")
	dir := c.String("dir")

	// TODO Maxime 11/5/2017
	// Check if the file passed in the CLI has the .go extension
	if file == "" && dir == "" {
		fmt.Println("Please specify a .go file to consume")
		return nil
	}

	// Handle directory
	if dir != "" {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, f := range files {
			err := handleFile(dir + "/" + f.Name())
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Handle file
	return handleFile(file)
}

// TODO Kristie 10/24/17
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
