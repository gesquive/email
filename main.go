package main

import "os"
import "fmt"
import "path/filepath"

var version = "v0.5.0"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("%s %s%s",
		filepath.Base(os.Args[0]),
		version,
		dirty)
	Execute(displayVersion)
}
