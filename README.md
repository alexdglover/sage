# sage
An alternative to Intuit's recently retired Mint application for tracking personal finances.

## Quickstart

1. Clone the repo locally
2. Run
   ```sh
   DROP_TABLES=true ADD_SAMPLE_DATA=true go run main.go
   ```

## Documentation

See [./docs](./docs) for full design documentation and user documentation.

## Caveats

Sage is designed to run locally on your machine. Data is provided in the format of CSV exports
from your financial institutions, so Sage doesn't need access to your bank accounts or
credentials. Data never leaves your machine, but this also means it is your responsibility
to store or backup the data used by Sage.
