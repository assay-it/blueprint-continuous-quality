/*

This example suite demonstrates ability of recursive behavior scenario.
The suite implement a sequence of nested interactions with endpoints.

See test/elementary.go for basic explanation about the suite structure.

*/

package suite

import (
	"github.com/assay-it/sdk-go/assay"
	c "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

/*

TestForEach proofs correctness of service behavior - the case involes multiple
recursive network operations. Here the case also show an ability to design
behavior without a dedicated scenario context type. The context is passed via
variables by references.
*/
func TestForEach() assay.Arrow {
	// seq variable holds the context of recursive scenario
	var seq List

	// Join composes a recursive I/O. Suite passes a reference to sequence of news
	// to each function as the context of execution.
	return assay.Join(
		// news checks quality of /api/news endpoint and fetches sequence of
		// articles.
		news(&seq),
		// So far, the case has successfully received the sequence of news into seq
		// variable. Next, for each element of sequence it fetches the news article.
		// Due to dynamic nature of the sequence, we need to build a loop
		//
		//   for _, id := range seq { ... }
		//
		// However, the sequence is materialized only when previous operation
		// succeeded. Therefore, the case builds a recursion with foreach function
		// and lifts it into category of HTTP I/O.
		c.FlatMap(func() assay.Arrow { return foreach(&seq) }),

		// assign the origin sequence into result of the case
		func(cat *assay.IOCat) *assay.IOCat {
			cat.HTTP.Recv.Payload = seq
			return cat
		},
	)
}

// News checks quality of /api/news endpoint and fetches sequence of news
// articles into seq variable
func news(seq *List) assay.Arrow {
	return http.Join(
		ø.GET("%s/news", host),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(seq),
	)
}

// foreach recursively iterate the list, for each step it returns HTTP I/O
// promise.
func foreach(seq *List) assay.Arrow {
	// Empty sequence causes termination of recursion.
	if len(*seq) == 0 {
		return nil
	}

	// each step of iterator fetches the head element from sequence and then
	// continue same function for the tail.
	return assay.Join(
		item((*seq)[0]),
		c.FlatMap(func() assay.Arrow {
			tl := (*seq)[1:]
			return foreach(&tl)
		}),
	)
}

// Item checks quality of /api/news/:id endpoint and proofs its correctens.
func item(expect News) assay.Arrow {
	var item News

	return http.Join(
		ø.GET("%s/news/%s", host, expect.ID),
		ƒ.Code(http.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
	).Then(
		c.Value(&item).Is(&expect),
	)
}
