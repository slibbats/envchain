package filter_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/envchain/internal/filter"
)

// buildLargeEnv creates an env map with n entries alternating between
// APP_ prefixed keys and DB_ prefixed keys.
func buildLargeEnv(n int) map[string]string {
	env := make(map[string]string, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			env[fmt.Sprintf("APP_VAR_%d", i)] = fmt.Sprintf("value_%d", i)
		} else {
			env[fmt.Sprintf("DB_VAR_%d", i)] = fmt.Sprintf("value_%d", i)
		}
	}
	return env
}

func BenchmarkApply_PrefixRule(b *testing.B) {
	env := buildLargeEnv(1000)
	f := filter.New().WithPrefix("APP_")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Apply(env)
	}
}

func BenchmarkApply_PatternRule(b *testing.B) {
	env := buildLargeEnv(1000)
	f, err := filter.New().WithPattern(`^APP_VAR_\d+$`)
	if err != nil {
		b.Fatalf("bad pattern: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Apply(env)
	}
}

func BenchmarkMatch_NoRules(b *testing.B) {
	f := filter.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Match("SOME_KEY")
	}
}
