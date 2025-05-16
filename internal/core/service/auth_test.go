package service_test

import (
	"context"
	"errors"
	"testing"

	"clean-arch/internal/core/domain"
	"clean-arch/internal/core/port/mock"
	"clean-arch/internal/core/service"
	"clean-arch/internal/core/util"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	type dependency struct {
		repo func(ctrl *gomock.Controller) *mock.MockUserRepository
		ts   func(ctrl *gomock.Controller) *mock.MockTokenService
	}

	type body struct {
		email    string
		password string
	}

	validPass := "password123"
	validHashedPassword, _ := util.HashPassword(validPass)

	invalidPass := "wrongpassword"
	invalidHashedPassword, _ := util.HashPassword(invalidPass + "123")

	// Define test cases
	tests := []struct {
		name string
		body
		dependency
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "when email is not found should and error is not ErrDataNotFound get error internal",
			body: body{
				email:    "poramin@treg.co.th",
				password: invalidPass,
			},
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) *mock.MockUserRepository {
					userRepo := mock.NewMockUserRepository(ctrl)
					userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
					return userRepo
				},
				ts: func(ctrl *gomock.Controller) *mock.MockTokenService {
					return nil
				},
			},
			expectErr:    true,
			expectErrMsg: domain.ErrInternal.Error(),
		},
		{
			name: "when email is not found should and error is ErrDataNotFound get error data not found",
			body: body{
				email:    "poramin@treg.co.th",
				password: invalidPass,
			},
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) *mock.MockUserRepository {
					userRepo := mock.NewMockUserRepository(ctrl)
					userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDataNotFound)
					return userRepo
				},
				ts: func(ctrl *gomock.Controller) *mock.MockTokenService {
					return nil
				},
			},
			expectErr:    true,
			expectErrMsg: domain.ErrDataNotFound.Error(),
		},
		{
			name: "when credentials is invalid should get error invalid credentials",
			body: body{
				email:    "poramin@treg.co.th",
				password: invalidPass,
			},
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) *mock.MockUserRepository {
					userRepo := mock.NewMockUserRepository(ctrl)
					userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
						Email:    "poramin@treg.co.th",
						Password: invalidHashedPassword,
					}, nil)
					return userRepo
				},
				ts: func(ctrl *gomock.Controller) *mock.MockTokenService {
					return nil
				},
			},
			expectErr:    true,
			expectErrMsg: domain.ErrInvalidCredentials.Error(),
		},
		{
			name: "when credentials is valid but create token error should get error",
			body: body{
				email:    "poramin@treg.co.th",
				password: validPass,
			},
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) *mock.MockUserRepository {
					userRepo := mock.NewMockUserRepository(ctrl)
					userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
						Email:    "poramin@treg.co.th",
						Password: validHashedPassword,
					}, nil)
					return userRepo
				},
				ts: func(ctrl *gomock.Controller) *mock.MockTokenService {
					tokenService := mock.NewMockTokenService(ctrl)
					tokenService.EXPECT().CreateToken(gomock.Any()).Return("", errors.New("token creation error"))
					return tokenService
				},
			},
			expectErr:    true,
			expectErrMsg: domain.ErrTokenCreationFailed.Error(),
		},
		{
			name: "when credentials is valid and token is created successfully should return token",
			body: body{
				email:    "poramin@treg.co.th",
				password: validPass,
			},
			dependency: dependency{
				repo: func(ctrl *gomock.Controller) *mock.MockUserRepository {
					userRepo := mock.NewMockUserRepository(ctrl)
					userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
						Email:    "poramin@treg.co.th",
						Password: validHashedPassword,
					}, nil)
					return userRepo
				},
				ts: func(ctrl *gomock.Controller) *mock.MockTokenService {
					tokenService := mock.NewMockTokenService(ctrl)
					tokenService.EXPECT().CreateToken(gomock.Any()).Return("valid_token", nil)
					return tokenService
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := service.NewAuth(test.dependency.repo(ctrl), test.dependency.ts(ctrl))
			defer service.ResetAuth()
			token, err := authService.Login(context.Background(), test.body.email, test.body.password)

			// Check if an error occurred
			if err != nil {
				// If an error was expected, check if it matches the expected error
				assert.True(t, test.expectErr, "Expected an error but got none")
				assert.Equal(t, test.expectErrMsg, err.Error(), "Expected an error")
				return
			}

			assert.Empty(t, err, "Expected no error but got one")
			assert.NotEmpty(t, token, "Expected a non-empty token")
		})
	}
}
