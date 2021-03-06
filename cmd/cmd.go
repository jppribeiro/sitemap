package cmd

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/url"
	"sitemap/internal"
)

type properties struct {
	parallelism uint
	maxDepth uint
	url string
}

type urlset struct {
	Url []internal.Url `xml:"url"`
}

func Start() {
	prop := config()

	fmt.Printf("Got config %v\n", prop)
	pool := internal.NewExecutor(prop.parallelism)

	u, err := url.Parse(prop.url)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mapper := internal.NewSiteMapper(u, prop.maxDepth, pool, make(chan internal.Url))

	err = pool.Queue(mapper)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Waiting\n")
	results := mapper.Wait()

	x, err := xml.MarshalIndent(urlset{Url: results}, "", "  ")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s\n", x)
}



func config() properties {
	p := flag.Uint("parallelism", 10, "Set max number of workers to transverse the website.")
	m := flag.Uint("max-depth", 1, "Max depth of the site tree to transverse.")

	flag.Parse()

	return properties{
		parallelism: *p,
		maxDepth : *m,
		url: flag.Arg(0),
	}
}
