package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name string `json:"Name"`
	Main struct {
		Kelvin     float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
	} `json:"main"`
	Weather interface {} `json:"weather"`
}

/*
type Weather struct {
	Forecast    string `json:"main"`
	Description string `json:"description"`
}
*/

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func query(city string) (WeatherData, error) {
	apiconfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return WeatherData{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?appid=" + apiconfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return WeatherData{}, err
	}

	defer resp.Body.Close()

	var d WeatherData

	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}

	return d, nil
}

func main() {
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternlServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}
