package cmd

import (
	"fmt"
	"os"
)

type Command interface {
	Execute(args []string)
}

// Help Command
type HelpCommand struct{}

func (c *HelpCommand) Execute(args []string) {
	fmt.Println(`
Redis Key Migrator

Commands:
  help          - Show this help message
  migrate       - Migrate keys from one Redis instance/cluster to another
  ... other commands
`)
	os.Exit(0)
}

// Command Factory
func GetCommand(commandName string) Command {
	switch commandName {
	case "help":
		return &HelpCommand{}
	case "migrate":
	default:
		fmt.Println("Unknown command. Type 'help' for usage.")
		os.Exit(1)
	}
	return nil
}
