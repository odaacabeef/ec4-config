package main

import (
	"io"
	"os"

	_ "embed"

	"cuelang.org/go/cue/cuecontext"
)

//go:embed schema.cue
var schemaSource string

type Config struct {
	Setups []Setup
}

type Setup struct {
	Name   string
	Groups []Group
}

type Group struct {
	Name     string
	Settings []Setting
}

type Setting struct {
	EC Encoder
	PB PushButton
}

type Encoder struct {
	Channel int
	Number  int
	Lower   int
	Upper   int
	Display string
	Type    string
	Mode    string
}

type PushButton struct {
	Channel int
	Number  int
	Lower   int
	Upper   int
	Display string
	Type    string
	Mode    string
}

func parseConfig() (cfg Config, err error) {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	ctx := cuecontext.New()
	s := ctx.CompileString(schemaSource)
	c := ctx.CompileBytes(stdin)
	err = s.Unify(c).Decode(&cfg)
	return
}
