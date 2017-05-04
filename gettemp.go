package main

import (
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"strings"
	"strconv"
	"fmt"
	"math"
	"time"
	"os"
)

// Convert HTML entity angle brackets to actual angle brackets
func formatXml(data string) (xml string) {
	xml = strings.Replace(data, "&lt;", "<", -1)
	xml = strings.Replace(xml, "&gt;", ">", -1)
	return
}

func convval(str string) (out float64) {
	str = strings.Replace(str, ",", ".", 1)
	str = strings.Trim(str, " ")
	out, _ = strconv.ParseFloat(str, 64)
	return
}

// Calculates wind chill based on temperature and wind speed.
// Formula taken from: http://www8.tfe.umu.se/weather-new/js/index.js.
func calcWindChill(temp, speed float64) float64 {
	expTmp := math.Pow(speed, 0.16)
	return 13.126667 +
	    0.6215*temp -
	    13.924748*expTmp +
	    0.4875195*temp*expTmp
}

// LOGGING

// Get the current date in ISO 8601.
func getDate() string {
	now := time.Now()
	return fmt.Sprintf("%s", now.Format("2006-01-02"))
}

// Log date and temp, semi-colon separated.
func logTemp(f *os.File, temp float64) {
	f.WriteString(fmt.Sprintf("%s;%.1f\n", getDate(), temp))
}

func getLogFile(fileName string) *os.File {
	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	f, err := os.OpenFile(fileName, flags, 0644)

	if err != nil {
		panic(err)
	}

	return f
}

func main() {

	type WeatherData struct {
		Temperature   string `xml:"root>tempmed"`
		WindSpeed	  string `xml:"root>vindh"`
		WindChill	  string `xml:"root>windChill"`
	}

	// Weather data URL
	const url = "http://www8.tfe.umu.se/WeatherWebService2012/Service.asmx/Aktuellavarden"

	// HTTP GET
	res, _ := http.Get(url)

	// Read all XML data
	bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	xmlStr := formatXml(string(bytes))

	// Parse XML
	var v WeatherData
	xml.Unmarshal([]byte(xmlStr), &v)

	temp := convval(v.Temperature)
	windSpeed := convval(v.WindSpeed)
	//windChill := convval(v.WindChill)
	windChill := calcWindChill(temp, windSpeed)

	// Output
	fmt.Printf("Current temp: %.1f°C\n", temp)

	if windSpeed < 1.2 {
		fmt.Println("Negligible wind chill.")
	} else {
		fmt.Printf("Feels like: %.1f°C\n", windChill)
	}

	fmt.Printf("Wind speed: %.1f m/s\n", windSpeed)

	// log date and time to a file provided on the commandline.
	if len(os.Args) > 1 {
		logFile := getLogFile(os.Args[1])
		defer logFile.Close()

		logTemp(logFile, temp)
	}
}
