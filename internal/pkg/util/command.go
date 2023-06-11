package util

import (
    "fmt"
    "os"
)

func PrintErrorWithExit(err error) {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    os.Exit(1)
}
