package models

type Store struct {
	Payload struct {
		StoreId            int64  `json:"store_id"`
		StoreName          string `json:"store_name"`
		StoreLocationId    int64  `json:"store_location_id"`
		StoreServiceId     int64  `json:"store_service_id"`
		StoreStockId       int64  `json:"store_stock_id"`
		StoreStatus        int64  `json:"store_status"`
		AccountId          int64  `json:"account_id"`
		LocationId         int64  `json:"location_id"`
		StoreConfig        int64  `json:"store_config"`
		Properties         string `json:"properties,omitempty"`
		StoreInvoiceSeries string `json:"store_invoice_series,omitempty"`
		ConfigZone         int64  `json:"config_zone"`
		Deleted            string `json:"__deleted,omitempty"`
	} `json:"payload"`
}

type StoreProperties struct {
	PackedLabel struct {
		Zone interface{} `json:"zone"`
	} `json:"packed_label,omitempty"`
}
