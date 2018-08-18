// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import "github.com/gobuffalo/packr"

// static assets, embedded during compilation
var assets packr.Box

func init() {
	assets = packr.NewBox("assets")
}
