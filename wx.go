// wx - Get NOAA weather.
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/cli"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Body requests url and returns body.
func Body(url string) (body []byte) {

	// request url
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	// read Body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return
}

// LatLon retrieves LatLon from SOAP zipcode request.
func LatLon(zip string) (latlon []string) {

	// checks
	if len(zip) != 5 {
		fmt.Printf("Zipcode is not five numbers: %s\n", zip)
		os.Exit(2)
	}

	// some data
	zipreq := "http://graphical.weather.gov/xml/SOAP_server/ndfdXMLclient.php?listZipCodeList=%s"
	type LL struct {
		LatLon string `xml:"latLonList"`
	}

	// request LatLon in XML
	url := fmt.Sprintf(zipreq, zip)
	body := Body(url)

	// parse LatLon from XML
	var ll LL
	if err := xml.Unmarshal(body, &ll); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	latlon = strings.Split(ll.LatLon, ",")

	return
}

type Day struct {
	Value string `xml:"period-name,attr"`
}
type Cond struct {
	Summary    string `xml:"weather-summary,attr"`
	Visibility string `xml:"value>visibility"` // miles
}
type Data struct {
	Humidity      string   `xml:"humidity>value"` // relative
	Condition     []Cond   `xml:"weather>weather-conditions"`
	WindDirection string   `xml:"direction>value"`   // degrees true
	WindSpeed     []string `xml:"wind-speed>value"`  // gust, sustained, knots
	Pressure      string   `xml:"pressure>value"`    // Barometric, inches of mercury
	Temperature   []string `xml:"temperature>value"` // Farenheit
	Description   []string `xml:"wordedForecast>text"`
}
type Wx struct {
	Place     string `xml:"data>location>description"`
	Place2    string `xml:"data>location>area-description"`
	Time      string `xml:"head>product>creation-date"`
	HalfDay   []Day  `xml:"data>time-layout>start-valid-time"`
	Parameter Data   `xml:"data>parameters"`
}

// Weather retrieves NOAA weather updates for a LatLon coordinates.
func Weather(ll []string) (wx Wx) {

	// some data
	freq := "http://forecast.weather.gov/MapClick.php?lat=%s&lon=%s&unit=0&lg=english&FcstType=dwml"

	// request forecast in XML
	url := fmt.Sprintf(freq, ll[0], ll[1])
	body := Body(url)
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	// decode XML
	if err := decoder.Decode(&wx); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return
}

// PlaceTime prints the location and time of the ZIP code provided
func PlaceTime(wx Wx) {
	place := ""
	if wx.Place != "" {
		place = wx.Place
	} else {
		place = wx.Place2
	}
	fmt.Printf("%s @ %s\n", place, wx.Time)
}

// Forecast prints a NOAA weather forecast
func Forecast(wx Wx) {
	PlaceTime(wx)
	t := len(wx.Parameter.Description)
	for i, Time := range wx.HalfDay[0:t] {
		fmt.Printf("%s: %s\n", Time.Value, wx.Parameter.Description[i])
	}
}

// Current prints NOAA current weather
func Current(wx Wx) {
	PlaceTime(wx)
	t := len(wx.Parameter.Temperature)
	fmt.Printf("  Summary: %s\n", wx.Parameter.Condition[t-2].Summary)
	fmt.Printf("  Temperature (F): %s\n", wx.Parameter.Temperature[t-2])
	fmt.Printf("  Dew Point (F): %s\n", wx.Parameter.Temperature[t-1])
	fmt.Printf("  Humidity: %s\n", wx.Parameter.Humidity)
	fmt.Printf("  Visibility: %s\n", wx.Parameter.Condition[t-1].Visibility)
	fmt.Printf("  Wind Gust (max, knots): %s\n", wx.Parameter.WindSpeed[0])
	fmt.Printf("  Wind Sustained (knots): %s\n", wx.Parameter.WindSpeed[1])
	fmt.Printf("  Pressure (in): %s\n", wx.Parameter.Pressure)
}

func main() {
	// custom help templates
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}}command ZIPCODE

VERSION:
   {{.Version}}

AUTHOR(S): 
   {{range .Authors}}{{ . }}
   {{end}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
	cli.CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   wx {{.Name}}{{if .Flags}} [command options]{{end}} ZIPCODE
`

	app := cli.NewApp()
	app.Name = "wx"
	app.Version = "0.1.0"
	app.Usage = "Get NOAA weather."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Edward A. Roualdes",
			Email: "eroualdes@csuchico.edu",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "forecast",
			Aliases: []string{"f"},
			Usage:   "7 day NOAA weather forecast",
			Action: func(c *cli.Context) {
				zipcode := "95926" // default location
				if len(c.Args()) > 0 {
					zipcode = c.Args()[0]
				}
				ll := LatLon(zipcode)
				wx := Weather(ll)
				Forecast(wx)
			},
		},
		{
			Name:    "current",
			Aliases: []string{"c"},
			Usage:   "current NOAA weather",
			Action: func(c *cli.Context) {
				zipcode := "95926" // default location
				if len(c.Args()) > 0 {
					zipcode = c.Args()[0]
				}
				ll := LatLon(zipcode)
				wx := Weather(ll)
				Current(wx)
			},
		},
	}

	app.Run(os.Args)
}
