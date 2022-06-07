package entity

// Really naive representation but good enough for our purposes
type Product struct {
	ID     uint64   `json:"id"`
	Name   string   `json:"name"`
	Price  float64  `json:"price"`
	Images *[]Image `json:"images"`
}
