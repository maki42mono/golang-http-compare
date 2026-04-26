package parser

import "context"

type TestCase func(ctx context.Context, links []string, dep int) (any, error)
