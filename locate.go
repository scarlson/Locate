package locate

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	url   = "http://www.telize.com/geoip"
	ipurl = "http://freegeoip.net/json/"
)

type inputLocation struct {
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	Asn            string  `json:"asn"`
	Offset         string  `json:"offset"`
	Ip             string  `json:"ip"`
	Area_code      string  `json:"area_code"`
	Continent_code string  `json:"continent_code"`
	Dma_code       string  `json:"dma_code"`
	City           string  `json:"city"`
	Timezone       string  `json:"timezone"`
	Timezone2      string  `json:"time_zone"`
	State          string  `json:"region"`
	Isp            string  `json:"isp"`
	Postal_code    string  `json:"postal_code"`
	Zip_code       string  `json:"zip_code"`
	Country        string  `json:"country"`
	Country_name   string  `json:"country_name"`
	Country_code   string  `json:"country_code"`
	Country_code3  string  `json:"country_code3"`
	Region         string  `json:"region"`
	Region_code    string  `json:"region_code"`
	Region_name    string  `json:"region_name"`
}

func (self *inputLocation) sanitize() *Location {
	loc := &Location{}
	loc.Longitude = self.Longitude
	loc.Latitude = self.Latitude
	loc.City = self.City
	if self.Region == "" {
		loc.State = self.Region_name
	} else {
		loc.State = self.Region
	}
	loc.StateCode = self.Region_code
	if self.Country == "" {
		loc.Country = self.Country_name
	} else {
		loc.Country = self.Country
	}
	loc.CountryCode = self.Country_code
	loc.Ip = self.Ip
	if self.Postal_code == "" {
		loc.ZipCode = self.Zip_code
	} else {
		loc.ZipCode = self.Postal_code
	}
	if self.Timezone2 == "" {
		loc.Timezone = self.Timezone
	} else {
		loc.Timezone = self.Timezone2
	}
	return loc
}

type Location struct {
	Longitude   float64
	Latitude    float64
	City        string
	ZipCode     string
	Timezone    string
	Country     string
	CountryCode string
	State       string
	StateCode   string
	Ip          string
}

func (l *Location) Json() (string, error) {
	bytes, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func WhereAmI() (*Location, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	loc := &inputLocation{}
	err = json.Unmarshal(bytes, loc)
	if err != nil {
		return nil, err
	}
	return loc.sanitize(), nil
}

func WhereIsThis(ip string) (*Location, error) {
	resp, err := http.Get(ipurl + ip)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	loc := &inputLocation{}
	err = json.Unmarshal(bytes, loc)
	if err != nil {
		return nil, err
	}
	return loc.sanitize(), nil
}

func main() {
	ip := flag.String("ip", "", "IP Address")
	format := flag.String("format", "text", "Output json/text")
	flag.Parse()
	loc := &Location{}
	var err error
	if *ip == "" {
		loc, err = WhereAmI()
		if err != nil {
			panic(err)
		}
	} else {
		loc, err = WhereIsThis(*ip)
		if err != nil {
			panic(err)
		}
	}
	if *format == "json" {
		jsonloc, err := json.Marshal(loc)
		if err != nil {
			panic(err)
		}
		jl := string(jsonloc)
		fmt.Println(jl)
	} else {
		fmt.Println(loc.Latitude, loc.Longitude)
	}
}
