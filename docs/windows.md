# Windows first setup
1. Install Scoop from https://scoop.sh/.
- Scoop is a package manager for Windows that automates installing programs and manages your PATH for you.

----

2. Install Git, Go, and PostgreSQL.
You can omit any software you already have installed.
```powershell
scoop install git go postgresql
```

----

3. Start Postgres.
```powershell
postgres -D $env:PGDATA
```
- The `-D` option specifies the data directory (where your Postgres databases are stored).
- `PGDATA` is an environment variable created by Scoop that holds your data directory path.

----

4. Clone and enter the Duit Git repository.
```powershell
git clone https://github.com/TwirlySeal/duit.git
cd .\duit
```

----

5. Run `setup.sql`. This creates a database called 'duit' and inserts some test data.
```powershell
psql -U postgres -f .\setup.sql
```
- `psql` is a command-line interface for sending queries to Postgres.
- `-U postgres` connects as the default user `postgres`, which was created by Scoop.
- The `-f` option specifies a file to read queries from.

----

6. Set the `POSTGRES_URL` environment variable. This is used because different platforms and setups will need different URLs to connect to the database.
```powershell
$env:POSTGRES_URL = "postgres://postgres@localhost/duit"
```

Optional: Save this step as a script.
  1. Create a `scripts` directory in the root of the project (it is in `.gitignore`).
  2. Inside `scripts`, write the command to a file called `postgres_url.ps1`.
  3. Run the script like this:
```powershell
. .\scripts\postgres_url.ps1
```

----

7. Enter the 'server' directory and run the Go server.
```powershell
cd .\server
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
