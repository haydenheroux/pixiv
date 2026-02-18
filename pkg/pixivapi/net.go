package pixivapi

import "net/http"

var client = http.Client{}

// setupGet creates an http GET request with the supplied headers.
func setupGet(url string, headers map[string]string) (*http.Client, *http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &client, nil, err
	}

	for header := range headers {
		req.Header.Add(header, headers[header])
	}
	return &client, req, nil
}
