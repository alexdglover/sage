# Development Setup

To set up your development environment for Sage:

## Prerequisites
- **Go 1.23.0** (required)
- **SQLite3** (for direct database access/testing)

## Install Dependencies
Run the following command to install Go dependencies:

```bash
go mod tidy
```

## Environment Variables
- `DROP_TABLES`: If set, deletes all data from all tables (use with caution)
- `ADD_SAMPLE_DATA`: If set, populates the database with sample data for testing

## Running the App Locally
Start the application with:

```bash
go run main.go
```

Then open your browser to [http://localhost:8080](http://localhost:8080).

## Database Access
You can interact directly with the database using:

```bash
sqlite3 sage.db
```

Or, if your database is elsewhere:

```bash
sqlite3 ./path/to/sage.db.file
```
