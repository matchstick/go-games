package gopark

/* crawl test and program lifted from the golang tour */
import (
	"fmt"
	"testing"
)

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
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func verifyDB(t *testing.T, db *safeDB) {
	var expected = map[string]elem{
		"https://golang.org/":         {"The Go Programming Language", 4, ""},
		"https://golang.org/pkg/":     {"Packages", 3, ""},
		"https://golang.org/cmd/":     {"", 2, "not found: https://golang.org/cmd/"},
		"https://golang.org/pkg/os/":  {"Package os", 1, ""},
		"https://golang.org/pkg/fmt/": {"Package fmt", 1, ""},
	}

	for k, e := range expected {
		testElem, ok := db.db[k]
		if !ok {
			t.Errorf("%s is not present.", k)
			continue
		}

		if e != testElem {
			t.Errorf("Elements don't match for key %s", k)
			continue
		}
	}
}

func TestCrawl(t *testing.T) {
	db := safeDB{db: make(map[string]elem)}
	Crawl("https://golang.org/", 4, fetcher, &db)
	verifyDB(t, &db)
}
