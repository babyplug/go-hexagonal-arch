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

func TestUserService_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *domain.User
	}

	tests := []struct {
		name         string
		args         args
		dependency   func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr    bool
		expectErrMsg string
		expectUser   *domain.User
	}{
		{
			name: "create user",
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:       "1",
					Name:     "John Doe",
					Email:    "poramin@treg.co.th",
					Password: "password",
				},
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepo
			},
			expectErr: false,
			expectUser: &domain.User{
				ID:       "1",
				Name:     "John Doe",
				Email:    "poramin@treg.co.th",
				Password: "password",
			},
		},
		{
			name: "when email already exists",
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Password: "password123",
				},
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "john.doe@example.com").Return(&domain.User{
					ID:       "1",
					Name:     "John Doe Doe",
					Email:    "john.doe@example.com",
					Password: "hashedpassword",
				}, nil)
				return userRepo
			},
			expectErr:    true,
			expectErrMsg: domain.ErrDuplicateEmail.Message,
			expectUser:   nil,
		},
		{
			name: "when password hashing fails",
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:       "3",
					Name:     "Error User",
					Email:    "error.user@example.com",
					Password: "123456790123456790123456790123456790123456790123456790123456790123456790123456909834091283409182304981203498120394820348",
				},
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "error.user@example.com").Return(nil, nil)
				return userRepo
			},
			expectErr:    true,
			expectErrMsg: "bcrypt: password length exceeds 72 bytes",
			expectUser:   nil,
		},
		{
			name: "when user creation fails",
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:       "4",
					Name:     "Success User",
					Email:    "success.user@example.com",
					Password: "password123",
				},
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "success.user@example.com").Return(nil, nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("mock err"))
				return userRepo
			},
			expectErr:    true,
			expectErrMsg: "mock err",
		},
		{
			name: "successful user creation",
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:       "4",
					Name:     "Success User",
					Email:    "success.user@example.com",
					Password: "password123",
				},
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "success.user@example.com").Return(nil, nil)
				userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return userRepo
			},
			expectErr: false,
			expectUser: &domain.User{
				ID:       "4",
				Name:     "Success User",
				Email:    "success.user@example.com",
				Password: "password123",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			err := userService.Create(test.args.ctx, test.args.user)
			// For each test, after calling the function under test:
			if test.expectErr {
				assert.Error(t, err)
				if test.expectErrMsg != "" {
					assert.Equal(t, test.expectErrMsg, err.Error())
				}
				return
			}
			assert.NoError(t, err)
			if test.expectUser != nil {
				assert.Equal(t, test.args.user.ID, test.expectUser.ID)
				assert.Equal(t, test.args.user.Name, test.expectUser.Name)
				assert.Equal(t, test.args.user.Email, test.expectUser.Email)
				// Password should be hashed and equal to the hashed password when compared with util.ComparePassword
				assert.Nil(t, util.ComparePassword(test.expectUser.Password, test.args.user.Password), "password should be  hashed")
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		dependency func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr  bool
		expectUser *domain.User
	}{
		{
			name: "get user by ID successfully",
			id:   "1",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByID(gomock.Any(), "1").Return(&domain.User{
					ID:    "1",
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}, nil)
				return userRepo
			},
			expectErr: false,
			expectUser: &domain.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		},
		{
			name: "user not found by ID",
			id:   "2",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByID(gomock.Any(), "2").Return(nil, errors.New("user not found"))
				return userRepo
			},
			expectErr:  true,
			expectUser: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			user, err := userService.GetByID(context.Background(), test.id)

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectUser, user)
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		dependency   func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr    bool
		expectErrMsg string
		expect       *domain.User
	}{
		{
			name:  "get by email successfully",
			email: "test@gmail.com",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@gmail.com").Return(&domain.User{ID: "1", Email: "test@gmail.com"}, nil)
				return userRepo
			},
			expectErr: false,
			expect: &domain.User{
				ID:    "1",
				Email: "test@gmail.com",
			},
		},
		{
			name:  "get by email fails",
			email: "test@gmail.com",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "test@gmail.com").Return(nil, errors.New("get by email error"))
				return userRepo
			},
			expectErr:    true,
			expectErrMsg: "get by email error",
			expect:       nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			user, err := userService.GetByEmail(context.Background(), test.email)

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expect, user)
		})
	}
}

func TestUserService_List(t *testing.T) {
	tests := []struct {
		name       string
		page       int64
		size       int64
		dependency func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr  bool
		expectList []*domain.User
	}{
		{
			name: "list users successfully",
			page: 0,
			size: 2,
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.User{
					{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe@example.com",
					},
					{
						ID:    "2",
						Name:  "Jane Doe",
						Email: "jane.doe@example.com",
					},
				}, nil)
				return userRepo
			},
			expectErr: false,
			expectList: []*domain.User{
				{
					ID:    "1",
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
				{
					ID:    "2",
					Name:  "Jane Doe",
					Email: "jane.doe@example.com",
				},
			},
		},
		{
			name: "error listing users",
			page: 0,
			size: 2,
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
				return userRepo
			},
			expectErr:  true,
			expectList: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			list, err := userService.List(context.Background(), test.page, test.size)

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectList, list)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	tests := []struct {
		name         string
		user         *domain.User
		dependency   func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr    bool
		expectErrMsg string
	}{
		{
			name: "when get by id failed",
			user: &domain.User{
				ID:    "1",
				Name:  "Updated Name",
				Email: "updated.email@example.com",
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByID(gomock.Any(), "1").Return(nil, errors.New("mock err"))
				return userRepo
			},
			expectErr:    true,
			expectErrMsg: domain.ErrDataNotFound.Error(),
		},
		{
			name: "update user successfully",
			user: &domain.User{
				ID:    "1",
				Name:  "Updated Name",
				Email: "updated.email@example.com",
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByID(gomock.Any(), "1").Return(&domain.User{
					ID:    "1",
					Name:  "Old Name",
					Email: "old.email@example.com",
				}, nil)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "updated.email@example.com").Return(nil, nil)
				userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				return userRepo
			},
			expectErr: false,
		},
		{
			name: "update fails due to email conflict",
			user: &domain.User{
				ID:    "1",
				Name:  "Updated Name",
				Email: "conflict.email@example.com",
			},
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().GetByID(gomock.Any(), "1").Return(&domain.User{
					ID:    "1",
					Name:  "Old Name",
					Email: "old.email@example.com",
				}, nil)
				userRepo.EXPECT().GetByEmail(gomock.Any(), "conflict.email@example.com").Return(&domain.User{}, nil)
				return userRepo
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			err := userService.Update(context.Background(), test.user)

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		dependency func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr  bool
	}{
		{
			name: "delete user successfully",
			id:   "1",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Delete(gomock.Any(), "1").Return(nil)
				return userRepo
			},
			expectErr: false,
		},
		{
			name: "delete user fails",
			id:   "2",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Delete(gomock.Any(), "2").Return(errors.New("delete error"))
				return userRepo
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			err := userService.Delete(context.Background(), test.id)

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserService_Count(t *testing.T) {
	tests := []struct {
		name        string
		dependency  func(ctrl *gomock.Controller) *mock.MockUserRepository
		expectErr   bool
		expectCount int
	}{
		{
			name: "count users successfully",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Count(gomock.Any()).Return(10, nil)
				return userRepo
			},
			expectErr:   false,
			expectCount: 10,
		},
		{
			name: "count users fails",
			dependency: func(ctrl *gomock.Controller) *mock.MockUserRepository {
				userRepo := mock.NewMockUserRepository(ctrl)
				userRepo.EXPECT().Count(gomock.Any()).Return(0, errors.New("count error"))
				return userRepo
			},
			expectErr:   true,
			expectCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := service.NewUser(test.dependency(ctrl))
			defer service.ResetUser()
			count, err := userService.Count(context.Background())

			if test.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expectCount, count)
		})
	}
}
