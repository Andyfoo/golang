package nodisk

import (
	"github.com/Andyfoo/golang/x/tools/internal/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
