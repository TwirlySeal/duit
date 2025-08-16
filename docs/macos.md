# macOS setup instructions
1. Install Homebrew from https://brew.sh/. Homebrew is a package manager for macOS that automates installing programs and manages your PATH for you.

2. Install Git, Go, and PostgreSQL.
You can omit any software you already have installed.
```zsh
brew install git go postgresql@17
```

3. Start Postgres.
```zsh
LC_ALL="C" postgres -D /opt/homebrew/var/postgresql@17
```
- `LC_ALL="C"` sets the locale environment variable for the current shell to one that works well with Postgres.
- The `-D` option specifies the data directory (where your Postgres databases are stored).

4. Optional: You might like to make a zsh alias for the Postgres start command. You can do this by opening `~/.zshrc` in a text editor and adding the following line.
```
alias start_postgres='LC_ALL="C" postgres -D /opt/homebrew/var/postgresql@17'
```

5. Clone and enter the Duit Git repository.
```zsh
git clone https://github.com/TwirlySeal/duit.git
cd duit
```

6. Run `setup.sql`. This creates a database called 'duit' and inserts some test data.
```zsh
psql postgres -f setup.sql
```
- `psql` is a command-line interface for sending queries to Postgres.
- The `-f` option specifies a file to read queries from.

7. Enter the 'server' directory and run the Go server.
```zsh
cd server
go run .
```

8. Visit `http://localhost:8080` in your browser to access the web app.
