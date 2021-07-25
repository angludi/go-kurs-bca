package scrapper

import "github.com/PuerkitoBio/goquery"

type ScrapperRepository interface {
	Scrapping(url string) (doc *goquery.Document, err error)
}