// The findcall command runs the findcall analyzer.
package main

import (
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/findcall"
	"github.com/Andyfoo/golang/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findcall.Analyzer) }
