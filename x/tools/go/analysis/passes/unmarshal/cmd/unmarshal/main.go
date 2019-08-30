// The unmarshal command runs the unmarshal analyzer.
package main

import (
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/unmarshal"
	"github.com/Andyfoo/golang/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(unmarshal.Analyzer) }
