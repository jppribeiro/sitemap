package internal

import (
	"fmt"
	"net/http"
	"sitemap/internal/parser"
	"time"
)

type SiteMapper struct {
	Url string
	Depth uint
	MaxDepth uint
	Pool *Executor
	Result chan Url

	expectedLeavesCounter chan int
	leafCounter chan int
}

func NewSiteMapper(url string, maxDepth uint, pool *Executor, resultChannel chan Url) *SiteMapper {
	return &SiteMapper{
		Url:                   url,
		Depth:                 0,
		MaxDepth:              maxDepth,
		Pool:                  pool,
		Result:                resultChannel,
		expectedLeavesCounter: make(chan int),
		leafCounter:           make(chan int),
	}
}

func newTaskFromSiteMapper(url string, depth uint, sm *SiteMapper) *SiteMapper {
	return &SiteMapper{
		Url:                   url,
		Depth:                 depth,
		MaxDepth:              sm.MaxDepth,
		Pool:                  sm.Pool,
		Result:                sm.Result,
		expectedLeavesCounter: sm.expectedLeavesCounter,
		leafCounter:           sm.leafCounter,
	}
}

func (sm *SiteMapper) Execute() {
	fmt.Printf("GET: %s", sm.Url)
	p, err := http.Get(sm.Url)

	fmt.Println("Success")
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	links := parser.GetLinks(*p)

	if sm.Depth == sm.MaxDepth - 1 {
		fmt.Printf("Expecting %d leaves", len(links))
		sm.expectedLeavesCounter <- len(links)
	}

	for _, l := range links {
		sm.Result <- Url{l}

		if sm.Depth < sm.MaxDepth {
			sm.Pool.Queue(newTaskFromSiteMapper(l, sm.Depth + 1, sm))
			continue
		}

		sm.leafCounter <- 1
	}
}

func (sm *SiteMapper) Wait() []Url {
	start := time.Now()

	var results []Url

	expectedLeaves := 0
	leaves := 0

	isStart := true

	for {
		if !isStart && expectedLeaves == leaves {
			return results
		}

		if time.Now().After(start.Add(time.Second * 10)) {
			return results
		}

		select {
		case url := <-sm.Result:
			fmt.Println("Append result")
			results = append(results, url)
		case n := <- sm.expectedLeavesCounter:
			fmt.Println("Expect leaves")
			expectedLeaves += n
		case n := <- sm.leafCounter:
			fmt.Println("Add leaves")
			isStart = false
			leaves += n
		}
	}
}