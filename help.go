package main

import (
	"fmt"
	"os"
	"strings"
)

// Program version
const Version = "1.3.4"

// Function to display help and list of commands
func ShowHelp() {
	readme, err := os.ReadFile("readme.txt")
	if err != nil {
		fmt.Println("Error reading readme.txt file:", err)
	} else {
		fmt.Println(string(readme))
	}
	fmt.Println("Available commands:")
	fmt.Println("  exit       - exit the program")
	fmt.Println("  clear/cls  - clear the console")
	fmt.Println("  history    - display command history")
	fmt.Println("  --help     - display help and information from readme.txt")
	fmt.Println("  --version  - display the program version")
}

// Function to display the program version
func ShowVersion() {
	fmt.Println("Program version:", Version)
}

// Function to handle commands
func HandleCommand(command string, history []string) (bool, bool) {
	switch strings.ToLower(command) {
	case "exit":
		return true, true // Exit the program, do not send to neural network
	case "clear", "cls":
		clearConsole()     // Clear the console
		return false, true // Do not send to neural network
	case "history":
		ShowHistory(history) // Show command history
		return false, true   // Do not send to neural network
	case "--help":
		ShowHelp()         // Display help
		return false, true // Do not send to neural network
	case "--version":
		ShowVersion()      // Display program version
		return false, true // Do not send to neural network
	default:
		// If the command starts with "--" and is unknown
		if strings.HasPrefix(command, "--") {
			fmt.Println("Unknown command. Enter '--help' for a list of commands.")
			return false, true // Do not send to neural network
		}
	}
	return false, false // Not a command, send to neural network
}

// Function to display command history
func ShowHistory(history []string) {
	fmt.Println("Command history:")
	for i, h := range history {
		fmt.Printf("%d: %s\n", i+1, h)
	}
}
