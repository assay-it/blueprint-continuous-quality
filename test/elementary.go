/*

Microservices have become a design style to evolve systems architecture in parallel,
implement stable and consistent interfaces. This architecture style brings additional
complexity and new problems. Once of them is validation of system behavior while its
components communicate over the network. We need an ability to quantitatively
evaluate and trade-off the architecture to ensure quality of the software solutions.

https://assay.it is a service that automatically performs a formal (objective)
proofs of the quality using Behavior as a Code paradigm. It connects cause-and-effect
(Given/When/Then) to the networking concepts (Input/Process/Output). The expected
behavior of each network component is declared using simple Golang program (suite).

Here is an example suite that illustrates an ability to apply unit-test like strategy
on quality assessment. The suite implements a test cases as function of the form:

   func TestXxx() gurl.Arrow

where Xxx is a unique name of test case. The case declares cause-and-effect:

↣ "Given" specifies the communication context and the known state of the expected behavior;

↣ "When" executes key actions about the interaction with remote component;

↣ "Then" observes output of remote component, validates its correctness and outputs results.

The service evaluates suites and its test cases sequentially one after another.

Let's look on the following example!

*/
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
//
func TestItemJSON() gurl.Arrow {
	var news News

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "2"),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&news),
		ƒ.Value(&news.ID).String("2"),
		ƒ.Value(&news.Title).String("Sed luctus tortor sit amet eros eleifend cursus."),
	)
}

//
//
func TestItemHTML() gurl.Arrow {
	var news []byte

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "2"),
		ø.Accept().Is("text/html"),
		ƒ.Code(200),
		ƒ.Served().Is("text/html"),
		ƒ.Bytes(&news),
		ƒ.Value(&news).Bytes([]byte("<h1>2: Sed luctus tortor sit amet eros eleifend cursus.</h1>")),
	)
}

//
//
func TestItemNotFound() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "9"),
		ƒ.Code(404),
	)
}

//
func main() {}
