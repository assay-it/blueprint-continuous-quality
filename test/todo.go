package main

import (
	"fmt"

	"github.com/assay-it/tk"
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

var host = tk.Env("HOST", "")

type TODO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TODOs []TODO

func (seq TODOs) Len() int                { return len(seq) }
func (seq TODOs) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq TODOs) Less(i, j int) bool      { return seq[i].ID < seq[j].ID }
func (seq TODOs) String(i int) string     { return seq[i].ID }
func (seq TODOs) Value(i int) interface{} { return seq[i] }

//
//
func TestList() gurl.Arrow {
	var seq TODOs
	item := TODO{ID: "#1", Title: "study assay.it"}

	return gurl.HTTP(
		ø.GET("https://%s/todo", host),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&seq),
		ƒ.Seq(&seq).Has(item.ID, item),
	)
}

//
//
func TestLookup() gurl.Arrow {
	var item TODO

	return gurl.HTTP(
		ø.GET("https://%s/todo/%s", host, item.ID),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
		ƒ.Value(&item.ID).String("#1"),
		ƒ.Value(&item.Title).String("study assay.it"),
	)
}

//
//
func TestNotFound() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/todo/%s", host, "unknown"),
		ƒ.Code(404),
	)
}

//
//
func TestLifeCycle() gurl.Arrow {
	origin := TODO{ID: "#2", Title: "Write Unit Tests"}
	remote := TODO{ID: "#2"}

	return gurl.Join(
		append(origin),
		lookup(&remote),
		ƒ.FMap(func() error {
			if origin.Title != remote.Title {
				return fmt.Errorf("origin item %v do not match remote %v", origin, remote)
			}
			return nil
		}),
		contain(origin),
		remove(remote),
	)
}

//
func contain(item TODO) gurl.Arrow {
	var seq TODOs

	return gurl.HTTP(
		ø.GET("https://%s/todo", host),
		ƒ.Code(200),
		ƒ.Recv(&seq),
		ƒ.Seq(&seq).Has(item.ID, item),
	)
}

//
func append(item TODO) gurl.Arrow {
	return gurl.HTTP(
		ø.POST("https://%s/todo", host),
		ø.ContentJSON(),
		ø.Send(item),
		ƒ.Code(200),
	)
}

//
func lookup(item *TODO) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/todo/%s", host, item.ID),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(item),
	)
}

//
func remove(item TODO) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/todo/%s", host, item.ID),
		ƒ.Code(200),
	)
}

//
func main() {}
