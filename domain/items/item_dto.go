package items

type Item struct {
	Id                string      `json:"id"`
	Seller            int64       `json:"seller"`
	Title             string      `json:"title"`
	Description       Description `json:"description"`
	Pictures          []Pictures  `json:"pictures"`
	Video             string      `json:"video"`
	Price             float32     `json:"price"`
	AvailableQuantity int         `json:"available_quantity"`
	SoldQuantity      int         `json:"sold_quantity"`
	Status            string      `json:"status"`
}

type PartialUpdateItem struct {
	Title             *string            `json:"title,omitempty"`
	Description       *UpdateDescription `json:"description,omitempty"`
	Pictures          *[]UpdatePictures  `json:"pictures,omitempty"`
	Video             *string            `json:"video,omitempty"`
	Price             *float32           `json:"price,omitempty"`
	AvailableQuantity *int               `json:"available_quantity,omitempty"`
	SoldQuantity      *int               `json:"sold_quantity,omitempty"`
	Status            *string            `json:"status,omitempty"`
}

type Description struct {
	PlainText string `json:"plain_text"`
	Html      string `json:"html"`
}

type UpdateDescription struct {
	PlainText *string `json:"plain_text"`
	Html      *string `json:"html"`
}

type Pictures struct {
	Id  int64  `json:"id"`
	Url string `json:"url"`
}

type UpdatePictures struct {
	Id  *int64  `json:"id"`
	Url *string `json:"url"`
}
