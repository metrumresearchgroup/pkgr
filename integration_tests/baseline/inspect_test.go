package baseline

import (
	"context"
	"encoding/json"
	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"testing"
)

func TestInspectDeps(t *testing.T) {
	testCmd := command.New()
	ctx := context.TODO()
	res, err := testCmd.Run(ctx, "pkgr", "inspect", "--deps", "--json")

	if err != nil {
		panic(err)
	}
	t.Run("get dependencies as json", func(t *testing.T) {
		g := goldie.New(t)
		g.Assert(t, "inspect-deps", res.Output)
	})

	// the situation that can arise is if log messages slip in, so the output would be some logrus message + json
	t.Run("dependencies are valid json", func(t *testing.T) {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal(res.Output, &jsonMap)
		if err != nil {
			t.Errorf("could not unmarshal dependency json with error %s", err)
		}
	})
}

func TestInspectReverseDeps(t *testing.T) {
	testCmd := command.New()
	ctx := context.TODO()
	res, err := testCmd.Run(ctx, "pkgr", "inspect", "--deps", "--reverse", "--json")
	if err != nil {
		panic(err)
	}
	t.Run("reverse dependencies are valid json", func(t *testing.T) {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal(res.Output, &jsonMap)
		if err != nil {
			t.Errorf("could not unmarshal reverse dependency json with error %s", err)
		}
	})
	// the situation that can arise is if log messages slip in, so the output would be some logrus message + json
	t.Run("get reverse dependencies as json", func(t *testing.T) {
		g := goldie.New(t)
		g.Assert(t, "inspect-reverse-deps", res.Output)
	})
}

