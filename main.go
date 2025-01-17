package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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

var genericNames = map[string]bool{
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

	var amountOfMalesNamedAsFemales = 0

	var amountOfPallywoodFromMOH = 0
	var amountOfPallywoodFromJudiciary = 0
	var amountOfPallywoodFromPublicSubmission = 0

	var amountOfWomen = 0
	var amountOfCivilians = 0
	var totalAmount = 0

	var amountOfSourcesFromMOH = 0
	var amountOfSourcesFromJudiciary = 0
	var amountOfSourcesFromPublicSubmission = 0

	for _, item := range data {

		totalAmount++

		if item.Source == "h" {
			amountOfSourcesFromMOH++
		}
		if item.Source == "c" {
			amountOfSourcesFromPublicSubmission++
		}
		if item.Source == "j" {
			amountOfSourcesFromJudiciary++
		}

		words := strings.Fields(item.Name)

		if item.Sex == "f" && len(words) > 0 && genericNames[words[0]] {
			fmt.Printf("Found record: ID: %s, Name: %s, Age: %d, Sex: %s, Dob: %s, Source: %s, EN_name: %s\n",
				item.Id, item.Name, item.Age, item.Sex, item.Dob, item.Source, item.En_name)

			amountOfMalesNamedAsFemales++

			if item.Source == "h" {
				amountOfPallywoodFromMOH++
			}
			if item.Source == "c" {
				amountOfPallywoodFromPublicSubmission++
			}
			if item.Source == "j" {
				amountOfPallywoodFromJudiciary++
			}
		}

		if item.Sex == "f" {
			amountOfWomen++
		}

		if item.Sex == "f" || (item.Sex == "m" && item.Age <= 10) {
			amountOfCivilians++
		}
	}

	fmt.Printf("This is the end of the PALLYWOOD names \n\n")

	fmt.Printf("This is the beginning of names without the population registry ID\n")

	// Define the regular expression pattern for exactly 9 digits
	pattern := `^\d{9}$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	var completeData = 0
	var incompleteData = 0

	for _, item := range data {
		if re.MatchString(item.Id) {
			completeData++
		} else {
			fmt.Printf("Found record: ID: %s, Name: %s, Age: %d, Sex: %s, Dob: %s, Source: %s, EN_name: %s\n",
				item.Id, item.Name, item.Age, item.Sex, item.Dob, item.Source, item.En_name)
			incompleteData++
		}
	}

	fmt.Printf("This is the end of names without the population registry ID\n\n")

	amountOfWomen = amountOfWomen - amountOfMalesNamedAsFemales

	var ProportionofWomenToTotal float32 = float32(amountOfWomen-amountOfMalesNamedAsFemales) / float32(totalAmount) * 100

	var ProportionOfCiviliansToTotal float32 = (float32(amountOfCivilians) - float32(amountOfMalesNamedAsFemales)) / float32(totalAmount) * 100

	fmt.Println()

	fmt.Printf("The amount of people with generic male names listed as having a Sex = f: %d\n\n", amountOfMalesNamedAsFemales)
	fmt.Printf("The amount of PALLYWOOD from the Ministry Of Health: %d\n", amountOfPallywoodFromMOH)
	fmt.Printf("The amount of PALLYWOOD from the Judicial Or House Committee: %d\n", amountOfPallywoodFromJudiciary)
	fmt.Printf("The amount of PALYYWOOD from the Public Submission: %d\n\n", amountOfPallywoodFromPublicSubmission)

	fmt.Printf("Total amount of people recorded: %d\n\n", totalAmount)

	fmt.Printf("Total amount of women recorded: %d\n", amountOfWomen)
	fmt.Printf("women:total = %f\n\n", ProportionofWomenToTotal)

	fmt.Printf("Total amount of civilians: %d\n", amountOfCivilians)
	fmt.Printf("civilian:total = %f\n\n", ProportionOfCiviliansToTotal)

	fmt.Printf("The margin of error for the total numbers: %f\n\n", float32(totalAmount)*0.02)

	fmt.Printf("Complete data: %d\n", completeData)
	fmt.Printf("Incomplete data: %d\n\n", incompleteData)

	fmt.Printf("The amount from the Ministry Of Health: %d\n", amountOfSourcesFromMOH)
	fmt.Printf("The amount from the Judicial Or House Committee: %d\n", amountOfSourcesFromJudiciary)
	fmt.Printf("The amount from the Public Submission: %d\n\n", amountOfSourcesFromPublicSubmission)

	fmt.Printf("MOH:Total = %f\n", float32(amountOfSourcesFromMOH)/float32(totalAmount)*100)
	fmt.Printf("Judicial Or House:Total = %f\n", float32(amountOfSourcesFromJudiciary)/float32(totalAmount)*100)
	fmt.Printf("Public:Total = %f\n", float32(amountOfSourcesFromPublicSubmission)/float32(totalAmount)*100)

	fmt.Printf("Estimated ratio of civilians to militants: %f\n", ProportionofWomenToTotal*2)
}
