# psutil-golang

A Go package for retrieving system and process information on Linux systems.

## Project Structure

```
.
├── go.mod
├── go.sum
├── psutils
│   ├── cpu.go
│   ├── erros.go
│   ├── internal.go
│   ├── internal_test.go
│   ├── mem.go
│   ├── proc.go
│   └── psutil_test.go
└── README.md
```

## Features

- CPU information retrieval
- Memory information retrieval
- Process information retrieval

## Installation

To use this package in your Go project, run:

```bash
go get github.com/codescalersinternships/psutil-golang-MohamedFadel
```

## Usage

Here's a basic example of how to use the package:

```go
package main

import (
    "fmt"
    "github.com/codescalersinternships/psutil-golang-MohamedFadel/psutils"
)

func main() {
    // Get CPU info
    cpuInfo, err := psutils.GetCPUInfo()
    if err != nil {
        fmt.Printf("Error getting CPU info: %v\n", err)
    } else {
        fmt.Printf("CPU Info: %+v\n", cpuInfo)
    }

    // Get Memory info
    memInfo, err := psutils.GetMemoryInfo()
    if err != nil {
        fmt.Printf("Error getting memory info: %v\n", err)
    } else {
        fmt.Printf("Memory Info: %+v\n", memInfo)
    }

    // Get Process info
    procInfo, err := psutils.GetProcessInfo()
    if err != nil {
        fmt.Printf("Error getting process info: %v\n", err)
    } else {
        fmt.Printf("Process Info: %+v\n", procInfo)
    }
}
```

## Testing

To run the tests, navigate to the project root and run:

```bash
go test ./...
```

## Dependencies

This project uses the following external dependency:

- github.com/stretchr/testify v1.9.0 (for testing)
