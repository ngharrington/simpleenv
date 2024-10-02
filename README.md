# simpleenv

Simple CLI to store and retrieve sets of environment variables and output "sourceable" text to export these environment variables. Meant to assist in VM initialization with init scripts, where secrets are required and no mature secrets infrastructure is present. This is a prototype and is not secure.

## Features

- **Write Environment Variables**: Store environment variables in a DigitalOcean Space.
- **Read Environment Variables**: Retrieve and output environment variables in a format suitable for sourcing in shell scripts.

## Prerequisites

- Go 1.22 or later
- DigitalOcean Space credentials set via environment variables:
  - `DO_SPACE_NAME`
  - `DO_SPACE_REGION`
  - `DO_ACCESS_KEY`
  - `DO_SECRET_KEY`

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Build the CLI tool:
   ```bash
   go build -o simpleenv
   ```


## Usage

### Write Environment Variables

To write environment variables to a DigitalOcean Space:

```bash
./simpleenv write --id <env-id> --vars KEY1=value1 --vars KEY2=value2
```

- `--id`: A unique identifier for the environment variables.
- `--vars`: Environment variables in `KEY=VALUE` format.

### Read Environment Variables

To read and output environment variables in a sourceable format:

```bash
source <(./simpleenv read --id <env-id> --source)
```

- `--id`: The unique identifier for the environment variables.
- `--source`: Outputs the variables in `export KEY=VALUE` format for sourcing in shell scripts.

## Testing

Run the tests using the Go testing tool:

```bash
go test ./...
```

To run integration tests you must have environment variables set as described above
```bash
go test -tag integration ./...
```

