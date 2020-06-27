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

//
//
//

var host = tk.Env("HOST", "")

func TestForEach() gurl.Arrow {
	var seq List

	return gurl.Join(
		elements(&seq),
		ƒ.FlatMap(func() gurl.Arrow {
			return foreach(&seq)
		}),
		func(x *gurl.IOCat) *gurl.IOCat {
			x.Body = seq
			return x
		},
	)
}

func elements(seq *List) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		ƒ.Code(200),
		ƒ.Recv(seq),
	)
}

func foreach(seq *List) gurl.Arrow {
	if len(*seq) == 0 {
		return nil
	}

	hd := (*seq)[0]

	return gurl.Join(
		lookup(hd),
		ƒ.FlatMap(func() gurl.Arrow {
			tl := (*seq)[1:]
			return foreach(&tl)
		}),
	)
}

func lookup(expect News) gurl.Arrow {
	var item News

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, expect.ID),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
		ƒ.Value(&item).Is(&expect),
	)
}

//
func main() {}
