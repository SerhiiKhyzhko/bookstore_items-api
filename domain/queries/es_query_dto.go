package queries

type EsQuery struct {
	SearchText        *string  `json:"search_text"`
	Status            *string  `json:"status"`
	Seller            *int64   `json:"seller"`
	MinPrice          *float64 `json:"min_price"`
	MaxPrice          *float64 `json:"max_price"`
	AvailableQuantity *int     `json:"available_quantity"`

	// Пагінація (Технічні поля)
	From *int `json:"from"` // Скільки пропустити (Offset)
	Size *int `json:"size"` // Скільки повернути (Limit)
}
