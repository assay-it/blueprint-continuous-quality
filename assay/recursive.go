/*

This example suite demonstrates ability of recursive behavior scenario.
The suite implement a sequence of nested interactions with micorservice.

See test/elementary.go for basic explanation about the suite structure.

*/
package assay

import (
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

/*

TestForEach proofs correctness of service behavior - the case involes multiple
recursive network operations. Here the case also show an ability to design behavior
without a dedicated scenario context type. The context is passed via variables by
references.
*/
func TestForEach() gurl.Arrow {
	// seq variable holds the context of recursive scenario
	var seq List

	// Join composes a recursive I/O. Suite passes a reference to sequence of news
	// to each function as the context of execution.
	return gurl.Join(
		// news checks quality of /api/news endpoint and fetches sequence of articles.
		news(&seq),
		// So far, the case has successfully received the sequence of news into seq variable.
		// Next, for each element of sequence it fetches the news article. Due to dynamic
		// nature of the sequence, we need to build a loop
		//
		//   for _, id := range seq { ... }
		//
		// However, the sequence is materialized only when previous operation succeeded.
		// Therefore, the case builds a recursion with foreach function and lifts it into
		// category of HTTP I/O.
		ƒ.FlatMap(func() gurl.Arrow {
			return foreach(&seq)
		}),
		// assign the origin sequence into result of the case
		func(x *gurl.IOCat) *gurl.IOCat {
			x.Body = seq
			return x
		},
	)
}

// News checks quality of /api/news endpoint and fetches sequence of news articles
// into seq variable
func news(seq *List) gurl.Arrow {
	return gurl.HTTP(
		ø.GET("%s/news", host),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(seq),
	)
}

// foreach recursively iterate the list, for each step it returns HTTP I/O promise.
func foreach(seq *List) gurl.Arrow {
	// Empty sequence causes termination of recursion.
	if len(*seq) == 0 {
		return nil
	}

	// each step of iterator fetches the head element from sequence and then
	// continue same function for the tail.
	return gurl.Join(
		item((*seq)[0]),
		ƒ.FlatMap(func() gurl.Arrow {
			tl := (*seq)[1:]
			return foreach(&tl)
		}),
	)
}

// Item checks quality of /api/news/:id endpoint and proofs its correctens.
func item(expect News) gurl.Arrow {
	var item News

	return gurl.HTTP(
		ø.GET("%s/news/%s", host, expect.ID),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		ƒ.Recv(&item),
		ƒ.Value(&item).Is(&expect),
	)
}
