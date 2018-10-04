package rcmd

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	"github.com/fatih/structtag"
	"github.com/thoas/go-funk"
)

// NewDefaultInstallArgs provides a set of sane default installation args
func NewDefaultInstallArgs() *InstallArgs {
	return &InstallArgs{
		WithKeepSource: true,
		NoMultiarch:    true,
	}
}

// CliArgs converts the InstallArgs struct to the proper cli args
// including only returning the relevant args
func (i *InstallArgs) CliArgs() []string {
	var args []string
	is := structs.New(i)
	nms := structs.Names(i)
	for _, n := range nms {
		fld := is.Field(n)
		if !fld.IsZero() {
			// ... and start using structtag by parsing the tag
			tag, _ := reflect.TypeOf(i).Elem().FieldByName(fld.Name())
			// ... and start using structtag by parsing the tag
			tags, err := structtag.Parse(string(tag.Tag))
			if err != nil {
				panic(err)
			}
			rcmd, err := tags.Get("rcmd")
			if fld.Kind() == reflect.String && funk.Contains(rcmd.Options, "fmt") {
				// format the tag name by injecting any value into the tag name
				// for example lib=%s and struct value is some/path -> lib=some/path
				rcmd.Name = fmt.Sprintf(rcmd.Name, fld.Value())
			}
			args = append(args, fmt.Sprintf("--%s", rcmd.Name))
		}
	}
	return args
}
