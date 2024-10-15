package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func main() {
	sb := &strings.Builder{}
	loadEnums(sb)
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String())
}

func loadEnums(sb *strings.Builder) {
	enums := []string{
		`CREATE TYPE gender AS ENUM (
			'MALE',
			'FEMALE'
		);`,
	}
	for _, enum := range enums {
		sb.WriteString(enum)
		sb.WriteString(";\n")
	}
}

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&entity.User{},
	}
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
}
