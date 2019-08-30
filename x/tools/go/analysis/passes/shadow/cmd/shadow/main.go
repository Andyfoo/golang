// The shadow command runs the shadow analyzer.
package main

import (
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/shadow"
	"github.com/Andyfoo/golang/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(shadow.Analyzer) }
