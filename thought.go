package main

import (
	"html/template"
	"time"
)

type Thought struct {
	TID      int64
	Title    string
	Date     time.Time
	HTML     template.HTML
	Markdown string
}
