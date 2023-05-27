package data

import "time"

type Book struct {
	ID int64 `json:"id"`
	//"-" means omit it from the marshalled value (struct tags)
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	// omitempty makes it optional
	Published int      `json:"pubished,omitempty"`
	Pages     int      `json:"pages,omitempty"`
	Genres    []string `json:"genres,omitempty"`
	Rating    float32  `json:"rating,omitempty"`
	Version   int32    `json:"-"`
}
