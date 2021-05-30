package configlib

import (
	"reflect"
	"testing"
)

func TestGetPackageCustomizationByName(t *testing.T) {
	// just enough config to confirm its not a default struct
	dplyrPkg := PkgConfig{
		Suggests: true,
	}
	type args struct {
		nm string
		c  Customizations
	}
	tests := []struct {
		name  string
		args  args
		want  PkgConfig
		want1 bool
	}{
		{
			name: "no customization",
			args: args{
				nm: "dplyr",
				c:  Customizations{},
			},
			want:  PkgConfig{},
			want1: false,
		},
		{
			name: "has customization",
			args: args{
				nm: "dplyr",
				c: Customizations{
					Packages: []map[string]PkgConfig{
						{
							"dplyr": dplyrPkg,
						},
					},
				},
			},
			want:  dplyrPkg,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetPackageCustomizationByName(tt.args.nm, tt.args.c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPackageCustomizationByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPackageCustomizationByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetRepoCustomizationByName(t *testing.T) {
	mpnRepo := RepoConfig{
		RepoSuffix: "SomeValue",
	}
	type args struct {
		nm string
		c  Customizations
	}
	tests := []struct {
		name  string
		args  args
		want  RepoConfig
		want1 bool
	}{
		{
			name: "no customization",
			args: args{
				nm: "MPN",
				c:  Customizations{},
			},
			want:  RepoConfig{},
			want1: false,
		},
		{
			name: "with customization",
			args: args{
				nm: "MPN",
				c: Customizations{
					Repos: []map[string]RepoConfig{
						{
							"MPN": mpnRepo,
						},
					},
				},
			},
			want:  mpnRepo,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetRepoCustomizationByName(tt.args.nm, tt.args.c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepoCustomizationByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetRepoCustomizationByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
