# Snnippetbox

A application to share texts and snippets.

## Project Structure

* root
  * [cmd](./cmd/README.md)
    * web
  * [internal](./internal/README.md)
  * [ui](./ui/README.md)
    * html
    * static

## Creating local DB

1. Start the DB container executing `docke-compose up -d`;
2. Create the database executing the script in `db/create-db.sql`;
3. Create the application user executing the script in `db/create-user.sql` (user: web | password: password);
4. Populate the db executing the script in `db/seed.sql`

## Running application

See all application configurations: `go run ./cmd/web -help`

### Local

1. Starting DB: `docker-compose up -d`;
2. Starting the web server: `go run ./cmd/web`
