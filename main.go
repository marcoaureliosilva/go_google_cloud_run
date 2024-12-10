package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func getCityByCEP(cep string) (string, error) {
	response, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", nil
	}

	var location map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&location); err != nil {
		return "", err
	}

	if location["erro"] != nil {
		return "", nil
	}

	return location["localidade"].(string), nil
}

func getWeather(city string) (WeatherResponse, error) {
	apiKey := "849f811824274dedb21192005241012"
	response, err := http.Get("https://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer response.Body.Close()

	var weatherData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
		return WeatherResponse{}, err
	}

	tempC := weatherData["current"].(map[string]interface{})["temp_c"].(float64)
	return WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273.15,
	}, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := mux.Vars(r)["cep"]

	if matched, _ := regexp.MatchString(`^\d{5}-?\d{3}$`, cep); !matched {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := getCityByCEP(cep)
	if err != nil || city == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather, err := getWeather(city)
	if err != nil {
		http.Error(w, "unable to retrieve weather", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", weatherHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
