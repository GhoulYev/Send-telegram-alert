package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
	Prefix  string `json:"prefix"`
}

func main() {
	message := flag.String("m", "", "Message to send")
	configPath := flag.String("c", "", "Path to configuration file")
	flag.Parse()

	if *message == "" {
		fmt.Println("Error: Message not specified. Use -m 'message text'")
		os.Exit(1)
	}

	finalConfigPath := *configPath
	if finalConfigPath == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		finalConfigPath = filepath.Join(usr.HomeDir, ".config", "sta", "config.json")
	}

	if _, err := os.Stat(finalConfigPath); os.IsNotExist(err) {
		if *configPath == "" {
			if err := createDefaultConfig(finalConfigPath); err != nil {
				fmt.Printf("Error creating config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Created template config: %s\nPlease fill it with your data!\n", finalConfigPath)
			os.Exit(0)
		} else {
			fmt.Printf("Error: Config file not found at: %s\n", finalConfigPath)
			os.Exit(1)
		}
	}

	config, err := loadConfig(finalConfigPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if config.Token == "" || config.ChatID == "" {
		fmt.Println("Error: Token or Chat ID are missing in config")
		os.Exit(1)
	}

	fullMessage := config.Prefix + " " + *message

	err = sendTelegram(config.Token, config.ChatID, fullMessage)
	if err != nil {
		fmt.Printf("Sending error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Message sent successfully!")
}

func createDefaultConfig(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	config := Config{
		Token:  "YOUR_BOT_TOKEN",
		ChatID: "YOUR_CHAT_ID",
		Prefix: "[Arch Desktop]",
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}

func loadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}

func sendTelegram(token, chatID, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.PostForm(
		apiURL,
		url.Values{
			"chat_id": {chatID},
			"text":    {text},
		},
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", body)
	}

	return nil
}