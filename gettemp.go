package main

import (
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "strings"
    "strconv"
    "fmt"
)

func main() {

    type WeatherData struct {
        Temperature string `xml:"root>tempmed"`
    }

    v := WeatherData{ Temperature: "none" }

    // Weather data URL
    const url = "http://www8.tfe.umu.se/WeatherWebService2012/Service.asmx/Aktuellavarden"

    // HTTP GET
    res, _ := http.Get(url)

    // Read all XML data
    bytes, _ := ioutil.ReadAll(res.Body)
    res.Body.Close()

    // Convert HTML entity angle brackets to actual angle brackets
    xmlStr := strings.Replace(string(bytes), "&lt;", "<", -1)
    xmlStr = strings.Replace(xmlStr, "&gt;", ">", -1)

    // Parse XML
    xml.Unmarshal([]byte(xmlStr), &v)

    // Prepare for conversion
    tempStr := strings.Replace(v.Temperature, ",", ".", 1)
    tempStr = strings.Trim(tempStr, " ")

    temp, _ := strconv.ParseFloat(tempStr, 32)

    // Output
    fmt.Printf("Current temp: %.1f°C\n", temp)
}
