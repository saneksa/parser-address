package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type CityList struct {
	XMLName xml.Name `xml:"root"`
	Items   []City   `xml:"item"`
}

type City struct {
	XMLName  xml.Name `xml:"item"`
	CityName string   `xml:"city,attr"`
	Street   string   `xml:"street,attr"`
	House    string   `xml:"house,attr"`
	Floor    string   `xml:"floor,attr"`
}

func parseXml(path string) CityList {
	xmlFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()

	var root CityList

	byteValue, _ := ioutil.ReadAll(xmlFile)
	err = xml.Unmarshal(byteValue, &root)

	if err != nil {
		log.Fatal(err)
	}
	return root
}

func getAmountFloor(cityList CityList) map[string]map[string]int {
	amountFloor := make(map[string]map[string]int)

	for _, city := range cityList.Items {
		if amountFloor[city.CityName] == nil {
			amountFloor[city.CityName] = map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
				"5": 0,
			}
		}

		amountFloor[city.CityName][city.Floor]++
	}
	return amountFloor
}

func findDuplicates(elements []City) map[City]int {
	temp := map[City]int{}
	res := make(map[City]int)

	for _, city := range elements {
		_, ok := temp[city]
		if ok {
			temp[city]++
			res[city] = temp[city]
		} else {
			temp[city] = 1
		}
	}
	return res
}

func main() {

	startTime := time.Now()
	fmt.Println("start parse ")

	cityList := parseXml("./address.xml")
	fmt.Println("end parse: ", time.Since(startTime))
	_ = getAmountFloor(cityList)
	fmt.Println("end count floor: ", time.Since(startTime))
	duplicates := findDuplicates(cityList.Items)
	fmt.Println("end find duplicates: ", time.Since(startTime))
	//fmt.Println(amountFloor)
	fmt.Println("Duplicates count: ", len(duplicates), "\n", duplicates)
	fmt.Println(time.Since(startTime))
}
