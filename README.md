# DirSnap

DirSnap is a lightweight Go CLI that **snapshots a directory’s file state** and later **detects changes** using content hashing — not timestamps.

## Why it’s useful
- Detect unexpected file changes
- Verify build or deployment outputs
- Audit folders for integrity
- Lightweight alternative to full backup/versioning tools

## Usage
```bash
# Create a snapshot
go run main.go snap ./mydir snapshot.json

# Compare directory against a snapshot
go run main.go diff ./mydir snapshot.json
