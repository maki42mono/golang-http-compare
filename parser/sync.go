package parser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type noResultErr struct{}

func (e *noResultErr) Error() string {
	return "No url were paesed"
}

var res map[string]int
var curDep int

func init() {
	res = make(map[string]int)
	curDep = 0
}

// https://chatgpt.com/g/g-p-678c0ac3548081918045c2cc6840396d/c/69d7dfcd-0eb0-838c-929b-f5cea6084d17

func SyncCase(ctx context.Context, links []string, dep int) (any, error) {
	if curDep == dep {
		curDep--
		return nil, nil
	}

	re := regexp.MustCompile(`https?://[^\s"'<>]+`)

	for i := range links {
		if res[links[i]] > 0 {
			res[links[i]]++
			fmt.Printf("%v was already met %v timed. Continue\n", links[i], res[links[i]])
			continue
		}

		resp, err := http.Get(links[i])
		if err != nil {
			fmt.Printf("Couldn't parse with get error %v. Continue\n", links[i])
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("Couldn't parse with resp error: %v: %v. Continue\n", resp.StatusCode, links[i])
			continue
		}

		html, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Couldn't parse with get error %v. Continue\n", links[i])
			continue
		}

		curDep++
		newLinks := re.FindAllString(string(html), -1)
		SyncCase(ctx, newLinks, dep)
	}

	return res, nil
}
