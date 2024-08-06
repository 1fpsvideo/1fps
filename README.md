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

## License

This is FSL software. For more information, visit [https://fsl.software/](https://fsl.software/).
