package errors

import (
	"github.com/Andyfoo/golang/x/tools/internal/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
