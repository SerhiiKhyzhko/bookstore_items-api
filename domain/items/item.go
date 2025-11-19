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
	Status            string      `json:"status"` //maybe bool?
}

type Description struct {
	PlainText string `json:"html"`
	Html      string `json:"html"`
}

type Pictures struct {
	Id  int64  `jsnon:"id"`
	Url string `json:"url"`
}
