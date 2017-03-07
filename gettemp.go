package main

import (
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "strings"
    "strconv"
    "fmt"
    "math"
)

func convprep(str string) (out string) {
    out = strings.Replace(str, ",", ".", 1)
    out = strings.Trim(out, " ")
    return
}

// Convert HTML entity angle brackets to actual angle brackets
func formatXml(data string) (xml string) {
    xml = strings.Replace(data, "&lt;", "<", -1)
    xml = strings.Replace(xml, "&gt;", ">", -1)
    return
}

// Implements the standard Wind Chill formula for Environment Canada.
// Values taken from http://www8.tfe.umu.se/weather-new/js/index.js.
func windChill(T, V float64) float64 {
    vtmp := math.Pow(V, 0.16)
    return 13.126667 + 0.6215 * T + vtmp * (0.4875195 * T - 13.924748)
}

func main() {

    type WeatherData struct {
        Temperature string `xml:"root>tempmed"`
        WindSpeed string `xml:"root>vindh"`
    }

    // Weather data URL
    const url = "http://www8.tfe.umu.se/WeatherWebService2012/Service.asmx/Aktuellavarden"

    // HTTP GET
    res, _ := http.Get(url)

    // Read all XML data
    bytes, _ := ioutil.ReadAll(res.Body)
    res.Body.Close()

    xmlStr := formatXml(string(bytes));

    // Parse XML
    v := WeatherData{ Temperature: "none" }
    xml.Unmarshal([]byte(xmlStr), &v)

    // Prepare for conversion
    tempStr := convprep(v.Temperature)
    windStr := convprep(v.WindSpeed)

    temp, _ := strconv.ParseFloat(tempStr, 32)
    wind, _ := strconv.ParseFloat(windStr, 32)

    Twc := windChill(temp, wind)

    // Output
    fmt.Printf("Current temp: %.1f°C\nFeels like: %.1f°C\n", temp, Twc)
}
