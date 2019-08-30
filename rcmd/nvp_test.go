package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name",
			"value",
			"append",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
		assert.Equal(t, tt.name, nvpList.Pairs[0].Name, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, tt.value, nvpList.Pairs[0].Value, fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestAppendNvp(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name",
			"value",
			"appendNvp",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.AppendNvp(tt.name + "=" + tt.value)
		assert.Equal(t, tt.name, nvpList.Pairs[0].Name, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, tt.value, nvpList.Pairs[0].Value, fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestGet(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name",
			"value",
			"get",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
		value, exists := nvpList.Get(tt.name)
		assert.Equal(t, tt.value, value, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, true, exists, fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestGetString(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name",
			"value",
			"getString",
		},
	}
	for _, tt := range tests {
		nvp := Nvp{
			Name:  tt.name,
			Value: tt.value,
		}
		assert.Equal(t, tt.name+"="+tt.value, nvp.GetString(tt.name), fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestGetNvp(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name",
			"value",
			"getNvp",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
		nvp, exists := nvpList.GetNvp(tt.name)
		assert.Equal(t, tt.value, nvp.Value, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, tt.name, nvp.Name, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, true, exists, fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestUpdate(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name1",
			"value1",
			"update",
		},
		{
			"name2",
			"value2",
			"update",
		},
		{
			"name3",
			"value3",
			"update",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
	}

	for _, tt := range tests {
		nvpList.Update(tt.name, tt.value+"_"+tt.context)
	}

	for _, tt := range tests {
		value, exists := nvpList.Get(tt.name)
		assert.Equal(t, true, exists, fmt.Sprintf("Fail: %s", tt.context))
		assert.Equal(t, tt.value+"_"+tt.context, value, fmt.Sprintf("Fail: %s", tt.context))
	}
}

func TestRemove(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name1",
			"value1",
			"remove",
		},
		{
			"name2",
			"value2",
			"remove",
		},
		{
			"name3",
			"value3",
			"remove",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
	}
	nvpList.Remove("name2")
	assert.Equal(t, 2, len(nvpList.Pairs), fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name1", nvpList.Pairs[0].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value1", nvpList.Pairs[0].Value, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name3", nvpList.Pairs[1].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value3", nvpList.Pairs[1].Value, fmt.Sprintf("Fail: remove"))
}

func TestRemove_First(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name1",
			"value1",
			"remove",
		},
		{
			"name2",
			"value2",
			"remove",
		},
		{
			"name3",
			"value3",
			"remove",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
	}
	nvpList.Remove("name1")
	assert.Equal(t, 2, len(nvpList.Pairs), fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name2", nvpList.Pairs[0].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value2", nvpList.Pairs[0].Value, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name3", nvpList.Pairs[1].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value3", nvpList.Pairs[1].Value, fmt.Sprintf("Fail: remove"))
}

func TestRemove_Last(t *testing.T) {
	var tests = []struct {
		name    string
		value   string
		context string
	}{
		{
			"name1",
			"value1",
			"remove",
		},
		{
			"name2",
			"value2",
			"remove",
		},
		{
			"name3",
			"value3",
			"remove",
		},
	}
	var nvpList NvpList
	for _, tt := range tests {
		nvpList.Append(tt.name, tt.value)
	}
	nvpList.Remove("name3")
	assert.Equal(t, 2, len(nvpList.Pairs), fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name1", nvpList.Pairs[0].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value1", nvpList.Pairs[0].Value, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "name2", nvpList.Pairs[1].Name, fmt.Sprintf("Fail: remove"))
	assert.Equal(t, "value2", nvpList.Pairs[1].Value, fmt.Sprintf("Fail: remove"))
}
