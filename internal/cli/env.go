package cli

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type appEnv struct {
	url             string
	parallelWorkers int
	outputFile      string
	maxDepth        int
	verbose         bool
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("sitemap-gen", flag.PanicOnError)
	fl.StringVar(&app.url, "url", "", "site url for sitemap generation")
	fl.IntVar(&app.parallelWorkers, "parallel", 3, "number of parallel workers to navigate through site (Default 3)")
	fl.StringVar(&app.outputFile, "output-file", "./temp.xml", "output file path")
	fl.IntVar(&app.maxDepth, "max-depth", 3, "max depth of url navigation recursion")
	fl.BoolVar(&app.verbose, "verbose", true, "display detailed processing information")
	fl.Parse(args)

	if err := app.validate(); err != nil {
		fl.Usage()
		return err
	}

	return nil
}

func (app *appEnv) validate() error {
	u, err := url.Parse(app.url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Fprintln(os.Stderr, "the provided url is not valid")
		return flag.ErrHelp
	}

	fileExtension := filepath.Ext(app.outputFile)
	if fileExtension != ".xml" {
		fmt.Fprintln(os.Stderr, "File extension ins't equal to .xml")
		return flag.ErrHelp
	}

	if app.parallelWorkers < 1 {
		fmt.Fprintln(os.Stderr, "Number of parallel workers cant be smaller than 1")
		return flag.ErrHelp
	}

	if app.maxDepth < 1 {
		fmt.Fprintln(os.Stderr, "Maximum Depth cant be smaller than 1")
		return flag.ErrHelp
	}

	return nil
}
