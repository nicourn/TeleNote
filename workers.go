package main

import (
	"strconv"
	"strings"
	"time"
)
const (
	layout = "2006-01-02 15:04:05"
)

func RemindWorker(users *UserList, chanel chan SendData){
	l, err := time.LoadLocation("Europe/Kiev")
	err_handler(err)
	for true{
		users.mu.Lock()
		for _, user := range users.users{
			for _, note := range user.Notes {
				for _, alert := range note.Get_info()["alerts"]{
					t, err := time.ParseInLocation(layout, alert, l)
					err_handler(err)
					if time.Until(t).Seconds() < 30 && time.Until(t).Seconds() > -30{
						chanel <- SendData{user.ID, note.Get_note()}
						// TODO Отправить список след алертов 
					}
				}
			}
		}
		users.mu.Unlock()
		time.Sleep(60 * time.Second)
	}
}

func ParseTime(text string) time.Time {
	times := strings.ReplaceAll(text, ":", " ")
	times = strings.ReplaceAll(times, "-", " ")

	remtime := make([]int, 0) 
	for _, d := range strings.Split(times, " "){
		num, err := strconv.Atoi(d)
		err_handler(err)
		remtime = append(remtime, num)
	}

	l, err := time.LoadLocation("Europe/Kiev")
	err_handler(err)
	now := time.Now()

	switch len(remtime) {
	case 2: 
		return time.Date(now.Year(), now.Month(), now.Day(), remtime[0], remtime[1], 0, 0, l)
	
	case 3: 
		return time.Date(now.Year(), now.Month(), remtime[0], remtime[1], remtime[2], 0, 0, l)

	case 4: 
		return time.Date(now.Year(), time.Month(remtime[0]), remtime[1], remtime[2], remtime[3], 0, 0, l)
	}
	return time.Now().Add(24 * time.Hour)
}