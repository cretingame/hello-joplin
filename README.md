# Joplin Fuse

Joplin is an amazing note taking application. I just need an interface to edit and search my notes without using the current user interface. So, I want to instantiate a filesystem that represent my Joplin notes.

## Features

- Mount Joplin notes as a read-only filesystem
- Read notes directly from your terminal or file explorer
- Simple command-line interface

## Installation

To install Joplin Fuse, use the following command:

```bash
go install github.com/cretingame/joplin-fuse@latest
```

## Usage

```bash
joplin-fuse [MOUNTING POINT]
```

For example:

```bash
joplin-fuse ~/JoplinMount
```

This will mount your Joplin notes at the specified mount point.

## Requirements

- FUSE installed and configured on your system
- Joplin desktop app or server with API access enabled

## API Configuration

Joplin Fuse connects to the Joplin API. Make sure the API is enabled in your Joplin settings.

Default configuration assumes the API is available at http://127.0.0.1:41184. If you're using a different host or port, you may need to set environment variables or provide configuration options (update this section as appropriate to your implementation).

You can check the API status by visiting:

http://127.0.0.1:41184/ping

If the response is `"JoplinClipperServer"`, the API is running.

## Building

Clone the repository and install dependencies:

```bash
git clone https://github.com/yourusername/joplin-fuse.git
cd joplin-fuse
go build
```

You can then run the binary:

```bash
./joplin-fuse [MOUNTING POINT]
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

GNU Affero General Public License v3.0 or later
