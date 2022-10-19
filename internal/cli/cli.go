package cli

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/Mihai22125/sitemap/pkg/sitemap"
	"github.com/Mihai22125/sitemap/pkg/workerpool"
)

func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}

	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}

	return 0
}

func (app *appEnv) run() error {
	wp := workerpool.New(app.parallelWorkers)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx)

	wp.GenerateFromJob(generateJob(PageJob{Url: app.url, Depth: 1}))

	seen := make(map[string]bool)
	var processed int

	file, err := os.Create(app.outputFile)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	writer.Write([]byte("<urlset>\n"))

	start := time.Now()
	defer func() {
		writer.Write([]byte("\n</urlset>"))
		writer.Flush()
		file.Close()
		if app.verbose {
			fmt.Fprintf(os.Stderr, "\nTime finished sitemap %s\n", time.Since(start))
		}
	}()

	for {
		if app.verbose {
			fmt.Fprintf(os.Stderr, "\rURLs Found: %5d\t\t Pages Processed: %5d\t\t Queue: %5d", processed+wp.GetQueueSize(), processed, wp.GetQueueSize())
		}
		if wp.Working == 0 {
			wp.CloseJobsChannel()
		}
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			wp.Working--

			if r.Err != nil {
				continue
			}

			page := r.Value.(sitemap.Page)
			data, err := xml.MarshalIndent(page, " ", "  ")
			if err != nil {
				if app.verbose {
					fmt.Fprintf(os.Stderr, "An error occured while Marshling to XML: %v\n", err)
				}
			}
			writer.Write(data)
			processed++

			if page.Depth < app.maxDepth {

				for _, link := range page.Links {
					if !seen[link] {
						seen[link] = true
						wp.GenerateFromJob(generateJob(PageJob{Url: link, Depth: page.Depth + 1}))
					}
				}
			}

		case <-wp.Done:
			return nil
		default:
		}
	}
}
