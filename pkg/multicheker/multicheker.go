// Package multicheker provides a Go analysis tool that integrates multiple static analysis checkers
// into a single linter. It uses a combination of built-in and third-party checkers to perform a variety
// of checks on Go code.
//
// The `MultichekerLinter` struct is designed to run all the analyzers in a single pass, providing
// an easy way to utilize multiple static analysis checks simultaneously.
//
// The linter includes checks for common Go coding issues, such as race conditions, incorrect function
// usage, unnecessary type assertions, unhandled errors, unused variables, and more. It integrates
// checkers from the Go community, as well as custom analyzers, like the ones from the `honnef.co` tools
// and the `Zrossiz/go-metrics` package.
//
// Key functionality:
//   - Initializes multiple static analyzers that perform various types of checks.
//   - Executes all analyzers with the `multichecker` tool to aggregate results into a unified output.
package multicheker

import (
	"github.com/Zrossiz/go-metrics/pkg/noexit"
	"github.com/gostaticanalysis/builtinprint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/timakin/bodyclose/passes/bodyclose"
)

type MultichekerLinter struct {
	// checkers contains a list of analysis.Analyzer instances
	// that will be executed during the analysis process.
	checkers []*analysis.Analyzer
}

// New creates a new instance of MultichekerLinter with a preconfigured
// list of analyzers from the Go tools, community tools, and custom analyzers.
//
// This function initializes the linter with common static analysis checks,
// including checks for race conditions, code formatting, performance issues,
// and potential bugs.
//
// It returns a new instance of MultichekerLinter.
func New() MultichekerLinter {
	checkers := []*analysis.Analyzer{
		// Go analysis passes
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}

	// Add analyzers from external packages like honnef.co and staticcheck
	for _, analyzers := range [][]*lint.Analyzer{
		simple.Analyzers,
		stylecheck.Analyzers,
		quickfix.Analyzers,
		staticcheck.Analyzers,
	} {
		for _, v := range analyzers {
			checkers = append(checkers, v.Analyzer)
		}
	}

	// Add custom and other third-party analyzers
	checkers = append(checkers, bodyclose.Analyzer)
	checkers = append(checkers, builtinprint.Analyzer)
	checkers = append(checkers, noexit.Analyzer)

	return MultichekerLinter{checkers}
}

// Run executes all the analyzers that have been added to the MultichekerLinter.
// This method uses the multichecker tool to run each analysis pass and aggregates
// the results into a unified output, which can be displayed to the user.
func (m MultichekerLinter) Run() {
	// Run all the analyzers using multichecker.Main
	multichecker.Main()
}
