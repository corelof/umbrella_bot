package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type dailyForecast struct {
	Time              uint64  `json:"time"`
	PrecipType        string  `json:"precipType"`
	PrecipProbability float64 `json:"precipProbability"`
}

type dailyForecastSet struct {
	Data []dailyForecast `json:"data"`
}

type forecast struct {
	Daily dailyForecastSet `json:"daily"`
}

func getTodayForecast() (forecast, error) {
	url := fmt.Sprintf("%s%s/%s,%s", os.Getenv("api_url"), os.Getenv("darksky_api_key"), os.Getenv("latitude"), os.Getenv("longtitude"))
	resp, err := http.Get(url)
	if err != nil {
		return forecast{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return forecast{}, err
	}
	var fc forecast
	err = json.Unmarshal(body, &fc)
	return fc, err
}

func TodayForecast() (string, int) {
	fc, err := getTodayForecast()
	if err != nil {
		log.Println(err)
		return "rain", 100
	}
	fmt.Println(fc.Daily.Data[0].PrecipProbability)
	return fc.Daily.Data[0].PrecipType, int(fc.Daily.Data[0].PrecipProbability * 100)
}
