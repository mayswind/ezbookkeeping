# ezBookkeeping dev environment setup (Git Bash)

## Build prerequisites (Windows + Git Bash)

The project uses `github.com/mattn/go-sqlite3` which requires CGO.
On Windows you need an MSYS2 GCC on PATH **for the spawned cgo child process**.

### One-time setup
1. Install MSYS2 from https://www.msys2.org/ (if not already).
2. In MSYS2 UCRT64 shell: `pacman -S mingw-w64-ucrt-x86_64-gcc`
3. Before **every** Go build in Git Bash, run:
   ```bash
   export PATH="/c/msys64/ucrt64/bin:$PATH"
   export CGO_ENABLED=1
   export CC=gcc
   ```
   (Or source this file: `source DEV_SETUP.md`)

## Build commands
```bash
# Backend (Go)
go build -o ezbookkeeping.exe .
go test ./...

# Frontend (Vue)
npm install
npm run dev      # dev server with HMR
npm run build    # production build
```

## Run
```bash
cp conf/ezbookkeeping.ini.example conf/ezbookkeeping.ini  # edit DB path if needed
./ezbookkeeping.exe serve
```
