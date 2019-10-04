package pbrtparser

import (
	"fmt"
	"strings"
)

type useCmd struct {
	Cmd
	What string `json:"what"` // use what?
}

// ObjectInstanceCmd stores parameters of object using command
//   ObjectInstance "foo"
type ObjectInstanceCmd useCmd

// NamedMaterialCmd stores parameters of material using command
//   NamedMaterial "myplastic"
type NamedMaterialCmd useCmd

func isUseCmd(rawCommand string) bool {
	return strings.HasPrefix(rawCommand, "NamedMaterial") ||
		strings.HasPrefix(rawCommand, "ObjectInstance")
}

func parseUseCmd(rawCommand string) (interface{}, error) {
	tokens := toTokens(rawCommand)
	if len(tokens) != 4 {
		return nil, fmt.Errorf("pbrtparser.usecmd.parseUseCmd: Error format")
	}

	class := tokens[0]
	err := assureForm(tokens[1:4], `"`, `"`)
	if err != nil {
		return nil, err
	}

	cmd := useCmd{What: tokens[2]}
	cmd.CmdType = class
	switch class {
	case "NamedMaterial":
		return NamedMaterialCmd(cmd), nil
	case "ObjectInstanceCmd":
		return ObjectInstanceCmd(cmd), nil
	default:
		panic("Error Use Cmd: " + class)
	}
	return nil, nil
}
