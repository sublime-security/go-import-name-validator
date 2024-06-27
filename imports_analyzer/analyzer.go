package imports_analyzer

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func GetAnalyzer(requiredNamingPairs, forbidden *StringSliceFlag) (*analysis.Analyzer, error) {
	ci := ConfiguredInspector{
		RequiredNamesRef:  requiredNamingPairs,
		ForbiddenPathsRef: forbidden,
	}

	return &analysis.Analyzer{
		Name:             "inspect_imports",
		Doc:              "checks for violations of user defined import constraints",
		Run:              ci.run,
		RunDespiteErrors: true,
		Requires:         []*analysis.Analyzer{},
	}, nil
}

func ParseRequiredNamePairs(values []string) ([]RequiredImportName, error) {
	var toReturn []RequiredImportName
	for _, val := range values {
		splitVal := strings.Split(val, "=")

		if splitLen := len(splitVal); splitLen == 1 {
			toReturn = append(toReturn, RequiredImportName{
				Path: splitVal[0],
				Name: "",
			})
		} else if splitLen == 2 {
			toReturn = append(toReturn, RequiredImportName{
				Path: splitVal[0],
				Name: splitVal[1],
			})
		} else {
			return nil, fmt.Errorf("unexpected required import format %s. Should be <path>=<?name>", val)
		}
	}

	return toReturn, nil
}

type RequiredImportName struct {
	Path string
	Name string
}

type ConfiguredInspector struct {
	RequiredNamesRef  *StringSliceFlag
	ForbiddenPathsRef *StringSliceFlag
}

func (ci ConfiguredInspector) run(pass *analysis.Pass) (interface{}, error) {
	requiredImportNaming, err := ParseRequiredNamePairs(*ci.RequiredNamesRef)
	if err != nil {
		return nil, err
	}

	forbiddenPaths := *ci.ForbiddenPathsRef

	for _, f := range pass.Files {
		fileName := pass.Fset.Position(f.Package).Filename

		// TODO: Make this an argument like -skip-files and pass in CI job instead
		if strings.HasSuffix(fileName, "internal/clients/google/google_mocks.go") {
			continue
		}

		for _, imp := range f.Imports {
			name := ""
			if imp.Name != nil {
				name = imp.Name.Name
			}
			path := strings.Trim(imp.Path.Value, `""`)

			// Go builds tests as separate binaries, and it looks like imports the package being tested as _test.
			// These can be skipped.
			if name == "_test" {
				continue
			}

			report := func(diagnosticMessage, fixMessage, newText string) {
				// Unclear where `fixMessage` ever appears in output
				var fixes []analysis.SuggestedFix
				if fixMessage != "" {
					fixes = append(fixes, analysis.SuggestedFix{
						Message: fixMessage,
						TextEdits: []analysis.TextEdit{{
							Pos:     imp.Pos(),
							End:     imp.End(),
							NewText: []byte(newText),
						}},
					})
				}

				pass.Report(analysis.Diagnostic{
					Pos:            imp.Pos(),
					End:            imp.End(),
					Message:        diagnosticMessage,
					Category:       "import-inspect",
					SuggestedFixes: fixes,
				})
			}

			for _, fp := range forbiddenPaths {
				if fp == path {
					report(
						fmt.Sprintf("import path %s is forbidden", path),
						fmt.Sprintf("import path %s must be deleted", path),
						"")
				}
			}

			for _, requirement := range requiredImportNaming {
				if path == requirement.Path && name != requirement.Name {
					report(
						fmt.Sprintf(`import path "%s" cannot be named "%s"`, path, name),
						fmt.Sprintf(`set import path "%s" name to "%s"`, path, requirement.Name),
						fmt.Sprintf(`%s "%s"`, requirement.Name, path))
				}

				if requirement.Name != "" && name == requirement.Name && path != requirement.Path {
					report(
						fmt.Sprintf(`import named "%s" must have "%s" for the path`, name, requirement.Path),
						// Most likely the import is correct and an alternate name should be used.
						"",
						"")
				}
			}
		}
	}

	return nil, nil
}
