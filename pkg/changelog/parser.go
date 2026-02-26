package changelog

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// Load reads and parses a YAML changelog from the given path.
func Load(path string) (*Changelog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening changelog: %w", err)
	}
	defer func() { _ = f.Close() }()
	return LoadFromReader(f)
}

// LoadFromReader parses a YAML changelog from a reader.
func LoadFromReader(r io.Reader) (*Changelog, error) {
	var c Changelog
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("decoding YAML: %w", err)
	}
	if errs := Validate(&c); len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		return nil, fmt.Errorf("validation failed:\n  %s", strings.Join(msgs, "\n  "))
	}
	return &c, nil
}

// Validate checks a Changelog for structural and semantic errors.
// An optional Config can be passed to control category validation.
func Validate(c *Changelog, cfg ...*Config) []ValidationError {
	var errs []ValidationError

	if strings.TrimSpace(c.Project) == "" {
		errs = append(errs, ValidationError{Field: "project", Message: "must not be empty"})
	}

	var allowed []string
	if len(cfg) > 0 && cfg[0] != nil {
		allowed = cfg[0].AllowedCategories()
	} else {
		allowed = DefaultCategories
	}

	allowedSet := make(map[string]bool, len(allowed))
	for _, cat := range allowed {
		allowedSet[cat] = true
	}

	seen := map[string]bool{}
	unreleasedCount := 0

	for i, v := range c.Versions {
		prefix := fmt.Sprintf("versions[%d]", i)

		if strings.TrimSpace(v.Version) == "" {
			errs = append(errs, ValidationError{Field: prefix + ".version", Message: "must not be empty"})
			continue
		}

		normalized := NormalizeVersion(v.Version)
		if seen[normalized] {
			errs = append(errs, ValidationError{
				Field:   prefix + ".version",
				Message: fmt.Sprintf("duplicate version %q", v.Version),
			})
		}
		seen[normalized] = true

		if v.IsUnreleased() {
			unreleasedCount++
			if unreleasedCount > 1 {
				errs = append(errs, ValidationError{
					Field:   prefix + ".version",
					Message: "only one unreleased version allowed",
				})
			}
		} else {
			if v.Date == "" {
				errs = append(errs, ValidationError{
					Field:   prefix + ".date",
					Message: "date required for released versions",
				})
			} else if !dateRegex.MatchString(v.Date) {
				errs = append(errs, ValidationError{
					Field:   prefix + ".date",
					Message: fmt.Sprintf("invalid date format %q, expected YYYY-MM-DD", v.Date),
				})
			}
		}

		if v.IsEmpty() && v.Internal.IsEmpty() && !v.IsUnreleased() {
			errs = append(errs, ValidationError{
				Field:   prefix,
				Message: "must have at least one entry",
			})
		}

		// Validate public categories
		errs = append(errs, validateChanges(v.Public, prefix, allowed, allowedSet)...)

		// Validate internal categories
		errs = append(errs, validateChanges(v.Internal, prefix+".internal", allowed, allowedSet)...)
	}

	return errs
}

// validateChanges checks entries within a Changes value.
func validateChanges(changes Changes, prefix string, allowed []string, allowedSet map[string]bool) []ValidationError {
	var errs []ValidationError
	for _, cat := range changes.Categories {
		// Check category is allowed (only when allowed is non-nil = strict mode)
		if allowed != nil && !allowedSet[cat.Name] {
			errs = append(errs, ValidationError{
				Field:   prefix + "." + cat.Name,
				Message: fmt.Sprintf("unknown category %q", cat.Name),
			})
		}
		for j, entry := range cat.Entries {
			if strings.TrimSpace(entry) == "" {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("%s.%s[%d]", prefix, cat.Name, j),
					Message: "entry must not be empty",
				})
			}
		}
	}
	return errs
}

// Save marshals a Changelog to YAML and writes it to the given path.
func Save(c *Changelog, path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// NormalizeVersion strips a leading "v" and lowercases the version string.
func NormalizeVersion(version string) string {
	v := strings.ToLower(version)
	return strings.TrimPrefix(v, "v")
}
