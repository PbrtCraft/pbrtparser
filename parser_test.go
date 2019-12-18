package pbrtparser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestXYZ(t *testing.T) {
	sp, err := NewCmdsParser("test/xyz.pbrt")
	if err != nil {
		t.Error(err)
		return
	}
	defer sp.Close()

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}

	// Translate 150 0 20
	// Scale 150 0 20
	// Rotate 180 1 0 0
	want := []interface{}{
		&XYZCmd{
			Cmd: Cmd{"Translate"},
			X:   150,
			Y:   0,
			Z:   20,
		},
		&XYZCmd{
			Cmd: Cmd{"Scale"},
			X:   150,
			Y:   0,
			Z:   20,
		},
		&RotateCmd{
			Cmd:   Cmd{"Rotate"},
			Angle: 180,
			X:     1,
			Y:     0,
			Z:     0,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestClass(t *testing.T) {
	sp, err := NewCmdsParser("test/class.pbrt")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}
	defer sp.Close()

	// Camera "perspective" "float fov" [80]
	// Integrator "path"
	want := []interface{}{
		&ClassCmd{
			Cmd:  Cmd{"Camera"},
			Name: "perspective",
			Params: []*Param{
				&Param{
					Name: "fov",
					Type: "float",
					Val:  []float64{80},
				},
			},
		},
		&ClassCmd{
			Cmd:    Cmd{"Integrator"},
			Name:   "path",
			Params: []*Param{},
		},
		&ClassCmd{
			Cmd:  Cmd{"Film"},
			Name: "image",
			Params: []*Param{
				&Param{
					Name: "xr",
					Type: "integer",
					Val:  []int{600},
				},
				&Param{
					Name: "filename",
					Type: "string",
					Val:  "test.exr",
				},
			},
		},
		&TextureCmd{
			Cmd:   Cmd{"Texture"},
			Name:  "name",
			Type:  "type",
			Class: "class",
			Params: []*Param{
				&Param{
					Name: "f",
					Type: "bool",
					Val:  true,
				},
				&Param{
					Name: "ok",
					Type: "bool",
					Val:  false,
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestUse(t *testing.T) {
	sp, err := NewCmdsParser("test/use.pbrt")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}
	defer sp.Close()

	// NamedMaterial "mat1"
	want := []interface{}{
		&UseCmd{
			Cmd:  Cmd{"NamedMaterial"},
			What: "mat1",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestLookAt(t *testing.T) {
	sp, err := NewCmdsParser("test/lookat.pbrt")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}
	defer sp.Close()

	// LookAt 3 4 1.5 .5 .5 0 0 0 1
	want := []interface{}{
		&LookAtCmd{
			Cmd:  Cmd{"LookAt"},
			Vals: []float64{3, 4, 1.5, 0.5, 0.5, 0, 0, 0, 1},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestBlock(t *testing.T) {
	sp, err := NewCmdsParser("test/block.pbrt")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}

	/* AttributeBegin
	       Material "matte" "color Kd" [0 0 0]
	       AttributeBegin
	           Translate 150 0  20
	           Shape "sphere" "float radius" [3]
	       AttributeEnd
	   AttributeEnd
	*/
	want := []interface{}{
		&BlockCmd{
			Cmd: Cmd{"Attribute"},
			Cmds: []interface{}{
				&ClassCmd{
					Cmd:  Cmd{"Material"},
					Name: "matte",
					Params: []*Param{&Param{
						Name: "Kd",
						Type: "color",
						Val:  []float64{0, 0, 0},
					}},
				},
				&BlockCmd{
					Cmd: Cmd{"Attribute"},
					Cmds: []interface{}{
						&XYZCmd{
							Cmd: Cmd{"Translate"},
							X:   150,
							Y:   0,
							Z:   20,
						},
						&ClassCmd{
							Cmd:  Cmd{"Shape"},
							Name: "sphere",
							Params: []*Param{&Param{
								Name: "radius",
								Type: "float",
								Val:  []float64{3},
							}},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestInclude(t *testing.T) {
	sp, err := NewCmdsParser("test/include.pbrt")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := sp.ParseCmds()
	if err != nil {
		t.Error(err)
		return
	}
	defer sp.Close()

	// Include "use.pbrt"
	want := []interface{}{
		&IncludeCmd{
			Cmd:      Cmd{"Include"},
			Filename: "use.pbrt",
			Cmds: []interface{}{
				&UseCmd{
					Cmd:  Cmd{"NamedMaterial"},
					What: "mat1",
				},
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
