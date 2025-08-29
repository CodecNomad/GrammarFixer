// Package projects provides functionality to manage and store grammar correction projects.
package projects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Project represents a saved grammar correction session
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	OriginalText  string    `json:"original_text"`
	CorrectedText string    `json:"corrected_text"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProjectManager handles project storage and retrieval
type ProjectManager struct {
	projectsDir string
}

// NewProjectManager creates a new project manager
func NewProjectManager() (*ProjectManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get home directory: %w", err)
	}

	projectsDir := filepath.Join(homeDir, ".grammarfixer", "projects")
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create projects directory: %w", err)
	}

	return &ProjectManager{projectsDir: projectsDir}, nil
}

// SaveProject saves a project to storage
func (pm *ProjectManager) SaveProject(originalText, correctedText, name string) (*Project, error) {
	project := &Project{
		ID:            fmt.Sprintf("proj_%d", time.Now().Unix()),
		Name:          name,
		OriginalText:  originalText,
		CorrectedText: correctedText,
		CreatedAt:     time.Now(),
	}

	if project.Name == "" {
		// Generate a name from the first few words of original text
		words := []rune(originalText)
		if len(words) > 30 {
			project.Name = string(words[:30]) + "..."
		} else {
			project.Name = originalText
		}
	}

	fileName := fmt.Sprintf("%s.json", project.ID)
	filePath := filepath.Join(pm.projectsDir, fileName)

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("unable to marshal project: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("unable to save project: %w", err)
	}

	return project, nil
}

// ListProjects returns all saved projects
func (pm *ProjectManager) ListProjects() ([]*Project, error) {
	files, err := ioutil.ReadDir(pm.projectsDir)
	if err != nil {
		return nil, fmt.Errorf("unable to read projects directory: %w", err)
	}

	var projects []*Project
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(pm.projectsDir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue // Skip files we can't read
		}

		var project Project
		if err := json.Unmarshal(data, &project); err != nil {
			continue // Skip invalid JSON files
		}

		projects = append(projects, &project)
	}

	return projects, nil
}

// GetProject retrieves a specific project by ID
func (pm *ProjectManager) GetProject(id string) (*Project, error) {
	fileName := fmt.Sprintf("%s.json", id)
	filePath := filepath.Join(pm.projectsDir, fileName)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("unable to parse project: %w", err)
	}

	return &project, nil
}
