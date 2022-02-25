package main

import (
	"fmt"

	"haydenheroux.xyz/pixivapi"
)

func main() {
	illustrations, _ := pixivapi.GetTopIllustrations()
	for n, illust := range illustrations {
		fmt.Printf("Downloading %d / %d...\n", n, len(illustrations))
		pixivapi.DownloadIllustration(illust, "test-"+illust.Title+".jpg")
	}
}
