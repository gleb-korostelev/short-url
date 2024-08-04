// Package main demonstrates the usage of the multichecker from the go/analysis package,
// combining various standard and third-party analyzers along with a custom-built analyzer
// to enforce specific coding practices.

// This program builds a multichecker with a variety of analyzers from the
// golang.org/x/tools/go/analysis/passes package and the staticcheck.io suite,
// and includes a custom analyzer that prohibits calls to os.Exit within the main function.
package main

import (
	"go/ast"

	"github.com/fatih/errwrap/errwrap"
	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// Collect standard analyzers from golang.org/x/tools
	allAnalyzers := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		lostcancel.Analyzer,
		loopclosure.Analyzer,
		copylock.Analyzer,
		unreachable.Analyzer,
		assign.Analyzer,
		asmdecl.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		reflectvaluecompare.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		goone.Analyzer,
		errwrap.Analyzer,
	}

	for _, a := range staticcheck.Analyzers {
		if isClassSA(a.Analyzer.Name) {
			allAnalyzers = append(allAnalyzers, a.Analyzer)
		}
	}

	// Custom analyzer, prohibiting os.Exit in main
	allAnalyzers = append(allAnalyzers, osExitAnalyzer)

	// Create multichecker
	multichecker.Main(allAnalyzers...)
}

func isClassSA(name string) bool {
	return name[0:2] == "SA"
}

// osExitAnalyzer defines a custom analysis tool that flags the use of os.Exit in the main function.
// It's designed to enforce policies against abrupt application termination strategies in main packages,
// encouraging proper error handling and resource management.
var osExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "prohibits use of os.Exit in the main function",
	Run:  runOsExitAnalyzer,
}

// runOsExitAnalyzer is the analysis function that inspects Go code to find and report
// any calls to os.Exit within the main function. It uses AST inspection to identify such calls
// and reports them as errors, thus enforcing coding standards that avoid using os.Exit
// for controlling program termination in main
func runOsExitAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := fun.X.(*ast.Ident); ok && pkgIdent.Name == "os" && fun.Sel.Name == "Exit" {
					pass.Reportf(callExpr.Pos(), "call to os.Exit in main is prohibited")
				}
			}
			return true
		})
	}
	return nil, nil
}
