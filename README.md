# Live Reload Server in Go

A minimal live reload server built with Go. It listens for changes to `.html`, `.css`, and `.js` files in the `static` directory, and automatically refreshes the browser when any of these files are saved.

## Features

- Serve static files (`.html`, `.css`, `.js`).
- Automatically detect changes to files in the `static` directory.
- Refresh the browser upon file changes.

## Prerequisites

- Go (version 1.16 or above)

## Usage

- Run the Go server:

- Copy code

`
go run main.go

`

- Open http://localhost:8080 in your browser.

The browser will automatically reload when you save changes to any .html, .css, or .js file in the static directory.

## Future Goals

### 1. Add WebSocket Support

- Replace the client-side polling mechanism with WebSocket connections to improve efficiency.
- WebSockets will allow real-time communication between the server and browser, ensuring instant reloads without needing to repeatedly request the `/reload` endpoint.

### 2. Asset Minification

- Implement a build step to automatically minify static assets (HTML, CSS, and JavaScript).
- This will reduce file sizes, improve load times, and optimize performance for the browser.

### 3. Gzip Compression

- Add Gzip compression to the server responses to further reduce the size of files being sent over HTTP.
- Compressing static assets before serving them will result in faster delivery and better user experience, especially for large files.
