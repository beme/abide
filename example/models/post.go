package models

import (
	"fmt"
)

// Post is a post.
type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (p *Post) String() string {
	return fmt.Sprintf("%s\n%s", p.Title, p.Body)
}
