package entity

type CartProduct struct {
	Product  *Product `json:"product"`
	Quantity uint64   `json:"quantity"`
}

type Cart struct {
	ID       uint64         `json:"id"`
	Products []*CartProduct `json:"products"`
}
