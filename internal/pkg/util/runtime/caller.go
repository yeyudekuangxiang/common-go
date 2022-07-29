package runtime

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

func Callers(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])

	s := strings.Builder{}
	for _, pc := range pcs[0:n] {
		f := errors.Frame(pc)
		s.WriteString(fmt.Sprintf("\n%+v", f))
	}
	return s.String()
}
