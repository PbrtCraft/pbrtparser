package main

import (
	"encoding/json"
	"flag"
	"fmt"

	pp "github.com/PbrtCraft/pbrtparser"
)

func main() {
	pbrtFilename := flag.String("file", "test.pbrt", "Parsing pbrt file")
	flag.Parse()

	sp, err := pp.NewCmdsParser(*pbrtFilename)
	if err != nil {
		panic(err)
	}
	defer sp.Close()

	cmds, err := sp.ParseCmds()
	if err != nil {
		panic(err)
	}

	bs, err := json.Marshal(cmds)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
}
