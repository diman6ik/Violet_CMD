package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// Structure for configuration file
type Config struct {
	HttpPost         string `json:"HttpPost"`
	Limit            int    `json:"Limit"`
	AnimationStyle   string `json:"AnimationStyle"`
	UserName         string `json:"UserName"`
	InterlocutorName string `json:"InterlocutorName"`
	Protocol         string `json:"Protocol"` // New parameter
}

// Structure for the AI request
type Prompt struct {
	Prompt    string   `json:"prompt"`
	Stopwords []string `json:"stop"`
	Limit     int      `json:"n_predict"`
	Cache     bool     `json:"cache_prompt"`
}

//go:embed config.json
var embeddedConfig []byte

//go:embed prompt.txt
var embeddedPrompt []byte

var history []string // Request history

// Function to load configuration from embedded file
func loadConfig() (*Config, error) {
	var config Config
	err := json.Unmarshal(embeddedConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded configuration file: %v", err)
	}
	return &config, nil
}

// Function to load text from embedded file
func loadPrompt() (string, error) {
	return string(embeddedPrompt), nil
}

// Function to clear the console
func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Configuration load error:", err)
		return
	}

	// Check Protocol value
	if config.Protocol != "http" && config.Protocol != "https" {
		fmt.Println("Error: Protocol must be 'http' or 'https'")
		return
	}

	// Load prompt text
	promptText, err := loadPrompt()
	if err != nil {
		fmt.Println("Prompt text load error:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Chat with %s. Type 'exit' to stop. (--help)\n", config.InterlocutorName)

	for {
		fmt.Printf("%s: ", config.UserName)
		scanner.Scan()
		input := scanner.Text()

		// Process commands with HandleCommand
		exit, isCommand := HandleCommand(input, history)
		if exit {
			break
		}

		// Skip sending to AI if it’s a command
		if isCommand {
			continue
		}

		// Add request to history
		history = append(history, input)

		// Create AI request
		prompt := &Prompt{
			Prompt:    fmt.Sprintf("%s\n%s: %s\n%s:", promptText, config.UserName, input, config.InterlocutorName),
			Stopwords: []string{fmt.Sprintf("\n%s:", config.UserName)},
			Limit:     config.Limit,
			Cache:     true,
		}

		var wg sync.WaitGroup
		wg.Add(1)
		stopChan := make(chan struct{})

		// Get the desired animation via GetAnimation
		animation := GetAnimation(config.AnimationStyle)
		go animation.Start(&wg, stopChan)

		// Substitute protocol in URL
		url := fmt.Sprintf("%s://%s/completion", config.Protocol, config.HttpPost)

		// Send POST request to AI server
		jsonPrompt, _ := json.Marshal(prompt)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPrompt))
		if err != nil {
			fmt.Println("Error:", err)
			close(stopChan)
			wg.Wait()
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		reply, ok := response["content"].(string)
		if ok {
			close(stopChan)
			wg.Wait()

			fmt.Print("\r                                                                                   \r")
			fmt.Printf("%s: %s\n", config.InterlocutorName, strings.Trim(reply, " \n\t"))
		} else {
			fmt.Println("Error: no response from AI")
		}
	}
}
