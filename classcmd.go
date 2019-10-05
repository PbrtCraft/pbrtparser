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
				fmt.Errorf("pbrtparser.classmd.newParam: %s", err)
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
				fmt.Errorf("pbrtparser.classmd: %s", err)
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

type classCmd struct {
	Cmd
	Name   string   `json:"name"`
	Params []*Param `json:"params"`
}

// CameraCmd stores parameters of camera command:
//   Camera "perspective" "float fov" [39]
type CameraCmd classCmd

// MaterialCmd stores parameters of material command:
//   Material "matte" "color Kd" [0 0 0]
type MaterialCmd classCmd

// ShapeCmd stores parameters of shape command:
//   Shape "sphere" "float radius" [3]
type ShapeCmd classCmd

// FilmCmd stores parameters of film command:
//   Film "image" "integer xresolution" [700] "integer yresolution" [700]
type FilmCmd classCmd

// SamplerCmd stores parameters of sampler command:
//   Sampler "halton" "integer pixelsamples" [8]
type SamplerCmd classCmd

// IntegratorCmd stores parameters of integrator command:
//   Integrator "path"
type IntegratorCmd classCmd

// LightSourceCmd stores parameters of integrator command:
//   LightSource "point" "rgb I" [ .5 .5 .5 ]
type LightSourceCmd classCmd

// AreaLightSourceCmd stores parameters of integrator command:
//   AreaLightSource "diffuse" "rgb L" [ .5 .5 .5 ]
type AreaLightSourceCmd classCmd

// MakeNamedMediumCmd stores parameters of make named medium command:
//   MakeNamedMedium "mymedium" "string type" "homogeneous" "rgb sigma_s" [100 100 100]
type MakeNamedMediumCmd classCmd

// MakeNamedMaterialCmd store parameters of make named material command:
//   MakeNamedMaterial "myplastic" "string type" "plastic" "float roughness" [0.1]
type MakeNamedMaterialCmd classCmd

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

func parseClassCmd(rawCommand string) (interface{}, error) {
	tokens := toTokens(rawCommand)
	cmd := classCmd{}

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
	switch class {
	case "Camera":
		return CameraCmd(cmd), nil
	case "Shape":
		return ShapeCmd(cmd), nil
	case "Material":
		return MaterialCmd(cmd), nil
	case "Film":
		return FilmCmd(cmd), nil
	case "Sampler":
		return SamplerCmd(cmd), nil
	case "Integrator":
		return IntegratorCmd(cmd), nil
	case "LightSource":
		return LightSourceCmd(cmd), nil
	case "AreaLightSource":
		return AreaLightSourceCmd(cmd), nil
	case "MakeNamedMedium":
		return MakeNamedMediumCmd(cmd), nil
	case "MakeNamedMaterial":
		return MakeNamedMaterialCmd(cmd), nil
	default:
		return nil, errors.New("Class name " + class + " no match")
	}
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
	return cmd, nil
}
