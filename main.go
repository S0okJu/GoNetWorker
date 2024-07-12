package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Works []Work `json:"works"`
}

type Work struct {
	Url     string    `json:"url"`
	Port    int       `json:"port"`
	Request []Request `json:"request"`
}

type Request struct {
	Path   string            `json:"path"`
	Method string            `json:"method"`
	Param  map[string]string `json:"param"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <go_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	config := Config{Works: []Work{{Request: []Request{}}}}
	var currentRequest *Request

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CommentGroup:
			for _, comment := range x.List {
				fmt.Printf("comment: %v\n", comment.Text)
				if strings.HasPrefix(comment.Text, "// @NewNetWorker") {
					parts := strings.Fields(comment.Text)
					config.Works[0].Url = parts[3]
					config.Works[0].Port = parseInt(parts[5])
				} else if strings.HasPrefix(comment.Text, "// @NetWorker") {
					currentRequest = &Request{Param: make(map[string]string)}
					parts := strings.Fields(comment.Text)
					for i := 2; i < len(parts); i += 2 {
						switch parts[i] {
						case "path":
							currentRequest.Path = parts[i+1]
						case "method":
							currentRequest.Method = parts[i+1]
						}
					}
					config.Works[0].Request = append(config.Works[0].Request, *currentRequest)
				}
			}
		case *ast.StructType:
			if currentRequest != nil {
				for _, field := range x.Fields.List {
					tag := field.Tag.Value
					jsonTag := strings.Split(strings.Trim(tag, "`"), "\"")[1]
					currentRequest.Param[jsonTag] = field.Type.(*ast.Ident).Name
				}
				currentRequest = nil
			}
		}
		return true
	})

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Create /docs directory if it doesn't exist
	docsDir := filepath.Join(".", "docs")
	err = os.MkdirAll(docsDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating /docs directory: %v\n", err)
		os.Exit(1)
	}

	// Write config.json in /docs directory
	configPath := filepath.Join(docsDir, "gonetworker.config.json")
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing config.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("config.json generated successfully in %s\n", configPath)
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
