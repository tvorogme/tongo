//go:build ignore

package main

import (
	"go/format"
	"os"
	"strings"

	"github.com/tonkeeper/tongo/abi/parser"
)

func main() {
	scheme, err := os.ReadFile("known.xml")
	if err != nil {
		panic(err)
	}
	interfaces, err := parser.ParseInterface(scheme)
	if err != nil {
		panic(err)
	}

	gen := parser.NewGenerator(nil, "test")

	err = gen.RegisterInterfaces(interfaces)
	if err != nil {
		panic(err)
	}

	types := gen.CollectedTypes()
	msgDecoder := gen.GenerateMsgDecoder()

	getMethods, err := gen.GetMethods()
	if err != nil {
		panic(err)
	}
	invocationOrder, err := gen.RenderInvocationOrderList()
	if err != nil {
		panic(err)
	}

	builder := strings.Builder{}
	_, err = builder.WriteString(`package abi
// Code autogenerated. DO NOT EDIT. 

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

`)
	if err != nil {
		panic(err)
	}
	_, err = builder.WriteString(types)
	if err != nil {
		panic(err)
	}
	_, err = builder.WriteString(msgDecoder)
	if err != nil {
		panic(err)
	}
	_, err = builder.WriteString(getMethods)
	if err != nil {
		panic(err)
	}
	_, err = builder.WriteString(invocationOrder)
	if err != nil {
		panic(err)
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	f, err := os.Create("generated.go")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(code)
	if err != nil {
		panic(err)
	}
}
