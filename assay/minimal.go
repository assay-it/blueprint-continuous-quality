/*

A minimal suite example.
See test/elementary.go for details and explanation of each statement in this program.

*/
package assay

import (
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

// TestNews endpoint
func TestNews() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/news", host),
		ƒ.Code(gurl.StatusCodeOK),
	)
}
