package usecase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	client *auth.Client
}

// func NewFirebaseService(projectID, credentialsJSON string) (*FirebaseService, error) {
func NewFirebaseService(credentialsJSON string) (*FirebaseService, error) {
	ctx := context.Background()

	opt := option.WithCredentialsJSON([]byte(credentialsJSON))
	// config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseService{client: client}, nil
}

func (s *FirebaseService) VerifyIDToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()
	return s.client.VerifyIDToken(ctx, idToken)
}

func (s *FirebaseService) CreateUser(email, password string) (*auth.UserRecord, error) {
	ctx := context.Background()
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)

	return s.client.CreateUser(ctx, params)
}

func (s *FirebaseService) GetUserByEmail(email string) (*auth.UserRecord, error) {
	ctx := context.Background()
	return s.client.GetUserByEmail(ctx, email)
}

func (s *FirebaseService) CreateCustomToken(uid string, claims map[string]interface{}) (string, error) {
	ctx := context.Background()

	return s.client.CustomToken(ctx, uid)
}
