package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestJWTService_GenerateJWT(t *testing.T) {
	type args struct {
		id       uuid.UUID
		username string
		secret   string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		checkFunc func(*testing.T, *TokenResult)
	}{
		{
			name: "success - generates valid token",
			args: args{
				id:       uuid.New(),
				username: "testuser",
				secret:   "supersecretkey",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, result *TokenResult) {
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.True(t, result.ExpiresIn <= 3600 && result.ExpiresIn > 3590, "ExpiresIn is not within expected range")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(tt.args.secret)

			got, err := service.GenerateJWT(tt.args.id, tt.args.username)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				tt.checkFunc(t, got)

				parts := strings.Split(got.AccessToken, ".")
				assert.Len(t, parts, 3, "JWT should have 3 parts")
			}
		})
	}
}
