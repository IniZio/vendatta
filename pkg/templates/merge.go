package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// TemplateData represents the merged template data
type TemplateData struct {
	Plugins map[string]*PluginData `yaml:"plugins"`
}

type PluginData struct {
	Skills   map[string]interface{} `yaml:"skills"`
	Rules    map[string]interface{} `yaml:"rules"`
	Commands map[string]interface{} `yaml:"commands"`
}

type PluginManifest struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
	Conditions  []struct {
		File string `yaml:"file"`
	} `yaml:"conditions,omitempty"`
}

func (m *Manager) Merge(baseDir string, enabledPlugins []string, extends []string) (*TemplateData, error) {
	data := &TemplateData{
		Plugins: make(map[string]*PluginData),
	}

	// Load extends as base configs
	if err := m.loadExtends(baseDir, extends, data); err != nil {
		return nil, fmt.Errorf("failed to load extends: %w", err)
	}

	// Load remote template repos
	templateReposDir := filepath.Join(baseDir, "remotes")
	if err := m.loadTemplateRepos(templateReposDir, data); err != nil {
		return nil, fmt.Errorf("failed to load template repos: %w", err)
	}

	// Load local plugin templates
	projectRoot := filepath.Dir(baseDir)
	pluginsDir := filepath.Join(baseDir, "plugins")
	if err := m.loadPluginTemplates(pluginsDir, projectRoot, data); err != nil {
		return nil, fmt.Errorf("failed to load plugin templates: %w", err)
	}

	// Load base templates last to ensure local rules take precedence
	baseTemplatesDir := filepath.Join(baseDir, "templates")
	basePlugin := m.getOrCreatePlugin(data, "base")
	if err := m.loadTemplatesFromDir(baseTemplatesDir, basePlugin); err != nil {
		return nil, fmt.Errorf("failed to load base templates: %w", err)
	}

	agentsDir := filepath.Join(baseDir, "agents")
	if err := m.applyAgentOverrides(agentsDir, data); err != nil {
		return nil, fmt.Errorf("failed to load agent overrides: %w", err)
	}

	return data, nil
}

func (m *Manager) loadExtends(baseDir string, extends []string, data *TemplateData) error {
	for _, extend := range extends {
		parts := strings.Split(extend, "/")
		if len(parts) != 2 {
			continue
		}
		// TODO: Implement fetching from GitHub
	}
	return nil
}

func (m *Manager) loadPluginManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest PluginManifest
	err = yaml.Unmarshal(data, &manifest)
	return &manifest, err
}

func (m *Manager) checkPluginConditions(projectRoot string, manifest *PluginManifest) bool {
	if len(manifest.Conditions) == 0 {
		return true // No conditions, always enabled
	}
	for _, cond := range manifest.Conditions {
		if _, err := os.Stat(filepath.Join(projectRoot, cond.File)); err == nil {
			return true // At least one condition met
		}
	}
	return false
}

func (m *Manager) applyAgentOverrides(agentsDir string, data *TemplateData) error {
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		agentName := entry.Name()
		agentDir := filepath.Join(agentsDir, agentName)

		overridePlugin := m.getOrCreatePlugin(data, "override")

		if err := m.applyOverrideForType(agentDir, "rules", overridePlugin.Rules); err != nil {
			return fmt.Errorf("failed to apply %s rules overrides: %w", agentName, err)
		}
		if err := m.applyOverrideForType(agentDir, "skills", overridePlugin.Skills); err != nil {
			return fmt.Errorf("failed to apply %s skills overrides: %w", agentName, err)
		}
		if err := m.applyOverrideForType(agentDir, "commands", overridePlugin.Commands); err != nil {
			return fmt.Errorf("failed to apply %s commands overrides: %w", agentName, err)
		}
	}

	return nil
}

func (m *Manager) applyOverrideForType(agentDir, templateType string, target map[string]interface{}) error {
	overrideDir := filepath.Join(agentDir, templateType)
	if _, err := os.Stat(overrideDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(overrideDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		isMd := strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".mdc")
		if !isMd {
			return nil
		}

		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if len(strings.TrimSpace(string(content))) == 0 {
			delete(target, filename)
			return nil
		}

		var overrideData map[string]interface{}
		frontmatter, mdContent := parseMarkdown(content)
		ruleData := make(map[string]interface{})
		for k, v := range frontmatter {
			ruleData[k] = v
		}
		ruleData["content"] = mdContent
		overrideData = map[string]interface{}{
			filename: ruleData,
		}

		for key, value := range overrideData {
			target[key] = value
		}

		return nil
	})
}
