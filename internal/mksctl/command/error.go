package command

import (
	"fmt"
	"io"
)

// PrintError prints error properly to the writter
func PrintError(w io.Writer, err error) {
	fmt.Fprintf(w, "%s\n", err)
}
