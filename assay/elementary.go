/*

Microservices have become a design style to evolve systems architecture in parallel,
implement stable and consistent interfaces. This architecture style brings additional
complexity and new problems. One of them is the assessment of system behavior while its
components communicate over the network - like integration testing but for distributed
environment. We need an ability to quantitatively evaluate and trade-off architecture
to ensure quality of the solutions.

https://assay.it is designed to perform a formal (objective) proofs of the quality
using Behavior as a Code contracts. It connects cause-and-effect (Given/When/Then) to
the networking concepts (Input/Process/Output). The expected behavior of each network
component is declared using simple Golang program.

Here is an example suite that illustrates an ability to apply contract testing of
quality assessment strategy. The contract implements a test cases as function of the form:

  func TestAbc() gurl.Arrow

where Abc is a unique name of test case. Each case declares cause-and-effect:

↣ "Given" specifies the communication context and the known state of the expected behavior;

↣ "When" executes key actions about the interaction with remote component;

↣ "Then" observes output of remote component, validates its correctness and outputs results.

The service evaluates suites and its test cases sequentially one after another.

Let's look on the following example!

*/

// Package test is a standard Golang declaration. It groups set of logically related contracts.
package assay

/*

Standard Golang import declaration.

However, assay.it restricts usage of some package.
Please check https://assay.it/doc for details of allowed packages.
We are constantly looking for your feedback, please open an issue to us.
*/
import (

	//
	// the toolkit for test suites development that provides various helper api.

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
		higher-order-functions and its compositions.
		See https://assay.it/doc/core for details about api
	*/

	// gurl.HTTP builds higher-order HTTP closure, so called gurl.Arrow, from
	// primitive elements.
	return gurl.HTTP(
		// module ø (gurl/http/send) defines function to declare HTTP request.
		// See the https://assay.it/doc/core for details about module ø

		// declares HTTP method and destination URL
		ø.GET("https://%s/news", host),

		/*
			"Then" observes output of remote component, validates its correctness and outputs results.
		*/

		// module ƒ (gurl/http/recv) defines function to validate correctness of HTTP response.
		// Each ƒ constrain might terminate execution of consequent ƒ's if it expectation fails.
		// See the https://assay.it/doc/core for details about module ƒ

		// requires HTTP Status Code to be 200 OK
		ƒ.Code(gurl.StatusCodeOK),
		// requites HTTP Header to be Content-Type: application/json
		ƒ.ServedJSON(),
		// requires a remote peer responds with List data type.
		// ƒ.Recv unmarshal JSON to the variable seq
		ƒ.Recv(&seq),
		// requires that expected element is present in the sequence.
		// Note: all received values are always passed by reference.
		ƒ.Seq(&seq).Has(expect.ID, expect),
	)
}

/*

TestNewsHTML checks quality of /api/news endpoint. The test case ensures
that api returns a sequence of news as HTML document
*/
func TestNewsHTML() gurl.Arrow {
	// Here the test case do not expect any particular value.
	// It just declares desired HTTP input and output.
	// Thus, we have omitted declaration of variables.
	return gurl.HTTP(
		ø.GET("https://%s/news", host),
		// output HTTP Header Accept: text/html
		ø.Accept().Is("text/html"),

		// requires HTTP Status Code to be 200 OK
		ƒ.Code(gurl.StatusCodeOK),
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
		ø.GET("https://%s/news/%s", host, "2"),

		ƒ.Code(gurl.StatusCodeOK),
		ƒ.ServedJSON(),
		// requires a remote peer responds with News data type.
		// ƒ.Recv unmarshal JSON to the variable news.
		ƒ.Recv(&news),
		// ƒ.Value is a helper function to assert received value(s)
		// Note: all received values are always passed by reference into assert functions
		ƒ.Value(&news.ID).String("2"),
		ƒ.Value(&news.Title).String("Sed luctus tortor sit amet eros eleifend cursus."),
	)
}

/*

TestItemHTML proofs correctens of example news article endpoint.
*/
func TestItemHTML() gurl.Arrow {
	// suite expects octet stream from api, here we declare a placeholder for the data.
	var news []byte
	// a constant value of expected response
	expect := []byte("<h1>2: Sed luctus tortor sit amet eros eleifend cursus.</h1>")

	return gurl.HTTP(
		ø.GET("https://%s/news/%s", host, "2"),
		ø.Accept().Is("text/html"),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Served().Is("text/html"),
		// ƒ.Bytes consumes response of remote peer into byte buffer.
		ƒ.Bytes(&news),
		// ƒ.Value assert received value.
		ƒ.Value(&news).Bytes(expect),
	)
}

/*

TestItemNotFound proofs correctens of example news article endpoint.
*/
func TestItemNotFound() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/news/%s", host, "9"),
		ƒ.Code(gurl.StatusCodeNotFound),
	)
}
