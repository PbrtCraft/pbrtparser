package main

import (
	"encoding/json"
	"fmt"

	pp "github.com/PbrtCraft/pbrtparser"
)

func main() {
	sp, _ := pp.NewCmdsParser("test.pbrt")
	cmds, _ := sp.ParseCmds()
	bs, _ := json.Marshal(cmds)
	fmt.Println(string(bs))
}
