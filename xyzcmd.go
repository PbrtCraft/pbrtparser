package pbrtparser

import (
	"fmt"
	"strconv"
	"strings"
)

type xyzCmd struct {
	Cmd
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// TranslateCmd stores the parameters of lookat command:
//   Translate x y z
type TranslateCmd xyzCmd

// ScaleCmd stores the parameters of lookat command:
//   Scale x y z
type ScaleCmd xyzCmd

func isXYZCmd(rawCommand string) bool {
	for _, prefix := range []string{
		"Translate",
		"Scale",
	} {
		if strings.HasPrefix(rawCommand, prefix) {
			return true
		}
	}
	return false
}

func parseXYZCmd(rawCommand string) (interface{}, error) {
	parts := strings.Fields(rawCommand)
	if len(parts) != 4 {
		return nil, fmt.Errorf("pbrtparser.parseXYZCmd: Error form")
	}

	var err error
	cmd := xyzCmd{}
	cmd.X, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseXYZCmd: %s", err)
	}
	cmd.Y, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseXYZCmd: %s", err)
	}
	cmd.Z, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseXYZCmd: %s", err)
	}

	cmd.CmdType = parts[0]
	if parts[0] == "Translate" {
		return TranslateCmd(cmd), nil
	} else if parts[0] == "Scale" {
		return ScaleCmd(cmd), nil
	}
	return nil, fmt.Errorf("pbrtparser.parseXYZCmd: No xyzCmd type match")
}

// RotateCmd stores the parameters of lookat command:
//   Rotate angle x y z
type RotateCmd struct {
	Cmd
	Angle float64 `json:"angle"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
}

func parseRotateCmd(rawCommand string) (interface{}, error) {
	parts := strings.Fields(rawCommand)
	if len(parts) != 5 {
		return nil, fmt.Errorf("pbrtparser.parseRotateCmd: Error form")
	}

	var err error
	cmd := RotateCmd{}
	cmd.Angle, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseRotateCmd: %s", err)
	}
	cmd.X, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseRotateCmd: %s", err)
	}
	cmd.Y, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseRotateCmd: %s", err)
	}
	cmd.Z, err = strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.parseRotateCmd: %s", err)
	}
	cmd.CmdType = "Rotate"
	return cmd, nil
}
