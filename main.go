package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

type variableFlag map[string]string

func (v variableFlag) String() string {
	var pairs []string
	for key, value := range v {
		pairs = append(pairs, fmt.Sprintf("%s:%s", key, value))
	}
	return strings.Join(pairs, ", ")
}

func (v variableFlag) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("variable must be in format 'key:value', got: %s", value)
	}
	v[parts[0]] = parts[1]
	return nil
}

func main() {
	var dir string
	var conn string
	variables := make(variableFlag)

	flag.StringVar(&dir, "dir", "", "Directory containing SQL files")
	flag.StringVar(&conn, "conn", "", "MSSQL connection string")
	flag.Var(variables, "variable", "Variable replacement in format 'old:new' (can be used multiple times)")
	flag.Parse()

	if dir == "" {
		log.Fatal("Directory parameter is required. Use --dir flag")
	}
	if conn == "" {
		log.Fatal("Connection string parameter is required. Use --conn flag")
	}

	db, err := sql.Open("mssql", conn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Printf("Connected to database successfully\n")
	fmt.Printf("Processing SQL files in directory: %s\n", dir)

	if len(variables) > 0 {
		fmt.Printf("Variable replacements:\n")
		for old, new := range variables {
			fmt.Printf("  %s -> %s\n", old, new)
		}
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".sql" {
				fmt.Printf("Executing: %s\n", path)
				if err := executeSQLFile(db, path, variables); err != nil {
					fmt.Printf("❌ Error executing %s: %v\n", path, err)
				} else {
					fmt.Printf("✅ Successfully executed: %s\n", path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	fmt.Println("All SQL files processed successfully")
}

func executeSQLFile(db *sql.DB, filePath string, variables variableFlag) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	sqlContent := string(content)
	if strings.TrimSpace(sqlContent) == "" {
		return nil
	}

	for old, new := range variables {
		sqlContent = strings.ReplaceAll(sqlContent, old, new)
	}

	_, err = db.Exec(sqlContent)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %v", err)
	}

	return nil
}
