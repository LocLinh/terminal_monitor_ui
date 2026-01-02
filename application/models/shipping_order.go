package models

type ShippingOrderMessage struct {
	Payload ShippingOrder `json:"payload"`
}

type ShippingOrder struct {
	Deleted              string `json:"__deleted"`
	CreatedTime          int64  `json:"created_time"`
	DeliveriedTime       int64  `json:"deliveried_time"`
	ExpectedDeliveryTime int64  `json:"expected_delivery_time"`
	PackedTime           int64  `json:"packed_time"`
	PartnerOrderID       int64  `json:"partner_order_id"`
	PaymentReceipt       int64  `json:"payment_receipt"`
	PickupStoreID        int64  `json:"pickup_store_id"`
	Properties           string `json:"properties"`
	ReceivingTime        int64  `json:"receiving_time"`
	SoAddress            string `json:"so_address"`
	SoCity               int64  `json:"so_city"`
	SoCod                int64  `json:"so_cod"`
	SoComissionType      int64  `json:"so_comission_type"`
	SoComissionValue     int64  `json:"so_comission_value"`
	SoCompanyID          int64  `json:"so_company_id"`
	SoCtime              int64  `json:"so_ctime"`
	SoCustomerCode       string `json:"so_customer_code"`
	SoCustomerID         int64  `json:"so_customer_id"`
	SoDistrict           int64  `json:"so_district"`
	SoFee                int64  `json:"so_fee"`
	SoID                 int64  `json:"so_id"`
	SoItemContent        string `json:"so_item_content"`
	SoName               string `json:"so_name"`
	SoNote               string `json:"so_note"`
	SoOrderCode          int64  `json:"so_order_code"`
	SoPackedCode         int64  `json:"so_packed_code"`
	SoPhone              string `json:"so_phone"`
	SoPriority           int64  `json:"so_priority"`
	SoShipperID          int64  `json:"so_shipper_id"`
	SoSourceType         int64  `json:"so_source_type"`
	SoStatus             int64  `json:"so_status"`
	SoUtime              int64  `json:"so_utime"`
	SoWard               int64  `json:"so_ward"`
}
