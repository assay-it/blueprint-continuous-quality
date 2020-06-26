package main

import (
	"github.com/assay-it/tk"
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

//
//
//

type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type List []News

func (seq List) Len() int                { return len(seq) }
func (seq List) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq List) Less(i, j int) bool      { return seq[i].ID < seq[j].ID }
func (seq List) String(i int) string     { return seq[i].ID }
func (seq List) Value(i int) interface{} { return seq[i] }

//
//
//

var host = tk.Env("HOST", "")

//
//
func TestNewsJSON() gurl.Arrow {
	var seq List
	item := News{ID: "2", Title: "Sed luctus tortor sit amet eros eleifend cursus."}

	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&seq),
		ƒ.Seq(&seq).Has(item.ID, item),
	)
}

//
//
func TestNewsHTML() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		ø.Accept().Is("text/html"),
		ƒ.Code(200),
		ƒ.Served().Is("text/html"),
	)
}

//
func main() {}
