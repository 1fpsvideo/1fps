# 1fps Client

This is the client component of the 1fps Screen Share application. The server part can be found at https://github.com/1fpsvideo/server.

## Overview

1fps Client is a Go application that captures screenshots, encrypts them, and uploads them to the server at 1 frame per second (fps). It also includes real-time cursor position tracking.

## Features

- Screen capture at 1 fps
- End-to-end encryption of screenshots
- Real-time cursor position tracking
- WebSocket-based communication for cursor updates

## Requirements

- Go (Golang)

## Usage

You can run the client application in one of the following ways:

1. Build and run:
   ```shell
   go build 1fps.go
   ./1fps
   ```

2. Run directly either:
   ```shell
   go run github.com/1fpsvideo/1fps@latest
   ```
   or replace `latest` with specific version, example `v0.1.11`, from [tags](https://github.com/1fpsvideo/1fps/tags).

Note: Windows users should scroll down to the Windows section for specific compilation steps.

## Linux Users

Linux users might need to install additional dependencies for the screen capture and cursor tracking functionality to work correctly. Specifically, you may need to install the `libxtst-dev` package:

For Ubuntu or Debian-based distributions:

```shell
sudo apt install libxtst-dev
```

For other distributions, the package name(s) might be slightly different. Please refer to your distribution's package management system.

For more detailed information about dependencies and installation instructions for different Linux distributions, you can check the RobotGo library documentation:

https://github.com/go-vgo/robotgo?tab=readme-ov-file#ubuntu

## Windows Users

### Cross-compilation from Linux (Recommended)

Due to the complexity of Windows tooling for CGO-based projects, we recommend building Windows executables from Linux using cross-compilation:

1. Install MinGW-w64 cross-compiler:
   ```bash
   sudo apt update
   sudo apt install -y gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64
   ```

2. Build the Windows executable:
   ```bash
   export GOOS=windows
   export GOARCH=amd64
   export CGO_ENABLED=1
   export CC=x86_64-w64-mingw32-gcc
   export CXX=x86_64-w64-mingw32-g++
   go build -ldflags="-w -s -extldflags '-static'" -o 1fps.exe
   ```

This will produce a statically-linked `1fps.exe` that can run on Windows without additional dependencies.

### Native Windows Compilation

Native Windows compilation with CGO is challenging. If you have working instructions for compiling this project natively on Windows, please contribute them via a pull request or issue.

For now, you can try:
- Using WSL2 with the Linux cross-compilation method above
- Running the pre-built executable if available in releases
- Using `go run github.com/1fpsvideo/1fps@latest` with proper MinGW setup (may not work due to CGO issues)

## Development

For local development, you need to create a `.env` file in the root directory of the project. The contents of the `.env` file should be:

```
ENV=development
```

This configuration sets the environment to development mode, which may enable certain debugging features or use local server addresses.

## License

This is FSL software. For more information, visit [https://fsl.software/](https://fsl.software/).
