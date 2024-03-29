// The nilness command applies the github.com/Andyfoo/golang/x/tools/go/analysis/passes/nilness
// analysis to the specified packages of Go source code.
package main

import (
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/nilness"
	"github.com/Andyfoo/golang/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nilness.Analyzer) }
