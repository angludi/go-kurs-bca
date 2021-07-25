package models

import "time"

type ExchangeRate struct {
	Date         time.Time 	`gorm:"column:date"`
	Symbol 		 string 	`gorm:"column:symbol"`
	ERateBuy     float64 	`gorm:"column:erate_buy"`
	ERateSell    float64 	`gorm:"column:erate_sell"`
	TTBuy        float64	`gorm:"column:tt_buy"`
	TTSell       float64	`gorm:"column:tt_sell"`
	BNBuy        float64	`gorm:"column:bn_buy"`
	BNSell       float64 	`gorm:"column:bn_sell"`
}

type RateSelectParameter struct {
	Code string
	StartDate string
	EndDate string
}

type RateDeleteParameter struct {
	Code string
	Date string
}

type RateRequestDetail struct {
	Jual float64 `json:"jual"`
	Beli float64 `json:"beli"`
}

type RateRequest struct {
	Symbol    string   `json:"symbol"`
	Date      string   `json:"date"`
	ERate     RateRequestDetail `json:"e_rate"`
	TTCounter RateRequestDetail `json:"tt_counter"`
	BankNotes RateRequestDetail `json:"bank_notes"`
}