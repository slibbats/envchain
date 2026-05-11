// Package filter provides composable key-matching rules for environment
// variable maps.
//
// A Filter is built by chaining one or more rules:
//
//	f := filter.New().
//		WithPrefix("APP_").
//		WithSuffix("_SECRET")
//
//	matched := f.Apply(env)   // keep only matching keys
//	rest    := f.Exclude(env) // keep only non-matching keys
//
// Rules are evaluated with OR semantics: a key matches if it satisfies
// at least one rule. When no rules are registered every key matches,
// making the filter a transparent pass-through.
//
// Supported rule types:
//
//   - WithPrefix  – simple string prefix check
//   - WithSuffix  – simple string suffix check
//   - WithPattern – compiled regular expression (returns error on bad syntax)
//
// Filter is safe to copy by value after construction but should not be
// mutated concurrently.
package filter
