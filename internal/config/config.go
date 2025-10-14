package config

type Config struct {
	Port                     string `env:"PORT,default=8080"`
	OpenAI_API_KEY           string `env:"OPENAI_API_KEY,required"`
	GEMINI_API_KEY           string `env:"GEMINI_API_KEY,required"`
	DATABASE_URL             string `env:"DATABASE_URL,required"`
	GOOGLE_CLIENT_SECRET     string `env:"GOOGLE_CLIENT_SECRET, required"`
	GOOGLE_CLIENT_ID         string `env:"GOOGLE_CLIENT_ID,required"`
	FIREBASE_SERVICE_ACCOUNT string `env:"FIREBASE_SERVICE_ACCOUNT,required"`
	FIREBASE_PROJECT_ID      string `env:"FIREBASE_PROJECT_ID,required"`
	FIREBASE_CREDENTIALS     string `env:"FIREBASE_CREDENTIALS,required"`
}

// DATABASE CONFIG SHOULD BE A MAP OF DIFFERENT DATABASES, POSTGRESQL PGVECTOR MONGODB AND MONGODB VECTOR

var GoogleScope = []string{
	"openid", "email", "profile", "https://www.googleapis.com/auth/drive.file", "https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/documents", "https://www.googleapis.com/auth/presentations", "https://www.googleapis.com/auth/forms.body", "https://www.googleapis.com/auth/forms.responses.readonly", "https://www.googleapis.com/auth/classroom.courses", "https://www.googleapis.com/auth/classroom.rosters", "https://www.googleapis.com/auth/classroom.coursework.students", "https://www.googleapis.com/auth/classroom.coursework.me", "https://www.googleapis.com/auth/classroom.announcements", "https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events",
}
