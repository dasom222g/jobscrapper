package main

import (
	"fmt"

	"github.com/dasom222g/jobscrapper/scrapper"
)

var baseUrl string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	pageLength := scrapper.GetPages(baseUrl)
	fmt.Println(pageLength)
	for i := 0; i < pageLength; i++ {
		scrapper.GetPage(baseUrl, i)
	}
}
