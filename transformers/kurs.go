package transformers

import "go-kurs-bca/models"

func (res *Transformer) RateTransformer(data models.ExchangeRate) {
	res.Data = assignRate(data)
}

func (res *CollectionTransformer) RateCollection(datas []models.ExchangeRate) {
	for _, data := range datas {
		res.Data = append(res.Data, assignRate(data))
	}
}

type JualBeli struct {
	Jual float64 `json:"jual"`
	Beli float64 `json:"beli"`
}

type ExchangeRate struct {
	Symbol    string   `json:"symbol"`
	Date      string   `json:"date"`
	ERate     JualBeli `json:"e_rate"`
	TTCounter JualBeli `json:"tt_counter"`
	BankNotes JualBeli `json:"bank_notes"`
}

func assignRate(req models.ExchangeRate) interface{} {
	var res ExchangeRate

	res.Symbol = req.Symbol
	res.Date = req.Date.Format("2006-01-02")
	res.ERate.Beli = req.ERateBuy
	res.ERate.Jual = req.ERateSell
	res.TTCounter.Beli = req.TTBuy
	res.TTCounter.Jual = req.TTSell
	res.BankNotes.Beli = req.BNBuy
	res.BankNotes.Jual = req.BNSell

	return res
}
