// The lostcancel command applies the github.com/Andyfoo/golang/x/tools/go/analysis/passes/lostcancel
// analysis to the specified packages of Go source code.
package main

import (
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/lostcancel"
	"github.com/Andyfoo/golang/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(lostcancel.Analyzer) }
