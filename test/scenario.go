package main

import (
	"sort"

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
//

type Scenario struct {
	Head string
	Tail string
}

func (s *Scenario) News() gurl.Arrow {
	var seq List

	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		ƒ.Code(200),
		ƒ.Recv(&seq),
		ƒ.FlatMap(func() gurl.Arrow {
			sort.Sort(seq)
			s.Head = seq[0].ID
			s.Tail = seq[len(seq)-1].ID
			return nil
		}),
	)
}

func (s *Scenario) Item(id *string) gurl.Arrow {
	var news News

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, id),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&news),
	)
}

func TestScenario() gurl.Arrow {
	s := &Scenario{}
	return gurl.Join(
		s.News(),
		s.Item(&s.Head),
		s.Item(&s.Tail),
	)
}

func main() {}
