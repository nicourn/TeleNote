package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func ParseTime(text string) time.Time {
	times := strings.ReplaceAll(text, ":", " ")
	times = strings.ReplaceAll(times, "-", " ")

	remtime := make([]int, 0) 
	for _, d := range strings.Split(times, " "){
		num, err := strconv.Atoi(d)
		ErrHandler(err)
		remtime = append(remtime, num)
	}

	l, err := time.LoadLocation("Europe/Kiev")
	ErrHandler(err)
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

func ErrHandler(err error){
	if err != nil{
		panic(err)
	}
}

type Config struct {
	Token string
}

func ParseConfig() Config {
	data, err := ioutil.ReadFile("./config.json")
	ErrHandler(err)
	var c Config
	err = json.Unmarshal(data, &c)
	ErrHandler(err)
	fmt.Println(c)
	return c 
}