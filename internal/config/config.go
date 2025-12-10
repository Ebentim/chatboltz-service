package config

type Config struct {
	Port                     string `env:"PORT,default=8080"`
	OPENAI_API_KEY           string `env:"OPENAI_API_KEY,required"`
	GEMINI_API_KEY           string `env:"GEMINI_API_KEY,required"`
	GOOGLE_API_KEY           string `env:"GOOGLE_API_KEY,required"` // For Google AI/Vertex AI
	DATABASE_URL             string `env:"DATABASE_URL,required"`
	GOOGLE_CLIENT_SECRET     string `env:"GOOGLE_CLIENT_SECRET, required"`
	GOOGLE_CLIENT_ID         string `env:"GOOGLE_CLIENT_ID,required"`
	FIREBASE_SERVICE_ACCOUNT string `env:"FIREBASE_SERVICE_ACCOUNT,required"`
	GCM_KEY                  string `env:"GCM_KEY,required"`
	JWT_SECRET               string `env:"JWT_SECRET,required"`
	COHERE_API_KEY           string `env:"COHERE_API_KEY,required"`
	GROQ_API_KEY             string `env:"GROQ_API_KEY,required"`
	PINECONE_API_KEY         string `env:"PINECONE_API_KEY"`
	PINECONE_INDEX_NAME      string `env:"PINECONE_INDEX_NAME,default=agent-knowledge"`
	VECTOR_DB_TYPE           string `env:"VECTOR_DB_TYPE,default=pgvector"`
	SMTP_HOST                string `env:"SMTP_HOST,required"`
	SMTP_PORT                string `env:"SMTP_PORT,required"`
	SMTP_USER                string `env:"SMTP_USER,required"`
	SMTP_PASS                string `env:"SMTP_PASS,required"`
	OTP_SECRET               string `env:"OTP_SECRET,required"`
	ENABLE_ORCHESTRATION     bool   `env:"ENABLE_ORCHESTRATION,default=false"`
	// DispatcherDeliveryTimeoutMS controls how long the in-memory dispatcher will
	// wait for a subscriber to accept an event before considering it dropped.
	// Value is in milliseconds.
	DispatcherDeliveryTimeoutMS int `env:"DISPATCHER_DELIVERY_TIMEOUT_MS,default=100"`
}

// Vector DB Types
const (
	VectorDBPgVector = "pgvector"
	VectorDBPinecone = "pinecone"
)

var GoogleScope = []string{
	"openid", "email", "profile", "https://www.googleapis.com/auth/drive.file", "https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/documents", "https://www.googleapis.com/auth/presentations", "https://www.googleapis.com/auth/forms.body", "https://www.googleapis.com/auth/forms.responses.readonly", "https://www.googleapis.com/auth/classroom.courses", "https://www.googleapis.com/auth/classroom.rosters", "https://www.googleapis.com/auth/classroom.coursework.students", "https://www.googleapis.com/auth/classroom.coursework.me", "https://www.googleapis.com/auth/classroom.announcements", "https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/calendar.events",
}
