package main

/*
func main() {
	// Get API key from environment
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable is required")
	}

	// Example 1: Direct provider creation
	fmt.Println("=== Direct Groq Provider Usage ===")
	provider, err := aiprovider.NewGroqAIClient(apiKey, "https://api.groq.com/openai/v1")
	if err != nil {
		log.Fatal(err)
	}

	conversation := aiprovider.Conversation{
		Messages: []aiprovider.Message{
			{Role: aiprovider.RoleSystem, Content: "You are a helpful assistant."},
			{Role: aiprovider.RoleUser, Content: "What is the capital of France?"},
		},
	}

	config := map[string]interface{}{
		"model":       "llama-3.3-70b-versatile",
		"temperature": 0.7,
		"max_tokens":  100,
	}

	response, err := provider.CompleteConversation(conversation, config)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response: %s\n\n", response)
	}

	// Example 2: Factory-based creation
	fmt.Println("=== Factory-based Provider Usage ===")
	factory := aiprovider.NewProviderFactory()

	groqProvider, err := factory.CreateProvider(entity.ProviderConfig{
		Provider: entity.Groq,
		APIKey:   apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Test capabilities
	caps := groqProvider.GetCapabilities()
	fmt.Printf("Capabilities: Text=%t, Voice=%t, Vision=%t\n", caps.Text, caps.Voice, caps.Vision)

	// Example 3: Agent-based usage (production pattern)
	fmt.Println("=== Agent-based Usage ===")
	agent := entity.Agent{
		AiProvider: "groq",
		AiModel:    "llama-3.1-8b-instant", // Faster model
	}

	agentProvider, err := factory.GetProviderFromAgent(agent, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	quickConversation := aiprovider.Conversation{
		Messages: []aiprovider.Message{
			{Role: aiprovider.RoleUser, Content: "Write a haiku about AI"},
		},
	}

	quickConfig := map[string]interface{}{
		"model":       agent.AiModel,
		"temperature": 0.9,
	}

	quickResponse, err := agentProvider.CompleteConversation(quickConversation, quickConfig)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Haiku Response:\n%s\n\n", quickResponse)
	}

	// Example 4: Streaming response
	fmt.Println("=== Streaming Response ===")
	fmt.Print("Streaming response: ")

	streamConfig := map[string]interface{}{
		"model":       "llama-3.1-8b-instant",
		"temperature": 0.8,
	}

	streamConversation := aiprovider.Conversation{
		Messages: []aiprovider.Message{
			{Role: aiprovider.RoleUser, Content: "Count from 1 to 5 with explanations"},
		},
	}

	err = agentProvider.CompleteConversationStream(streamConversation, streamConfig, func(chunk string, done bool) error {
		if done {
			fmt.Println("\n[Stream complete]")
			return nil
		}
		fmt.Print(chunk)
		return nil
	})

	if err != nil {
		log.Printf("Streaming error: %v", err)
	}

	fmt.Println("\n=== Examples Complete ===")
}
*/
