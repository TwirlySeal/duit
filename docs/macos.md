# macOS first setup
1. Install Homebrew from https://brew.sh/. Homebrew is a package manager for macOS that automates installing programs and manages your PATH for you.

----

2. Install Git, Go, and PostgreSQL.
You can omit any software you already have installed.
```zsh
brew install git go postgresql@17
```

----

3. Start Postgres.
```zsh
LC_ALL="C" postgres -D /opt/homebrew/var/postgresql@17
```
- `LC_ALL="C"` sets the locale environment variable for the current shell to one that works well with Postgres.
- The `-D` option specifies the data directory (where your Postgres databases are stored).

Optional: Make a simpler zsh alias for this command. You can do this by opening `~/.zshrc` in a text editor and adding the following line. Now you can just type `start_postgres`.
```
alias start_postgres='LC_ALL="C" postgres -D /opt/homebrew/var/postgresql@17'
```

----

4. Clone and enter the Duit Git repository.
```zsh
git clone https://github.com/TwirlySeal/duit.git
cd duit
```

----

5. Run `setup.sql`. This creates a database called 'duit' and inserts some test data.
```zsh
psql postgres -f setup.sql
```
- `psql` is a command-line interface for sending queries to Postgres.
- The `-f` option specifies a file to read queries from.

----

6. Set the `POSTGRES_URL` environment variable. This is used because different platforms and setups will need different URLs to connect to the database.
```zsh
export POSTGRES_URL="postgres://localhost/duit"
```

Optional: Save this command as a script in the `scripts` directory (which has been added to `.gitignore`). Write it into a file called `postgres_url` then run it like this:
```zsh
source postgres_url
```

----

7. Enter the 'server' directory and run the Go server.
```zsh
cd server
go run .
```

----

8. Visit `http://localhost:8080` in your browser to access the web app.

# Running the project
After first setup, these are the only steps needed to run the project.
1. Start Postgres
2. Set the `POSTGRES_URL` environment variable
3. Run the Go server
4. Visit the web app in your browser.
