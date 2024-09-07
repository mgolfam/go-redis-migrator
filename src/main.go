package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	rhelper "github.com/integrii/go-redis-migrator/src/rhelper"
)

type CommandHandler interface {
	Execute()
}

// MigrateCommand implements the CommandHandler interface for migrating Redis keys
type MigrateCommand struct {
	sourceHosts         []string
	destinationHosts    []string
	sourcePassword      string
	destinationPassword string
	keyFilter           string
	keyFile             string
}

func (cmd *MigrateCommand) Execute() {
	sourceHandler := rhelper.NewRedisHandler(cmd.sourceHosts, cmd.sourcePassword)
	destHandler := rhelper.NewRedisHandler(cmd.destinationHosts, cmd.destinationPassword)

	// Migrate all databases from source to destination
	sourceHandler.MigrateAllDatabases(destHandler)
}

// HelpCommand implements the CommandHandler interface for displaying help information
type HelpCommand struct{}

func (cmd *HelpCommand) Execute() {
	fmt.Println("Usage: redis-migrator [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  migrate   Migrate keys between Redis instances")
	fmt.Println("  help      Show this help message")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func main() {
	// Define common flags
	sourceHosts := flag.String("sourceHosts", "", "Source Redis hosts (comma-separated)")
	destinationHosts := flag.String("destinationHosts", "", "Destination Redis hosts (comma-separated)")
	sourcePassword := flag.String("sourcePassword", "", "Password for the source Redis instance (optional)")
	destinationPassword := flag.String("destinationPassword", "", "Password for the destination Redis instance (optional)")
	keyFilter := flag.String("keyFilter", "*", "Filter for Redis keys to migrate")
	keyFile := flag.String("keyFile", "", "File containing list of keys to migrate (optional)")

	// Parse command-line arguments
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Error: command is required")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "migrate":
		migrateCmd := &MigrateCommand{
			sourceHosts:         strings.Split(*sourceHosts, ","),
			destinationHosts:    strings.Split(*destinationHosts, ","),
			sourcePassword:      *sourcePassword,
			destinationPassword: *destinationPassword,
			keyFilter:           *keyFilter,
			keyFile:             *keyFile,
		}
		migrateCmd.Execute()
	case "help":
		helpCmd := &HelpCommand{}
		helpCmd.Execute()
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
