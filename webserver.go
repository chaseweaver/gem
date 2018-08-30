// Chase Weaver

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/chilts/sid"
	"github.com/chyeh/pubip"
)

// Client contains vairables used for html page
type Client struct {
	ClientID string
	ClientIP string
	Text     string
}

func main() {
	http.HandleFunc("/", homePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// homePage(http.ResponseWriter, *http.Request)
// Handles loading and parsing of HTML page
func homePage(w http.ResponseWriter, r *http.Request) {

	// Parse HTML page template
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print(err)
	}

	// Get outward ClientIP
	ip, err := pubip.Get()
	if err != nil {
		log.Println(err)
	}

	// Execute template and with variables
	err = t.Execute(w, Client{
		ClientID: sid.Id(),
		ClientIP: fmt.Sprintf("%v", ip),
	})
	if err != nil {
		log.Print(err)
	}
}
