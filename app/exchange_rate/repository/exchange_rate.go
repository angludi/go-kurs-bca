package repository

import (
	"fmt"
	RateInterface "go-kurs-bca/app/exchange_rate"
	"go-kurs-bca/models"
	"gorm.io/gorm"
	"strings"
)

type RateRepository struct {
	DB *gorm.DB
}

func NewRateRepository(db *gorm.DB) RateInterface.IRateRepository {
	return &RateRepository{DB: db}
}

func (r *RateRepository) Create(rates []models.ExchangeRate) (err error) {
	sqlQuery := `INSERT IGNORE INTO exchange_rates (date, symbol, erate_buy, erate_sell, tt_buy, tt_sell, bn_buy, bn_sell, created_at)
				 VALUES`

	valueString := []string{}
	valueArgs   := []interface{}{}
	for _, rate := range rates {
		valueString = append(valueString, "(?, ?, ?, ?, ?, ?, ?, ?, NOW())")
		valueArgs = append(valueArgs, rate.Date.Format("2006-01-02"))
		valueArgs = append(valueArgs, rate.Symbol)
		valueArgs = append(valueArgs, rate.ERateBuy)
		valueArgs = append(valueArgs, rate.ERateSell)
		valueArgs = append(valueArgs, rate.TTBuy)
		valueArgs = append(valueArgs, rate.TTSell)
		valueArgs = append(valueArgs, rate.BNBuy)
		valueArgs = append(valueArgs, rate.BNSell)
	}

	sqlStat := fmt.Sprintf("%s %s", sqlQuery, strings.Join(valueString, ","))

	tx := r.DB.Begin()

	if err = tx.Exec(sqlStat, valueArgs...).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (r *RateRepository) Select(param models.RateSelectParameter) (rates []models.ExchangeRate, err error) {
	tx := r.DB

	if  param.Code != "" {
		tx = tx.Where("symbol = ?", param.Code)
	}

	if param.StartDate != "" {
		tx = tx.Where("date >= ? and date <= ?", param.StartDate, param.EndDate)
	}

	err = tx.Find(&rates).Error

	return
}

func (r *RateRepository) Update(rate models.ExchangeRate) (err error) {
	tx := r.DB.Begin()

	if err = tx.Where("date = ? AND symbol = ?", rate.Date.Format("2006-01-02"), rate.Symbol).Updates(&rate).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}

func (r *RateRepository) Delete(param models.RateDeleteParameter) (err error) {
	sqlQuery := `DELETE FROM exchange_rates WHERE 1 `

	if len(param.Code) > 0 {
		sqlQuery += `AND symbol = "` + param.Code + `" `
	}

	if len(param.Date) > 0 {
		sqlQuery += `AND date = "` + param.Date + `" `
	}

	tx := r.DB.Begin()

	if err = tx.Exec(sqlQuery).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return
}