// Package snapshot provides point-in-time capture and comparison of resolved
// environment variable states produced by envchain.
//
// # Capturing a snapshot
//
// After resolving environment variables via the resolver package, a snapshot
// can be created and persisted to disk:
//
//	snap := snapshot.New(resolvedEnv, layerNames)
//	if err := snap.Save(".envchain.snapshot.json"); err != nil {
//		log.Fatal(err)
//	}
//
// # Loading and diffing snapshots
//
// Two snapshots can be compared to understand what changed between runs or
// environment promotions:
//
//	base, _ := snapshot.Load("snap-before.json")
//	next, _ := snapshot.Load("snap-after.json")
//	for _, entry := range snapshot.Diff(base, next) {
//		fmt.Printf("%s [%s]: %q -> %q\n", entry.Key, entry.Status, entry.Old, entry.New)
//	}
//
// Snapshot files are written with mode 0600 to reduce the risk of accidental
// secret exposure on shared systems.
package snapshot
