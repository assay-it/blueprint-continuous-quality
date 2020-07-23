/*

A minimal suite example.
See test/elementary.go for details and explanation of each statement in this program.

*/
package main

import (
	"fmt"

	"github.com/assay-it/tk"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

var host = fmt.Sprintf("v%s.%s", tk.Env("BUILD_ID", ""), tk.Env("CONFIG_DOMAIN", ""))

// TestNews endpoint
func TestNews() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/news", host),
		ƒ.Code(gurl.StatusCodeOK),
	)
}

func main() {}
