package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Klinik PKP API")
	fmt.Println("==============")
	fmt.Println()
	fmt.Println("SERVER")
	fmt.Println("> go run cmd/server/main.go           - Start the API server")
	fmt.Println()
	fmt.Println("MIGRATION")
	fmt.Println("> go run cmd/migrate/main.go          - Run all migrations")
	fmt.Println("> go run cmd/migrate/main.go user     - Run specific migration")
	fmt.Println("> go run cmd/migrate/main.go --help   - Show migration help")
	fmt.Println()
	fmt.Println("DATABASE SEEDING")
	fmt.Println("> go run cmd/seed/main.go             - Run all seeders")
	fmt.Println("> go run cmd/seed/main.go province    - Run specific seeder")
	fmt.Println("> go run cmd/seed/main.go --help      - Show seeder help")
	fmt.Println()
	fmt.Println("HOT RELOADING (FOR DEVELOPMENT, ADJUST air.toml AS NEEDED)")
	fmt.Println("> air")

	os.Exit(0)
}
