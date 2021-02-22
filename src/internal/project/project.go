package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/home"
	"github.com/Samasource/jen/src/internal/spec"
	"gopkg.in/yaml.v2"
)

// GetProjectDir returns the project's root dir. It finds it by looking for the jen.yaml file
// in current working dir and then walking up the directory structure until it reaches the
// volume's root dir. If it doesn't find it, it returns an empty string.
func GetProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		path := filepath.Join(dir, constant.ProjectFileName)
		if helpers.PathExists(path) {
			return dir, nil
		}
		if dir == "/" {
			return "", nil
		}
		dir = filepath.Dir(dir)
	}
}

// Project represents the configuration file in a project's root dir
type Project struct {
	Version       string
	TemplateName  string
	Vars          map[string]interface{}
	Dir           string   `yaml:"-"`
	OverridenVars []string `yaml:"-"`
}

// Save saves project file into given project directory
func (p Project) Save() error {
	p.Version = constant.ProjectFileVersion
	doc, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	path := filepath.Join(p.Dir, constant.ProjectFileName)
	return ioutil.WriteFile(path, doc, os.ModePerm)
}

// Load loads the project file from given project directory
func Load(dir string) (*Project, error) {
	specFilePath := filepath.Join(dir, constant.ProjectFileName)
	buf, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, fmt.Errorf("loading project file: %w", err)
	}
	var project Project
	err = yaml.Unmarshal(buf, &project)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling project file yaml: %w", err)
	}

	if project.Version != constant.ProjectFileVersion {
		return nil, fmt.Errorf("unsupported jen project file version %s (expected %s)", project.Version, constant.ProjectFileVersion)
	}

	project.Dir = dir
	return &project, nil
}

var varOverrideRegexp = regexp.MustCompile(`^(\w+)=(.*)$`)

// LoadOrCreate loads current project file and, if it doesn't
// exists, prompts user whether to create it.
func LoadOrCreate(templateName string, skipConfirm bool, varOverrides []string) (*Project, error) {
	projectDir, err := GetProjectDir()
	if err != nil {
		return nil, err
	}
	if projectDir == "" {
		if !skipConfirm {
			err := confirmCreateProject()
			if err != nil {
				return nil, err
			}
		}
		proj := Project{}
		if err := proj.Save(); err != nil {
			return nil, err
		}
	}

	proj, err := Load(projectDir)
	if err != nil {
		return nil, err
	}

	if templateName != "" {
		proj.TemplateName = templateName
		if err := proj.Save(); err != nil {
			return nil, err
		}
	}

	templatesDir, err := home.GetTemplatesDir()
	if err != nil {
		return nil, err
	}
	if proj.TemplateName == "" {
		proj.TemplateName, err = promptTemplate(templatesDir)
		if err != nil {
			return nil, fmt.Errorf("prompting for template: %w", err)
		}
		if err := proj.Save(); err != nil {
			return nil, err
		}
	}

	// Apply command-line variable overrides
	for _, entry := range varOverrides {
		submatch := varOverrideRegexp.FindStringSubmatch(entry)
		if submatch == nil {
			return nil, fmt.Errorf("failed to parse set variable %q", entry)
		}
		name := submatch[1]
		value := submatch[2]
		proj.Vars[name] = value
		proj.OverridenVars = append(proj.OverridenVars, name)
	}
	if len(varOverrides) > 0 {
		if err := proj.Save(); err != nil {
			return nil, err
		}
	}

	return proj, nil
}

func confirmCreateProject() error {
	var result bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Jen project not found. Do you want to initialize current directory as your project root?",
		Default: false,
	}, &result)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("cancelled by user")
	}
	return nil
}

func promptTemplate(templatesDir string) (string, error) {
	// Read templates dir
	infos, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return "", fmt.Errorf("reading templates directory %q: %w", templatesDir, err)
	}

	// Build list of choices
	var templates []string
	var titles []string
	for _, info := range infos {
		template := info.Name()
		if strings.HasPrefix(template, ".") {
			continue
		}
		templateDir := filepath.Join(templatesDir, template)
		spec, err := spec.Load(templateDir)
		if err != nil {
			return "", err
		}
		templates = append(templates, template)
		titles = append(titles, fmt.Sprintf("%s - %s", template, spec.Description))
	}

	// Any templates found?
	if len(templates) == 0 {
		return "", fmt.Errorf("no templates found in %q", templatesDir)
	}

	// Prompt
	prompt := &survey.Select{
		Message: "Select template",
		Options: titles,
	}
	var index int
	if err := survey.AskOne(prompt, &index); err != nil {
		return "", err
	}

	return templates[index], nil
}

// GetTemplateDir returns the path to this project's template
func (p Project) GetTemplateDir() (string, error) {
	templatesDir, err := home.GetTemplatesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(templatesDir, p.TemplateName), nil
}
