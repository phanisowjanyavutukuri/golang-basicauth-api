package main

import (
    "encoding/json"
    _"io"
    "log"
    "net/http"
	"encoding/base64"
    "strings"

      _ "fmt"
    "io/ioutil"
    "rest/data"
    "github.com/gorilla/mux"
)
type Person struct {
        ID        string   `json:"id,omitempty"`
        Firstname string   `json:"first,omitempty"`
        Lastname  string   `json:"lastname,omitempty"`
        Address   *Address `json:"address,omitempty"`
}
type ResponsePerson struct {
        ID        string   `json:"id,omitempty"`
        Firstname string   `json:"first,omitempty"`
        Lastname  string   `json:"lastname,omitempty"`
        Address   *Address `json:"address,omitempty"`
        Message   string   `json:"message"`
}

type Address struct {
        City  string `json:"city,omitempty"`
        State string `json:"state,omitempty"`
}

var people []Person

func main() {
	router := mux.NewRouter()

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
        people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
        people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})

    // public views
   router.HandleFunc("/", GetPeople)

    // private views
   router.HandleFunc("/people", PostOnly(BasicAuth(CreatePerson)))
    router.HandleFunc("/person/{id}", GetOnly(BasicAuth(GetPerson)))

    log.Fatal(http.ListenAndServe(":8080", router))
}




func GetPeople(w http.ResponseWriter, r *http.Request) {
        buffer, _ := json.Marshal(people)
        w.Write(buffer)
        data.DisplayAll()
}
func GetPerson(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        for _, items := range people {
                if items.ID == params["id"] {
                        json.NewEncoder(w).Encode(items)
                }
        }
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
        var personDetails Person
        buffer, _ := ioutil.ReadAll(r.Body)
        json.Unmarshal(buffer, &personDetails)
        responsePerson := Person{
                ID:        personDetails.ID,
                Firstname: personDetails.Firstname,
                Lastname:  personDetails.Lastname,
                Address:   personDetails.Address,
        }
        buffer, _ = json.Marshal(responsePerson)
      data.Insert(buffer)
      w.Write(buffer)
}



type handler func(w http.ResponseWriter, r *http.Request)

func BasicAuth(pass handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {

        auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

        if len(auth) != 2 || auth[0] != "Basic" {
            http.Error(w, "authorization failed", http.StatusUnauthorized)
            return
        }

        payload, _ := base64.StdEncoding.DecodeString(auth[1])
        pair := strings.SplitN(string(payload), ":", 2)

        if len(pair) != 2 || !validate(pair[0], pair[1]) {
            http.Error(w, "authorization failed", http.StatusUnauthorized)
            return
        }

        pass(w, r)
    }
}

func validate(username, password string) bool {
    if username == "test" && password == "test" {
        return true
    }
    return false
}





func GetOnly(h handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            h(w, r)
            return
        }
        http.Error(w, "get only", http.StatusMethodNotAllowed)
    }
}

func PostOnly(h handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            h(w, r)
            return
        }
        http.Error(w, "post only", http.StatusMethodNotAllowed)
    }
}
