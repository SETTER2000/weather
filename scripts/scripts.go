package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/addetz/go-temperature-checker/apis"
	"honnef.co/go/js/dom"
)

//go:embed cities
var citiesInput string

func main() {
	cities := strings.Split(citiesInput, "\n")
	document := dom.GetWindow().Document()
	populateCities(document, cities)

	selection := document.GetElementByID("citiesDropdown")
	selection.AddEventListener("change", false, func(e dom.Event) {
		selEl := selection.(*dom.HTMLSelectElement)
		index := selEl.SelectedIndex
		selectedText := selEl.Options()[index].Text
		city := strings.Split(selectedText, ",")[0]
		go func(callback func(document dom.Document, bs *apis.BackendResponse)) {
			resetWeather(document)
			resp, err := http.Get(fmt.Sprintf("/weather/%s", city))
			if err != nil {
				log.Fatal(err)
			}
			bs, err := apis.NewBackendResponse(resp)
			if err != nil {
				log.Fatal(err)
			}
			callback(document, bs)

		}(weatherCallback)
	})
}

func populateCities(document dom.Document, cities []string) {
	selEl := document.GetElementByID("citiesDropdown").(*dom.HTMLSelectElement)
	for _, c := range cities {
		o := document.CreateElement("option")
		o.SetTextContent(c)
		selEl.AppendChild(o)
	}
}

func weatherCallback(document dom.Document, bs *apis.BackendResponse) {
	weatherContainer := document.GetElementByID("weatherResult")
	weatherContainer.Class().Remove("d-none")
	title := document.GetElementByID("weatherTitle").(*dom.HTMLHeadingElement)
	title.SetTextContent(fmt.Sprintf("Weather forecast for %s", bs.CityName))
	desc := document.GetElementByID("weatherDesc").(*dom.HTMLHeadingElement)
	desc.SetTextContent(fmt.Sprintf("Description: %s", bs.Description))
	temp := document.GetElementByID("weatherTemp").(*dom.HTMLHeadingElement)
	temp.SetTextContent(fmt.Sprintf("Temperature: %s C", bs.Temp))
	feelsLike := document.GetElementByID("weatherFeels").(*dom.HTMLHeadingElement)
	feelsLike.SetTextContent(fmt.Sprintf("Feels Like: %s C", bs.FeelsLike))
	switch {
	case strings.Contains(bs.Description, "clear"):
		img := document.GetElementByID("sunImg")
		img.Class().Remove("d-none")
	case strings.Contains(bs.Description, "clouds"):
		img := document.GetElementByID("cloudsImg")
		img.Class().Remove("d-none")
	case strings.Contains(bs.Description, "rain"):
		img := document.GetElementByID("rainImg")
		img.Class().Remove("d-none")
	case strings.Contains(bs.Description, "storm"):
		img := document.GetElementByID("stormImg")
		img.Class().Remove("d-none")
	case strings.Contains(bs.Description, "snow"):
		img := document.GetElementByID("snowImg")
		img.Class().Remove("d-none")
	case strings.Contains(bs.Description, "mist"):
		img := document.GetElementByID("mistImg")
		img.Class().Remove("d-none")
	}
}

func resetWeather(document dom.Document) {
	weatherContainer := document.GetElementByID("weatherResult")
	weatherContainer.Class().Add("d-none")
	document.GetElementByID("sunImg").Class().Add("d-none")
	document.GetElementByID("cloudsImg").Class().Add("d-none")
	document.GetElementByID("rainImg").Class().Add("d-none")
	document.GetElementByID("stormImg").Class().Add("d-none")
	document.GetElementByID("snowImg").Class().Add("d-none")
	document.GetElementByID("mistImg").Class().Add("d-none")
}
