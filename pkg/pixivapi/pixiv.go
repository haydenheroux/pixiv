package pixivapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// GetTopIllustrations spoofs Pixiv's AJAX requests to query the top illustrations
// at the time of the request.
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

func GetSearchIllustrations(query string) ([]PixivIllustration, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/search/artworks/%s", query)
	client, req, _ := setupGet(url, defaultHeaders)
	res, err := client.Do(req)
	if err != nil {
		return make([]PixivIllustration, 0), err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return make([]PixivIllustration, 0), err
	}
	var body PixivSearchResponseBody
	json.Unmarshal([]byte(content), &body)
	illustrations := body.MangaIllustrations.Illusts.IllustrationList
	return illustrations, nil
}
