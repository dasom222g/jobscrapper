package scrape

import (
	"fmt"

	"github.com/dasom222g/jobscrapper/scrape/scrapper"
)

func Scrape(keyword string) {
	baseUrl := "https://kr.indeed.com/jobs?q=" + keyword + "&limit=50"
	pageLength := scrapper.GetPages(baseUrl)
	var totalJobs []scrapper.ExtractJob

	// channel
	mainChannel := make(chan []scrapper.ExtractJob)
	for i := 0; i < pageLength; i++ {
		go scrapper.GetPage(baseUrl, i, mainChannel)
	}
	for i := 0; i < pageLength; i++ {
		extractedJobs := <-mainChannel
		totalJobs = append(totalJobs, extractedJobs...)
	}
	scrapper.WriteCsv(keyword, totalJobs)
	fmt.Println("Done.", len(totalJobs))
}
