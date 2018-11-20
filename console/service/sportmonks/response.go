package sportmonks

type Country struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Extra Extra `json:"extra"`
}

type Extra struct {
	Continent string `json:"continent"`
	ISO string `json:"iso"`
}

type CountryResponse struct {
	CountryList []Country `json:"data"`
}
