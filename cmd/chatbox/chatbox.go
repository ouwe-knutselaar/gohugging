package main

import (
	"bufio"
	"flag"
	"fmt"
	"gohugging/pkg/gohugging"
	"log"
	"os"
)

func main() {

	debug := flag.Bool("d", false, "enable debugging")
	flag.Parse()

	// Load configuration
	configData := GetConfig()

	gh, err := gohugging.New(configData)
	if err != nil {
		log.Fatalf("Failed to create GoHugging instance: %v", err)
	}
	if debug != nil && *debug {
		gh.EnableDebugging()
	}

	loop := true
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type your message (type 'exit' to quit):")
	for loop {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if handleCommand(input, gh) {
			continue
		}

		response, err := gh.SendChatMessage(input)
		if err != nil {
			log.Printf("Failed to send chat message: %v", err)
			continue
		}

		fmt.Println("Response:", response)
		fmt.Println()
	}
	fmt.Println("Exiting chat.")

}

func GetConfig() []byte {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configData, err := os.ReadFile(fmt.Sprintf("%s%c.clai.%c%s", home, os.PathSeparator, os.PathSeparator, "huggingface.yaml"))
	if err != nil {
		panic(err)
	}

	return configData
}

func handleCommand(input string, gh *gohugging.GoHugging) bool {
	switch input {
	case "/exit":
		fmt.Println("Exiting chat.")
		os.Exit(0)
	case "/history":
		if len(gh.Context) == 0 {
			fmt.Println("No chat history.")
		} else {
			for i, msg := range gh.Context {
				fmt.Printf("%d: [%s] %s\n", i, msg.Role, msg.Content)
			}
		}
		return true
	case "/clear":
		gh.Clear()
		fmt.Println("Chat history cleared.")
		return true
	case "/help":
		fmt.Println(`
Available commands:
/exit    - Exit the chat
/history - Show chat history
/clear   - Clear chat history
/help    - Show this help message`)
		return true
	}
	return false
}
