package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/yetone/magpie/internal/appdir"
)

func codexConfigPath() string {
	return filepath.Join(appdir.CodexHome(), "config.toml")
}

// codexConfig is the part of Codex's config.toml magpie imports: the fields
// that pick a provider and a model, the providers themselves, and the
// profiles that pick one each. The TOML parser reads it, so a value is what
// TOML says it is — an escape decoded, a multi-line string one string — and a
// table is the one its header names: [model_providers."my.relay"] is the
// provider "my.relay", while [model_providers.my.relay] is a table inside the
// provider "my".
type codexConfig struct {
	Model            string                   `toml:"model"`
	ModelProvider    string                   `toml:"model_provider"`
	ModelCatalogJSON string                   `toml:"model_catalog_json"`
	ModelProviders   map[string]codexProvider `toml:"model_providers"`
	Profiles         map[string]codexProfile  `toml:"profiles"`
}

// codexProvider is one [model_providers.<id>] table: what magpie carries into
// the provider it adds. Its http_headers are read here whether the file writes
// them inline ({ "X-Relay-Key" = "abc" }) or as a table of their own.
type codexProvider struct {
	Name        string            `toml:"name"`
	BaseURL     string            `toml:"base_url"`
	WireAPI     string            `toml:"wire_api"`
	BearerToken string            `toml:"experimental_bearer_token"`
	HTTPHeaders map[string]string `toml:"http_headers"`
}

// codexProfile is one [profiles.<name>] table.
type codexProfile struct {
	ModelProvider    string `toml:"model_provider"`
	Model            string `toml:"model"`
	ModelCatalogJSON string `toml:"model_catalog_json"`
}

// parseCodexConfig reads Codex's config text. A text TOML's parser refuses is
// an error: nothing in it says which providers it holds, so none is imported.
func parseCodexConfig(b []byte) (codexConfig, error) {
	var cfg codexConfig
	err := toml.Unmarshal(b, &cfg)
	return cfg, err
}

func readCodexConfig(path string) ([]AppImport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg, err := parseCodexConfig(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var ids []string
	for id := range cfg.ModelProviders {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	var profiles []string
	for name := range cfg.Profiles {
		profiles = append(profiles, name)
	}
	slices.Sort(profiles)
	var out []AppImport
	for _, id := range ids {
		mp := cfg.ModelProviders[id]
		name := mp.Name
		if name == "" {
			name = id
		}
		it := AppImport{Ref: id, From: "config.toml"}
		if id == "magpie" {
			it.Provider, it.Skip = Provider{Name: name}, "it points at magpie itself"
			out = append(out, it)
			continue
		}
		var models []string
		if cfg.ModelProvider == id {
			models = append(models, cfg.Model)
			models = append(models, codexImportModels(path, cfg.ModelCatalogJSON)...)
		}
		for _, profileName := range profiles {
			profile := cfg.Profiles[profileName]
			providerID := profile.ModelProvider
			if providerID == "" {
				providerID = cfg.ModelProvider
			}
			if providerID != id {
				continue
			}
			models = append(models, profile.Model)
			catalogPath := profile.ModelCatalogJSON
			if catalogPath == "" && cfg.ModelProvider == id {
				catalogPath = cfg.ModelCatalogJSON
			}
			models = append(models, codexImportModels(path, catalogPath)...)
		}
		var eps endpoints
		if mp.WireAPI == "chat" {
			eps.chat = mp.BaseURL
		} else {
			eps.responses = mp.BaseURL
		}
		it.Provider, it.Skip = imported(name, mp.BearerToken, eps, cleanList(models))
		if it.Skip != "" {
			it.Provider = Provider{Name: name}
		} else {
			it.Provider.Headers = mp.HTTPHeaders
		}
		out = append(out, it)
	}
	return out, nil
}

func codexImportModels(configPath, catalogPath string) []string {
	if catalogPath == "" {
		return nil
	}
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(catalogPath, "~/") {
		catalogPath = filepath.Join(home, strings.TrimPrefix(catalogPath, "~/"))
	} else if !filepath.IsAbs(catalogPath) {
		catalogPath = filepath.Join(filepath.Dir(configPath), catalogPath)
	}
	if filepath.Clean(catalogPath) == filepath.Join(appdir.CodexHomeIn(home), "magpie-models.json") {
		return nil
	}
	b, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil
	}
	var list struct {
		Models []struct {
			Slug       string `json:"slug"`
			Visibility string `json:"visibility"`
		} `json:"models"`
	}
	if json.Unmarshal(b, &list) != nil {
		return nil
	}
	var models []string
	for _, m := range list.Models {
		if m.Slug != "" && m.Visibility != "hide" {
			models = append(models, m.Slug)
		}
	}
	return models
}
