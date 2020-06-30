package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/path"
)

//
// TODO is a type of elements our application deals with
type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

//
// tToDo implements REST API to managed the list of TODO items
type tNews struct {
	list map[string]string
}

//
// lambda factory
func main() {
	api := tNews{list: map[string]string{
		"1": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"2": "Sed luctus tortor sit amet eros eleifend cursus.",
		"3": "Proin volutpat leo eu dui tristique, sit amet aliquet diam molestie.",
		"4": "In in odio vel velit commodo ultrices.",
		"5": "Nulla quis neque pulvinar, mollis libero in, varius libero.",
	}}

	lambda.Start(
		µ.Serve(
			api.NewsHTML(),
			api.NewsJSON(),
			api.ItemHTML(),
			api.ItemJSON(),
		),
	)
}

//
// NewsHTML endpoint
//
// GET /news
func (todo *tNews) NewsHTML() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("news")),
		µ.Header(header.Is("Accept", "text/html")),
		µ.FMap(func() error {
			seq := []string{}

			seq = append(seq, "<ul>")
			for id, title := range todo.list {
				seq = append(seq, fmt.Sprintf("<li>%s: %s</li>", id, title))
			}
			seq = append(seq, "</ul>")

			return µ.Ok().Text(strings.Join(seq, "\n")).With("Content-Type", "text/html")
		}),
	)
}

//
// NewsJSON endpoint
//
// GET /news
func (todo *tNews) NewsJSON() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("news")),
		µ.FMap(func() error {
			seq := []News{}

			for id, title := range todo.list {
				seq = append(seq, News{ID: id, Title: title})
			}

			return µ.Ok().JSON(seq)
		}),
	)
}

//
// ItemHTML endpoint
//
// GET /news/:id
func (todo *tNews) ItemHTML() µ.Endpoint {
	var id string

	return µ.GET(
		µ.Path(path.Is("news"), path.String(&id)),
		µ.Header(header.Is("Accept", "text/html")),
		µ.FMap(func() error {
			item, ok := todo.list[id]
			if !ok {
				return µ.NotFound(fmt.Errorf(id))
			}

			return µ.Ok().Text(fmt.Sprintf("<h1>%s: %s</h1>", id, item)).With("Content-Type", "text/html")
		}),
	)
}

//
// ItemJSON endpoint
//
// GET /news/:id
func (todo *tNews) ItemJSON() µ.Endpoint {
	var id string

	return µ.GET(
		µ.Path(path.Is("news"), path.String(&id)),
		µ.FMap(func() error {
			item, ok := todo.list[id]
			if !ok {
				return µ.NotFound(fmt.Errorf(id))
			}

			return µ.Ok().JSON(News{ID: id, Title: item})
		}),
	)
}
