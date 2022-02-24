package pixivapi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func dateTimeFromString(date string) (string, string, string, string, string, string) {
	splitDateTime := strings.Split(date, "T")
	dateFragment := splitDateTime[0]
	dateSlice := strings.Split(dateFragment, "-")
	year := dateSlice[0]
	month := dateSlice[1]
	day := dateSlice[2]
	timeFragment := splitDateTime[1]
	justTimeFragment := strings.Split(timeFragment, "+")[0]
	timeSlice := strings.Split(justTimeFragment, ":")
	hour := timeSlice[0]
	minute := timeSlice[1]
	second := timeSlice[2]
	return year, month, day, hour, minute, second
}

func GetIllustrationDownloadURL(illustration PixivIllustration) string {
	id := illustration.Id
	date := illustration.UpdateDate
	year, month, day, hour, minute, second := dateTimeFromString(date)
	url := fmt.Sprintf("https://i.pximg.net/img-master/img/%s/%s/%s/%s/%s/%s/%s_p0_master1200.jpg", year, month, day, hour, minute, second, id)
	return url
}

func DownloadIllustration(illustration PixivIllustration) (int64, error) {
	url := GetIllustrationDownloadURL(illustration)
	fileName := "test-" + illustration.Title
	output, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	downloadHeaders := map[string]string{
		"Referer": "https://www.pixiv.net",
	}
	client, req, err := setupGet(url, downloadHeaders)
	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	n, err := io.Copy(output, res.Body)
	if err != nil {
		return n, err
	}

	return n, nil
}
