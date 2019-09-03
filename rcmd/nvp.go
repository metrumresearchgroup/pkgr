package rcmd

import (
	"fmt"
	"strings"
)

// Append a name and value pair to the list as an Nvp object
func (list *NvpList) Append(name, value string) {
	list.Pairs = append(list.Pairs, Nvp{Name: strings.Trim(name, " "), Value: strings.Trim(value, " ")})
}

// AppendNvp append a string of name=value pair to the list as an Nvp object
func (list *NvpList) AppendNvp(nvp string) {
	b := strings.Split(nvp, "=")
	if len(b) == 2 {
		list.Append(b[0], b[1])
	}
}

// Get a value by name
func (list *NvpList) Get(name string) (value string, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair.Value, true
		}
	}
	return "", false
}

// GetNvp an nvp by name
func (list *NvpList) GetNvp(name string) (nvp Nvp, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair, true
		}
	}
	return Nvp{}, false
}

// Remove by name
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
func (nvp *Nvp) GetString(name string) (value string) {
	return fmt.Sprintf("%s=%s", nvp.Name, nvp.Value)
}
