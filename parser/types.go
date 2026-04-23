package parser

import "context"

type TestCase func(ctx context.Context) (any, error)
