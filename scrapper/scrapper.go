package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

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

func GetPage(url string, index int) {
	fmt.Println("===================================")
	resultUrl := url + "&start=" + strconv.Itoa(index*50)
	res, err := http.Get(resultUrl)
	checkError(err)
	checkCode(res)
	defer res.Body.Close()

	// Load the html file
	document, err := goquery.NewDocumentFromReader(res.Body)
	document.Find(".tapItem").Each(func(index int, selection *goquery.Selection) {
		fmt.Println("index", index)
		card, err := selection.Html()
		checkError(err)

		fmt.Println("card", card)
	})
}
