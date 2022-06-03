package internal

import (
	"fmt"
	"net/http"
	"sitemap/internal/parser"
)

func Transverse(url string, maxDepth int, currentDepth int) []Url {
	fmt.Println(maxDepth)
	fmt.Println(currentDepth)
	if currentDepth > maxDepth {
		return nil
	}

	p, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	links := parser.GetLinks(*p)

	var childrenUrl []Url

	for _, l := range links {
		childTree := Transverse(l, maxDepth, currentDepth + 1)

		childrenUrl = append(childrenUrl, Url{l})

		childrenUrl = append(childrenUrl, childTree...)
	}

	return append([]Url{{url}}, childrenUrl...)
}
