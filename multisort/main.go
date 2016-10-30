package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
)

type Data struct {
	UserId int
	Id     int
	Tİtle  string
	Body   string
}

func areYouAwake(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am awake!")
}

func getData(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	v := []Data{}
	err = decoder.Decode(&v)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := json.Marshal(v)

	if (err != nil) {
		panic(err)
	}

	w.Write(b)
}

func main() {
	http.HandleFunc("/areyouawake", areYouAwake)
	http.HandleFunc("/getdata", getData)
	http.ListenAndServe(":8080", nil)
}
