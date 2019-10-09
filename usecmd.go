package pbrtparser

import (
	"fmt"
	"strings"
)

// UseCmd stores parameters of use command
//   COMMANDTYPE "What"
type UseCmd struct {
	Cmd
	What string `json:"what"` // use what?
}

func isUseCmd(rawCommand string) bool {
	return strings.HasPrefix(rawCommand, "NamedMaterial") ||
		strings.HasPrefix(rawCommand, "ObjectInstance")
}

func parseUseCmd(rawCommand string) (*UseCmd, error) {
	tokens := toTokens(rawCommand)
	if len(tokens) != 4 {
		return nil, fmt.Errorf("pbrtparser.usecmd.parseUseCmd: Error format")
	}

	err := assureForm(tokens[1:4], `"`, `"`)
	if err != nil {
		return nil, err
	}

	cmd := UseCmd{
		Cmd:  Cmd{tokens[0]},
		What: tokens[2],
	}
	return &cmd, nil
}
