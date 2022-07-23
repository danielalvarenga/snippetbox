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
2. Fetch dependencies: `go mod download`
3. Starting the web server: `go run ./cmd/web`

## Useful Go commands often used

* Download dependencies: `go mod download`
* Remove all unused packages in _go.mod_ and _go.sum_: `go mod tidy`
* Verify is downloaded packages wasn't modified unexpectedly: `go mod verify`
* Adding dependency with the latest release for the major version v1: `go get github.com/foo/bar@v1`
* Adding dependency with the specific version: `go get github.com/foo/bar@v1.2.3`
* Upgrading to the latest minor or patch release: `go get -u github.com/foo/bar`
* Upgrading to specific version: `go get -u github.com/foo/bar@v2.0.0`
