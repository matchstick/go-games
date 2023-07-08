/* Package gopark is a place to play and learn about golang.
 * This webcrawler is an implementation of a solution from the GoLang Tour
 */
package gopark

import (
	"sync"
)

// Fetcher returns the body of URL and
// a slice of URLs found on that page.
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type elem struct {
	body   string
	count  int
	errMsg string
}

// Helper Type to all routine safe map access
type safeDB struct {
	m  sync.Mutex
	db map[string]elem
}

func (db *safeDB) isFetchedLocked(url string) bool {
	elem, present := db.db[url]
	if !present {
		return false
	}
	elem.count++
	db.db[url] = elem
	return true
}

func (db *safeDB) add(url string, body string, err error) {
	db.m.Lock()
	defer db.m.Unlock()

	// There is a race where we can have two routines
	// calling this on top of each other. It's a corner case but one we need to
	// handle.
	if db.isFetchedLocked(url) {
		return
	}

	if err != nil {
		db.db[url] = elem{body, 1, err.Error()}
	} else {
		db.db[url] = elem{body, 1, ""}
	}
}

func (db *safeDB) isFetched(url string) bool {
	db.m.Lock()
	defer db.m.Unlock()
	return db.isFetchedLocked(url)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, db *safeDB) {
	var wg sync.WaitGroup
	if depth <= 0 {
		return
	}

	if db.isFetched(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	db.add(url, body, err)

	if err != nil {
		return
	}

	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			Crawl(u, depth-1, fetcher, db)
		}(u)
	}
	wg.Wait()
}
