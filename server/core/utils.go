package core

import "fmt"

func Url(uri string, port int, path string) string {
	baseURI := uri
	if port != 0 {
		baseURI = fmt.Sprintf("%s:%d", baseURI, port)
	}
	return baseURI + path
}
