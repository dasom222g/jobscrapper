package main

import (
	"fmt"

	"github.com/dasom222g/jobscrapper/scrapper"
)

var baseUrl string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
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
	scrapper.WriteCsv(totalJobs)
	fmt.Println("Done.", len(totalJobs))
	// for _, job := range totalJobs {
	// 	fmt.Println("job", job)
	// }
}
