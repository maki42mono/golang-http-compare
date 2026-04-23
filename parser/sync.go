package parser

import (
	"context"
	"fmt"
	"time"
)

func SyncCall(ctx context.Context) (any, error) {
	for i := 1; i <= 3; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Context was canselled before I finished. So I finish with an error =))0")
			return nil, ctx.Err()
		case <-time.After(300 * time.Millisecond):
			fmt.Printf("step %d finished\n", i)
			i++
		}
	}

	return "Sync was done!", nil
}
