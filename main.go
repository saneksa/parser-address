package main

import (
	"bufio"
	"fmt"
	"sync"
	//"github.com/olekukonko/tablewriter"
	xmlparser "github.com/tamerh/xml-stream-parser"
	"log"
	"os"
	"time"
)

type City struct {
	CityName string
	Street   string
	House    string
	Floor    string
}

func parseXml(path string, duplCh, floorCh chan []City, wg *sync.WaitGroup) {
	defer wg.Done()
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
		xmlAttr := Xml.Attrs

		arr = append(arr, City{
			xmlAttr["city"],
			xmlAttr["street"],
			xmlAttr["house"],
			xmlAttr["floor"],
		})
	}

	defer xmlFile.Close()

	fmt.Println("end parse", time.Since(startTime))

	duplCh <- arr
	floorCh <- arr
}

func getAmountFloor(floorCh chan []City, wg *sync.WaitGroup) {
	defer wg.Done()
	cityList := <-floorCh
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

func findDuplicates(duplCh chan []City, wg *sync.WaitGroup) {
	defer wg.Done()
	cityList := <-duplCh
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

	wg.Add(3)
	floorCh := make(chan []City, 1)
	duplCh := make(chan []City, 1)

	go parseXml("./address.xml", duplCh, floorCh, wg)
	go findDuplicates(duplCh, wg)
	go getAmountFloor(floorCh, wg)

	wg.Wait()
	fmt.Println("total time work program: ", time.Since(startTotalTime))

}
