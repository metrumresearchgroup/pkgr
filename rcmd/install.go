package rcmd

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	"github.com/fatih/structtag"
	"github.com/thoas/go-funk"
)

// type InstallArgs struct {
// 	Clean          bool `rcmd:"clean"`
// 	Preclean       bool `rcmd:"preclean"`
// 	Debug          bool `rcmd:"debug"`
// 	NoConfigure    bool `rcmd:"no-configure"`
// 	Example        bool `rcmd:"example"`
// 	Fake           bool `rcmd:"fake"`
// 	Build          bool `rcmd:"build"`
// 	InstallTests   bool `rcmd:"install-tests"`
// 	NoMultiarch    bool `rcmd:"no-multiarch"`
// 	WithKeepSource bool `rcmd:"with-keep.source"`
// 	ByteCompile    bool `rcmd:"byte-compile"`
// 	NoTestLoad     bool `rcmd:"no-test-load"`
// 	NoCleanOnError bool `rcmd:"no-clean-on-error"`
// 	//set
// 	Library string `rcmd:"library=%s,fmt"`
// }
func NewDefaultInstallArgs() InstallArgs {
	return InstallArgs{
		WithKeepSource: true,
		NoMultiarch:    true,
	}
}

func (i InstallArgs) CliArgs() []string {
	var args []string
	is := structs.New(i)
	nms := structs.Names(i)
	for _, n := range nms {
		fld := is.Field(n)
		if !fld.IsZero() {
			// ... and start using structtag by parsing the tag
			tag, _ := reflect.TypeOf(i).FieldByName(fld.Name())
			// ... and start using structtag by parsing the tag
			tags, err := structtag.Parse(string(tag.Tag))
			if err != nil {
				panic(err)
			}
			rcmd, err := tags.Get("rcmd")
			fmt.Println(rcmd.Name)
			fmt.Println(rcmd.Options)
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
