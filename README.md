# sqlc-gen-r: Generate Type-Safe R Code from SQL

[![Awesome](https://awesome.re/badge.svg)](https://awesome.re)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

**sqlc-gen-r** is a `sqlc` plugin that supercharges your R data workflows by generating type-safe R functions directly from your SQL queries! Say goodbye to manual data wrangling and potential type-related bugs, and hello to a streamlined, robust R data access layer.

## ✨ Features

* **Type Safety:** Generate R functions with automatic type checks for your SQL queries, powered by r-lib standalone type checks.
* **SQL Injection Prevention:** Uses `glue::glue_sql` to create safe SQL queries and prevent vulnerabilities.
* **Easy Integration:** Seamlessly integrates with `sqlc` for a streamlined workflow.
* **Github Actions Ready:** Examples include automated code generation with Github Actions.
* **Customizable Templates:** Modify the code generation to fit your specific needs.

## 🚀 Getting Started

### Prerequisites

* **sqlc**: Make sure you have `sqlc` installed. See the [official sqlc documentation](https://sqlc.dev/docs/overview/install) for installation instructions.

### Installation

1. **Configure `sqlc.yaml`:** Add a plugin to your `sqlc.yaml`

    ```yaml
    version: "2"
    plugins:
      - name: sqlc-gen-r
        wasm:
            url: https://github.com/ghandic/sqlc-gen-r/releases/download/v0.1.0/main.wasm
            sha256: 3ffd8a5272cf0ff2452f9fd34ed6388bc0c8e0a27d859494a3a93b3f0151f619
    sql:
     - engine: "sqlite"
       queries: "queries.sql"
       schema: "schema.sql"
       codegen:
       - out: out
         plugin: sqlc-gen-r
         options:
           filename: "db.R"
    ```

2. **Create `schema.sql`:** Define your database schema in a `schema.sql` file. Example:

    ```sql
    CREATE TABLE authors (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        bio TEXT
    );
    ```

3. **Write SQL Queries in `query.sql`:**
    Add your SQL queries. Remember to use the `-- name:` annotation to name your queries.

    ```sql
    -- name: GetAuthor :one
    SELECT id, name, bio FROM authors WHERE id = ?;
    ```

### Usage

Run `sqlc generate` in your project directory. This will generate the R code specified in `db.R`.

## ⚙️ Configuration Options

| Option        | Description                                                     | Default Value |
| :------------ | :-------------------------------------------------------------- | :------------ |
| `filename`   | The name of the output R file.                               | `"db.R"`      |

## 📜 Examples

### Basic Example

Let's create a database to store author information, and then query this.

1. **Setup:**

    * Create `schema.sql`:

        ```sql
        CREATE TABLE authors (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            bio TEXT
        );
        ```

    * Create `query.sql`:

        ```sql
        -- name: GetAuthor :one
        SELECT id, name, bio FROM authors WHERE id = ?;
        ```

    * Ensure your `sqlc.yaml` is configured as above (specifying the wasm plugin etc)

2. **Generate the R code:**

    Run `sqlc generate` in your terminal.  This will generate a file called `db.R` in the project.

3. **Use the Generated Code in R:**

    ```r
    library(DBI)
    library(RSQLite)

    # Establish a database connection
    conn <- dbConnect(SQLite(), "example.db") # Replace "example.db" with your database path

    # Source the generated R code
    source("db.R")

    # Execute the generated function
    author <- GetAuthor(conn, 1)  # Get author with ID 1

    # Print the results
    print(author)

    # Disconnect from the database
    dbDisconnect(conn)
    ```

### Github Actions Integration

Here's how to automate code generation using Github Actions:

1. **Create a Workflow File:**
    Create a workflow file in `.github/workflows/generate_r.yaml`:

    ```yaml
    name: Generate R Code

    on:
        push:
            branches: [ main ]
        pull_request:
            branches: [ main ]

    jobs:
        generate:
            runs-on: ubuntu-latest
            steps:
            - uses: actions/checkout@v4
            - uses: sqlc-dev/setup-sqlc@v4
              with:
                sqlc-version: '1.18.0'

            - name: Generate R code
              run: (cd example && sqlc generate)

            - name: Commit changes (if any)
              run: |
                git config --local user.email "actions@github.com"
                git config --local user.name "GitHub Actions"
                git add example/db.R
                git diff --exit-code || git commit -m "Regenerate R code"
                git push
    ```

2. **Github Action setup:** The workflow will do the following

    * Check out your code, setup Go and install `sqlc`.
    * Run the `sqlc generate` code.
    * If there are code updates, the code will update the `db.R` and push back to the repo!

## 🤝 Contributing

Contributions are welcome! Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.
If you notice a bug or some missing code, please submit an issue or a pull request.

## 📝 License

This project is licensed under the MIT License.
