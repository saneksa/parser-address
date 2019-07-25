package main

import (
	"bufio"
	"fmt"
	xmlparser "github.com/tamerh/xml-stream-parser"
	"log"
	"os"
	"sync"
	"time"
)

type City struct {
	CityName string
	Street   string
	House    string
	Floor    string
}

func parseXml(path string) []City {
	startTime := time.Now()
	fmt.Println("start parse ")
	xmlFile, err := os.Open(path)

	var arr []City

	if err != nil {
		log.Fatal(err)
	}

	br := bufio.NewReaderSize(xmlFile, 65536)
	parser := xmlparser.NewXMLParser(br, "item")

	for Xml := range parser.Stream() {
		arr = append(arr, City{
			Xml.Attrs["city"],
			Xml.Attrs["street"],
			Xml.Attrs["house"],
			Xml.Attrs["floor"],
		})
	}

	defer xmlFile.Close()

	fmt.Println("end parse", time.Since(startTime))

	return arr
}

func getAmountFloor(cityList []City, wg *sync.WaitGroup) {
	defer wg.Done()
	amountFloor := make(map[string]map[string]int)
	startTime := time.Now()
	fmt.Println("start count floor")

	for _, city := range cityList {
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

	defer fmt.Println("end count floor: ", time.Since(startTime))
	defer fmt.Println("result ", amountFloor)
}

func findDuplicates(cityList []City, wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	fmt.Println("start find duplicates")

	temp := make(map[City]int)
	duplicates := make(map[City]int)

	for _, city := range cityList {
		_, ok := temp[city]
		if ok {
			temp[city]++
			duplicates[city] = temp[city]
		} else {
			temp[city] = 1
		}
	}
	defer fmt.Println("end find duplicates: ", time.Since(startTime))
	defer fmt.Println("Duplicates count: ", duplicates)

}

func main() {
	startTotalTime := time.Now()
	wg := new(sync.WaitGroup)
	cityList := parseXml("./address.xml")
	wg.Add(2)

	go findDuplicates(cityList, wg)
	go getAmountFloor(cityList, wg)

	wg.Wait()
	fmt.Println("total time work program: ", time.Since(startTotalTime))

}
