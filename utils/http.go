package utils

import (
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(method, url string, body io.Reader) (*http.Response, error) {
	client := new(http.Client)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("gonews:v%s (by /u/drgarcia1986)", Version))
	return client.Do(req)
}
