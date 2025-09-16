# Windows first setup
1. Install Powershell from WinGet. This is distinct from [Windows Powershell](https://learn.microsoft.com/en-us/powershell/scripting/what-is-windows-powershell) which is included with Windows.
```powershell
winget install Microsoft.PowerShell
```

----

2. Install [Scoop](https://scoop.sh).

----

3. Install Git, Go, and PostgreSQL. You can omit any software you already have installed.
```powershell
scoop install git go postgresql
```

----

4. Start Postgres.
```powershell
postgres
```

----

5. In a new terminal, clone and enter the Duit Git repository.
```powershell
git clone https://github.com/TwirlySeal/duit.git
cd .\duit
```

----

6. Run `setup.sql` to create the `duit` database and insert mock data.
```powershell
psql -U postgres -f .\setup.sql
```

----

7. Make a script to set the `POSTGRES_URL` environment variable.

a. Create a directory called 'scripts' in the root of the project (`scripts/` is in `.gitignore`)

b. Add a file inside called 'postgres_url.ps1' with the following content:
```powershell
$env:POSTGRES_URL = "postgres://postgres@localhost/duit"
```
c. Run the script.
```powershell
. .\scripts\postgres_url.ps1
```

----

8. Enter the 'server' directory and run the Go server.
```powershell
cd .\server
go run .
```

----

9. Open `http://localhost:8080` in a browser to access the web app.

# Running the project
After first setup, these are the only steps needed to run the project.

1. Start Postgres
2. Run the `postgres_url.ps1` script
3. Run the Go server
4. Open the web app in a browser
