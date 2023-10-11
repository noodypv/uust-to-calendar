package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

var (
	sFlag *string
)

type Calendar struct {
	events *[]Event
}

type Event struct {
	Date      string
	Type      string
	Name      string
	TimeStart string
	TimeEnd   string
	Auditory  string
	Teacher   string
}

func main() {
	sFlag = flag.String("u", "", "Link to your schedule page")

	flag.Parse()
	log.Println(*sFlag)
	u, err := url.Parse(*sFlag)
	if err != nil {
		log.Println(err)
		return
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Println(err)
		return
	}

	scheduleID, ok := m["schedule_semestr_id"]
	if !ok {
		log.Println("Invalid URL.")
		return
	}

	groupID, ok := m["student_group_id"]
	if !ok {
		log.Println("Invalid URL")
		return
	}

	log.Println(scheduleID[0], groupID[0])
	doc := getPage(scheduleID[0], groupID[0], 0)

	weeks := findMaxWeek(doc)
	log.Println("Парсим недель:", weeks)

	es := []Event{}

	for w := 1; w <= weeks; w++ {
		log.Println("Parsing week:", w)
		var prevWeekDay string
		doc = getPage(scheduleID[0], groupID[0], w)
		doc.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
			var startEndTime, weekday string
			e := Event{}
			s.Find("td").Each(func(j int, curr *goquery.Selection) {
				switch j {
				case 0:
					weekday = curr.Find("p").Text()
					if weekday != "" && weekday != prevWeekDay {
						prevWeekDay = weekday
					}

					weekday = prevWeekDay
				case 1:
					startEndTime = curr.Find("p").Text()
				case 2:
					e.Name = curr.Find("p").Text()
				case 3:
					e.Type = curr.Find("p").Text()
				case 4:
					e.Teacher = curr.Find("p").Text()
				case 5:
					e.Auditory = curr.Find("p").Text()
				case 6:
					re := regexp.MustCompile(`(\d+:\d+)`)
					times := re.FindAllString(startEndTime, -1)

					if len(times) < 2 {
						return
					}

					e.TimeStart = times[0] + "00"
					e.TimeEnd = times[1] + "00"

					rg := regexp.MustCompile(`(\d[^+]+)`)
					e.Date = rg.FindString(weekday)

					es = append(es, e)

				}
			})
		})
	}

	c := Calendar{
		events: &es,
	}

	c.CreateCalendar()
}

func findMaxWeek(doc *goquery.Document) int {
	num := doc.Find(".col-lg-10").Find("option").Last().Text()

	integer, err := strconv.Atoi(num)
	if err != nil {
		return -1
	}

	return integer
}

func getPage(semester, group string, week int) *goquery.Document {
	str := fmt.Sprintf("https://isu.ugatu.su/api/new_schedule_api/?schedule_semestr_id=%s&WhatShow=1&student_group_id=%s&weeks=%d", semester, group, week)

	client := &http.Client{}
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer res.Body.Close()

	utf8, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (c *Calendar) CreateCalendar() {
	f, err := os.Create("calendar.ics")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintln(f, "BEGIN:VCALENDAR")

	for _, e := range *c.events {
		fmt.Fprintln(f, "BEGIN:VEVENT")

		summary := e.Auditory + " " + e.Type + " " + e.Name + " " + e.Teacher
		fmt.Fprint(f, "SUMMARY:")
		fmt.Fprintln(f, summary)

		e.TimeStart = strings.ReplaceAll(e.TimeStart, ":", "")
		e.TimeEnd = strings.ReplaceAll(e.TimeEnd, ":", "")

		date, _ := time.Parse("02.01.2006", e.Date)
		e.Date = date.Format("2006.01.02")

		e.Date = strings.ReplaceAll(e.Date, ".", "")

		fmt.Fprintf(f, "DTSTART;VALUE=DATE-TIME:%sT%s", e.Date, e.TimeStart)
		fmt.Fprint(f, "\n")
		fmt.Fprintf(f, "DTEND;VALUE=DATE-TIME:%sT%s", e.Date, e.TimeEnd)
		fmt.Fprint(f, "\n")
		fmt.Fprintln(f, "END:VEVENT")
	}

	fmt.Fprint(f, "END:VCALENDAR")
}
