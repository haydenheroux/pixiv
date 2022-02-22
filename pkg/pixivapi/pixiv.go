package pixivapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func setupGet(url string, headers map[string]string) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return client, nil, err
	}

	for header := range headers {
		req.Header.Add(header, headers[header])
	}
	return client, req, nil
}

func GetTopIllustrations() ([]PixivIllustration, error) {

	client, req, _ := setupGet("https://www.pixiv.net/ajax/top/illust?mode=all&lang=en", defaultHeaders)
	res, err := client.Do(req)
	if err != nil {
		return make([]PixivIllustration, 0), err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return make([]PixivIllustration, 0), err
	}
	var body PixivResponseBody
	json.Unmarshal([]byte(content), &body)
	illustrations := body.Thumbnails.Illusts.IllustrationList
	return illustrations, nil
}
