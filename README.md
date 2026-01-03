# Gator
### RSS feed aggregator in Go!

it's a CLI tool that allows users to:

1. Add RSS feeds from across the internet to be collected
2. Store the collected posts in a PostgreSQL database
3. Follow and unfollow RSS feeds that other users have added
4. View summaries of the aggregated posts in the terminal, with a link to the full post

## Requirements
- [golang](https://github.com/hardmanhimself/gator?tab=readme-ov-file) (version 1.16 or higher)
- [postgresql](https://www.postgresql.org/) 15+

## Installation
run the command:
```
go install github.com/ArashPoorazam/Gator
```

## Run postgres 
macOS with [brew](https://brew.sh/):
```
brew install postgresql@15
```
Linux / WSL (Debian). Here are the [docs from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql), but simply:
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
Ensure the installation worked. The psql command-line utility is the default client for Postgres. Use it to make sure you're on version 15+ of Postgres:
```
psql --version
```
(Linux only) Update postgres password:
```
sudo passwd postgres
```
Enter a password, and be sure you won't forget it. You can just use something easy like postgres.

Start the Postgres server in the background:

Mac: `brew services start postgresql@15`\
Linux: `sudo service postgresql start`

## Configuration

Create a configuration file at `~/.gatorconfig.json` in your home directory with the following structure:
```json
{
  "db_url": "postgres://username:password@localhost:5432/database_name?sslmode=disable",
  "current_user_name": "your_username"
}
```

### Configuration Parameters:

- **db_url**: Your PostgreSQL connection string in the format:
```
  postgres://[username]:[password]@[host]:[port]/[database]?sslmode=disable
```
- **current_user_name**: Your default username (can be changed later)

### Example Configuration:
```json
{
  "db_url": "postgres://postgres:mypassword@localhost:5432/gator?sslmode=disable",
  "current_user_name": "john"
}
```

## commands

> ⚠️ **Important:** it is not a finished project more commands will be added

### User Management

**Register a new user:**
```bash
gator register <username>
```

**Login as an existing user:**
```bash
gator login <username>
```

**List all registered users:**
```bash
gator users
```
Shows all registered users and highlights the currently logged-in user.

### Feed Management

**Add a new feed:**
```bash
gator addfeed <feed_name> <feed_url>
```
Example: `gator addfeed "Tech News" https://example.com/feed.rss`

**List all feeds:**
```bash
gator feeds
```
Displays all feeds that have been added to the system.

**Follow a feed:**
```bash
gator follow <feed_url>
```

**Unfollow a feed:**
```bash
gator unfollow <feed_url>
```

**View feeds you're following:**
```bash
gator following
```
Shows all feeds that the current user is following.

### Content Aggregation

**Start the aggregator:**
```bash
gator agg <time_interval>
```

The `agg` command starts a background process that continuously fetches and updates RSS feeds at the specified time interval. 

**Time interval format:**
- `30s` - every 30 seconds
- `5m` - every 5 minutes
- `1h` - every 1 hour

**How it works:**
- Fetches the oldest feed in the queue based on last fetch time
- Marks the feed as fetched with the current timestamp
- Parses the RSS feed and extracts all posts/items
- Saves new posts to the database with title, URL, description, and publication date
- Repeats the process at the specified interval

This ensures your feeds are always up-to-date and posts are continuously aggregated into your local database.

**Browse saved posts:**
```bash
gator browse
```
Displays posts that have been saved to the database from your followed feeds.

### Database Management

**Reset the database:**
```bash
gator reset
```
⚠️ **Warning:** This command deletes all data from the database. Useful for testing purposes.

## Usage Example
```bash
# Register a new user
gator register alice

# Add some feeds
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator addfeed "Go Blog" https://blog.golang.org/feed.atom

# Follow feeds
gator follow https://news.ycombinator.com/rss

# Start aggregating feeds every 1 minute
gator agg 1m

# Browse collected posts
gator browse
```
