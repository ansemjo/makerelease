package main

import "github.com/gobuffalo/packr"

// static assets, embedded during compilation
var assets packr.Box

func init() {
	assets = packr.NewBox("assets")
}
