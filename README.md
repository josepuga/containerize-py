# CONTAINERIZE PYTHON PROJECT

A simple yet effective utility that generates the necessary container files to create an image for the specified Python project.

## Features:
- By default, it creates images using `python-alpine`, but you can specify any `python-slim-bullseye`, `python-slim-bookworm`, etc.
- Automatically selects the Python version from the Virtual Environment, if present, or falls back to the system’s Python version.

## Installation:
- Executables for Linux and Windows are available in the `bin/` directory. They do not require additional libraries.
- If you wish to clone the repository and compile it yourself, it is recommended to use the `build.sh` script.

## Considerations:
- It has not been extensively tested, so proper functionality is not guaranteed for large projects with numerous dependencies.
- The generated files are `Containerfile` and `.containerignore`. Both are recognized by Docker (taking precedence over `Dockerfile` and `.dockerignore`) as they follow the new non-vendor-specific standard.
- A `main.py` and a `requirements.txt` file are required.

## Example Usage:
Although it has several options, the simplest way is to generate the files in the project directory and use the default Alpine image.
```bash
$ cd /path/to/my_project
$ containerize-py .
```
This will start the process and generate the necessary files.
```text
Checking directories...
Checking for main.py...
Checking for requirements.txt...
Getting Python project version...
Writing Ignore File...
Writing Container File...
Done!.
Check that no files containing sensitive information (passwords, API keys, etc.) are present in your project. If there are any, make sure to include them in .containerignore
```
Two files will be generated. `Containerfile`:
```Dockerfile
# Generated with containerize-py v0.1.0
FROM python:3.13.0-alpine
WORKDIR /app
COPY . .
RUN pip install --no-cache-dir -r requirements.txt
CMD ["python3", "main.py"]
```
And another to exclude unnecessary files, `.containerignore`:
```text
# Generated with containerize-py v0.1.0
# Virtual environments
.venv/
venv/
.env/
env/

# Cache
__pycache__/
*.pyc
*.pyo

# 
Containerfile
.containerignore
```
