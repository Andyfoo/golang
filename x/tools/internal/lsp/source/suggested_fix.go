package source

import (
	"go/token"

	"github.com/Andyfoo/golang/x/tools/go/analysis"
	"github.com/Andyfoo/golang/x/tools/internal/lsp/diff"
	"github.com/Andyfoo/golang/x/tools/internal/span"
)

func getCodeActions(fset *token.FileSet, diag analysis.Diagnostic) ([]SuggestedFixes, error) {
	var cas []SuggestedFixes
	for _, fix := range diag.SuggestedFixes {
		var ca SuggestedFixes
		ca.Title = fix.Message
		for _, te := range fix.TextEdits {
			span, err := span.NewRange(fset, te.Pos, te.End).Span()
			if err != nil {
				return nil, err
			}
			ca.Edits = append(ca.Edits, diff.TextEdit{Span: span, NewText: string(te.NewText)})
		}
		cas = append(cas, ca)
	}
	return cas, nil
}
