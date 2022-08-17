package models

type Cart struct {
	ID       string         `json:"id"`
	Products []*CartProduct `json:"products"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    *string `json:"image_url"`
}

type CartProduct struct {
	CartID    string  `json:"cart_id"`
	ProductID string  `json:"product_id"`
	Quantity  uint    `json:"quantity"`
	Product   Product `json:"product"`
}
