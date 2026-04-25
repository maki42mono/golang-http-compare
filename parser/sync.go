package parser

import (
	"context"
	"fmt"
	"time"
)

func SyncCall(ctx context.Context) (any, error) {
	for i := 1; i <= 8; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Context was cancelled before I finished. So I finish with an error =))0")
			return nil, ctx.Err()
		case <-time.After(1000 * time.Millisecond):
			fmt.Printf("step %d finished\n", i)
		}
	}

	return "Sync was done!", nil
}
