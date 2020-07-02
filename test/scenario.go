/*

Here is an example suite that shows an ability to execute a behavioral scenario with
quality assessment. The suite implement a sequence of dependent interactions with
remote microservice. Each consequent interaction requires input from previous one.

See test/elementary.go for basic explanation about the suite structure.

*/
package main

import (
	"fmt"
	"sort"

	"github.com/assay-it/tk"
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// News a type used by the example application.
type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// List is a sequence of news, a core type of example application.
type List []News

// Value and other functions implements sort.Interface and gurl.Ord interfaces
func (seq List) Value(i int) interface{} { return seq[i] }
func (seq List) Len() int                { return len(seq) }
func (seq List) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq List) Less(i, j int) bool      { return seq[i].ID < seq[j].ID }
func (seq List) String(i int) string     { return seq[i].ID }

var host = fmt.Sprintf("v%s.%s", tk.Env("BUILD_ID", ""), tk.Env("HOST", ""))

/*

Scenario is a typed context to depict sequence of related interactions with
remote component(s). It is just a convenient container to group related
networking I/O and organize intermediate variables.
*/
type Scenario struct {
	Head string
	Last string
}

/*

News checks quality of /api/news endpoint, fetches sequence of News articles and
stores head and last elements of sequence to the memory (They are used in following
assessments).

Please note that Scenario is a pointer receiver type. It ensures that the case is able
to lift networking outputs with type-safe manner back to scenario context.
*/
func (s *Scenario) news() gurl.Arrow {
	var seq List

	return gurl.HTTP(
		ø.GET("https://%s/news", host),
		ƒ.Code(200),
		ƒ.Recv(&seq),
		// So far, the case has successfully received the sequence of news into seq variable.
		// We need to write a small function that extracts first and last elements from
		// sequence. A following traditional coding style would not work
		//
		//   var seq List
		//   gurl.HTTP(...)
		//   s.Head = seq[0].ID
		//
		// The sequence is not defined yet at the moment when gurl.HTTP(...) returns. It only
		// returns a "promise" of HTTP I/O which is materialized later. Therefore, any
		// computation have to be lifted-and-composed with this promise. ƒ.FMap does it.
		// ƒ.FMap takes a closure and applies it to the results of network communication.
		// In this example, the closure sort received sequence and fetches first and last
		// elements into Scenario context.
		ƒ.FMap(func() error {
			sort.Sort(seq)
			s.Head = seq[0].ID
			s.Last = seq[len(seq)-1].ID
			return nil
		}),
	)
}

/*

Item checks quality of /api/news/:id endpoint and proofs its correctens.
The function have to take pointers to context and an article id.
Pointers ensure a correct value is used during the evaluation of the case.
*/
func (s *Scenario) item(id *string) gurl.Arrow {
	var news News

	return gurl.HTTP(
		ø.GET("https://%s/news/%s", host, id),
		ƒ.Code(200),
		ƒ.ServedJSON(),
		ƒ.Recv(&news),
	)
}

/*

TestScenario proofs correctness of service behavior - the case involes multiple
consequent network operations. Similarly to other cases, this function is defined as

  func TestAbc() gurl.Arrow

Please note, this module implements a single test cases func TestScenario() gurl.Arrow
Other functions are used to compose this scenario.
*/
func TestScenario() gurl.Arrow {
	// Creates a scenario context, declare all given facts...
	s := &Scenario{}

	// gurl.Join composes elementary communication into high-order scenario.
	// The formal definition of Join is (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
	// It returns HTTP I/O HoC as the result
	return gurl.Join(
		s.news(),
		// Here, the suite parametrises item quality check with different values
		// from the context of the scenario. The values are passed by reference.
		s.item(&s.Head),
		s.item(&s.Last),
	)
}

func main() {}
