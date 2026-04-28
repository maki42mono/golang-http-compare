package parser

import (
	"context"
	"net/http"
	"strings"
	"sync"

	goquery "github.com/PuerkitoBio/goquery"
)

var async_res map[string]*ParsedLink

func init() {
	async_res = make(map[string]*ParsedLink)
}

func craw(ctx context.Context, links []string, dep int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	if dep == 0 {
		return
	}

	for i := range links {
		mu.Lock()
		parsedLink, ok := async_res[links[i]]
		var linkCopy *ParsedLink
		if ok && parsedLink.Count > 0 {
			parsedLink.Count++
			mu.Unlock()
			continue
		}
		if !ok {
			parsedLink := ParsedLink{}
			parsedLink.Link = links[i]
			parsedLink.Deep = dep
			async_res[links[i]] = &parsedLink
			linkCopy = &parsedLink
		} else {
			linkCopy = parsedLink
		}
		linkCopy.Count++
		mu.Unlock()

		resp, err := http.Get(links[i])
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		mu.Lock()
		if resp.StatusCode != 200 {
			linkCopy.OK = false
			mu.Unlock()
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			linkCopy.OK = false
			mu.Unlock()
			continue
		}
		linkCopy.OK = true
		mu.Unlock()

		var newLinks []string
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			val, ok := s.Attr("href")
			if !ok {
				return
			}
			newLinks = append(newLinks, val)
		})

		// fmt.Printf("Parsing %v\n", links[i])
		for j := range newLinks {
			if !isPageURL(newLinks[j]) {
				continue
			}
			newLink := make([]string, 1)
			newLink[0] = newLinks[j]
			if strings.HasPrefix(newLink[0], "mailto:") {
				continue
			}
			if !strings.HasPrefix(newLink[0], "http") {
				newLink[0] = links[i] + newLink[0]
			}
			wg.Add(1)
			go craw(ctx, newLink, dep-1, wg, mu)
		}
	}
}

func AsyncCase(ctx context.Context, links []string, dep int) (map[string]*ParsedLink, error) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	craw(ctx, links, dep, &wg, &mu)
	wg.Wait()

	return async_res, nil
}
