package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var keyvalues map[string]string

func handleRequests() {
	// Create a Router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/keys", getAllKeys)
	router.HandleFunc("/key/{id}", upsertKey).Methods("POST")
	router.HandleFunc("/key/{id}", deleteKey).Methods("DELETE")
	router.HandleFunc("/key/{id}", getKey)

	fmt.Println("Running on port 10000")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home Page!")
	fmt.Println("Endpoint Hit: homePage")
}

func getAllKeys(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(keyvalues)
}

func getKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Fprintf(w, keyvalues[key])
}

func upsertKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	keyvalues[key] = string(reqBody)
	writeMapToFile()
	fmt.Fprintf(w, "Upsert success")
}

func deleteKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	delete(keyvalues, key)
	writeMapToFile()
	fmt.Fprintf(w, "Delete success")
}

func main() {
	fmt.Println("Simple Key Value store v1.0 (c)2021 Andrew Ferguson")
	keyvalues = readFileIntoMap()
	fmt.Println("Loaded values:")
	for key, element := range keyvalues {
		fmt.Println("Key:", key, "=>", "Value:", element)
	}
	handleRequests()
}

// Load values from File into Map
func readFileIntoMap() map[string]string {
	m := make(map[string]string)

	f, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		kv := strings.Split(scanner.Text(), "=")
		m[kv[0]] = kv[1]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m
}

// Write values from Map into File
func writeMapToFile() {
	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	s := ""
	for key, element := range keyvalues {
		s = s + key + "=" + element + "\n"
	}

	_, err2 := f.WriteString(s)

	if err2 != nil {
		log.Fatal(err2)
	}
}
