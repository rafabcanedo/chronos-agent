package mocks

import "context"

type MockChatService struct {
	HandleFn func(ctx context.Context, userID, message string) (string, error)
}

func (m *MockChatService) Handle(ctx context.Context, userID, message string) (string, error) {
	return m.HandleFn(ctx, userID, message)
}
