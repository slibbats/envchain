// Package profile provides named environment profile management for envchain.
//
// A profile is a named, ordered list of env-file layer paths that together
// define the full environment for a given context (e.g. "dev", "staging",
// "prod"). Profiles are declared in a JSON configuration file:
//
//	{
//	  "profiles": [
//	    { "name": "dev",  "layers": [".env", ".env.dev"] },
//	    { "name": "prod", "layers": [".env", ".env.prod"] }
//	  ]
//	}
//
// Layers are listed in ascending priority order; later entries override
// earlier ones when the resolver builds the final environment chain.
//
// The active profile is selected by name via the CLI flag or the
// ENVCHAIN_PROFILE environment variable as a fallback.
package profile
