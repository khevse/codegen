package childpkg

import "fmt"

// Interface comment
type Interface interface {
	fmt.Stringer
	// OtherMethod comment
	OtherMethod() any
}
