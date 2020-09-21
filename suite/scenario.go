/*

Here is an example suite that shows an ability to execute a behavioral scenario.
The suite implement a sequence of dependent interactions with microservice.
Each consequent interaction requires input from previous one.

See test/elementary.go for basic explanation about the suite structure.

*/

package suite

import (
	"sort"

	"github.com/assay-it/sdk-go/assay"
	c "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

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

News checks quality of /api/news endpoint, fetches sequence of News articles
and stores head and last elements of sequence to the memory (They are used in
following assessments).

Please note that Scenario is a pointer receiver type. It ensures that the case
is able to lift networking outputs with type-safe manner back to scenario
context.
*/
func (s *Scenario) news() assay.Arrow {
	var seq List

	return http.Join(
		ø.GET("%s/news", host),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&seq),
		// So far, the case has successfully received the sequence of news into seq
		// variable. We need to write a small function that extracts first and last
		// elements from sequence. A following traditional coding style would not
		// work
		//
		//   var seq List
		//   http.Join(...)
		//   s.Head = seq[0].ID
		//
		// The sequence is not "materialized" yet at the moment when http.Join(...)
		// returns. It only returns a "promise" of HTTP I/O which is materialized
		// later. Therefore, any computation have to be lifted-and-composed with
		// this promise. c.FMap does it. c.FMap takes a closure and applies it to
		// the results of network communication. In this example, the closure sort
		// received sequence and fetches first and last elements into Scenario
		// context.
	).Then(
		c.FMap(func() error {
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
func (s *Scenario) item(id *string) assay.Arrow {
	var news News

	return http.Join(
		ø.GET("%s/news/%s", host, id),
		ƒ.Code(http.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&news),
	)
}

/*

TestScenario proofs correctness of service behavior - the case involes multiple
consequent network operations. Similarly to other cases, this function is
defined as

  func TestAbc() assay.Arrow

Please note, this module implements a single test cases TestScenario,
Other functions are used to compose this scenario.
*/
func TestScenario() assay.Arrow {
	// Creates a scenario context, declare all given facts...
	s := &Scenario{}

	// assay.Join composes elementary communication into high-order scenario.
	// The formal definition of Join is (a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
	// It returns HTTP I/O HoC as the result
	return assay.Join(
		s.news(),
		// Here, the suite parametrises item quality check with different values
		// from the context of the scenario. The values are passed by reference.
		s.item(&s.Head),
		s.item(&s.Last),
	)
}
