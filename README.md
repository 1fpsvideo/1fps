# 1fps Client

This is the client component of the 1fps Screen Share application. The server part can be found at https://github.com/1fpsvideo/server.

## Overview

1fps Client is a Go application that captures screenshots, encrypts them, and uploads them to the server at 1 frame per second (fps). It also includes real-time cursor position tracking.

## Features

- Screen capture at 1 fps
- End-to-end encryption of screenshots
- Real-time cursor position tracking
- Automatic session creation and management
- WebSocket-based communication for cursor updates

## Requirements

- Go (Golang)

## Usage

You can run the client application in one of the following ways:

1. Build and run:
   ```
   go build 1fps.go
   ./1fps
   ```

2. Run directly:
   ```
   go run github.com/1fpsvideo/1fps@v0.1.1
   ```
   Replace `v0.1.1` with the latest version from [tags](https://github.com/1fpsvideo/1fps/tags).

Note: Windows users should scroll down to the Windows section for specific compilation steps.

## Linux Users

Linux users might need to install additional dependencies for the screen capture and cursor tracking functionality to work correctly. Specifically, you may need to install the `libxtst-dev` package:

For Ubuntu or Debian-based distributions:

```
sudo apt install libxtst-dev
```

For other distributions, the package name(s) might be slightly different. Please refer to your distribution's package management system.

For more detailed information about dependencies and installation instructions for different Linux distributions, you can check the RobotGo library documentation:

https://github.com/go-vgo/robotgo?tab=readme-ov-file#ubuntu

## Windows Users

Compiling on Windows requires a few additional steps. Please follow these instructions:

1. Install Golang, for example from https://webinstall.dev/golang/
2. Install the GCC compiler pack from https://github.com/skeeto/w64devkit/releases
   - Download the exe file, which will automatically unpack (probably to your Downloads folder)
   - Run w64devkit.exe
3. In the w64devkit terminal, type:
   ```
   go env -w CGO_ENABLED=1
   ```
4. Run the main command from the 1fps.video website. It's better to copy the command directly from the website or use the latest version from the tags:
   ```
   go run github.com/1fpsvideo/1fps@v0.1.1
   ```

Please note that these steps are necessary until we produce binaries for Windows. We understand that compiling on Windows has been challenging for various software projects. We're currently in alpha, so please check back later for easier installation options with pre-compiled binaries.

## License

This is FSL software. For more information, visit [https://fsl.software/](https://fsl.software/).
