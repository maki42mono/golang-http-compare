package speedtest

import (
	"fmt"
	"speedtest/parser"
)

var registry = map[string]parser.TestCase{
	"sync": parser.SyncCall,
}

func main() {
	fmt.Println("Test1")
}
