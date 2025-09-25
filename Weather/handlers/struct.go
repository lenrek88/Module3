package handlers

type WeeklyResponse struct {
	Cod     string     `json:"cod"`
	Message any        `json:"message"`
	Cnt     int        `json:"cnt"`
	City    City       `json:"city"`
	List    []Forecast `json:"list"`
}

type City struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Country    string      `json:"country"`
	Population int         `json:"population"`
	Timezone   int         `json:"timezone"`
	Sunrise    int64       `json:"sunrise"`
	Sunset     int64       `json:"sunset"`
	Coord      Coordinates `json:"coord"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Forecast struct {
	Dt         int64     `json:"dt"`
	DtTxt      string    `json:"dt_txt"`
	Main       MainInfo  `json:"main"`
	Weather    []Weather `json:"weather"`
	Clouds     Clouds    `json:"clouds"`
	Wind       Wind      `json:"wind"`
	Visibility int       `json:"visibility"`
	Pop        float64   `json:"pop"`
	Sys        Sys       `json:"sys"`
	Rain       *Rain     `json:"rain,omitempty"`
}

type MainInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Sys struct {
	Pod string `json:"pod"`
}

type Rain struct {
	ThreeH float64 `json:"3h"`
}

type SimplifiedForecast struct {
	DateTime     string  `json:"дата"`
	Temperature  float64 `json:"температура"`
	FeelsLike    float64 `json:"ощущается_как"`
	Description  string  `json:"погода"`
	Humidity     int     `json:"влажность"`
	WindSpeed    float64 `json:"скорость_ветра"`
	ChanceOfRain float64 `json:"вероятность_дождя"`
}

type SimplifiedResponse struct {
	City      string               `json:"город"`
	Forecasts []SimplifiedForecast `json:"прогноз"`
}

type TodayResponse struct {
	Coord      any            `json:"coord"`
	Weather    []WeatherToday `json:"weather"`
	Base       string         `json:"base"`
	Main       MainToday      `json:"main"`
	Visibility int            `json:"visibility"`
	Wind       WindToday      `json:"wind"`
	Clouds     map[string]int `json:"clouds"`
	Dt         int            `json:"dt"`
	Sys        SysToday       `json:"sys"`
	Timezone   int            `json:"timezone"`
	Id         int            `json:"id"`
	Name       string         `json:"name"`
	Code       int            `json:"code"`
}

type WindToday struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type MainToday struct {
	FeelsLike float64 `json:"feels_like"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	Temp      float64 `json:"temp"`
	TempMax   float64 `json:"temp_max"`
	TempMin   float64 `json:"temp_min"`
}
type WeatherToday struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type SysToday struct {
	Type    int    `json:"type"`
	Id      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}
