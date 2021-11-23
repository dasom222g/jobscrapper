package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type ExtractJob struct {
	id          string
	title       string
	companyName string
	location    string
	salary      string
	summary     string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Request failed with status", res.StatusCode)
	}
}

func GetPages(url string) int {
	res, err := http.Get(url)

	// 유효성 체크
	checkError(err)
	checkCode((res))
	defer res.Body.Close()

	// Load the html file
	response, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	var pageLength = 0
	response.Find(".pagination").Each(func(index int, selection *goquery.Selection) {
		pageLength = selection.Find("a").Length()
	})
	return pageLength
}

func GetPage(url string, index int) []ExtractJob {
	var jobs = []ExtractJob{}
	resultUrl := url + "&start=" + strconv.Itoa(index*50)
	res, err := http.Get(resultUrl)
	checkError(err)
	checkCode(res)
	defer res.Body.Close()
	fmt.Println("resultUrl!!!!", resultUrl)

	// Load the html file
	document, err := goquery.NewDocumentFromReader(res.Body)
	document.Find(".tapItem").Each(func(index int, card *goquery.Selection) {
		job := extractedJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractedJob(card *goquery.Selection) ExtractJob {
	id, _ := card.Attr("data-jk")
	title := card.Find(".jobTitle>span").Text()
	companyName := card.Find(".companyName").Text()
	location := card.Find(".companyLocation").Text()
	summary := card.Find(".job-snippet").Text()
	return ExtractJob{
		id:          id,
		title:       title,
		companyName: companyName,
		location:    location,
		summary:     summary,
	}
}

func cleanString(str string) {
	fmt.Println("str", str)
}
