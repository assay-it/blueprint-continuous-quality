/*

A minimal example of valid test suite.
See test/elementary.go for details and explanation of each constrain.

*/
package main

import (
	"github.com/assay-it/tk"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

var host = tk.Env("HOST", "")

// TestNews endpoint
func TestNews() gurl.Arrow {
	return gurl.HTTP(
		ø.GET("https://%s/api/news", host),
		ƒ.Code(200),
	)
}

func main() {}
