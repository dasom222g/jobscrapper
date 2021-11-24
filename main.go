package main

import (
	"fmt"

	"github.com/dasom222g/jobscrapper/scrapper"
)

var baseUrl string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	pageLength := scrapper.GetPages(baseUrl)
	var totalJobs []scrapper.ExtractJob
	for i := 0; i < pageLength; i++ {
		extractedJobs := scrapper.GetPage(baseUrl, i)
		totalJobs = append(totalJobs, extractedJobs...)
	}
	fmt.Println("Total jobs", len(totalJobs))
	for _, job := range totalJobs {
		fmt.Println("job", job)
	}
}
