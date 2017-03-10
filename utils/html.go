package utils

import (
	"io/ioutil"
	"regexp"
	"strings"

	readability "github.com/mauidude/go-readability"
)

var (
	htmlCleanRe  = regexp.MustCompile(`<[^>]*>`)
	htmlReplaces = []struct {
		from string
		to   string
	}{
		{"&#39;", "'"},
	}
)

func GetPreview(url string) (string, error) {
	resp, err := MakeRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return cleanHTML(string(body))
}

func cleanHTML(html string) (string, error) {
	doc, err := readability.NewDocument(html)
	if err != nil {
		return "", err
	}

	content := doc.Content()
	for _, r := range htmlReplaces {
		content = strings.Replace(content, r.from, r.to, -1)
	}

	return htmlCleanRe.ReplaceAllString(content, ""), nil
}
