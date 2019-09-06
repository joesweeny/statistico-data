package process

import (
	spClient "github.com/statistico/sportmonks-go-client"
	"testing"
)

func TestProcess(t *testing.T) {
	//repo := new(stMock.CountryRepository)
	//
	//server := stMock.HttpClient(func(req *http.Request) *http.Response {
	//	assert.Equal(t, req.URL.String(), "http://example.com/api/v2.0/countries?api_token=my-key&page=1")
	//	b, _ := json.Marshal(countryResponse())
	//	return &http.Response{
	//		StatusCode: 200,
	//		Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
	//	}
	//})
	//
	//client := spClient.Client{
	//	Client:  server,
	//	BaseURL: "http://example.com",
	//	ApiKey:  "my-key",
	//}
	//
	//processor := CountryProcessor{
	//	repository: repo,
	//	factory:    sportmonks.CountryFactory{Clock: clockwork.NewFakeClock()},
	//	client:     &client,
	//	logger:     log.New(ioutil.Discard, "", 0),
	//}
	//
	//t.Run("inserts new country", func(t *testing.T) {
	//	done := make(chan bool)
	//
	//	repo.On("GetById", 1).Return(&app.Country{}, errors.New("not Found"))
	//	repo.On("Insert", mock.Anything).Return(nil)
	//	repo.AssertNotCalled(t, "Update", mock.Anything)
	//	processor.Process("country", "", done)
	//})
	//
	//t.Run("updates existing country", func(t *testing.T) {
	//	done := make(chan bool)
	//
	//	c := stMock.Country(1)
	//	repo.On("GetById", 1).Return(c, nil)
	//	repo.On("Update", &c).Return(nil)
	//	repo.MethodCalled("Update", &c)
	//	repo.AssertNotCalled(t, "Insert", mock.Anything)
	//	processor.Process("country", "", done)
	//})
}

func countryResponse() spClient.CountriesResponse {
	c := clientCountry()

	m := spClient.Meta{}
	m.Pagination.Total = 1
	m.Pagination.Count = 1
	m.Pagination.PerPage = 1
	m.Pagination.CurrentPage = 1
	m.Pagination.TotalPages = 1

	res := spClient.CountriesResponse{}
	res.Data = append(res.Data, c)
	res.Meta = m

	return res
}

func clientCountry() spClient.Country {
	return spClient.Country{
		ID:   1,
		Name: "Brazil",
		Extra: struct {
			Continent   string      `json:"continent"`
			SubRegion   string      `json:"sub_region"`
			WorldRegion string      `json:"world_region"`
			Fifa        interface{} `json:"fifa,string"`
			ISO         string      `json:"iso"`
			Longitude   string      `json:"longitude"`
			Latitude    string      `json:"latitude"`
		}{
			Continent:   "South America",
			SubRegion:   "South America",
			WorldRegion: "South America",
			Fifa:        "BRA",
			ISO:         "BRA",
		},
	}
}
