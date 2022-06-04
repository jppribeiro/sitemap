package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"sitemap/internal/parser"
	"strings"
	"sync/atomic"
	"time"
)

var client = http.Client{Timeout: time.Second * 10}

type SiteMapper struct {
	Url *url.URL
	Depth uint
	MaxDepth uint
	Pool *Executor
	Result chan Url

	node chan int32
}

func NewSiteMapper(url *url.URL, maxDepth uint, pool *Executor, resultChannel chan Url) *SiteMapper {
	return &SiteMapper{
		Url:                   url,
		Depth:                 0,
		MaxDepth:              maxDepth,
		Pool:                  pool,
		Result:                resultChannel,

		node: make(chan int32),
	}
}

func newTaskFromSiteMapper(url *url.URL, depth uint, sm *SiteMapper) *SiteMapper {
	return &SiteMapper{
		Url:                   url,
		Depth:                 depth,
		MaxDepth:              sm.MaxDepth,
		Pool:                  sm.Pool,
		Result:                sm.Result,

		node: sm.node,
	}
}

func (sm *SiteMapper) Execute(r uint) {
	sm.node <- 1

	defer func() {
		sm.node <- -1
	}()

	//fmt.Printf("GET: %s\n", sm.Url)

	p, err := http.Get(sm.Url.String())

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Pass current url to the results channel
	sm.Result <- Url{p.Request.URL.String()}

	if sm.Depth >= sm.MaxDepth {
		return
	}

	links := parser.GetLinks(*p)

	for _, l := range links {
		fmt.Println(l.String())
		if isHost(p.Request.URL.Host, l.Host) {
			fmt.Printf("Not on host: %s, %s\n", l.String(), p.Request.URL.Host)
			continue
		}

		task := newTaskFromSiteMapper(l, sm.Depth + 1, sm)

		err = sm.Pool.Queue(task)

		if err != nil {
			fmt.Println("Pool is starving.")
			task.Execute(r)
		}
	}
}

func (sm *SiteMapper) Wait() []Url {
	start := time.Now()

	var visited = make(map[string]int)

	var results []Url

	var nodes int32

	isStart := false

	for {
		if time.Now().Sub(start) > time.Second * 10 {
			return results
		}

		if isStart && nodes == 0 {
			return results
		}

		select {
		case u := <-sm.Result:
			if visited[u.Loc] == 0 {
				visited[u.Loc] += 1
				results = append(results, u)
			}

			break
		case n := <- sm.node:
			isStart = true
			atomic.AddInt32(&nodes, n)
		}
	}
}

func isHost(h string, c string) bool {
	splitHost := strings.Split(h, ".")

	h = fmt.Sprintf("%s.%s", splitHost[len(splitHost)-2], splitHost[len(splitHost)-1])

	return h == c
}