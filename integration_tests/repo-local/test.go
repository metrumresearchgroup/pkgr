// +build linux darwin

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/metrumresearchgroup/command"
)

func main() {
	env := os.Environ()

	// this is where we manipulate the environment

	r, err := command.New(command.WithEnv(env)).Run(context.Background(), "pkgr", "plan")
	if err != nil {
		fmt.Printf("FAIL: %s", err)
	}

	fmt.Print(r.Output)
	fmt.Println("---") // this is just a temporary thing as I work through it; yaml doc separator.

	r, err = command.New(command.WithEnv(env)).Run(context.Background(), "pkgr", "install")
	if err != nil {
		fmt.Printf("FAIL: %s", err)
	}
	fmt.Print(r.Output)
}
