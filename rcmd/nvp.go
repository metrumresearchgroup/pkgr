package rcmd

import (
	"fmt"
	"strings"
)

// TODO: generally, about this type, a lot of the value comes from the format
//  of a list of things. Honestly, the environment functions in the os package
//  are a bit hard to work with, but most of what is being described in this file
//  could be expressed as a []string as it is in the os package.
// TODO: what we're really proposing is a way to edit the environment for later
//  insertion into a sub-shell. That manipulation can be done much easier/cheaper
//  with the natural format found in os.



// Append a name and value pair to the list as an Nvp object
// TODO: Append is called by AppendNvp.
// TODO: Append is used 4x in configure.go
func (list *NvpList) Append(name, value string) {
	list.Pairs = append(list.Pairs, Nvp{Name: strings.Trim(name, " "), Value: strings.Trim(value, " ")})
}

// AppendNvp append a string of name=value pair to the list as an Nvp object
// TODO: AppendNvp is used 1x in configure.go
func (list *NvpList) AppendNvp(nvp string) {
	b := strings.Split(nvp, "=")
	if len(b) == 2 {
		list.Append(b[0], b[1])
	}
}

// Get a value by name
// TODO: Get is used 2x in configure.go
func (list *NvpList) Get(name string) (value string, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair.Value, true
		}
	}
	return "", false
}

// GetNvp an nvp by name
// TODO: GetNvp is only used in tests
func (list *NvpList) GetNvp(name string) (nvp Nvp, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair, true
		}
	}
	return Nvp{}, false
}

// Remove by name
// TODO: Remove is only used in tests
func (list *NvpList) Remove(name string) {
	n := -1
	for i, pair := range list.Pairs {
		if name == pair.Name {
			n = i
			break
		}
	}
	if n >= 0 {
		list.Pairs = append(list.Pairs[:n], list.Pairs[n+1:]...)
	}
}

// Update a value by name
// TODO: Update is only used in tests
func (list *NvpList) Update(name string, value string) (nvp Nvp, exists bool) {
	n := -1
	for i, pair := range list.Pairs {
		if name == pair.Name {
			n = i
		}
	}

	if n >= 0 {
		list.Pairs[n].Value = value
		return list.Pairs[n], true
	}
	return Nvp{}, false
}

// GetString returns a string as name=value
// TODO: GetString is used once in configure.go
func (nvp *Nvp) GetString(name string) (value string) {
	return fmt.Sprintf("%s=%s", nvp.Name, nvp.Value)
}
