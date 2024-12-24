package main

const CONTAINER_FILE = "Containerfile"
const IGNORE_FILE = ".containerignore"
const IGNORE_CONTENT = `

`
type Project struct {
    PythonVersion string
    HasRequirements bool
    From string
}

