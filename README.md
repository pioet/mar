<p align="center">
  <img src="docs/images/mar-logo.png" alt="Mar Logo" width="160">
</p>

<h1 align="center">
  <strong>mar</strong>
</h1>

**Mar** is a fast and efficient CLI tool for managing and navigating bookmarks. Built with **Golang**, it offers blazing fast performance.  It supports **not only web URLs but also local file paths**, and can export bookmarks to browser-compatible files (such as Chrome/Edge). With the use of **tag**, you can quickly access your frequently used bookmarks at lightning speed.

## Features

- **Quick access** — Organize and access bookmarks by tags (alias)
- **Auto-complete** — Intelligent suggestion of bookmark titles by scraping webpage metadata  
- **Multi-type Support** — Manage web URLs and local file paths  
- **Compatibility** — Export browser-importable bookmark files (HTML format)  
- **Fast** — Written in Golang for near-instant command execution  

## Installation

```bash
go install github.com/pioet/mar@latest
```

## Usage

| Command       | Description                         |
|---------------|-----------------------------------|
| `add`         | Add a new bookmark                 |
| `clear`       | Delete all saved bookmarks         |
| `edit`        | Modify an existing bookmark        |
| `export`      | Export bookmarks to a file         |
| `get`         | Retrieve a bookmark's URI          |
| `import`      | Import bookmarks from a plaintext file |
| `list`        | List all saved bookmarks           |
| `reindex`     | Reassign sequential IDs to bookmarks |
| `rm`          | Remove bookmarks by ID or tag      |
| `search`      | Search bookmarks by multiple keywords |
| `show`        | Show bookmarks by ID or tag        |
