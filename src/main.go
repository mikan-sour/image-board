package main

import (
	"github.com/jedzeins/image-board/src/cmd"
)

func main() {
	// lshc := healthcheck.NewHealthcheck("http://localstack:4566")
	// err := lshc.DoWhile()
	// if err != nil {
	// 	panic(err)
	// }

	cmd.Execute()
}
