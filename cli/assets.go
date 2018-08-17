package main

// TODO: move step to makefile
//go:generate bash -c "mkdir -p assets && tar cvf assets/context.tar -C ../ dockerfile makerelease.sh"

import "github.com/gobuffalo/packr"

// static assets, embedded during compilation
var assets packr.Box

func init() {
	assets = packr.NewBox("assets")
}
