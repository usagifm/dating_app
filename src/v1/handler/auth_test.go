package handler

import (
	"context"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	mock_atomic "github.com/usagifm/dating-app/lib/atomic/mock"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/v1/contract"
	mock_handler "github.com/usagifm/dating-app/src/v1/handler/mock"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	os.Chdir("../../../")

	app.Init(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	atomicMock := mock_atomic.NewMockAtomicSessionProvider(ctrl)
	sAuthMock := mock_handler.NewMockAuthService(ctrl)
	sDatingMock := mock_handler.NewMockDatingService(ctrl)
	// atomicSess := mock_atomic.NewMockAtomicSession(ctrl)
	// atomicSessCtx := atomic.NewAtomicSessionContext(context.Background(), atomicSess)

	type mockFields struct {
		atomicSession *mock_atomic.MockAtomicSessionProvider
		sAuthMock     *mock_handler.MockAuthService
	}

	mockParam := contract.SignUpRequest{
		Name:            faker.Name(),
		Email:           faker.Email(),
		Password:        faker.Password(),
		Gender:          faker.Gender(),
		Age:             20,
		Bio:             faker.Paragraph(),
		PhotoUrl:        faker.Paragraph(),
		PreferredGender: faker.Gender(),
		MinAge:          18,
		MaxAge:          40,
	}

	mocks := mockFields{
		atomicSession: atomicMock,
		sAuthMock:     sAuthMock,
	}

	type args struct {
		ctx    context.Context
		params contract.SignUpRequest
	}

	tests := []struct {
		name      string
		parameter map[string]string
		mockFunc  func(mock mockFields, arg args)
		args      args
		want      *error
		wantErr   bool
	}{
		{
			name: "failed to sign up",
			args: args{
				ctx:    context.Background(),
				params: mockParam,
			},
			mockFunc: func(mock mockFields, arg args) {
				// mockParam, _ := contract.ValidateAndBuildSignUpRequest(&http.Request{})
				mock.sAuthMock.EXPECT().SignUp(gomock.Any(), mockParam).Return(assert.AnError)
			},
			want:    &assert.AnError,
			wantErr: true,
		},
		{
			name: "success sign up",
			args: args{
				ctx:    context.Background(),
				params: mockParam,
			},
			mockFunc: func(mock mockFields, arg args) {
				// mockParam, _ := contract.ValidateAndBuildSignUpRequest(&http.Request{})
				mock.sAuthMock.EXPECT().SignUp(gomock.Any(), mockParam).Return(nil)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			h := NewDatingAppHandler(sAuthMock, sDatingMock)
			h.sAuth.SignUp(tt.args.ctx, tt.args.params)

		})
	}

}
