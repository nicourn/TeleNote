package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	tgbot "github.com/Syfaro/telegram-bot-api"
)
type SendData struct{
	ID int
	Data string
}

type UserList struct{
	mu sync.Mutex
	users []*User
}

func NewUserList() *UserList{
	return &UserList{sync.Mutex{}, make([]*User, 0)}
}

type User struct{
	ID int
	Name string
	Status string
	Chanel chan SendData `json: ""`
	Notes []Noteble
}

type MarshaledUser struct {
	ID int `json:"id"`
	Notes []Noteble `json:"notes"`
}

func NewUser(ID int, Name string, Chanel chan SendData) *User{
	return &User{ID, Name, "", Chanel, make([]Noteble, 0)}
}

func (u *User) AddNote(n Noteble){
	u.Notes = append(u.Notes, n)
}

func (u *User) Save(){
	b, err := json.Marshal(MarshaledUser{u.ID, u.Notes})
	err_handler(err)
	filename := fmt.Sprintf("%v.json", u.ID)
	os.Chmod(filename, 0777)
	err_handler(ioutil.WriteFile(filename, b, 0777))
}


func (u *User) MessageHandler(message tgbot.Update){
	switch u.Status{
	case "change name":
		if message.Message.Text != ""{
			u.Name = message.Message.Text
			u.Chanel <- SendData{ u.ID, fmt.Sprintf("Теперь вас зовут %v", message.Message.Text) }
		}
		return
	case "new note": 
		u.AddNote(ParseNote(message.Message.Text))
		u.Status = ""
	case "new remind1":
		note := ParseNote(message.Message.Text)
		u.AddNote(NewReminder(note.Name, note.Text))
		u.Chanel <- SendData{ u.ID, "Отправте напоминания или '.' что бы закончить" }
		u.Status = "new remind2"
	case "new remind2":
		if message.Message.Text == "."{
			u.Status = ""
			return
		}
		u.Notes[len(u.Notes) - 1].Add_info(map[string]string{"alert": ParseTime(message.Message.Text).Format("2006-01-02 15:04:05") })
	}

	switch message.Message.Text {
	case "/new_note":
		u.Status = "new note"
		u.Chanel <- SendData{u.ID, fmt.Sprintf("Напишите заметку: ")}
	
	case "/new_remind":
		u.Status = "new remind1"
		u.Chanel <- SendData{u.ID, fmt.Sprintf("Напишите напоминание: ")}


	case "/get_notes":
		for i, note := range u.Notes{
			u.Chanel <- SendData{u.ID,  fmt.Sprintf("%v)", i + 1) + note.Get_note()}
		}
	} 
	u.Save()
}
