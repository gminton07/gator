# Gator: Blog Aggregator

## Required
- Go CLI: Main programming language
- Postgres: Local database backend

## Go Install
Compile and install the gator package. 
```bash
go install gator
```

Usage:
```bash
gator <command> [args]
```

## Setup
### Config file: ~/.gatorconfig.json

Fields:
- "db_url": Postgres URI for accessing the gator database
- "current_user_name": Populated by program

### SQL Schema
SQL schema (stored in sql/schema) can be handled with the goose CLI tool. It can be installed from:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Usage (using the same connection_string from the config file):
```bash
goose postgres <connection_string> <up/down>
```

## Commands
- addfeed <name string> <url string> - Add RSS feed to database and create feed_follow entry for current user
- agg - Run feed aggregation script and save posts
- browse [limit int = 2] - Shows first limit posts for current user
- login <user string> - Change which user is currently activated
- help - List all subcommands
- register <user string> - Add a user to database
- reset - Remove all users from database
- unfollow <url string> - Remove feed_follow entry

## Tables
- users
- feeds
- feed_follows
- posts