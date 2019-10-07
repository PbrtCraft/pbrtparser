package pbrtparser

import (
	"testing"
)

func TestParser(t *testing.T) {
	sp, err := NewCmdsParser("test/xyz.pbrt")

	if err != nil {
		t.Error(err)
		return
	}

	cmds, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}

	if len(cmds) != 1 {
		t.Error("Command parse error")
		return
	}

	cmd, ok := cmds[0].(TranslateCmd)
	if !ok {
		t.Error("Command type error")
		return
	}

	if cmd.X != 150 || cmd.Y != 0 || cmd.Z != 20 {
		t.Error("Params parse error")
		return
	}
}
