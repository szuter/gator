
## Getting Started

Gator is a Go-based application designed to manage RSS feeds and user posts. It leverages SQL for database interactions and provides a structured approach to handle various functionalities such as creating posts, fetching posts for users, and managing feeds.

### Prerequisites

- Go 1.21 or later
- PostgreSQL

### Installation

```sh
go install github.com/szuter/gator@latest
```

### Configuration

Create a `.gatorconfig.json` file in your home directory:

```json
{
    "database": {
        "host": "localhost",
        "port": 5432,
        "name": "gator_db",
        "user": "your_username",
        "password": "your_password"
    }
}
```

### Basic Commands

- `gator addfeed name url` - Create new feed
- `gator browse` - List all following posts
- `gator feeds` - List all feeds