package parser

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
)

type ParsedLink struct {
	Link  string
	Deep  int
	Count int
	OK    bool
}

func (p *ParsedLink) String() string {
	data, _ := json.MarshalIndent(p, "", "")
	return string(data)
}

func (p *ParsedLink) GetRow() []string {
	// res := make([]string, 4)
	return []string{p.Link, strconv.Itoa(p.Deep), strconv.Itoa(p.Count), func(v bool) string {
		if v {
			return "ok"
		}
		return "err"
	}(p.OK)}
}

var res map[string]*ParsedLink

func init() {
	res = make(map[string]*ParsedLink)
}

func isPageURL(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}

	path := parsed.Path

	// no extension → likely page
	if !strings.Contains(path, ".") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))

	bad := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".svg":  true,
		".pdf":  true,
		".zip":  true,
		".mp4":  true,
		".mp3":  true,
		".css":  true,
		".js":   true,
	}

	return !bad[ext]
}

// https://chatgpt.com/g/g-p-678c0ac3548081918045c2cc6840396d/c/69d7dfcd-0eb0-838c-929b-f5cea6084d17

func SyncCase(ctx context.Context, links []string, dep int) (map[string]*ParsedLink, error) {
	if dep == 0 {
		return nil, nil
	}

	for i := range links {
		parsedLink, ok := res[links[i]]
		var linkCopy *ParsedLink
		if ok && parsedLink.Count > 0 {
			parsedLink.Count++
			continue
		}
		if !ok {
			parsedLink := ParsedLink{}
			parsedLink.Link = links[i]
			parsedLink.Deep = dep
			res[links[i]] = &parsedLink
			linkCopy = &parsedLink
		} else {
			linkCopy = parsedLink
		}
		linkCopy.Count++

		resp, err := http.Get(links[i])
		if err != nil {
			// fmt.Printf("Couldn't parse with get error %v. Continue\n", links[i])
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			linkCopy.OK = false
			// fmt.Printf("Couldn't parse with resp error: %v: %v. Contingitue\n", resp.StatusCode, links[i])
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			linkCopy.OK = false
			// fmt.Printf("Couldn't parse with get error %v. Continue\n", links[i])
			continue
		}
		linkCopy.OK = true

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

			SyncCase(ctx, newLink, dep-1)
		}

	}

	return res, nil
}
