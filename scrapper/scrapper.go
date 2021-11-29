package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
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

var header = []string{"id", "title", "companyName", "location", "salary", "summary"}

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
	fmt.Println("requesting! ", resultUrl)

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
	salary := card.Find(".salary-snippet").Text()
	summary := card.Find(".job-snippet").Text()
	return ExtractJob{
		id:          id,
		title:       title,
		companyName: companyName,
		location:    location,
		salary:      salary,
		summary:     summary,
	}
}

func cleanString(str string) {
	fmt.Println("str", str)
}

func WriteCsv(jobs []ExtractJob) {
	file, err := os.Create("jobs.csv") // file생성
	checkError(err)

	w := csv.NewWriter(file) // writer생성
	defer w.Flush()          // file에 데이터 입력 및 저장

	// 데이터 작성
	wErr := w.Write(header)
	checkError(wErr)

	for _, job := range jobs {
		link := "https://kr.indeed.com/viewjob?jk=" + job.id
		jobSlice := []string{link, job.title, job.companyName, job.location, job.salary, job.summary}
		jobWErr := w.Write(jobSlice)
		checkError(jobWErr)
	}
}
