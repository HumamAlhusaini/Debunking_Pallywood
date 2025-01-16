package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

type gaza_list struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Dob     string `json:"dob"`
	Sex     string `json:"sex"`
	Age     int    `json:"age"`
	Source  string `json:"source"`
	En_name string `json:"en_name"`
}

var data []*gaza_list

func main() {
	list, err := http.Get("https://data.techforpalestine.org/api/v2/killed-in-gaza.json")
	if err != nil {
		fmt.Println("We've had an error in getting JSON!")
	} else {
		fmt.Println("success!")
	}

	fmt.Println("Content-Type:", list.Header.Get("Content-Type"))

	defer list.Body.Close()

	body, err := io.ReadAll(list.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if utf8.Valid(body) {
		fmt.Println("The response body is valid UTF-8 encoded.")
	} else {
		fmt.Println("The response body is NOT UTF-8 encoded.")
	}

	var count = 0

	genericNames := map[string]bool{
		"محمد":    true,
		"أحمد":    true,
		"علي":     true,
		"حسن":     true,
		"عبد":     true,
		"محمود":   true,
		"يوسف":    true,
		"إبراهيم": true,
		"خالد":    true,
		"سعيد":    true,
		"عمر":     true,
		"إسماعيل": true,
		"جمال":    true,
		"صالح":    true,
		"سامي":    true,
		"نادر":    true,
		"رشيد":    true,
		"حاتم":    true,
		"يحيى":    true,
		"زياد":    true,
		"ماهر":    true,
		"حسين":    true,
		"فارس":    true,
		"رامي":    true,
		"أمين":    true,
		"كمال":    true,
		"وائل":    true,
		"علاء":    true,
		"حمزة":    true,
		"شريف":    true,
		"منير":    true,
		"عادل":    true,
		"باسم":    true,
		"طارق":    true,
		"سامر":    true,
		"نبيل":    true,
		"هيثم":    true,
		"عدنان":   true,
		"مازن":    true,
		"عمار":    true,
		"مالك":    true,
		"إياد":    true,
		"أنس":     true,
		"مراد":    true,
		"سيف":     true,
		"بدر":     true,
		"راشد":    true,
		"منصور":   true,
		"عصام":    true,
		"زكريا":   true,
		"غسان":    true,
		"مصطفى":   true,
		"رائد":    true,
		"سامح":    true,
		"نزار":    true,
		"وهيب":    true,
		"أيمن":    true,
	}

	for _, item := range data {
		words := strings.Fields(item.Name)

		if item.Sex == "f" && len(words) > 0 && genericNames[words[0]] {
			fmt.Printf("Found record: ID: %s, Name: %s, Age: %d, Sex: %s, Dob: %s, Source: %s, EN_name: %s\n",
				item.Id, item.Name, item.Age, item.Sex, item.Dob, item.Source, item.En_name)

			count++
		}
	}

	fmt.Println(count)
}
