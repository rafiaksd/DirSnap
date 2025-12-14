package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}

type Snapshot struct {
	Files map[string]FileInfo `json:"files"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:")
		fmt.Println("  snap <directory> <snapshot.json>")
		fmt.Println("  diff <directory> <snapshot.json>")
		os.Exit(1)
	}

	command := os.Args[1]
	dir := os.Args[2]
	snapFile := os.Args[3]

	switch command {
	case "snap":
		snapshot, err := createSnapshot(dir)
		if err != nil {
			panic(err)
		}
		saveSnapshot(snapshot, snapFile)
		fmt.Println("Snapshot saved to", snapFile)

	case "diff":
		oldSnap, err := loadSnapshot(snapFile)
		if err != nil {
			panic(err)
		}
		newSnap, err := createSnapshot(dir)
		if err != nil {
			panic(err)
		}
		diffSnapshots(oldSnap, newSnap)

	default:
		fmt.Println("Unknown command:", command)
	}
}

func createSnapshot(root string) (*Snapshot, error) {
	snap := &Snapshot{
		Files: make(map[string]FileInfo),
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(root, path)
		snap.Files[rel] = FileInfo{
			Path: rel,
			Size: info.Size(),
			Hash: hash,
		}
		return nil
	})

	return snap, err
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func saveSnapshot(snap *Snapshot, filename string) {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

func loadSnapshot(filename string) (*Snapshot, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	err = json.Unmarshal(data, &snap)
	return &snap, err
}

func diffSnapshots(oldSnap, newSnap *Snapshot) {
	seen := make(map[string]bool)

	fmt.Println("🔍 Differences:\n")

	for path, oldFile := range oldSnap.Files {
		newFile, exists := newSnap.Files[path]
		seen[path] = true

		if !exists {
			fmt.Println("❌ Deleted:", path)
			continue
		}

		if oldFile.Hash != newFile.Hash {
			fmt.Println("✏️ Modified:", path)
		}
	}

	for path := range newSnap.Files {
		if !seen[path] {
			fmt.Println("🆕 Added:", path)
		}
	}
}
