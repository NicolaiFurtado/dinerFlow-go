package utils

import (
	"fmt"
	"os"
)

func Dd(v any) {
	fmt.Printf("%#v\n", v)
	os.Exit(1)
}
