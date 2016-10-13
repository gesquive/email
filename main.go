package main

import "fmt"

var version = "v0.5.1"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("email %s%s",
		version,
		dirty)
	Execute(displayVersion)
}
