package main

import (
	"os"
	"runtime"

	"github.com/Mihai22125/sitemap/internal/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(cli.CLI(os.Args[1:]))
}
