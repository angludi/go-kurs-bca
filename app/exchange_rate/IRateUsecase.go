package exchange_rate

import "go-kurs-bca/models"

type IRateUsecase interface {
	Indexing() error
	Create(rate models.ExchangeRate) error
	Update(rate models.ExchangeRate) error
	GetRates(code, startdate, enddate string) ([]models.ExchangeRate, error)
	DeleteByDate(date string) error
}
