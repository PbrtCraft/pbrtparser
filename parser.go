package pbrtparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Cmd is the base type of all types of commands
type Cmd struct {
	CmdType string `json:"cmd_type"`
}

// IncludeCmd stores parameters of include command and
//   the cmds in the file:
//   Include "geometry/killeroo.pbrt"
type IncludeCmd struct {
	Cmd
	Filename string        `json:"filename"`
	Cmds     []interface{} `json:"cmds"`
}

func parseIncludeCmd(rawCommand string) (*IncludeCmd, error) {
	tokens := toTokens(rawCommand)
	if len(tokens) != 4 {
		return nil, ErrClassCommandForm
	}

	if tokens[1] != `"` || tokens[3] != `"` {
		return nil, ErrClassCommandForm
	}

	inc := IncludeCmd{Filename: tokens[2]}
	inc.CmdType = "Include"
	return &inc, nil
}

func (inc *IncludeCmd) resolve(dir string) error {
	sp, err := NewCmdsParser(path.Join(dir, inc.Filename))
	if err != nil {
		return fmt.Errorf("pbrtparser.IncludeCmd.resolve: %s", err)
	}
	defer sp.Close()
	cmds, err := sp.ParseCmds()
	if err != nil {
		return fmt.Errorf("pbrtparser.IncludeCmd.resolve: %s", err)
	}
	inc.Cmds = cmds
	return nil
}

// BlockCmd stores parameters of block command and
//   the cmds in the block:
//   BLOCKBegin
//     ...
//	 BLOCKEnd
type BlockCmd struct {
	Cmd
	Cmds []interface{} `json:"cmds"`
}

// CmdsParser parse the cmd in a file
// Usage:
//   sp, _ := pp.NewCmdsParser("test.pbrt")
//   cmds, _ := sp.ParseCmds()
//   sp.Close()
type CmdsParser struct {
	scanner  *bufio.Scanner
	filename string
	dir      string
	file     *os.File
	cmds     []interface{}
	prevLine string
}

// NewCmdsParser prepare scanner for file and returns CmdsParser
func NewCmdsParser(filename string) (*CmdsParser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("pbrtparser.NewCmdsParser: %s", err)
	}

	return &CmdsParser{
		scanner:  bufio.NewScanner(file),
		filename: filename,
		file:     file,
		dir:      filepath.Dir(filename),
	}, nil
}

// Close close the file opened in CmdsParser
func (sp *CmdsParser) Close() {
	sp.file.Close()
}

func (sp *CmdsParser) nextRawCommand() (string, error) {
	curCommand := sp.prevLine
	sp.prevLine = ""
	for {
		res := sp.scanner.Scan()
		line := strings.TrimLeft(sp.scanner.Text(), " 	")
		if !res {
			if curCommand == "" {
				return "", io.EOF
			}
			break
		}
		if line == "" {
			continue
		}
		if line[:1] == `"` {
			curCommand += " " + line
		} else {
			sp.prevLine = line
			break
		}
	}
	return curCommand, nil
}

func (sp *CmdsParser) rawToCommand(rawCommand string) (interface{}, error) {
	if rawCommand == "" {
		// Empty Command
		return nil, nil
	}

	if strings.HasPrefix(rawCommand, "LookAt") {
		return parseLookAtCmd(rawCommand)
	} else if strings.HasPrefix(rawCommand, "Rotate") {
		return parseRotateCmd(rawCommand)
	} else if strings.HasPrefix(rawCommand, "Texture") {
		return parseTextureCmd(rawCommand)
	} else if isUseCmd(rawCommand) {
		return parseUseCmd(rawCommand)
	} else if isClassCmd(rawCommand) {
		return parseClassCmd(rawCommand)
	} else if isXYZCmd(rawCommand) {
		return parseXYZCmd(rawCommand)
	} else if strings.HasPrefix(rawCommand, "Include") {
		inc, err := parseIncludeCmd(rawCommand)
		if err != nil {
			return nil, err
		}
		err = inc.resolve(sp.dir)
		if err != nil {
			return nil, err
		}
		return inc, nil
	} else if strings.HasPrefix(rawCommand, "AttributeBegin") ||
		strings.HasPrefix(rawCommand, "WorldBegin") ||
		strings.HasPrefix(rawCommand, "ObjectBegin") ||
		strings.HasPrefix(rawCommand, "TransformBegin") {
		attrCmd := BlockCmd{
			Cmds: []interface{}{},
		}

		// Remove "Begin" from Prefix
		blockName := strings.Split(rawCommand, " 	")[0]
		blockName = blockName[:len(blockName)-5]

		attrCmd.CmdType = blockName
		for {
			rawCommand, err := sp.nextRawCommand()
			if err != nil {
				return nil, errors.New("EOF occur in block")
			}
			if strings.HasPrefix(rawCommand, blockName+"End") {
				break
			}
			cmd, err := sp.rawToCommand(rawCommand)
			if err != nil {
				return nil,
					fmt.Errorf("File: %s\nCommand: %s\nErr: %s", sp.filename, rawCommand, err)
			}
			if cmd == nil {
				continue
			}
			attrCmd.Cmds = append(attrCmd.Cmds, cmd)
		}
		return &attrCmd, nil
	}
	return nil, errors.New("Command no match:\n" + rawCommand)
}

// ParseCmds return the cmds in the file of CmdsParser
func (sp *CmdsParser) ParseCmds() ([]interface{}, error) {
	cmds := []interface{}{}
	for {
		rawCommand, err := sp.nextRawCommand()
		if err == io.EOF {
			break
		}

		cmd, err := sp.rawToCommand(rawCommand)
		if err != nil {
			return nil,
				fmt.Errorf("File: %s\nCommand: %s\nErr: %s", sp.filename, rawCommand, err)
		}
		if cmd == nil {
			continue
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}
