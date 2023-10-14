package apis

import (
	"encoding/json"
	"io"
	"net/http"
)

type BackendResponse struct {
	Message     string `json:"message"`
	CityName    string `json:"city_name"`
	FeelsLike   string `json:"feels_like"`
	Temp        string `json:"temp"`
	Description string `json:"description"`
}

func NewBackendResponse(resp *http.Response) (*BackendResponse, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bs := new(BackendResponse)
	if err := json.Unmarshal(body, &bs); err != nil {
		return nil, err
	}

	return bs, nil
}
