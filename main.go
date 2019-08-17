package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Dictionary struct {
	Word string
	Explanation interface{}
}

var (
	dictionary map[string][]Dictionary
	keys []string
)

func main() {
	jsonFile, err := os.Open("dictionary.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened dictionary.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	dictionary = make(map[string][]Dictionary)

	for word, explanation := range result {
		dictionary[strings.ToUpper(string(word[0]))] = append(dictionary[strings.ToUpper(string(word[0]))], Dictionary{word, explanation})
		if IfExists(strings.ToUpper(string(word[0]))) {
			keys = append(keys, strings.ToUpper(string(word[0])))
		}
	}

	sort.Strings(keys)

	r := mux.NewRouter()
	r.HandleFunc("/letter/{letter}", LetterHandler)
	r.HandleFunc("/all", GetAll)
	r.HandleFunc("/keys", GetKeys)

	fmt.Println("Done reading. Service is ready at 127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}

func IfExists(key string) bool {
	for _, k := range keys {
		if k == key {
			return false
		}
	}
	return true
}

func GetKeys(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(keys)
}

func GetAll(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dictionary)
}

func LetterHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dictionary[strings.ToUpper(vars["letter"])])
}
