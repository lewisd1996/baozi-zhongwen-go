package util

import (
	"fmt"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     error
	}{
		{
			name:     "password is too short",
			password: "short",
			want:     fmt.Errorf("password must be at least 8 characters long"),
		},
		{
			name:     "password is all lowercase",
			password: "lowercase",
			want:     fmt.Errorf("password must contain at least one uppercase letter"),
		},
		{
			name:     "password is all uppercase",
			password: "UPPERCASE",
			want:     fmt.Errorf("password must contain at least one lowercase letter"),
		},
		{
			name:     "password contains no special characters",
			password: "Password123",
			want:     fmt.Errorf("password must contain at least one special character"),
		},
		{
			name:     "password contains no numbers",
			password: "Password!",
			want:     fmt.Errorf("password must contain at least one number"),
		},
		{
			name:     "password is valid",
			password: "ValidPassword123!",
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePassword(tt.password)
			if (got != nil && tt.want != nil && got.Error() != tt.want.Error()) || (got != nil && tt.want == nil) || (got == nil && tt.want != nil) {
				t.Errorf("ValidatePassword(%v) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}
