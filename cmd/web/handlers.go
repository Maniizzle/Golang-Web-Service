package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Home page")

}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "View a single book")

}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Home page")

}
