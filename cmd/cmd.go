package cmd

import (
	"encoding/xml"
	"flag"
	"fmt"
	"sitemap/internal"
)

type properties struct {
	parallelism int
	maxDepth int
	url string
}

type urlset struct {
	Url []internal.Url `xml:"url"`
}

func Execute() {
	prop := config()

	sitemap := urlset{Url: internal.Transverse(prop.url, prop.maxDepth, 0)}

	fmt.Printf("%+v\n", sitemap)
	x, err := xml.MarshalIndent(sitemap, "", "  ")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s\n", x)
}

func config() properties {
	p := flag.Int("parallelism", 1, "Set max number of workers to transverse the website.")
	m := flag.Int("max-depth", 1, "Max depth of the site tree to transverse.")

	flag.Parse()

	return properties{
		parallelism: *p,
		maxDepth : *m,
		url: flag.Arg(0),
	}
}
