package pbrtparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrClassCommandForm stands for format of command wrong
	ErrClassCommandForm = errors.New("Token form wrong")

	// ErrClassEmpty stands for that a name of class is empty
	ErrClassEmpty = errors.New("Empty class")
)

// Param stores param in class command
type Param struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Val  interface{} `json:"val"`
}

func newParam(name, typ, rawVal string) (*Param, error) {
	switch typ {
	case "string", "texture":
		return &Param{
			Name: name,
			Type: typ,
			Val:  rawVal,
		}, nil
	case "bool":
		var val bool
		if rawVal == "true" {
			val = true
		} else if rawVal == "false" {
			val = false
		} else {
			return nil, fmt.Errorf("pbrtparser.classcmd.newParam: bool value=%s", rawVal)
		}
		return &Param{
			Name: name,
			Type: typ,
			Val:  val,
		}, nil
	case "integer":
		vals := []int{}
		for _, s := range strings.Fields(rawVal) {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("pbrtparser.classmd.newParam: %s", err)
			}
			vals = append(vals, i)
		}
		return &Param{
			Name: name,
			Type: typ,
			Val:  vals,
		}, nil

	case "float", "color", "rgb", "xyz", "spectrum",
		"point", "point2", "point3",
		"vector", "vector2", "vector3",
		"blackbody", "normal":
		vals := []float64{}
		for _, s := range strings.Fields(rawVal) {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("pbrtparser.classmd: %s", err)
			}
			vals = append(vals, f)
		}
		return &Param{
			Name: name,
			Type: typ,
			Val:  vals,
		}, nil
	default:
		return nil, errors.New(typ + " param not match")
	}
}

func parseParamList(tokens []string) ([]*Param, error) {
	ll := len(tokens)
	var err error
	params := []*Param{}
	for lx := 0; lx < ll; lx += 6 {
		// ", ParamName Type, ", [, Raw Value, ]
		err = assureForm(tokens[lx:lx+3], `"`, `"`)
		if err != nil {
			return nil, err
		}
		def := strings.Fields(tokens[lx+1])
		if len(def) != 2 {
			return nil, fmt.Errorf("pbrtparser.parseParamList: Error def of param")
		}
		if def[0] == "string" || def[0] == "bool" || def[0] == "texture" {
			err = assureForm(tokens[lx+3:lx+6], `"`, `"`)
		} else {
			err = assureForm(tokens[lx+3:lx+6], `[`, `]`)
		}
		if err != nil {
			return nil, err
		}
		param, err := newParam(def[1], def[0], tokens[lx+4])
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

// ClassCmd stores parameters of class command
//   COMMANDTYPE "Name" "Par1.Type1 Par1.Name1" [Par1.Vals] ...
type ClassCmd struct {
	Cmd
	Name   string   `json:"name"`
	Params []*Param `json:"params"`
}

func isClassCmd(rawCommand string) bool {
	for _, prefix := range []string{
		"Camera",
		"Material",
		"Shape",
		"Film",
		"Sampler",
		"Integrator",
		"LightSource",
		"AreaLightSource",
		"MakeNamedMedium",
		"MakeNamedMaterial",
	} {
		if strings.HasPrefix(rawCommand, prefix) {
			return true
		}
	}
	return false
}

func parseClassCmd(rawCommand string) (*ClassCmd, error) {
	tokens := toTokens(rawCommand)
	cmd := ClassCmd{}

	if len(tokens) == 0 {
		return nil, ErrClassEmpty
	}

	class := tokens[0]
	tokens = tokens[1:]
	ll := len(tokens)
	if ll == 0 {
		return nil, fmt.Errorf("pbrtparse.parseClass: Empty Command")
	}

	// Schema: "NAME" "ParamName Type" [Values] ...
	// Len = 3 + 6x
	if (ll-3)%6 != 0 {
		return nil, ErrClassCommandForm
	}

	var err error
	err = assureForm(tokens[:3], `"`, `"`)
	if err != nil {
		return nil, err
	}
	cmd.Name = tokens[1]
	cmd.Params, err = parseParamList(tokens[3:])
	if err != nil {
		return nil, err
	}

	cmd.CmdType = class
	return &cmd, nil
}

// TextureCmd stores parameters of texture command:
type TextureCmd struct {
	Cmd
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Class  string   `json:"class"`
	Params []*Param `json:"params"`
}

func parseTextureCmd(rawCommand string) (interface{}, error) {
	tokens := toTokens(rawCommand)

	if len(tokens) == 0 {
		return nil, ErrClassEmpty
	}

	tokens = tokens[1:]
	ll := len(tokens)
	if ll == 0 {
		return nil, fmt.Errorf("pbrtparse.parseClass: Empty Command")
	}

	// Schema: "NAME" "Type" "Class" "ParamName Type" [Values] ...
	// Len = 9 + 6x
	if ll < 9 || (ll-9)%6 != 0 {
		return nil, ErrClassCommandForm
	}

	var err error
	for i := 0; i < 9; i += 3 {
		err = assureForm(tokens[i:i+3], `"`, `"`)
		if err != nil {
			return nil, err
		}
	}
	cmd := TextureCmd{
		Name:  tokens[1],
		Type:  tokens[4],
		Class: tokens[7],
	}
	cmd.CmdType = "Texture"
	cmd.Params, err = parseParamList(tokens[9:])
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}
