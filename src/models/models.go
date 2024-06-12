package models

type GoldPriceResponse struct {
	ResponseData      GoldPriceData `json:"response_data"`
	ResponseMessage   string        `json:"response_message"`
	ResponseTimestamp string        `json:"response_timestamp"`
}

type GoldPriceData struct {
	BarBuy        string `json:"bar_buy"`
	BarSell       string `json:"bar_sell"`
	OrnamentBuy   string `json:"ornament_buy"`
	OrnamentSell  string `json:"ornament_sell"`
	StatusChange  string `json:"status_change"`
	TodayChange   string `json:"today_change"`
	UpdatedDate   string `json:"updated_date"`
	UpdatedTime   string `json:"updated_time"`
	UpdateTheTime string `json:"updated_the_time"`
}
