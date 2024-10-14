package auth

import "testing"

func TestJWTToken(t *testing.T) {
	tests := []struct {
		name      string
		userId    int
		expectErr bool
	}{
		{name: "valid case", userId: 1, expectErr: false},     // Valid case
		{name: "invalid user id", userId: 0, expectErr: true}, // Invalid userId
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := CreateToken(tt.userId)
			if err != nil {
				t.Errorf("got error creating token for userId %d: %v", tt.userId, err)
				return
			}

			if tokenString == "" {
				t.Errorf("expected token to not be empty for userId %d", tt.userId)
			}

			token, err := VerifyToken(tokenString)
			if err != nil {
				t.Errorf("error verifying token for userId %d: %v", tt.userId, err)
				return
			}

			resultUserId, err := GetUserIDFromToken(token)
			if err != nil {
				if !tt.expectErr {
					t.Errorf("error getting userId from token for userId %d: %v", tt.userId, err)
				}
				return
			}

			// Check if the userId matches the original
			if resultUserId != tt.userId {
				t.Errorf("result userId %d does not match initial userId %d", resultUserId, tt.userId)
			}
		})
	}
}
