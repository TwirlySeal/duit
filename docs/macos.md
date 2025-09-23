# macOS first setup
1. Install [Homebrew](https://brew.sh/).

----

2. Install Git, Go, PostgreSQL, and Caddy. You can omit any software you already have installed.
```zsh
brew install git go postgresql@17 caddy
```

----

3. Add the following lines to `~/.zshrc` and restart your shell.
```zsh
export PATH="/opt/homebrew/opt/postgresql@17/bin:$PATH"
export PGDATA="/opt/homebrew/var/postgresql@17"
alias pg_start='LC_ALL="C" postgres'
```

----

4. Start Postgres.
```zsh
pg_start
```

----

5. In a new terminal, clone and enter the Duit Git repository.
```zsh
git clone https://github.com/TwirlySeal/duit.git
cd duit
```

----

6. Run `setup.sql` to create the 'duit' database and insert mock data.
```zsh
psql postgres -f setup.sql
```

----

7. Make a script to set the `POSTGRES_URL` environment variable.

a. Create a directory called 'scripts' in the root of the project (`scripts/` is in `.gitignore`)

b. Add a file inside called 'postgres_url' with the following content:
```zsh
export POSTGRES_URL="postgres://localhost/duit"
```
c. Make the file executable.
```zsh
chmod +x scripts/postgres_url
```
d. Run the script.
```zsh
source scripts/postgres_url
```

----

8. Enter the 'server' directory and run the Go server.
```zsh
cd server
go run .
```

----

9. Start Caddy from the root directory in a new terminal.
```zsh
caddy run
```

----

10. Open `http://localhost/1` in a browser to access the web app.

# Running the project
After first setup, these are the only steps needed to run the project.

1. Start Postgres
2. Run the `postgres_url` script
3. Start the Go server
4. Start Caddy
5. Open the web app in a browser
