package main

import (
	"fmt"
)

type Post struct {
	Title string
	Body  string
}

func (p *Post) String() string {
	return fmt.Sprintf("%s\n%s", p.Title, p.Body)
}
