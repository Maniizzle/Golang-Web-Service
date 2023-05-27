package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"reading.maniizzle.io/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "environment: %s\n", "dev")
	// fmt.Fprintf(w, "version: %s\n", "1.0.0")

}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// books := []data.Book{
		// 	{
		// 		ID:        1,
		// 		CreatedAt: time.Now(),
		// 		Title:     "Echeos in the light",
		// 		Published: 2019,
		// 		Pages:     300,
		// 		Genres:    []string{"fiction", "thriller"},
		// 		Rating:    4.5,
		// 		Version:   1,
		// 	},
		// 	{
		// 		ID:        1,
		// 		CreatedAt: time.Now(),
		// 		Title:     "Deep waters",
		// 		Published: 2020,
		// 		Pages:     300,
		// 		Genres:    []string{"mystery", "documentary"},
		// 		Rating:    4.7,
		// 		Version:   1,
		// 	},
		// }
		books, err := app.models.Books.GetAll()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//fmt.Fprintln(w, "Display a list of books on the reading")
	}

	if r.Method == http.MethodPost {
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"pubished"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float32  `json:"rating"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		err = json.Unmarshal(body, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		book := &data.Book{
			Title:     input.Title,
			Published: input.Published,
			Pages:     input.Pages,
			Genres:    input.Genres,
			Rating:    input.Rating,
		}

		if err := app.models.Books.Insert(book); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/books/%d", book.ID))

		err = app.writeJSON(w, http.StatusCreated, envelope{"books": book}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		//fmt.Fprintf(w, "%v\n", input)
	}
}

func (app *application) getUpdateDeleteBooksHanler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {

	// }
	switch r.Method {
	case http.MethodGet:
		app.getBook(w, r)
	case http.MethodPut:
		app.updateBook(w, r)
	case http.MethodDelete:
		app.deleteBook(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "BAd reqeuest", http.StatusBadRequest)
		return
	}

	book, err := app.models.Books.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	// book := data.Book{
	// 	ID:        idInt,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Echeos in the light",
	// 	Published: 2019,
	// 	Pages:     300,
	// 	Genres:    []string{"fiction", "thriller"},
	// 	Rating:    4.5,
	// 	Version:   1,
	// }

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintln(w, "Return the etails of book with Id: %d", idInt)

}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	book, err := app.models.Books.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"pubished"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	// when working with static data
	// book := data.Book{
	// 	ID:        idInt,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Echeos in the light",
	// 	Published: 2019,
	// 	Pages:     300,
	// 	Genres:    []string{"fiction", "thriller"},
	// 	Rating:    4.5,
	// 	Version:   1,
	// }

	err = app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Published != nil {
		book.Published = *input.Published
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if len(input.Genres) > 0 {
		book.Title = *input.Title
	}
	if input.Rating != nil {
		book.Rating = *input.Rating
	}

	err = app.models.Books.Update(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "%v\n", book)

}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "BAd reqeuest", http.StatusBadRequest)
	}

	err = app.models.Books.Delete(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "book deleted succesfully"}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintln(w, "deleting  the details of book with Id: %d", idInt)

}
