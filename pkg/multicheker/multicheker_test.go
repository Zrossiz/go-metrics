package multicheker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/passes/assign"
)

func TestNew_LinterInitialization(t *testing.T) {
	linter := New()

	assert.NotNil(t, linter)
	assert.True(t, len(linter.checkers) > 0, "Linter should contain analyzers")

	found := false
	for _, checker := range linter.checkers {
		if checker == assign.Analyzer {
			found = true
			break
		}
	}
	assert.True(t, found, "MultichekerLinter should include the assign analyzer")
}
