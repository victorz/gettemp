package main

import (
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "strings"
    "strconv"
    "fmt"
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

func main() {

    type WeatherData struct {
        Temperature string `xml:"root>tempmed"`
        WindSpeed string `xml:"root>vindh"`
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

    temp := convval(v.Temperature)
    windSpeed := convval(v.WindSpeed)
    windChill := convval(v.WindChill)

    // Output
    fmt.Printf("Current temp: %.1f°C\n", temp);

    if windSpeed < 1.2 {
        fmt.Println("Negligible wind chill.");
    } else {
        fmt.Printf("Feels like: %.1f°C\n", windChill);
    }

    fmt.Printf("Wind speed: %.1f m/s\n", windSpeed);
}
