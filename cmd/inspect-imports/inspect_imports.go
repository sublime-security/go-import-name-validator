package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sublime-security/go-import-name-validator/imports_analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var (
		requiredNames  imports_analyzer.StringSliceFlag
		forbiddenPaths imports_analyzer.StringSliceFlag
	)

	flag.Var(&requiredNames, "require-name", "A pair of import paths to required name. Name may be omitted to require an unnamed import. E.g. github.com/pkg/errors=pErrors or errors=")
	flag.Var(&forbiddenPaths, "forbidden", "An import path which is forbidden.")

	analyzer, err := imports_analyzer.GetAnalyzer(&requiredNames, &forbiddenPaths)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Main calls flag.Parse after adding its own flags
	singlechecker.Main(analyzer)
}
