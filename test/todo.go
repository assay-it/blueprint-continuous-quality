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
//

var host = tk.Env("HOST", "")

//
// Elementary Tests - ...
//

//
//
func TestList() gurl.Arrow {
	var seq TODOs
	item := TODO{ID: "1", Title: "study assay.it"}

	return gurl.HTTP(
		ø.GET("https://%s/api/todo", host),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&seq),
		ƒ.Seq(&seq).Has(item.ID, &item),
	)
}

//
//
func TestLookup() gurl.Arrow {
	var item TODO

	return gurl.HTTP(
		ø.GET("https://%s/api/todo/%s", host, "1"),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
		ƒ.Value(&item.ID).String("1"),
		ƒ.Value(&item.Title).String("study assay.it"),
	)
}

//
//
func TestUpdate() gurl.Arrow {
	item := TODO{ID: "3", Title: "write test suites for your app"}
	recv := TODO{}

	return gurl.HTTP(
		ø.PUT("https://%s/api/todo/%s", host, item.ID),
		ø.ContentJSON(),
		ø.Send(item),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&recv),
		ƒ.Value(&recv).Is(&item),
	)
}

//
//
func TestNotFound() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/todo/%s", host, "unknown"),
		ƒ.Code(404),
	)
}

//
// Elementary Chain - ...
//

//
//
func TestLifeCycle() gurl.Arrow {
	item := TODO{ID: "4", Title: "have fun!"}

	return gurl.Join(
		append(item),
		lookup(item),
		contain(item),
		remove(item),
	)
}

//
func contain(item TODO) gurl.Arrow {
	var seq TODOs

	return gurl.HTTP(
		ø.GET("https://%s/api/todo", host),
		ƒ.Code(200),
		ƒ.Recv(&seq),
		ƒ.Seq(&seq).Has(item.ID, item),
	)
}

//
func append(item TODO) gurl.Arrow {
	return gurl.HTTP(
		ø.POST("https://%s/api/todo", host),
		ø.ContentJSON(),
		ø.Send(item),
		ƒ.Code(200),
	)
}

//
func lookup(expect TODO) gurl.Arrow {
	var item TODO

	return gurl.HTTP(
		ø.GET("https://%s/api/todo/%s", host, expect.ID),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
		ƒ.Value(&item).Is(expect),
	)
}

//
func remove(item TODO) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/todo/%s", host, item.ID),
		ƒ.Code(200),
	)
}

//
// Advance Chain - ...
//

func TestForEach() gurl.Arrow {
	var seq TODOs

	return gurl.Join(
		elements(&seq),
		foreach(&seq),
	)
}

func elements(seq *TODOs) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/todo", host),
		ƒ.Code(200),
		ƒ.Recv(seq),
	)
}

func foreach(seq *TODOs) gurl.Arrow {
	if len(*seq) == 0 {
		return nil
	}

	return gurl.Join(
		lookup((*seq)[0]),
		ƒ.FlatMap(func() gurl.Arrow {
			tail := (*seq)[1:]
			return foreach(&tail)
		}),
	)
}

//
func main() {}
