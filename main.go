package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Data struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

var (
	datas  []Data
	nextId int
)

const fileName = "data.json"

func main() {
	Loader()

	http.HandleFunc("/", Handler)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/remove", RemoveTask)

	log.Println("Sever at port: http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalln("ERROR WHILE OPENING THE PAGE: ", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("Template/index.html")
	if err != nil {
		http.Error(w, "ERROR WHILE PARING THE FILE: ", http.StatusBadRequest)
		return
	}
	err = temp.Execute(w, datas)
	if err != nil {
		http.Error(w, "ERROR WHILE EXECUTING THE TEMPLATE: ", http.StatusInternalServerError)
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		itemName := r.FormValue("description")
		itemAmount := r.FormValue("amount")
		if itemName == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		amount, err := strconv.Atoi(itemAmount)
		if err != nil {
			http.Error(w, "ERROR WHILE CONVERTING: ", http.StatusBadRequest)
			return
		}

		data := Data{
			ID:          nextId,
			Date:        time.Now().Format(time.Stamp),
			Description: itemName,
			Amount:      amount,
		}
		datas = append(datas, data)
		nextId++
		Saver()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Saver() {
	data, err := json.MarshalIndent(datas, "", " ")
	if err != nil {
		log.Fatalln("ERROR WHILE MARSHALLING DATA: ", err)
		return
	}
	err = os.WriteFile(fileName, data, 0o644)
	if err != nil {
		log.Fatalln("ERROR WHILE WRITING THE FILE: ", err)
		return
	}
}

func Loader() {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			datas = []Data{}
			nextId = 1
			return
		}
		log.Fatalln("ERROR WHILE READING THE FILE: ", err)
	}
	err = json.Unmarshal(file, &datas)
	if err != nil {
		log.Fatalln("ERROR WHILE UNMARSHALING THE DATA: ", err)
		return
	}
	for _, data := range datas {
		if data.ID >= nextId {
			nextId = data.ID + 1
		}
	}
}

func RemoveTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		dataId := r.FormValue("id")
		id, err := strconv.Atoi(dataId)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		for i, data := range datas {
			if data.ID == id {
				datas = append(datas[:i], datas[i+1:]...)
			}
		}
		Saver()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
