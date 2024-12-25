package main

// By José Puga. 2024. GPL3 License
import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	_ "text/template"
)

var version = "[unknown]" // Get from Git tag

const DEFAULT_IMAGE = "alpine"

const CONTAINER_FILE = "Containerfile"
const CONTAINER_FILE_CONTENT = `# Generated with containerize-py {{.Version}}
FROM python:{{.PythonVersion}}-{{.FromImage}}
WORKDIR /app
COPY {{.SourcePath}} .
RUN pip install --no-cache-dir -r requirements.txt
CMD ["python3", "main.py"]
`

const IGNORE_FILE = ".containerignore"
const IGNORE_FILE_CONTENT = `# Generated with containerize-py {{.Version}}
# Virtual enviroments
.venv/
venv/
.env/
env/

# Cache
__pycache__/
*.pyc
*.pyo

# 
{{.ContainerFile}}
{{.IgnoreFile}}
`

const NOPARAM_MSG = `containerize-py {{.Version}}
Creates the {{.ContainerFile}} and {{.IgnoreFile}} from a python project.
This files can be used to create/launch a Container with docker or podman.
Syntax:
   containerize-py [--from=python_image] <project_dir> [output_dir]

--from=python_image : By default the image is alpine. There are more available, some are:
 
  · alpine: The lighest, but can lead to compatibility issues with glibc libraries dependent. Popular libraries (e.g., pandas, numpy, scipy) may need manual compilation due to native dependencies.
  · slim-bookworm: Good for new projects.
  · slim-bullseye: For production projects requiring stability.
  · slim-buster: Is end of support. Maybe good for some very old projects.

<project_dir> : Your python project location

[output_dir]  : Container files location. If empty use same place as project.

If there is any virtual enviroment directory the python version of that enviroment will be used, if not, the system version instead.
`

const NOREQUIREMENTS_MSG = `requirements.txt not found. 
You can create it in 2 ways:
1. Using pipreqs (recommended)
  pipreqs /path/your_project

2. Using pip freeze (the Virtual Enviroment must be active):
  pip freeze > requirements.txt
`
const NOMAINPY_MSG = `main.py not found.
It's mandatory to use this file to launch your project.
`

type Containerize struct {
	Version       string
	SourcePath    string
	OutputPath    string
	PythonVersion string
	FromImage     string
	ContainerFile string
	IgnoreFile    string
}

func main() {
	var err error

	cont := Containerize{
		Version:       version,
		ContainerFile: CONTAINER_FILE,
		IgnoreFile:    IGNORE_FILE,
	}

	// Define flags
	fromFlag := flag.String("from", "alpine", "Base Python image to use for the container (e.g., alpine, slim-bookworm, slim-bullseye)")
	flag.Usage = func() {
		t := template.Must(template.New("usage").Parse(NOPARAM_MSG))
		t.Execute(os.Stderr, cont)
	}

	// Parse the flags
	flag.Parse()

	// Retrieve the positional arguments
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Get working directories
	sourcePath := args[0]
	outputPath := sourcePath
	if len(args) == 2 {
		outputPath = args[1]
	}

	// Set the last fields of the Containerize struct
	cont.SourcePath = sourcePath
	cont.OutputPath = outputPath
	cont.FromImage = *fromFlag

	// Check valid directories
	fmt.Println("Checking directories...")
	if !isDirectory(cont.SourcePath) {
		fmt.Fprintf(os.Stderr, "%s is not a directory.\n", cont.SourcePath)
		os.Exit(1)
	}
	if !isDirectory(cont.OutputPath) {
		fmt.Fprintf(os.Stderr, "%s is not a directory.\n", cont.OutputPath)
		os.Exit(1)
	}

	// Check requirements.txt
	fmt.Println("Checking for requirements.txt...")
	if !fileExists(filepath.Join(cont.SourcePath, "requirements.txt")) {
		fmt.Fprint(os.Stderr, NOREQUIREMENTS_MSG)
		os.Exit(1)
	}

	// Check main.py
	fmt.Println("Checking for requirements.txt...")
	if !fileExists(filepath.Join(cont.SourcePath, "main.py")) {
		fmt.Fprint(os.Stderr, NOMAINPY_MSG)
		os.Exit(1)
	}

	// Check project python version
	fmt.Println("Getting Python project version...")
	cont.PythonVersion, err = getPythonVersion(cont.SourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get Python version: %v\n", err)
		os.Exit(1)
	}

	// Write ignore file
	fmt.Printf("Writing Ignore File...\n")
	t := template.Must(template.New("").Parse(IGNORE_FILE_CONTENT))
	path := filepath.Join(cont.OutputPath, cont.IgnoreFile)
	err = writeFileWithTemplate(path, t, cont)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	// Write Container file
	fmt.Printf("Writing Container File...\n")
	t = template.Must(template.New("").Parse(CONTAINER_FILE_CONTENT))
	path = filepath.Join(cont.OutputPath, cont.ContainerFile)
	err = writeFileWithTemplate(path, t, cont)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done!.")
	fmt.Printf("Check that no files containing sensitive information (passwords, API keys, etc.) are present in your project. If there are any, make sure to include them in %s\n", cont.IgnoreFile)
}

func writeFileWithTemplate(path string, t *template.Template, ttags any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = t.Execute(f, ttags)
	return err

}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func getPythonVersion(projectPath string) (string, error) {
	// Some VE location...
	envDirs := []string{"venv", ".venv", "env", ".env"}

	// Look for virtual enviroments
	for _, envDir := range envDirs {
		venvPath := filepath.Join(projectPath, envDir)
		pythonExecPath := filepath.Join(venvPath, "bin", "python")
		if os.PathSeparator == '\\' { // if Windows... (hack!)
			pythonExecPath = filepath.Join(venvPath, "Scripts", "python.exe")
		}
		// Check if the executable exists in the VE
		if _, err := os.Stat(pythonExecPath); err == nil {
			return execPythonVersion(pythonExecPath)
		}
	}
	return execPythonVersion("python")
}

// Run the executable python and returns it version. Used by getPythonVersion
func execPythonVersion(pythonExecPath string) (string, error) {
	cmd := exec.Command(pythonExecPath, "--version")

	// Captures output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // `python --version` sends the output to stderr in some versions

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute %s: %v", pythonExecPath, err)
	}

	result := strings.TrimSpace(out.String()) // Rips newline characters too
	result = strings.Split(result, " ")[1]    //TODO: Check for bad string?
	return result, nil

}
