package modelsorder

type Order struct {
	OrderUID    string `json:"order_uid" env-required:"true"`
	TrackNumber string `json:"track_number" env-required:"true"`
	Entry       string `json:"entry" env-required:"true"`
	Delivery
	Payment
	Items             []Items
	Locale            string `json:"locale" env-required:"true"`
	InternalSignature string `json:"internal_signature" env-required:"true"`
	CustomerID        string `json:"customer_id" env-required:"true"`
	DeliveryService   string `json:"delivery_service" env-required:"true"`
	ShardKey          string `json:"shardkey" env-required:"true"`
	SmID              uint16 `json:"sm_id" env-required:"true"`
	DateCreated       string `json:"date_created" env-required:"true"`
	OofShard          string `json:"oof_shard" env-required:"true"`
}

type Delivery struct {
	Name    string `json:"name" env-required:"true"`
	Phone   string `json:"phone" env-required:"true"`
	Zip     string `json:"zip" env-required:"true"`
	City    string `json:"city" env-required:"true"`
	Address string `json:"address" env-required:"true"`
	Region  string `json:"region" env-required:"true"`
	Email   string `json:"email" env-required:"true"`
}

type Payment struct {
	Transaction  string `json:"transaction" env-required:"true"`
	RequestID    string `json:"request_id" env-required:"true"`
	Currency     string `json:"currency" env-required:"true"`
	Provider     string `json:"provider" env-required:"true"`
	Amount       uint16 `json:"amount" env-required:"true"`
	PaymentDT    uint16 `json:"paymentDT" env-required:"true"`
	Bank         string `json:"bank" env-required:"true"`
	DeliveryCost uint16 `json:"delivery_cost" env-required:"true"`
	GoodsTotal   uint16 `json:"goods_total" env-required:"true"`
	CustomFee    uint16 `json:"custom_fee" env-required:"true"`
}

type Items struct {
	ChrtID      uint16 `json:"chrt_id" env-required:"true"`
	TrackNumber string `json:"track_number" env-required:"true"`
	Price       uint16 `json:"price" env-required:"true"`
	RID         string `json:"rid" env-required:"true"`
	Name        string `json:"name" env-required:"true"`
	Sale        uint16 `json:"sale" env-required:"true"`
	Size        uint16 `json:"size" env-required:"true"`
	TotalPrice  uint16 `json:"total_price" env-required:"true"`
	NmID        uint16 `json:"nmID" env-required:"true"`
	Brand       string `json:"brand" env-required:"true"`
	Status      uint16 `json:"status" env-required:"true"`
}
