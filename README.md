# SQL Helper

A simple Go application that executes SQL files from a directory against an MSSQL database.

## Usage

```bash
go run main.go --dir /path/to/sql/files --conn "server=localhost;user id=sa;password=yourpassword;database=yourdatabase"
```

## Parameters

- `--dir`: Directory containing SQL files (.sql extension)
- `--conn`: MSSQL connection string

## Example

```bash
go run main.go --dir ./stored_procedures --conn "server=myserver;user id=myuser;password=mypass;database=mydb"
```

The application will traverse the specified directory and execute all `.sql` files against the MSSQL database.