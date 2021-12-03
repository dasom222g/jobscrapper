package main

import (

	// "github.com/dasom222g/jobscrapper/scrape"
	"os"
	"strings"

	"github.com/dasom222g/jobscrapper/scrape"
	"github.com/dasom222g/jobscrapper/scrape/scrapper"
	"github.com/labstack/echo"
)

func removeFile(fileName string) {
	err := os.Remove(fileName)
	scrapper.CheckError(err)
}
func handleHome(c echo.Context) error {
	return c.File("pages/home.html")
}
func handleScrape(c echo.Context) error {
	keyword := strings.ToLower(scrapper.CleanString(c.FormValue("keword")))
	fileName := keyword + "_jobs.csv"
	defer removeFile(fileName)

	scrape.Scrape(keyword)
	return c.Attachment(fileName, fileName)
}

func main() {
	// scrape.Scrape("python")
	// scrape.Scrape("golang")
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
