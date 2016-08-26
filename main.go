package main

import "os"
import "fmt"
import "path/filepath"
import "github.com/gesquive/email/cmd"

var version = "v0.1.0"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("%s %s%s",
		filepath.Base(os.Args[0]),
		version,
		dirty)
	cmd.Execute(displayVersion)
}
