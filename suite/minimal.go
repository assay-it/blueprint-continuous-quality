/*

A minimal suite example.
See test/elementary.go for details and explanation of each statement
in this program.

*/

package suite

import (
	"github.com/assay-it/sdk-go/assay"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

// TestNews endpoint
func TestNews() assay.Arrow {
	return http.Join(
		ø.GET("%s/news", host),
		ƒ.Code(http.StatusCodeOK),
	)
}
