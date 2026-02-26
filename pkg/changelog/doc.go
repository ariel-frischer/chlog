// Package changelog provides a Go library for YAML-first changelog management.
//
// It reads and writes CHANGELOG.yaml files following the Keep a Changelog
// convention, with support for validation, querying, rendering to Markdown,
// and programmatic release promotion.
//
// # Loading and querying
//
//	c, err := changelog.Load("CHANGELOG.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	latest := c.GetLatestRelease()
//	fmt.Println(latest.Version, latest.Date)
//
//	entries := c.GetLastN(5)
//	for _, e := range entries {
//		fmt.Printf("[%s] %s: %s\n", e.Version, e.Category, e.Text)
//	}
//
// # Rendering to Markdown
//
//	md := changelog.RenderMarkdown(c, changelog.RenderOptions{})
//	os.WriteFile("CHANGELOG.md", []byte(md), 0644)
//
// # Programmatic releases
//
//	if err := c.Release("2.0.0", "2024-06-01"); err != nil {
//		log.Fatal(err)
//	}
//	changelog.Save(c, "CHANGELOG.yaml")
//
// # Parsing from any reader
//
//	f, _ := os.Open("CHANGELOG.yaml")
//	c, err := changelog.LoadFromReader(f)
package changelog
