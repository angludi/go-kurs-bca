package usecase

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	RateInterface "go-kurs-bca/app/exchange_rate"
	IndexingInterface "go-kurs-bca/app/scrapper"
	"go-kurs-bca/models"
	"strconv"
	"strings"
	"time"
)

type RateUsecase struct {
	RateRepository     RateInterface.IRateRepository
	ScrapperRepository IndexingInterface.ScrapperRepository
}

func NewRateUsecase(k RateInterface.IRateRepository, s IndexingInterface.ScrapperRepository) RateInterface.IRateUsecase {
	return &RateUsecase{
		RateRepository:     k,
		ScrapperRepository: s,
	}
}

func (u *RateUsecase) Indexing() (err error) {
	var exchangeRates []models.ExchangeRate
	var codes []string
	var erbuys, ersells, ttbuys, ttsells, bnbuys, bnsells []float64

	html, err := u.ScrapperRepository.Scrapping("https://www.bca.co.id/en/informasi/kurs")
	if err != nil {
		return
	}

	// Date
	rateDateStr, _ := html.Find("div.o-kurs-refresh-wrapper span.desc-ref-kurs.refresh-date").Html()
	rateDate, _ := time.Parse("02 January 2006 15:04", rateDateStr)

	// Exchange Rates
	html.Find("table.m-table-kurs tbody p").Each(func(index int, item *goquery.Selection) {
		attr, _ := item.Attr("rate-type")

		valueStr := strings.TrimSpace(item.Text())
		valueStr = strings.Replace(valueStr, ",", "", 1)

		value, _ := strconv.ParseFloat(valueStr, 64)

		attrStr := fmt.Sprintf("%s", strings.ToLower(strings.TrimSpace(attr)))

		switch attrStr {
		case "erate-buy":
			erbuys = append(erbuys, value)
		case "erate-sell":
			ersells = append(ersells, value)
		case "tt-buy":
			ttbuys = append(ttbuys, value)
		case "tt-sell":
			ttsells = append(ttsells, value)
		case "bn-buy":
			bnbuys = append(bnbuys, value)
		case "bn-sell":
			bnsells = append(bnsells, value)
		default:
			codes = append(codes, strings.TrimSpace(item.Text()))
		}
	})

	// Mapping to struct
	for i := range codes {
		var rate models.ExchangeRate

		rate.Date = rateDate
		rate.Symbol = codes[i]
		rate.ERateBuy = erbuys[i]
		rate.ERateSell = ersells[i]
		rate.TTBuy = ttbuys[i]
		rate.TTSell = ttsells[i]
		rate.BNBuy = bnbuys[i]
		rate.BNSell = bnsells[i]

		exchangeRates = append(exchangeRates, rate)
	}

	if len(exchangeRates) > 0 {
		// Insert Into DB
		err = u.RateRepository.Create(exchangeRates)
		if err != nil {
			return
		}
	} else {
		err = errors.New("Can not fetch data")
	}

	return
}

func (u *RateUsecase) GetRates(code, startdate, enddate string) (rates []models.ExchangeRate, err error) {
	param := models.RateSelectParameter{
		Code: code,
		StartDate: startdate,
		EndDate: enddate,
	}
	rates, err = u.RateRepository.Select(param)

	return
}

func (u *RateUsecase) Create(rate models.ExchangeRate) (err error) {
	var rates []models.ExchangeRate
	rates = append(rates, rate)

	err = u.RateRepository.Create(rates)

	return
}

func (u *RateUsecase) Update(rate models.ExchangeRate) (err error) {
	err = u.RateRepository.Update(rate)

	return
}

func (u *RateUsecase) DeleteByDate(date string) (err error) {
	param := models.RateDeleteParameter{
		Date: date,
	}
	err = u.RateRepository.Delete(param)

	return
}