package provider

import (
	"errors"
	"math/rand"
	"time"
)

type EmailProvider interface {
	Send(email string, message string) error
}

type mockEmailProvider struct{}

func NewMockEmailProvider() EmailProvider {
	return &mockEmailProvider{}
}

func (m *mockEmailProvider) Send(email string,message string) error {
	time.Sleep(500 * time.Millisecond)

	if rand.Float32() < 0.3 {
		return errors.New("failed to connect to email provider server")
	}

	return nil
}