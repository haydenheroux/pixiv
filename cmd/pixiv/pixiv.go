package main

import (
	"haydenheroux.xyz/pixivapi"
)

func main() {
	illustrations, _ := pixivapi.GetTopIllustrations()
	for _, illust := range illustrations {
		pixivapi.DownloadIllustration(illust)
	}
}
