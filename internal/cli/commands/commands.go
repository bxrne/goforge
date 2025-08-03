package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bxrne/goforge/internal/templates"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	var template string
	var db string
	var auth string
	var goVersion string

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Go project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]

			// Auto-detect Go version if not provided
			if goVersion == "" {
				detectedVersion, err := detectGoVersion()
				if err != nil {
					fmt.Printf("Warning: Could not detect Go version, using default 1.21: %v\n", err)
					goVersion = "1.21"
				} else {
					goVersion = detectedVersion
				}
			}

			scaffoldProject(projectName, template, db, auth, goVersion)
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "api", "Template type (cli, api, grpc)")
	cmd.Flags().StringVarP(&db, "db", "d", "", "Database (sqllite, postgres, mongodb, mysql, oracle)")
	cmd.Flags().StringVarP(&auth, "auth", "a", "", "Authentication (jwt, oauth2)")
	cmd.Flags().StringVarP(&goVersion, "go-version", "g", "", "Go version to use (auto-detected if not specified)")

	return cmd
}

func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [package]",
		Short: "Add a package to the project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			packageName := args[0]
			fmt.Printf("Adding package: %s\n", packageName)
			// TODO: Implement package addition logic
		},
	}

	return cmd
}

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [resource]",
		Short: "Generate scaffolded resources",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resource := args[0]
			fmt.Printf("Generating resource: %s\n", resource)
			// TODO: Implement resource generation logic
		},
	}

	return cmd
}

func detectGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'go version': %w", err)
	}

	// Parse output like "go version go1.21.5 darwin/amd64"
	versionRegex := regexp.MustCompile(`go(\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return "", fmt.Errorf("could not parse Go version from: %s", strings.TrimSpace(string(output)))
	}

	return matches[1], nil
}

func scaffoldProject(projectName, template, db, auth, goVersion string) {
	// Initialize template engine
	engine, err := templates.NewEngine()
	if err != nil {
		fmt.Printf("Error initializing template engine: %v\n", err)
		return
	}

	// Create project directory
	projectPath := filepath.Join(".", projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	// Create basic project structure based on template type
	var dirs []string
	if template == "cli" {
		dirs = []string{
			filepath.Join(projectPath, "cmd", projectName),
			filepath.Join(projectPath, "internal", "cli"),
			filepath.Join(projectPath, "internal", "cli", "commands"),
			filepath.Join(projectPath, "pkg"),
		}
	} else {
		dirs = []string{
			filepath.Join(projectPath, "cmd", projectName),
			filepath.Join(projectPath, "pkg"),
			filepath.Join(projectPath, "internal", "api"),
			filepath.Join(projectPath, "configs"),
		}
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Prepare template data
	data := templates.TemplateData{
		ProjectName: projectName,
		Template:    template,
		Database:    db,
		Auth:        auth,
		GoVersion:   goVersion,
	}

	// Render templates based on project type
	if template == "api" {
		// Create go.mod
		if err := engine.RenderToFile("templates/api/go.mod.tmpl", filepath.Join(projectPath, "go.mod"), data); err != nil {
			fmt.Printf("Error creating go.mod: %v\n", err)
			return
		}

		// Create main.go
		mainPath := filepath.Join(projectPath, "cmd", projectName, "main.go")
		if err := engine.RenderToFile("templates/api/main.go.tmpl", mainPath, data); err != nil {
			fmt.Printf("Error creating main.go: %v\n", err)
			return
		}

		// Create .env file
		if err := engine.RenderToFile("templates/api/env.tmpl", filepath.Join(projectPath, ".env"), data); err != nil {
			fmt.Printf("Error creating .env: %v\n", err)
			return
		}

		// Create Dockerfile
		if err := engine.RenderToFile("templates/api/Dockerfile.tmpl", filepath.Join(projectPath, "Dockerfile"), data); err != nil {
			fmt.Printf("Error creating Dockerfile: %v\n", err)
			return
		}
	} else if template == "cli" {
		// Create go.mod
		if err := engine.RenderToFile("templates/cli/go.mod.tmpl", filepath.Join(projectPath, "go.mod"), data); err != nil {
			fmt.Printf("Error creating go.mod: %v\n", err)
			return
		}

		// Create main.go
		mainPath := filepath.Join(projectPath, "cmd", projectName, "main.go")
		if err := engine.RenderToFile("templates/cli/main.go.tmpl", mainPath, data); err != nil {
			fmt.Printf("Error creating main.go: %v\n", err)
			return
		}

		// Create CLI root command
		rootPath := filepath.Join(projectPath, "internal", "cli", "root.go")
		if err := engine.RenderToFile("templates/cli/root.go.tmpl", rootPath, data); err != nil {
			fmt.Printf("Error creating root.go: %v\n", err)
			return
		}

		// Create CLI commands
		commandsPath := filepath.Join(projectPath, "internal", "cli", "commands", "commands.go")
		if err := engine.RenderToFile("templates/cli/commands.go.tmpl", commandsPath, data); err != nil {
			fmt.Printf("Error creating commands.go: %v\n", err)
			return
		}

		// Create README.md
		if err := engine.RenderToFile("templates/cli/README.md.tmpl", filepath.Join(projectPath, "README.md"), data); err != nil {
			fmt.Printf("Error creating README.md: %v\n", err)
			return
		}

		// Create Makefile
		if err := engine.RenderToFile("templates/cli/Makefile.tmpl", filepath.Join(projectPath, "Makefile"), data); err != nil {
			fmt.Printf("Error creating Makefile: %v\n", err)
			return
		}
	}

	fmt.Printf("Project %s created successfully with template: %s, db: %s, auth: %s\n", projectName, template, db, auth)
}
