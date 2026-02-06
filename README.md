# GO JSON Server

A lightweight and Dependency-free Fake JSON API Server implemented in Go.
Converted into Go and inspired by the original [`JSON-Server`](https://github.com/typicode/json-server).

## This project was built as a learning exercise during my stay at **Kodehode**

## Features

- **Zero External dependencies** - Only utilizing the GO standard Library
- **Automatic JSON DB creation** - No need to make any directories or files, automatically creates a json file in (`data/db.json`)
- **Dynamic Collections** - (Created on first POST request to that collection's name)
- **Full CRUD API**
- **Middleware** - Logging & CORS
- **CORS Support** (simple, permissive, json-server style)

---

## Installation

Clone or Fork the repo. No need to install anything else - Just Build and run:

```
go run ./cmd/jsonserver
```

The server will start with a default empty JSON database at `data/db.json`
If the file or directory doesn't exist, the code will do it for you.

---

## API Endpoints

| Method | Path               | Description                                         |
| ------ | ------------------ | --------------------------------------------------- |
| GET    | /{collection}      | Lists all entries in a collection                   |
| GET    | /{collection}/{id} | Get a single entry within a collection              |
| POST   | /{collection}      | Create a new entry (and collection if it's missing) |
| PUT    | /{collection}/{id} | Replace an entry                                    |
| PATCH  | /{collection}/{id} | Update an entry                                     |
| DELETE | /{collection}/{id} | Delete an entry                                     |

### Health Check

| Method | Path    | Description               |
| ------ | ------- | ------------------------- |
| GET    | /health | Returns `{"status":"ok"}` |

---

## CORS

CORS is enabled by default with permissive settings:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type`

Preflight(`OPTIONS`) requests are handled automatically.

---

## Project Structure

```
go-json-server/
├── cmd/
│   └── jsonserver/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── router.go
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── health.go
│   ├── db/
│   │   └── readwrite.go
│   ├── model/
│   │   └── data.go
│   └── service/
│       └── service.go
├── static/
│   └── index.html
└── data/
    └── db.json (auto-created)
```

## Credits

- Original Concept and API design by [`Typicode`](https://github.com/typicode)
- GO implementation written by [`Me`](https://github.com/OleKodehode) _(Admittedly with help from LLMs for breaking down some concepts and steps)_

---

## License

This is open source. Feel free to use it for learning, experimentation and modification.
It probably needs some modification either way if you want to run it on a server. This works first and foremost in a dev environment. _But_ it should be a good starting point?
