package cli

import (
	"context"

	"github.com/Mihai22125/sitemap/pkg/sitemap"
	"github.com/Mihai22125/sitemap/pkg/workerpool"
)

type PageJob struct {
	Url   string
	Depth int
}

func generateJob(page PageJob) workerpool.Job {
	return workerpool.Job{Descriptor: workerpool.JobDescriptor{ID: 1}, ExecFn: wrapper, Args: page}
}

func wrapper(ctx context.Context, pageJob interface{}) (interface{}, error) {
	return processPage(ctx, pageJob.(PageJob))
}

func processPage(ctx context.Context, pageJob PageJob) (sitemap.Page, error) {
	page, err := sitemap.ParsePage(pageJob.Url)
	if err != nil {
		return sitemap.Page{}, err
	}

	page.Depth = pageJob.Depth
	if page.Location != pageJob.Url {
		page.Depth = pageJob.Depth + 1
	}

	page.Priority = sitemap.PriorityMap[page.Depth]

	return page, nil
}
