package parser

import "context"

type TestCase func(ctx context.Context, links []string, dep int) (map[string]*ParsedLink, error)

type Printable interface {
	String() string
}

type CsvRow interface {
	GetRow() []string
}
