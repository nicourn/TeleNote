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
	Chanel chan SendData
	Notes []Noteble
}

type MarshaledUser struct {
	Status string `json:"status"`
	Notes []MarshNotes `json:"notes"`
}

func NewUser(ID int, Name string, Chanel chan SendData) *User{
	return &User{ID, Name, "", Chanel, make([]Noteble, 0)}
}

func (u *User) Load() *User{
	files, err := ioutil.ReadDir("./Data")
	if err != nil{
		fmt.Errorf(err.Error())
		return u
	}
	for _, file := range files{
		if file.Name() == fmt.Sprintf("%v.json", u.ID){
			filedata, err := ioutil.ReadFile(fmt.Sprintf("./Data/%v", file.Name()))
			err_handler(err)
			var buff MarshaledUser
			err_handler(json.Unmarshal(filedata, &buff))
			u.Status = buff.Status
			umnotes := make([]Noteble, 0)
			for _, note := range buff.Notes{
				switch note.Type{
				case "Note": nnote := NewNote("Note", "")
					nnote.Unmarsh(note)
					umnotes = append(umnotes, nnote)
				case "Reminder": nrem := NewReminder("Note", "")
					nrem.Unmarsh(note)
					umnotes = append(umnotes, nrem)
				}
			}
			u.Notes = umnotes
			fmt.Println("Опа, ну тут я все помню)")
			return u
		}
	}
	fmt.Println("Ну подгружать то нечего(")
	return u
}


func (u *User) AddNote(n Noteble){
	u.Notes = append(u.Notes, n)
}

func (u *User) Save(){
	os.Mkdir("./Data", 0777)
	mnotes := make([]MarshNotes, 0)
	for _, note := range u.Notes{
		mnotes = append(mnotes, note.Marsh())
	}
	b, err := json.Marshal(MarshaledUser{u.Status, mnotes})
	err_handler(err)
	filename := fmt.Sprintf("./Data/%v.json", u.ID)
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
			u.Chanel <- SendData{u.ID,  fmt.Sprintf("%v. ", i + 1) + note.Get_note()}
		}
	} 
	u.Save()
}
