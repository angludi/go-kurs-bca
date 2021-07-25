package repository

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	ScrapperInterface "go-kurs-bca/app/scrapper"
	"net/http"
)

type IndexingRepository struct {
}

func NewScrapperRepository() ScrapperInterface.ScrapperRepository {
	return &IndexingRepository{}
}

func (r *IndexingRepository) Scrapping(url string) (doc *goquery.Document, err error){
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	doc, err = goquery.NewDocumentFromReader(reader)
	return
}