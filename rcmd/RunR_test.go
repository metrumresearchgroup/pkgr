package rcmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRunRBatch(t *testing.T) {
	type args struct {
		fs      afero.Fs
		rs      RSettings
		cmdArgs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test R Version",
			args: args{
				fs: afero.NewOsFs(),
				rs: RSettings{
					Rpath: "/usr/local/bin/R",
				},
				cmdArgs: []string{
					"--version",
				},
			},
			want: []byte(`R version 3.5.3 (2019-03-11) -- "Great Truth"
Copyright (C) 2019 The R Foundation for Statistical Computing
Platform: x86_64-apple-darwin15.6.0 (64-bit)

R is free software and comes with ABSOLUTELY NO WARRANTY.
You are welcome to redistribute it under the terms of the
GNU General Public License versions 2 or 3.
For more information about these matters see
http://www.gnu.org/licenses/.

`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := RunRBatch(tt.args.fs, tt.args.rs, tt.args.cmdArgs)
			assert.Equal(t, nil, err, "error")

			n := bytes.Compare(got, tt.want)
			msg := fmt.Sprintf("\ngot<\n%v\n>\nwant<\n%v\n>", string(got), string(tt.want))
			assert.Equal(t, 0, n, msg)

		})
	}
}

func BenchmarkRunR(b *testing.B) {
	rs := NewRSettings("/usr/local/bin/R")
	fs := afero.NewOsFs()
	for n := 0; n < b.N; n++ {
		RunR(fs, "", rs, "version", "")
	}
}

func BenchmarkRunRBatch(b *testing.B) {
	rs := NewRSettings("/usr/local/bin/R")
	fs := afero.NewOsFs()
	for n := 0; n < b.N; n++ {
		RunRBatch(fs, rs, []string{"--version"})
	}
}
