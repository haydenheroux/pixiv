package pixivapi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// dateTimeFromString converts a string representing a date and time
// with the format YY-MM-DDTHH:MM:SS+TZHH:TZMM into its constituents, namely into year, month, day, hour, minute, and second components.
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

// GetIllustrationDownloadURL provides a URL containing the image contents
// related to the Pixiv illustration passed as an argument. Pixiv stores
// actual images on the pximg.net domain under the "i" subdomain, and within
// this subdomain there exists two routes which pertain the images themselves.
// The "img-original" route appears to store the images as uploaded by the artist,
// keeping the original format. The "img-master" route appears to store the images
// Pixiv generates to represent the image, such as alterations to the size.
//
// The "img-master" route is used within this function because all of the images
// within the scope of a basic search are saved in a JPEG format, which entails
// the elimination of dynamically finding the original format. While there is a
// strong argument to be made that the original image in the original format is the
// ideal way to view the image, for sake of programming this function cares only about
// JPEG images.
func GetIllustrationDownloadURL(illustration PixivIllustration) string {
	id := illustration.Id
	date := illustration.UpdateDate
	year, month, day, hour, minute, second := dateTimeFromString(date)
	url := fmt.Sprintf("https://i.pximg.net/img-master/img/%s/%s/%s/%s/%s/%s/%s_p0_master1200.jpg", year, month, day, hour, minute, second, id)
	return url
}

// DownloadIllustration downloads the Pixiv illustration passed as an argument to
// the location specified in the filePath argument.
func DownloadIllustration(illustration PixivIllustration, filePath string) (int64, error) {
	url := GetIllustrationDownloadURL(illustration)
	output, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	downloadHeaders := map[string]string{
		// When Pixiv requests the image from the pximg.net domain,
		// the only header pximg.net requires for access (e.g. not error 403)
		// is the Referer header referencing Pixiv.net
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
