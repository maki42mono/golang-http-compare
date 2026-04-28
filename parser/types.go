package parser

import (
	"context"
	"encoding/json"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type TestCase func(ctx context.Context, links []string, dep int) (map[string]*ParsedLink, error)

type ParsedLink struct {
	Link  string
	Deep  int
	Count int
	OK    bool
}

type Printable interface {
	String() string
}

type CsvRow interface {
	GetRow() []string
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
