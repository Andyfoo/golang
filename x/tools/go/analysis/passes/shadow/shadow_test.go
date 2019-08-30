package shadow_test

import (
	"testing"

	"github.com/Andyfoo/golang/x/tools/go/analysis/analysistest"
	"github.com/Andyfoo/golang/x/tools/go/analysis/passes/shadow"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, shadow.Analyzer, "a")
}
