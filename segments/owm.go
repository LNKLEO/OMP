package segments

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"

	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"

	"github.com/LNKLEO/OMP/log"
	"github.com/LNKLEO/OMP/properties"
)

type Owm struct {
	base

	Weather     string
	URL         string
	units       string
	UnitIcon    string
	Temperature float64
	WeatherIcon string
}

const (
	// APIKey openweathermap api key
	APIKey properties.Property = "api_key"
	// Location openweathermap location
	Location properties.Property = "location"
	// Units openweathermap units
	Units properties.Property = "units"
	// CacheKeyResponse key used when caching the response
	CacheKeyResponse string = "owm_response"
	// CacheKeyURL key used when caching the url responsible for the response
	CacheKeyURL string = "owm_url"

	OMPOWMAPIKey = "OMP_OWM_API_KEY"
)

type weather struct {
	ShortDescription string `json:"main"`
	Description      string `json:"description"`
	TypeID           string `json:"icon"`
}
type temperature struct {
	Value float64 `json:"temp"`
}

type owmDataResponse struct {
	Data        []weather `json:"weather"`
	temperature `json:"main"`
}

func (d *Owm) Enabled() bool {
	err := d.setStatus()

	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func (d *Owm) Template() string {
	return " {{ .Weather }} ({{ .Temperature }}{{ .UnitIcon }}) "
}

func (d *Owm) getResult() (*owmDataResponse, error) {
	response := new(owmDataResponse)

	apikey := properties.OneOf(d.props, ".", APIKey, "apiKey")
	if len(apikey) == 0 {
		apikey = d.env.Getenv(OMPOWMAPIKey)
	}

	location := d.props.GetString(Location, "De Bilt,NL")
	location = url.QueryEscape(location)

	if len(apikey) == 0 || len(location) == 0 {
		return nil, errors.New("no api key or location found")
	}

	units := d.props.GetString(Units, "standard")
	httpTimeout := d.props.GetInt(properties.HTTPTimeout, properties.DefaultHTTPTimeout)

	d.URL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=%s&appid=%s", location, units, apikey)

	body, err := d.env.HTTPRequest(d.URL, nil, httpTimeout)
	if err != nil {
		return new(owmDataResponse), err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return new(owmDataResponse), err
	}

	return response, nil
}

func (d *Owm) setStatus() error {
	units := d.props.GetString(Units, "standard")

	q, err := d.getResult()
	if err != nil {
		return err
	}

	if len(q.Data) == 0 {
		return errors.New("No data found")
	}

	d.Temperature = math.Round(q.temperature.Value*2) / 2.0

	id := q.Data[0].TypeID
	icon := ""
	switch id {
	case "01n":
		icon = "󰖔"
	case "01d":
		icon = "󰖙"
	case "02n":
		icon = "󰼱"
	case "02d":
		icon = "󰖕"
	case "03n":
		icon = "󰖐"
	case "03d":
		icon = "󰖐"
	case "04n":
		icon = "󰼯"
	case "04d":
		icon = "󰼯"
	case "09n":
		icon = "󰖒"
	case "09d":
		icon = "󰖒"
	case "10n":
		icon = "󰖗"
	case "10d":
		icon = "󰼳"
	case "11n":
		icon = "󰖓"
	case "11d":
		icon = "󰼲"
	case "13n":
		icon = "󰖘"
	case "13d":
		icon = "󰼴"
	case "50n":
		icon = "󰖑"
	case "50d":
		icon = "󰖑"
	}
	d.Weather = cases.Title(lang.Und).String(q.Data[0].Description)
	d.WeatherIcon = icon
	d.units = units
	d.UnitIcon = "°"
	switch d.units {
	case "imperial":
		d.UnitIcon = "󰔅"
	case "metric":
		d.UnitIcon = "󰔄"
	case "standard":
		d.UnitIcon = "󰔆"
	case "":
		d.UnitIcon = "󰔆"
	}
	return nil
}
