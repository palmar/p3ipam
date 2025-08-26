package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"p3ipam/db"
)

const (
	version = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	action := os.Args[1]
	args := os.Args[2:]

	switch action {
	case "init":
		handleInit()
	case "version":
		fmt.Printf("p3ipam v%s\n", version)
	case "help", "-h", "--help":
		showHelp()
	case "add":
		handleAdd(args)
	case "list":
		handleList(args)
	case "delete":
		handleDelete(args)
	case "edit":
		handleEdit(args)
	case "ping":
		handlePing(args)
	case "search":
		handleSearch(args)
	default:
		fmt.Printf("Unknown action: %s\n", action)
		fmt.Println("Use 'help' to see available actions")
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("p3ipam - Lightweight IP Address Management Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  p3ipam <action> <object> <target> [--arguments]")
	fmt.Println("")
	fmt.Println("Actions:")
	fmt.Println("  init                    - Initialize the database")
	fmt.Println("  version                 - Show version information")
	fmt.Println("  help                    - Show this help message")
	fmt.Println("  add <object>            - Add a new object")
	fmt.Println("  list <object>           - List objects")
	fmt.Println("  delete <object> <id>    - Delete an object")
	fmt.Println("  edit <object> <id>      - Edit an object")
	fmt.Println("  ping <object> <target>  - Ping and discover hosts")
	fmt.Println("  search <query>          - Search across all objects")
	fmt.Println("")
	fmt.Println("Objects:")
	fmt.Println("  subnet                  - Network subnet (e.g., 192.168.1.0/24)")
	fmt.Println("  host                    - Network host (e.g., 192.168.1.1)")
	fmt.Println("")
	fmt.Println("Parent References:")
	fmt.Println("  --parent accepts: subnet name, ID, or CIDR notation")
	fmt.Println("  Example: --parent home-network, --parent ABC123, or --parent 192.168.1.0/24")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  p3ipam add subnet --cidr 192.168.1.0/24 --name home-network")
	fmt.Println("  p3ipam add host --parent home-network --address 192.168.1.1 --name router")
	fmt.Println("  p3ipam add host --parent 192.168.1.0/24 --address 192.168.1.2 --name server")
	fmt.Println("  p3ipam search 192.168.1")
	fmt.Println("  p3ipam ping subnet home-network")
}

func handleInit() {
	fmt.Println("=== p3ipam Database Initialization ===")
	fmt.Println()

	// Get database location
	defaultPath := db.GetDatabasePath()
	fmt.Printf("Database location [%s]: ", defaultPath)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dbLocation := strings.TrimSpace(scanner.Text())

	if dbLocation == "" {
		dbLocation = defaultPath
	}

	// Get database password (for future use)
	fmt.Print("Database password (optional, press Enter for none): ")
	scanner.Scan()
	dbPassword := strings.TrimSpace(scanner.Text())

	fmt.Println()
	fmt.Printf("Initializing database at: %s\n", dbLocation)
	if dbPassword != "" {
		fmt.Println("Password will be stored (encryption not yet implemented)")
	}
	fmt.Println()

	// Show environment variable note if custom path
	if dbLocation != defaultPath {
		fmt.Println("💡 Note: To use this custom database location in the future,")
		fmt.Printf("   set the environment variable: P3IPAM_DATADIR=%s\n", filepath.Dir(dbLocation))
		fmt.Println()
	}

	database, err := db.Connect(dbLocation)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := database.Init(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Database initialized successfully!")
	fmt.Printf("📁 Database file: %s\n", dbLocation)
}

func handleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Object type required")
		fmt.Println("Usage: p3ipam add <object> [--arguments]")
		os.Exit(1)
	}

	objectType := args[0]
	objectArgs := args[1:]

	switch objectType {
	case "subnet":
		handleAddSubnet(objectArgs)
	case "host":
		handleAddHost(objectArgs)
	default:
		fmt.Printf("Unknown object type: %s\n", objectType)
		fmt.Println("Supported types: subnet, host")
		os.Exit(1)
	}
}

func handleAddSubnet(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Subnet arguments required")
		fmt.Println("Usage: p3ipam add subnet --cidr <cidr> [--name <name>] [--parent <parent_id>] [--comment <comment>]")
		os.Exit(1)
	}

	var cidr, name, parentID, comment string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--cidr":
			if i+1 < len(args) {
				cidr = args[i+1]
				i++
			}
		case "--name":
			if i+1 < len(args) {
				name = args[i+1]
				i++
			}
		case "--parent":
			if i+1 < len(args) {
				parentID = args[i+1]
				i++
			}
		case "--comment":
			if i+1 < len(args) {
				comment = args[i+1]
				i++
			}
		}
	}

	if cidr == "" {
		fmt.Println("Error: --cidr is required")
		fmt.Println("Usage: p3ipam add subnet --cidr <cidr> [--name <name>] [--parent <parent_id>] [--comment <comment>]")
		os.Exit(1)
	}

	// Connect to database
	database, err := db.Connect(db.GetDatabasePath())
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Add subnet to database
	subnet, err := database.AddSubnet(cidr, name, parentID, comment)
	if err != nil {
		fmt.Printf("Error adding subnet: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Subnet added successfully!\n")
	fmt.Printf("   ID: %s\n", subnet.ID)
	fmt.Printf("   CIDR: %s\n", subnet.CIDR)
	if subnet.Name != "" {
		fmt.Printf("   Name: %s\n", subnet.Name)
	}
	if subnet.Comment != "" {
		fmt.Printf("   Comment: %s\n", subnet.Comment)
	}
}

func handleAddHost(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Host arguments required")
		fmt.Println("Usage: p3ipam add host --address <address> [--name <name>] [--parent <parent_id>] [--comment <comment>]")
		os.Exit(1)
	}

	var address, name, parentID, comment string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--address":
			if i+1 < len(args) {
				address = args[i+1]
				i++
			}
		case "--name":
			if i+1 < len(args) {
				name = args[i+1]
				i++
			}
		case "--parent":
			if i+1 < len(args) {
				parentID = args[i+1]
				i++
			}
		case "--comment":
			if i+1 < len(args) {
				comment = args[i+1]
				i++
			}
		}
	}

	if address == "" {
		fmt.Println("Error: --address is required")
		fmt.Println("Usage: p3ipam add host --address <address> [--name <name>] [--parent <parent_id>] [--comment <comment>]")
		os.Exit(1)
	}

	// Connect to database
	database, err := db.Connect(db.GetDatabasePath())
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Add host to database
	host, err := database.AddHost(address, name, parentID, comment)
	if err != nil {
		fmt.Printf("Error adding host: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Host added successfully!\n")
	fmt.Printf("   ID: %s\n", host.ID)
	fmt.Printf("   Address: %s\n", host.Address)
	if host.Name != "" {
		fmt.Printf("   Name: %s\n", host.Name)
	}
	if host.Comment != "" {
		fmt.Printf("   Comment: %s\n", host.Comment)
	}
}

func handleList(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Object type required")
		fmt.Println("Usage: p3ipam list <object>")
		os.Exit(1)
	}

	objectType := args[0]
	switch objectType {
	case "subnets":
		handleListSubnets()
	case "hosts":
		handleListHosts()
	case "discoveries":
		handleListDiscoveries()
	default:
		fmt.Printf("Unknown object type: %s\n", objectType)
		fmt.Println("Supported types: subnets, hosts, discoveries")
		os.Exit(1)
	}
}

func handleListSubnets() {
	// TODO: Implement subnet listing
	fmt.Println("Listing subnets... (not yet implemented)")
}

func handleListHosts() {
	// TODO: Implement host listing
	fmt.Println("Listing hosts... (not yet implemented)")
}

func handleListDiscoveries() {
	// TODO: Implement discovery listing
	fmt.Println("Listing discoveries... (not yet implemented)")
}

func handleDelete(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: Object type and ID required")
		fmt.Println("Usage: p3ipam delete <object> <id>")
		os.Exit(1)
	}

	objectType := args[0]
	objectID := args[1]

	switch objectType {
	case "subnet":
		handleDeleteSubnet(objectID)
	case "host":
		handleDeleteHost(objectID)
	default:
		fmt.Printf("Unknown object type: %s\n", objectType)
		fmt.Println("Supported types: subnet, host")
		os.Exit(1)
	}
}

func handleDeleteSubnet(id string) {
	// TODO: Implement subnet deletion
	fmt.Printf("Deleting subnet %s... (not yet implemented)\n", id)
}

func handleDeleteHost(id string) {
	// TODO: Implement host deletion
	fmt.Printf("Deleting host %s... (not yet implemented)\n", id)
}

func handleEdit(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: Object type and ID required")
		fmt.Println("Usage: p3ipam edit <object> <id> [--arguments]")
		os.Exit(1)
	}

	objectType := args[0]
	objectID := args[1]
	editArgs := args[2:]

	switch objectType {
	case "subnet":
		handleEditSubnet(objectID, editArgs)
	case "host":
		handleEditHost(objectID, editArgs)
	default:
		fmt.Printf("Unknown object type: %s\n", objectType)
		fmt.Println("Supported types: subnet, host")
		os.Exit(1)
	}
}

func handleEditSubnet(id string, args []string) {
	// TODO: Implement subnet editing
	fmt.Printf("Editing subnet %s... (not yet implemented)\n", id)
}

func handleEditHost(id string, args []string) {
	// TODO: Implement host editing
	fmt.Printf("Editing host %s... (not yet implemented)\n", id)
}

func handlePing(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: Object type and target required")
		fmt.Println("Usage: p3ipam ping <object> <target>")
		os.Exit(1)
	}

	objectType := args[0]
	target := args[1]

	switch objectType {
	case "subnet":
		handlePingSubnet(target)
	default:
		fmt.Printf("Unknown object type: %s\n", objectType)
		fmt.Println("Supported types: subnet")
		os.Exit(1)
	}
}

func handlePingSubnet(target string) {
	// TODO: Implement subnet ping
	fmt.Printf("Pinging subnet %s... (not yet implemented)\n", target)
}

func handleSearch(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Search query required")
		fmt.Println("Usage: p3ipam search <query>")
		os.Exit(1)
	}

	query := args[0]

	database, err := db.Connect(db.GetDatabasePath())
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	results, err := database.Search(query)
	if err != nil {
		fmt.Printf("Error searching database: %v\n", err)
		os.Exit(1)
	}

	displaySearchResults(results)
}

func displaySearchResults(results *db.SearchResults) {
	fmt.Printf("Search Results:\n\n")

	if len(results.Subnets) > 0 {
		fmt.Println("Subnets:")
		for _, subnet := range results.Subnets {
			fmt.Printf("  %s (%s) - %s\n", subnet.CIDR, subnet.ID, subnet.Name)
			if subnet.Comment != "" {
				fmt.Printf("    Comment: %s\n", subnet.Comment)
			}
		}
		fmt.Println()
	}

	if len(results.Hosts) > 0 {
		fmt.Println("Hosts:")
		for _, host := range results.Hosts {
			fmt.Printf("  %s (%s) - %s\n", host.Address, host.ID, host.Name)
			if host.Comment != "" {
				fmt.Printf("    Comment: %s\n", host.Comment)
			}
		}
		fmt.Println()
	}

	if len(results.Discoveries) > 0 {
		fmt.Println("Discoveries:")
		for _, discovery := range results.Discoveries {
			fmt.Printf("  %s (%s) - Status: %s\n", discovery.Address, discovery.ID, discovery.Status)
		}
		fmt.Println()
	}

	if len(results.Subnets) == 0 && len(results.Hosts) == 0 && len(results.Discoveries) == 0 {
		fmt.Println("No results found.")
	}
}
