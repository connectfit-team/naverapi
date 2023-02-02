package geocode

type Lang string

type FilterType string

type Response struct {
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage"`
	Meta         Meta      `json:"meta"`
	Addresses    []Address `json:"addresses"`
}

type Meta struct {
	TotalCount int64 `json:"totalCount"`
	Page       int64 `json:"page,omitempty"`
	Count      int64 `json:"count"`
}

type Address struct {
	RoadAddress     string           `json:"roadAddress"`
	JibunAddress    string           `json:"jibunAddress"`
	EnglishAddress  string           `json:"englishAddress"`
	X               string           `json:"x"`        // X 경도, Longitude
	Y               string           `json:"y"`        // Y 위도, Latitude
	Distance        float64          `json:"distance"` // when requested coordinate param
	AddressElements []AddressElement `json:"addressElements"`
}

type AddressElement struct {
	Types     []string `json:"types"`
	LongName  string   `json:"longName"`
	ShortName string   `json:"shortName"`
	Code      string   `json:"code"`
}
