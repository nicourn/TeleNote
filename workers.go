package main

import (
	"time"
)
const (
	layout = "2006-01-02 15:04:05"
)

func RemindWorker(users *UserList, chanel chan SendData){
	l, err := time.LoadLocation("Europe/Kiev")
	ErrHandler(err)
	for true{
		users.mu.Lock()
		for _, user := range users.users{
			for _, note := range user.Notes {
				for _, alert := range note.Get_info()["alerts"]{
					t, err := time.ParseInLocation(layout, alert, l)
					ErrHandler(err)
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
