package engine

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rampartfw/rampart/internal/model"
	"gopkg.in/yaml.v3"
)

// ResolveIncludes merges included policy files into the provided PolicySetYAML.
func ResolveIncludes(ps *model.PolicySetYAML, basePath string) error {
	return resolveIncludesRecursive(ps, basePath, 0)
}

func resolveIncludesRecursive(ps *model.PolicySetYAML, basePath string, depth int) error {
	if depth > 10 {
		return fmt.Errorf("include depth exceeded (max 10), possible circular include")
	}

	for _, inc := range ps.Includes {
		var data []byte
		var err error

		if inc.URL != "" {
			data, err = fetchURL(inc.URL)
		} else {
			path := inc.Path
			if !filepath.IsAbs(path) {
				path = filepath.Join(filepath.Dir(basePath), path)
			}
			data, err = os.ReadFile(path)
		}

		if err != nil {
			return fmt.Errorf("include %s%s: %w", inc.Path, inc.URL, err)
		}

		var included model.PolicySetYAML
		if err := yaml.Unmarshal(data, &included); err != nil {
			return fmt.Errorf("unmarshal include %s%s: %w", inc.Path, inc.URL, err)
		}

		// Recursively resolve includes for the newly included file
		if err := resolveIncludesRecursive(&included, inc.Path, depth+1); err != nil {
			return err
		}

		// Merge included policies
		ps.Policies = append(ps.Policies, included.Policies...)
	}

	// Clear includes to prevent re-processing
	ps.Includes = nil
	return nil
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
