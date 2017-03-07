package main

import (
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "strings"
    "strconv"
    "fmt"
)

// Prepares string before float conversion by replacing decimal comma with
// decimal point and trimming white space.
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

func main() {

    type WeatherData struct {
        Temperature string `xml:"root>tempmed"`
        WindChill string `xml:"root>windChill"`
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
    windStr := convprep(v.WindChill)

    temp, _ := strconv.ParseFloat(tempStr, 32)
    windChill, _ := strconv.ParseFloat(windStr, 32)

    // Output
    fmt.Printf("Current temp: %.1f°C\nFeels like: %.1f°C\n", temp, windChill)
}
