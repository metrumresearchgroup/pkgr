package main

import (
	"fmt"
	"reflect"

	"github.com/dpastoor/rpackagemanager/rcmd"
)

func main() {
	ia := rcmd.InstallArgs{Library: "some/path"}
	t := reflect.TypeOf(ia)

	// Get the type and kind of our user variable
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		tag := field.Tag.Get("rcmd")
		if field.Name == "Library" {
			tag = fmt.Sprintf(tag, ia.Library)
		}
		fmt.Printf("%d. %v (%v), tag: '--%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
	args := ia.CliArgs()
	fmt.Println(args)
}
