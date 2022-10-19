package sitemap

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var ErrHeaderValueNotFound = errors.New("last-moodified Header value not found")

var basePattern = regexp.MustCompile(`<base[\s\S]*?href="([^"]+)"[\s\S]*?>`)
var hrefPattern = regexp.MustCompile(`<a[\s\S]*?href="([^"]+)"[\s\S]*?>`)

func GetLastUpdatedDate(resp *http.Response) (time.Time, error) {
	lastModified := resp.Header.Get("last-modified")
	if lastModified == "" {
		return time.Time{}, ErrHeaderValueNotFound
	}

	date, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func GetLinks(resp *http.Response) ([]string, error) {
	hostname := resp.Request.URL.Hostname()
	root := resp.Request.URL.Scheme + "://" + resp.Request.URL.Hostname()

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	baseMatches := basePattern.FindStringSubmatch(string(html))

	baseUrl := ""
	if len(baseMatches) == 2 {
		baseUrl = baseMatches[1]
	}

	matches := hrefPattern.FindAllStringSubmatch(string(html), -1)

	var links []string

	for _, match := range matches {
		foundLink := SanitizeUrl(match[1])

		if strings.HasPrefix(foundLink, "//") {
		} else if strings.HasPrefix(foundLink, "/") {
			foundLink = root + foundLink
		} else if !strings.HasPrefix(strings.ToLower(foundLink), "http") {
			foundLink = baseUrl + foundLink
		}

		if isValidLink(foundLink, hostname) {
			links = append(links, foundLink)

		}
	}

	return links, nil
}

func doRequest(url string) (*http.Response, error) {
	return http.Get(url)

}

func extractData(resp *http.Response, URL string) (Page, error) {
	baseUrl, err := url.Parse(URL)
	if err != nil {
		return Page{}, err
	}

	if resp.Request.URL.Hostname() != baseUrl.Hostname() {
		return Page{}, errors.New("not same host")
	}

	lastModified, err := GetLastUpdatedDate(resp)
	if err != nil && err != ErrHeaderValueNotFound {
		return Page{}, err
	}

	links, err := GetLinks(resp)
	if err != nil {
		return Page{}, err
	}

	page := Page{
		Location:     URL,
		LastModified: &lastModified,
		Links:        links,
	}

	return page, nil
}

func ParsePage(URL string) (Page, error) {

	resp, err := doRequest(URL)
	if err != nil {
		return Page{}, err
	}

	return extractData(resp, URL)
}

func SanitizeUrl(link string) string {

	for _, fal := range FalseUrls {
		if strings.Contains(link, fal) {
			return ""
		}
	}

	link = strings.TrimSpace(link)
	tram := strings.Split(link, "#")[0]

	return tram
}

func isInternLink(link, root string) bool {

	url, err := url.Parse(link)
	if err != nil {
		return false
	}

	return url.Hostname() == root
}

func isStart(link, root string) bool {
	return strings.Compare(link, root) == 0
}

func isValidExtension(link string) bool {
	preparedLink := strings.TrimSuffix(strings.ToLower(link), "/")
	for _, extension := range Extensions {
		if strings.HasSuffix(preparedLink, extension) {
			return false
		}
	}
	return true
}

func isValidLink(link, root string) bool {

	if isInternLink(link, root) && !isStart(link, root) && isValidExtension(link) {
		return true
	}

	return false
}
