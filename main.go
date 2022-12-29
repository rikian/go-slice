package main

import (
	"go/slice/src"
)

func main() {
	src.NewListener().Run("0.0.0.0", "9091")
}
