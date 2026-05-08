// Package resolver builds a resolved environment Chain from a layered
// configuration. It orchestrates loading env files via the loader package and
// stacking them into a chain.Chain, optionally seeding the base layer from
// the current OS environment.
//
// Typical usage:
//
//	cfg := resolver.Config{
//		InjectOS: false,
//		Layers: []resolver.LayerConfig{
//			{Name: "base",    FilePath: ".env"},
//			{Name: "staging", FilePath: ".env.staging"},
//		},
//	}
//	c, err := resolver.Build(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	val, ok := c.Get("DATABASE_URL")
//
// Layers are applied in order: later layers override earlier ones for the
// same key, so secrets defined in higher-priority files (e.g. prod) will
// shadow defaults without modifying the base file.
package resolver
