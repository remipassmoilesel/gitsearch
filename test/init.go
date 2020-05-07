package test

import (
	"os"
	"path"
	"runtime"
)

// Here we ensure that tests will be run from root directory
// Use it like this:
//
//  import "gitlab.com/remipassmoilesel/gitsearch/test"
//  import _ "gitlab.com/remipassmoilesel/gitsearch/test"

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
