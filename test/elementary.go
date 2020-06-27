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

// each suite is always package main
package main

/*

Standard Golang import declaration.

However, assay.it restricts usage of package to the list of allowed one.
Please check doc.assay.it for details.
*/
import (
	//
	// the toolkit for test suites development that provides various helper api.
	"github.com/assay-it/tk"

	//
	// gurl library is a class of High Order Component which can do http requests
	// with few interesting property such as composition and laziness.
	// It implements a human-friendly syntax of HTTP communication and
	// Behavior as a Code paradigm. It connect cause-and-effect with the networking
	// primitives. Usage of gurl is a preferred approach for networking I/O.
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

/*

Golang type declaration.

It is possible to declare any types as part of the suite implementation.
*/

// News a type used by the example application. This type models a core data of
// the application and used by suites to validates correctness of outputs.
type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// List is a sequence of news, a core type of example application.
type List []News

// Value and other functions implements sort.Interface and gurl.Ord interfaces
// for the sequence. It allows asserting and validation of sequences
// (e.g. ƒ.Seq(&seq).Has(...) )
func (seq List) Value(i int) interface{} { return seq[i] }
func (seq List) Len() int                { return len(seq) }
func (seq List) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq List) Less(i, j int) bool      { return seq[i].ID < seq[j].ID }
func (seq List) String(i int) string     { return seq[i].ID }

/*

Suite constants and other global variables

*/

// Settings of assay.it allows developers to customize suite via environment
// variables. These variables are provided into your code during the assessment.
// This example application requires a HOST setting that declares a target SUT.
// Here the toolkit is used to read the value of environment variable.
var host = tk.Env("HOST", "")

/*

TestNewsJSON checks quality of /api/news endpoint. The test case ensures
that api returns a sequence of news and the sequence contains a mandatory
element.

Let's look in depth on the anatomy of test case
*/
func TestNewsJSON() gurl.Arrow {
	/*
		"Given" specifies the communication context and the known state of the expected behavior
	*/

	// seq variable holds a response of the endpoint under the test.
	var seq List
	// a constant value of expected response
	expect := News{ID: "2", Title: "Sed luctus tortor sit amet eros eleifend cursus."}

	/*
		"When" executes key actions about the interaction with remote component

		gurl library defines a rich techniques to hide the networking complexity using
		higher-order-functions and its compositions. See the doc.assay.it for details
		about api.
	*/

	// gurl.HTTP builds higher-order HTTP closure, so called gurl.Arrow, from
	// primitive elements.
	return gurl.HTTP(
		// module ø (gurl/http/send) defines function to specify request of HTTP protocol.
		// See the doc.assay.it for details about other function of module ø.

		// HTTP method and destination URL
		ø.GET("https://%s/api/news", host),

		/*
			"Then" observes output of remote component, validates its correctness and outputs results.
		*/

		// module ƒ (gurl/http/recv) defines function to validate correctness of HTTP protocol.
		// each ƒ constrain might terminate execution of consequent one if its constrain is not
		// valid. See the doc.assay.it for details about other function of module ƒ.

		// requires HTTP Status Code to be 200 OK
		ƒ.Code(200),
		// requites HTTP Header to be Content-Type: application/json
		ƒ.ServedJSON(),
		// requires a remote peer responds with List data type.
		// ƒ.Recv unmarshal JSON to the variable
		ƒ.Recv(&seq),
		// requires that expected element is present in the sequence.
		ƒ.Seq(&seq).Has(expect.ID, expect),
	)
}

/*

TestNewsHTML checks quality of /api/news endpoint. The test case ensures
that api returns a sequence of news as HTML document
*/
func TestNewsHTML() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		// output HTTP Header Accept: text/html
		ø.Accept().Is("text/html"),

		// requires HTTP Status Code to be 200 OK
		ƒ.Code(200),
		// requites HTTP Header to be Content-Type: text/html
		ƒ.Served().Is("text/html"),
	)
}

/*

TestItemJSON proofs correctens of example news article endpoint.
*/
func TestItemJSON() gurl.Arrow {
	// the response type MUST be News document
	var news News

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "2"),

		ƒ.Code(200),
		ƒ.ServedJSON(),
		// requires a remote peer responds with News data type.
		// ƒ.Recv unmarshal JSON to the variable
		ƒ.Recv(&news),
		// ƒ.Value is a helper function to assert received value(s)
		ƒ.Value(&news.ID).String("2"),
		ƒ.Value(&news.Title).String("Sed luctus tortor sit amet eros eleifend cursus."),
	)
}

/*

TestItemHTML proofs correctens of example news article endpoint.
*/
func TestItemHTML() gurl.Arrow {
	var news []byte

	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "2"),
		ø.Accept().Is("text/html"),
		ƒ.Code(200),
		ƒ.Served().Is("text/html"),
		// ƒ.Bytes consumes response of remote peer into bytes buffer
		ƒ.Bytes(&news),
		ƒ.Value(&news).Bytes([]byte("<h1>2: Sed luctus tortor sit amet eros eleifend cursus.</h1>")),
	)
}

/*

TestItemNotFound proofs correctens of example news article endpoint.
*/
func TestItemNotFound() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news/%s", host, "9"),
		ƒ.Code(404),
	)
}

// main function is a required element of main package.
// It have to be declared for each suite.
func main() {}
