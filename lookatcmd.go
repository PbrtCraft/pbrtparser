package pbrtparser

import (
	"strconv"
	"strings"
)

// LookAtCmd stores the parameters of lookat command:
//   LookAt 3 4 1.5 .5 .5 0 0 0 1
type LookAtCmd struct {
	Cmd
	Vals []float64 `json:"vals"`
}

func parseLookAtCmd(rawCommand string) (interface{}, error) {
	parts := strings.Fields(rawCommand)
	if len(parts) != 10 {
		return nil, ErrClassCommandForm
	}

	cmd := LookAtCmd{}
	cmd.CmdType = "LookAt"
	for i := 1; i <= 9; i++ {
		val, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return nil, ErrClassCommandForm
		}
		cmd.Vals = append(cmd.Vals, val)
	}
	return cmd, nil
}
