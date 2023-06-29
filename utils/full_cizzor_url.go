package utils

import (
	"fmt"
	"os"
)

func GetFullCizzorUrl(shortUrl string) string {
	baseUrl := os.Getenv("BASE_URL")
	url := baseUrl + "/" + shortUrl
	fmt.Println(url)
	return url
}
