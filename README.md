# shortener
Shortener is a tiny, self-hosted URL shortener written in Go.

## Installation
1. Clone the repository
```
git clone git@github.com:atomisadev/shortener.git
```
2. `cd` into the repository
```
cd shortener
```
3. Install the dependencies & Run the code
```
go get github.com/mattn/go-sqlite3
go run .
```

Make sure you have Go 1.13+ installed in order for it to work properly.

## Usage
```
# Actually shorten URLs (method: GET)
curl "http://localhost:8080/shorten?url=https://www.google.com"

# Visit a shortened URL (method: GET)
curl "http://localhost:8080/<alias>"

# Delete a shortened url (method: GET)
curl "http://localhost:8080/delete?alias=<alias>"
```
