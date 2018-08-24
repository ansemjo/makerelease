// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package assets

import "github.com/gobuffalo/packr"

// Box includes static assets, embedded during compilation.
// E.g. the Docker build context: `context.tar`.
var Box packr.Box

func init() {
	Box = packr.NewBox(".")
}
