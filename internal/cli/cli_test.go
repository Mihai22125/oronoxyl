package cli

import (
	"flag"
	"testing"
)

func TestFromArgs(t *testing.T) {
	var app appEnv

	args := []string{"-url", "http://example.com", "-output-file", "example.xml", "-parallel", "3", "-max-depth", "3"}

	err := app.fromArgs(args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if app.url != "http://example.com" {
		t.Errorf("Expected URL to be http://example.com, got %s", app.url)
	}
	if app.outputFile != "example.xml" {
		t.Errorf("Expected output to be example.xml, got %s", app.outputFile)
	}
	if app.parallelWorkers != 3 {
		t.Errorf("Expected parallel workers to be 3, got %d", app.parallelWorkers)
	}
	if app.maxDepth != 3 {
		t.Errorf("Expected max depth to be 3, got %d", app.maxDepth)
	}
}

func TestValidate(t *testing.T) {
	var app appEnv

	testData := []struct {
		args          []string
		expectedError error
	}{
		{[]string{"-url", "http://example.com", "-output-file", "example.xml", "-parallel", "3", "-max-depth", "3"}, nil},
		{[]string{"-url", "example.com", "-output-file", "example.xml", "-parallel", "3", "-max-depth", "3"}, flag.ErrHelp},
		{[]string{"-url", "http://example.com", "-output-file", "example.pdf", "-parallel", "3", "-max-depth", "3"}, flag.ErrHelp},
		{[]string{"-url", "http://example.com", "-output-file", "example.xml", "-parallel", "0", "-max-depth", "3"}, flag.ErrHelp},
		{[]string{"-url", "http://example.com", "-output-file", "example.xml", "-parallel", "1", "-max-depth", "0"}, flag.ErrHelp},
	}

	for _, test := range testData {
		err := app.fromArgs(test.args)
		if err != test.expectedError {
			t.Errorf("Expected error %v, got %v", test.expectedError, err)
		}
	}
}

func TestGenerateJob(t *testing.T) {
	job := generateJob(PageJob{Url: "http://example.com", Depth: 1})
	args := job.Args.(PageJob)

	if args.Url != "http://example.com" {
		t.Errorf("Expected URL to be http://example.com, got %s", args.Url)
	}
	if args.Depth != 1 {
		t.Errorf("Expected depth to be 1, got %d", args.Depth)
	}

}
