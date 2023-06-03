package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Book struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Published int      `json:"pubished"`
	Pages     int      `json:"pages,omitempty"`
	Genres    []string `json:"genres,omitempty"`
	Rating    float32  `json:"rating,omitempty"`
	Version   int32    `json:"-"`
}

type BookResponse struct {
	Book *Book `json:"book"`
}

type BooksResponse struct {
	Book *[]Book `json:"books"`
}

type ReadingListModel struct {
	Endpoint string
}

func (m *ReadingListModel) GetAll() (*[]Book, error) {

	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var booksResp BooksResponse
	err = json.Unmarshal(data, &booksResp)
	if err
}
