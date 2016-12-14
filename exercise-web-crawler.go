package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type WebData struct {
	body string
	urls []string
	depth int
	err  error
}

// Cache struct for urls
type Cache struct {
	data map[string]bool
	mux  sync.Mutex
}

func (c *Cache) Set(url string) {
	c.mux.Lock()
	c.data[url] = true
	c.mux.Unlock()
}

func (c *Cache) Exist(url string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.data[url]
	return ok
}

var cache = Cache{data: make(map[string]bool)}

func StartCrawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan WebData)

	go Crawl(url, depth, fetcher, ch)

	totalCount := 1

	for i := 0; i < totalCount; i ++ {
		webData := <- ch

		if webData.depth > 0 {
			for _, subUrl := range webData.urls {
				// number of url searched increases
				totalCount += 1
				// search sub url
				go Crawl(subUrl, webData.depth, fetcher, ch)
			}
		}
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, ch chan WebData) {
	if !cache.Exist(url) {
		body, urls, err := fetcher.Fetch(url)

		cache.Set(url)

		if err != nil {
			fmt.Println(err)
			// stop searching
			ch <- WebData{depth: 0}
		} else {
			fmt.Printf("found: %s %q\n", url, body)
			// searching in another depth
			ch <- WebData{body, urls, depth - 1, err}
		}
	} else {
		fmt.Printf("had been already fetched: %s\n", url)
		// stop searching
		ch <- WebData{depth: 0}
	}
	return
}

func main() {
	StartCrawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
