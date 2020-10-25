package main

import (
	"fmt"
	"strings"
	"time"
)

type Noteble interface{
	Get_note() string;
	Get_info() map[string][]string;
	Add_info(map[string]string) ;
}

type Note struct{
	Name string
	Text string
	Date string
}

func NewNote(Name, Text string) *Note {
	return &Note {Name, Text, time.Now().Format("2006-01-02 15:04:05")}
}

func ParseNote(text string) *Note {
	lines := strings.Split(text, "\n")
	if len(lines) > 1 && lines[0][0] == '.'{
		return NewNote(lines[0], strings.Join(lines[1:], ". "))
	} else {
		return NewNote("Note", text)
	}

}

func (n *Note) Get_note() string{
	return fmt.Sprintf("%v\n%v\nВремя создания: %v", n.Name, n.Text, n.Date)
}

func (n *Note) Get_info() map[string][]string{
	return map[string][]string{}
}

func (r *Note) Add_info(info map[string]string) {
	fmt.Println(info)
}


type Reminder struct{
	Note *Note
	Alert []string
	Date string
	Do bool
}

func NewReminder(Name, Text string) *Reminder{
	return &Reminder{ NewNote(Name, Text), make([]string, 0), "", false }
}


func (r *Reminder) Get_note() string{
	do := ""
	if r.Do {
		do = r.Date
	} else {
		do = "Не выполнено"
	}
	return fmt.Sprintf("%v\n%v\n%v", do, r.Note.Get_note(), r.Alert)
}

func (r *Reminder) Get_info() map[string][]string{
	return map[string][]string{"alerts": r.Alert}
}

func (r *Reminder) Add_info(info map[string]string) {
	if data, ok := info["alert"]; ok{
		r.Alert = append(r.Alert, data)
	}
}