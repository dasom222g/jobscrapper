package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

var header = []string{"Id", "Title", "CompanyName", "Location", "Salary", "Summary"}

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

func GetPage(url string, index int, mainChannel chan<- []ExtractJob) {
	var jobs = []ExtractJob{}
	resultUrl := url + "&start=" + strconv.Itoa(index*50)
	res, err := http.Get(resultUrl)
	checkError(err)
	checkCode(res)
	defer res.Body.Close()
	fmt.Println("requesting! ", resultUrl)

	// channel
	channel := make(chan ExtractJob)

	// Load the html file
	document, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	cards := document.Find(".tapItem")
	cards.Each(func(index int, card *goquery.Selection) {
		go extractedJob(card, channel)
	})

	for i := 0; i < cards.Length(); i++ {
		job := <-channel
		jobs = append(jobs, job)
	}

	mainChannel <- jobs
}

func extractedJob(card *goquery.Selection, channel chan<- ExtractJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle>span").Text())
	companyName := cleanString(card.Find(".companyName").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	channel <- ExtractJob{
		id:          id,
		title:       title,
		companyName: companyName,
		location:    location,
		salary:      salary,
		summary:     summary,
	}
}

func cleanString(str string) string {
	return strings.TrimSpace(str)
}

func WriteCsv(fileName string, jobs []ExtractJob) {
	file, err := os.Create(fileName + "_jobs.csv") // file생성
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
