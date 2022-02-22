package main

import (
	"fmt"

	"haydenheroux.xyz/pixivapi"
)

func main() {
	illustrations, _ := pixivapi.GetTopIllustrations()
	for _, illust := range illustrations {
		fmt.Printf("%+v\n", illust)
	}
}
