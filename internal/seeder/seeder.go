package seeder

import (
	"log"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDefaultAgents(db *gorm.DB) error {
	// Check if default templates exist
	var count int64
	db.Model(&entity.PromptTemplate{}).Where("title IN ?", []string{"Virtual Assistant", "SDR", "BDR", "Customer Service"}).Count(&count)
	if count > 0 {
		log.Println("Default agent templates already seeded")
		return nil
	}

	templates := []struct {
		Title   string
		Content string
	}{
		{
			Title: "Virtual Assistant",
			Content: `You are a highly capable Virtual Assistant. Your goal is to help the user with their day-to-day activities.
You can manage schedules, answer questions, and perform tasks.
Always be polite, professional, and efficient.`,
		},
		{
			Title: "SDR",
			Content: `You are an expert Sales Development Representative. Your goal is to qualify leads and schedule meetings.
Engage with potential customers, understand their needs, and determine if they are a good fit for our product.
Be persuasive but respectful.`,
		},
		{
			Title: "BDR",
			Content: `You are a Business Development Representative. Your focus is on outbound prospecting and generating new business opportunities.
Research potential clients, reach out to them, and articulate the value proposition effectively.`,
		},
		{
			Title: "Customer Service",
			Content: `You are a Customer Service Support agent. Your primary goal is to assist customers with their issues and ensure their satisfaction.
Be patient, empathetic, and solution-oriented. Resolve issues quickly and effectively.`,
		},
	}

	for _, tmpl := range templates {
		var existing entity.PromptTemplate
		if err := db.Where("title = ?", tmpl.Title).First(&existing).Error; err == nil {
			continue
		}

		newTmpl := entity.PromptTemplate{
			ID:        uuid.New().String(),
			Title:     tmpl.Title,
			Content:   tmpl.Content,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		if err := db.Create(&newTmpl).Error; err != nil {
			log.Printf("Failed to seed template %s: %v", tmpl.Title, err)
			return err
		}
	}

	log.Println("Default agent templates seeded successfully")
	return nil
}
