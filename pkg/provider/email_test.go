package provider

import (
	"testing"
)

func TestMockEmailProvider(t *testing.T) {
	provider := NewMockEmailProvider()

	for i := 1; i <= 5; i++ {
		err := provider.Send("mockemailfornow@mock.com","Test notification message")
		if err != nil {
			t.Logf("Attempt %d:  Failed with error: %v", i, err)
		} else {
			t.Logf("Attempt %d:  Sent successfully", i)
		}
	}
}