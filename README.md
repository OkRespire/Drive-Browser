# Drive Browser

A simple TUI file browser for your Google Drive.

## Prerequisites
- Go version 1.24.3 or newer (older versions may also work)
- A Google Cloud account

## Setup
This project requires Google OAuth2 credentials to access the Drive API.

1. Follow the [Authorize credentials for a desktop application](https://developers.google.com/workspace/drive/api/quickstart/go#authorize_credentials_for_a_desktop_application) guide.
2. Download your `credentials.json` file from Google Cloud Console.
3. Place the file in the project root directory.
4. Run the project with:
   ```bash
   go run main.go
   ```
