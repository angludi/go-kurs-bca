package exchange_rate

import "go-kurs-bca/models"

type IRateRepository interface {
	Create(rates []models.ExchangeRate) error
	Select(param models.RateSelectParameter) ([]models.ExchangeRate, error)
	Update(rate models.ExchangeRate) error
	Delete(param models.RateDeleteParameter) error
}
