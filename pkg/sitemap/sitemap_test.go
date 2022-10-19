package sitemap

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGetLastUpdatedDate(t *testing.T) {
	mockResponse := http.Response{}
	mockResponse.Header = make(http.Header)

	testTable := []struct {
		lastModified string

		expected    time.Time
		expectedErr error
	}{
		{"", time.Time{}, ErrHeaderValueNotFound},
		{"Tue, 15 Nov 1994 08:12:31 GMT", time.Date(1994, 11, 15, 8, 12, 31, 0, time.UTC), nil},
		{"Tue, 15 Nov 1994 08:12:31 INVALID", time.Time{}, errors.New(`parsing time "Tue, 15 Nov 1994 08:12:31 INVALID" as "Mon, 02 Jan 2006 15:04:05 MST": cannot parse "INVALID" as "MST"`)},
	}

	for _, test := range testTable {
		mockResponse.Header.Set("Last-Modified", test.lastModified)
		date, err := GetLastUpdatedDate(&mockResponse)
		if err != nil && err.Error() != test.expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", test.expectedErr, err)
		}
		date = date.UTC()
		if date.UTC() != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, date)
		}
	}
}

func TestGetLinks(t *testing.T) {
	mockResponse := http.Response{}
	mockResponse.Header = make(http.Header)

	testTable := []struct {
		hostname string
		root     string
		html     string
		expected []string
	}{
		{"", "", "", []string{}},
		{"example.com", "http://example.com", `<a href="/">Home</a>`, []string{"http://example.com/"}},
		{"example.com", "http://example.com", `<a href="//example.com">Home</a>`, []string{"//example.com"}},
		{"example.com", "http://www.example.com", `<a href="http://www.example.com/home">Home</a>`, []string{"http://www.example.com/home"}},
	}

	for _, test := range testTable {
		testUrl, _ := url.Parse(test.root)

		mockResponse.Request = &http.Request{URL: testUrl}
		mockResponse.Header.Set("Host", test.hostname)
		mockResponse.Header.Set("Root", test.root)
		mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(test.html))
		links, err := GetLinks(&mockResponse)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(links) != len(test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, links)
		}
		for i, link := range links {
			if link != test.expected[i] {
				t.Errorf("Expected %v, got %v", test.expected[i], link)
			}
		}
	}
}

func TestExtractData(t *testing.T) {
	mockResponse := http.Response{}
	mockResponse.Header = make(http.Header)

	testTable := []struct {
		hostname string
		root     string
		html     string
		expected Page
	}{
		{"", "", "", Page{LastModified: &time.Time{}}},
		{"example.com", "http://example.com", `<a href="/">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{"http://example.com/"}}},
		{"example.com", "http://example.com", `<a href="//example.com">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{"//example.com"}}},
		{"example.com", "http://www.example.com", `<a href="http://www.example.com/home">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{"http://www.example.com/home"}}},
		{"example.com", "http://www.example.com", `<a href="http://www.example.com/home.jpg">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{}}},
		{"example.com", "http://www.example.com", `<a href="mailto:http://www.example.com/home.jpg">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{}}},
		{"example.com", "http://www.example", `<head><base href="http://www.example.com/"></head><a href="home">Home</a>`, Page{LastModified: &time.Time{}, Links: []string{}}},
	}

	for _, test := range testTable {
		testUrl, _ := url.Parse(test.root)

		mockResponse.Request = &http.Request{URL: testUrl}
		mockResponse.Header.Set("Host", test.hostname)
		mockResponse.Header.Set("Root", test.root)
		mockResponse.Body = ioutil.NopCloser(bytes.NewBufferString(test.html))
		page, err := extractData(&mockResponse, test.root)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if page.Priority != test.expected.Priority {
			t.Errorf("Expected %v, got %v", test.expected.Priority, page.Priority)
		}

		if !page.LastModified.Equal(*test.expected.LastModified) {
			t.Errorf("Expected %v, got %v", test.expected.LastModified, page.LastModified)
		}

		if page.ChangeFrequency != test.expected.ChangeFrequency {
			t.Errorf("Expected %v, got %v", test.expected.ChangeFrequency, page.ChangeFrequency)
		}

		if len(page.Links) != len(test.expected.Links) {
			t.Errorf("Expected %v, got %v", test.expected.Links, page.Links)
			return
		}

		for i, link := range page.Links {
			if link != test.expected.Links[i] {
				t.Errorf("Expected %v, got %v", test.expected.Links[i], link)
			}
		}

	}
}

func TestFrequency_String(t *testing.T) {
	if Always.String() != "Always" {
		t.Errorf("Expected 'Always', got '%s'", Always.String())
	}
	if Hourly.String() != "Hourly" {
		t.Errorf("Expected 'Hourly', got '%s'", Hourly.String())
	}
	if Daily.String() != "Daily" {
		t.Errorf("Expected 'Daily', got '%s'", Daily.String())
	}
	if Weekly.String() != "Weekly" {
		t.Errorf("Expected 'Weekly', got '%s'", Weekly.String())
	}
	if Monthly.String() != "Monthly" {
		t.Errorf("Expected 'Monthly', got '%s'", Monthly.String())
	}
	if Yearly.String() != "Yearly" {
		t.Errorf("Expected 'Yearly', got '%s'", Yearly.String())
	}
	if Never.String() != "Never" {
		t.Errorf("Expected 'Never', got '%s'", Never.String())
	}
}

func TestParsePage(t *testing.T) {
	testTable := []struct {
		URL      string
		expected error
	}{
		{`https://www.example.com`, nil},
		{`www.example.com`, errors.New(`Get "www.example.com": unsupported protocol scheme ""`)},
	}

	for _, test := range testTable {
		_, err := ParsePage(test.URL)
		if err != nil && test.expected == nil {
			t.Errorf("Expected %v, got %v", test.expected, err)
			return
		}

		if err != nil && err.Error() != test.expected.Error() {
			t.Errorf("Expected %v, got %v", test.expected, err)
		}
	}
}
