# Pbrt Parser

![](https://travis-ci.org/PbrtCraft/pbrtparser.svg?branch=master)

## Example

```golang
package main

import (
	"encoding/json"
	"fmt"

	pp "github.com/PbrtCraft/pbrtparser"
)

func main() {
	sp, _ := pp.NewCmdsParser("test.pbrt")
	defer sp.Close()
	cmds, _ := sp.ParseCmds()
	bs, _ := json.Marshal(cmds)
	fmt.Println(string(bs))
}
```

### test.pbrt

```
LookAt 0 0 0 1 1 1 1 -1 0 
Rotate -5 0 0 1
Camera "perspective" "float fov" [80]
Film "image"
    "integer xresolution" [600] "integer yresolution" [600]
    "string filename" "test.exr"

Sampler "halton" "integer pixelsamples" [8]
Integrator "path"

WorldBegin

AttributeBegin
    Material "matte" "color Kd" [0 0 0]
    Translate 150 0  20
    AreaLightSource "area"  "color L" [1000 1000 1000] "integer nsamples" [128]
    Shape "sphere" "float radius" [3]
AttributeEnd

WorldEnd
```

### JSON Output

```json
[
  {
    "cmd_type": "LookAt",
    "vals": [0, 0, 0, 1, 1, 1, 1, -1, 0]
  },
  {
    "cmd_type": "Rotate", "angle": -5, "x": 0, "y": 0, "z": 1
  },
  {
    "cmd_type": "Camera",
    "name": "perspective",
    "params": [
      {
        "name": "fov",
        "type": "float",
        "val": [80]
      }
    ]
  },
  {
    "cmd_type": "Film",
    "name": "image",
    "params": [
      {
        "name": "xresolution",
        "type": "integer",
        "val": [600]
      },
      {
        "name": "yresolution",
        "type": "integer",
        "val": [600]
      },
      {
        "name": "filename",
        "type": "string",
        "val": "test.exr"
      }
    ]
  },
  {
    "cmd_type": "Sampler",
    "name": "halton",
    "params": [
      {
        "name": "pixelsamples",
        "type": "integer",
        "val": [8]
      }
    ]
  },
  {
    "cmd_type": "Integrator",
    "name": "path",
    "params": []
  },
  {
    "cmd_type": "World",
    "cmds": [
      {
        "cmd_type": "Attribute",
        "cmds": [
          {
            "cmd_type": "Material",
            "name": "matte",
            "params": [
              {
                "name": "Kd",
                "type": "color",
                "val": [0, 0, 0]
              }
            ]
          },
          {
            "cmd_type": "Translate", "x": 150, "y": 0, "z": 20
          },
          {
            "cmd_type": "AreaLightSource",
            "name": "area",
            "params": [
              {
                "name": "L",
                "type": "color",
                "val": [1000, 1000, 1000]
              },
              {
                "name": "nsamples",
                "type": "integer",
                "val": [128]
              }
            ]
          },
          {
            "cmd_type": "Shape",
            "name": "sphere",
            "params": [
              {
                "name": "radius",
                "type": "float",
                "val": [3]
              }
            ]
          }
        ]
      }
    ]
  }
]
```

