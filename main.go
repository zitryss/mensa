package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type item struct {
	name  string
	price string
}

func main() {
	file := downloadMenu()
	menu := readMenu(file)
	printMenu(menu)
}

func downloadMenu() io.ReadCloser {
	_, week := time.Now().ISOWeek()
	resp, err := http.Get("http://www.stwno.de/infomax/daten-extern/csv/UNI-P/" + strconv.Itoa(week) + ".csv")
	if err != nil {
		log.Fatalln(err)
	}
	return resp.Body
}

func readMenu(file io.ReadCloser) map[string][]item {
	defer file.Close()
	menu := make(map[string][]item)
	curDate := time.Now().Format("02.01.2006")
	r := csv.NewReader(charmap.ISO8859_1.NewDecoder().Reader(file))
	r.Comma = ';'
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		date := record[0]
		group := defineGroup(record[2])
		name := record[3]
		price := record[5]
		if curDate == date {
			for _, name := range splitFields(name) {
				menu[group] = append(menu[group], item{name, price})
			}
		}
	}
	return menu
}

func defineGroup(s string) string {
	switch s[0] {
	case 'S':
		return "Suppe"
	case 'H':
		return "Hauptgericht"
	case 'B':
		return "Beilage"
	case 'N':
		return "Nachtisch"
	default:
		return s
	}
}

func splitFields(s string) []string {
	tmp := regexp.MustCompile(`\s\(.+?\)`).Split(s, -1)
	return tmp[:len(tmp)-1]
}

func printMenu(menu map[string][]item) {
	maxLength := maxFieldLength(menu)
	for _, group := range []string{"Suppe", "Hauptgericht", "Beilage", "Nachtisch"} {
		items, ok := menu[group]
		if ok {
			fmt.Println(group)
			for _, item := range items {
				fmt.Print("	", item.name)
				for i := 0; i < maxLength-stringLength(item.name); i++ {
					fmt.Print(" ")
				}
				fmt.Println("	", item.price)
			}
			fmt.Println()
		}
	}
}

func maxFieldLength(menu map[string][]item) int {
	maxLength := 0
	for _, items := range menu {
		for _, item := range items {
			l := stringLength(item.name)
			if maxLength < l {
				maxLength = l
			}
		}
	}
	return maxLength
}

func stringLength(s string) int {
	counter := 0
	for range s {
		counter++
	}
	return counter
}
