package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	xmlparser "github.com/tamerh/xml-stream-parser"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type City struct {
	CityName,
	Street,
	House,
	Floor string
}

type TasksChan struct {
	floorCh,
	duplCh chan City
}

func parseXml(path string, tasksChan TasksChan, wg *sync.WaitGroup) {
	fmt.Println("start parse ")
	xmlFile, err := os.Open(path)
	defer wg.Done()
	defer close(tasksChan.duplCh)
	defer close(tasksChan.floorCh)

	if err != nil {
		log.Fatal(err)
	}

	defer xmlFile.Close()

	br := bufio.NewReaderSize(xmlFile, 64*1024)
	parser := xmlparser.NewXMLParser(br, "item")

	for Xml := range parser.Stream() {
		if Xml.Err != nil {
			log.Fatal(Xml.Err)
		}

		xmlAttr := Xml.Attrs

		city := City{
			xmlAttr["city"],
			xmlAttr["street"],
			xmlAttr["house"],
			xmlAttr["floor"],
		}

		tasksChan.floorCh <- city
		tasksChan.duplCh <- city
	}
}

func getAmountFloor(floorCh chan City, wg *sync.WaitGroup) {
	defer wg.Done()

	amountFloor := make(map[string]map[string]int)
	fmt.Println("start count floor")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Город", "1 этажей", "2 этажей", "3 этажей", "4 этажей", "5 этажей"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	for {
		if city, ok := <-floorCh; ok {
			if _, ok := amountFloor[city.CityName]; !ok {
				amountFloor[city.CityName] = map[string]int{
					"1": 0,
					"2": 0,
					"3": 0,
					"4": 0,
					"5": 0,
				}
			}

			amountFloor[city.CityName][city.Floor]++

		} else {
			break
		}
	}

	for val := range amountFloor {
		af := amountFloor[val]

		table.Append([]string{
			val,
			strconv.Itoa(af["1"]),
			strconv.Itoa(af["2"]),
			strconv.Itoa(af["3"]),
			strconv.Itoa(af["4"]),
			strconv.Itoa(af["5"]),
		})
	}

	defer table.Render()
	defer fmt.Println("Количество этажей в городах")
}

func findDuplicates(duplCh chan City, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("start find duplicates")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Город", "Улица", "№ Дома", "Этаж", "Повторов"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	temp := make(map[City]int)
	duplicates := make(map[City]int)

	for {
		if city, ok := <-duplCh; ok {
			_, ok := temp[city]
			if ok {
				temp[city]++
				duplicates[city] = temp[city]
			} else {
				temp[city] = 1
			}
		} else {
			break
		}
	}

	for val := range duplicates {
		table.Append([]string{
			val.CityName,
			val.Street,
			val.House,
			val.Floor,
			strconv.Itoa(duplicates[val]),
		})
	}

	defer table.Render()
	defer fmt.Println("Дубликаты")
}

func lineCounter(path string) (int, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	startTotalTime := time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(3)

	var path string
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Путь до *.xml файла не задан!")
		os.Exit(0)
	} else {
		path = args[1]
	}

	lineCount, err := lineCounter(path)
	if err != nil {
		log.Fatal(err)
	}

	tasksChan := TasksChan{
		make(chan City, lineCount),
		make(chan City, lineCount),
	}

	go parseXml(path, tasksChan, wg)
	go findDuplicates(tasksChan.duplCh, wg)
	go getAmountFloor(tasksChan.floorCh, wg)

	wg.Wait()
	fmt.Println("total time work program: ", time.Since(startTotalTime))
}
