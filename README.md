# Keeper

A lightweight, secure CLI password manager written in Go.

## Overview

Keeper is a command-line password management tool that focuses on simplicity and ease of use.

## Features

- Store passwords securely
- Generate strong, random passwords
- Retrieve stored passwords
- Remove passwords from storage
- Simple, intuitive command-line interface

## Installation

```bash
# Using go install
go install github.com/michaelyang12/keeper@latest

# Or download the binary directly
# [TODO: Add download instructions when available]
```

## Usage

```bash
# Store a new credential
keeper store github johnnysmith mypassword

# Generate a new password
keeper generate --length 16

# Retrieve a password
keeper get github

# Remove a password
keeper remove github

# List all stored passwords
keeper list
```

## Security

Keeper stores all passwords using AES-GCM encryption to ensure data security. Passwords are encrypted before being stored locally on your machine.

## Requirements

- Go 1.18 or higher
- Unix-based system or Windows

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[TODO: Add license]

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
