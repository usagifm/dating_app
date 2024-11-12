package auth_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	mock_atomic "github.com/usagifm/dating-app/lib/atomic/mock"
	"github.com/usagifm/dating-app/lib/helper"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
	"github.com/usagifm/dating-app/src/v1/contract"
	"github.com/usagifm/dating-app/src/v1/service/auth"
	mock_auth "github.com/usagifm/dating-app/src/v1/service/auth/mock"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	os.Chdir("../../../../")

	app.Init(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestSignUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	atomicMock := mock_atomic.NewMockAtomicSessionProvider(ctrl)
	rUser := mock_auth.NewMockUserRepository(ctrl)
	rUserPreference := mock_auth.NewMockUserReferenceRepository(ctrl)
	// atomicSess := mock_atomic.NewMockAtomicSession(ctrl)
	// atomicSessCtx := atomic.NewAtomicSessionContext(context.Background(), atomicSess)

	type mockFields struct {
		atomicSession   *mock_atomic.MockAtomicSessionProvider
		rUser           *mock_auth.MockUserRepository
		rUserPreference *mock_auth.MockUserReferenceRepository
	}

	mocks := mockFields{
		atomicSession:   atomicMock,
		rUser:           rUser,
		rUserPreference: rUserPreference,
	}

	type args struct {
		ctx    context.Context
		params contract.SignUpRequest
	}

	// var errorNil string

	// mockUser := entity.User{
	// 	Id:         1,
	// 	IsVerified: false,
	// 	Name:       faker.Name(),
	// 	Gender:     faker.Gender(),
	// 	Email:      faker.Email(),
	// 	Password:   faker.Password(),
	// 	Age:        20,
	// 	Bio:        faker.Paragraph(),
	// }

	mockCreateNewUser := entity.User{
		Name:   "faqqih",
		Gender: "male",
		Email:  "usagifm@gmail.com",
		Age:    20,
		Bio:    "test",
	}

	// mockCreateNewUserPreference := entity.UserPreference{
	// 	UserId:          2,
	// 	PreferredGender: faker.Gender(),
	// 	MinAge:          18,
	// 	MaxAge:          40,
	// }

	userEntity := &entity.User{}

	tests := []struct {
		name     string
		mockFunc func(mock mockFields, arg args)
		args     args
		want     *error
		wantErr  bool
	}{
		{
			name: "Error when trying to get user by email",
			args: args{
				ctx: context.Background(),
				params: contract.SignUpRequest{
					Name:            "faqqih",
					Gender:          "male",
					Email:           "usagifm@gmail.com",
					Age:             20,
					Bio:             "test",
					PreferredGender: "female",
					MinAge:          18,
					MaxAge:          40,
				},
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.rUser.EXPECT().GetUserByEmail(gomock.Any(), arg.params.Email).Return(nil, assert.AnError).Times(1)
				// mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
				// mock.rUser.EXPECT().GetUserByEmail(gomock.Any(), arg.params.Email).Return(userEntity, nil)
				// if userEntity != nil {
				// 	atomicSess.EXPECT().Rollback(gomock.Any()).Times(1)
				// }

				// // mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(atomicSessCtx, arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(entity.Template{}, assert.AnError)
				// // mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLog).Return(nil)
				// atomicSess.EXPECT().Rollback(gomock.Any()).Times(1)
			},
			wantErr: true,
		},
		{
			name: "User with the same email is already registered",
			args: args{
				ctx: context.Background(),
				params: contract.SignUpRequest{
					Name:            "faqqih",
					Gender:          "male",
					Email:           "usagifm@gmail.com",
					Age:             20,
					Bio:             "test",
					PreferredGender: "female",
					MinAge:          18,
					MaxAge:          40,
				},
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.rUser.EXPECT().GetUserByEmail(gomock.Any(), arg.params.Email).Return(userEntity, nil).Times(1)
				// mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
				// mock.rUser.EXPECT().GetUserByEmail(gomock.Any(), arg.params.Email).Return(userEntity, nil)
				// if userEntity != nil {
				// 	atomicSess.EXPECT().Rollback(gomock.Any()).Times(1)
				// }

				// // mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(atomicSessCtx, arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(entity.Template{}, assert.AnError)
				// // mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLog).Return(nil)
				// atomicSess.EXPECT().Rollback(gomock.Any()).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Password hashing fails",
			args: args{
				ctx: context.Background(),
				params: contract.SignUpRequest{
					Name:            "faqqih",
					Gender:          "male",
					Email:           "usagifm@gmail.com",
					Age:             20,
					Bio:             "test",
					PreferredGender: "female",
					MinAge:          18,
					MaxAge:          40,
				},
			},
			mockFunc: func(mock mockFields, arg args) {
				// Simulate user does not exist
				mock.rUser.EXPECT().GetUserByEmail(gomock.Any(), arg.params.Email).Return(nil, sql.ErrNoRows).Times(1)

				hashed, _ := helper.HashPassword(arg.params.Password)

				mockCreateNewUser.Password = hashed
				mock.rUser.EXPECT().CreateNewUser(gomock.Any(), mockCreateNewUser).Return(nil, assert.AnError)

			},
			wantErr: true, // Expect error due to hashing failure
		},
		// {
		// 	name: "Failed to send SMS ",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		params: contract.Delivery{
		// 			Channel:   "sms",
		// 			Service:   "otp",
		// 			Purpose:   "registration",
		// 			Recipient: "0812345678",
		// 			Params: []contract.Param{
		// 				{Name: "otp", Value: "12312312"},
		// 				{Name: "unique_id", Value: "Wfg1Ha52"},
		// 			},
		// 		},
		// 	},
		// 	mockFunc: func(mock mockFields, arg args) {
		// 		mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
		// 		mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(gomock.Any(), arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(mockGetSMSTemplate, nil)
		// 		mock.rSMS.EXPECT().SendSMSWithSME(atomicSessCtx, arg.params.Recipient, mockGetSMSTemplate.MessageTemplate.String, arg.params.Params).Return(assert.AnError)
		// 		mockCreateSendLogSMS.Recipient = arg.params.Recipient
		// 		mockCreateSendLogSMS.Error = sql.NullString{String: assert.AnError.Error(), Valid: true}
		// 		mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLogSMS).Return(nil)
		// 		atomicSess.EXPECT().Commit(gomock.Any()).Times(1)
		// 	},
		// 	want:    &assert.AnError,
		// 	wantErr: false,
		// },
		// {
		// 	name: "Failed to send Email ",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		params: contract.Delivery{
		// 			Channel:   "email",
		// 			Service:   "academy",
		// 			Purpose:   "submission",
		// 			Recipient: "test@gmail.com",
		// 		},
		// 	},
		// 	mockFunc: func(mock mockFields, arg args) {
		// 		mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
		// 		mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(gomock.Any(), arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(mockGetEmailTemplate, nil)
		// 		mock.rEmail.EXPECT().SendEmailWithMailgun(atomicSessCtx, mockGetEmailTemplate.SenderEmail.String, mockGetEmailTemplate.SenderName.String, arg.params.Recipient, mockGetEmailTemplate.Cc.String, mockGetEmailTemplate.EmailSubject.String, mockGetEmailTemplate.EmailTemplate.String, arg.params.Params).Return(assert.AnError)
		// 		mockCreateSendLogEmail.Recipient = arg.params.Recipient
		// 		mockCreateSendLogEmail.Error = sql.NullString{String: assert.AnError.Error(), Valid: true}
		// 		mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLogEmail).Return(nil)
		// 		atomicSess.EXPECT().Commit(gomock.Any()).Times(1)
		// 	},
		// 	want:    &assert.AnError,
		// 	wantErr: false,
		// },
		// {
		// 	name: "Success to send Email but failed to create send log",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		params: contract.Delivery{
		// 			Channel:   "email",
		// 			Service:   "academy",
		// 			Purpose:   "submission",
		// 			Recipient: "test@gmail.com",
		// 		},
		// 	},
		// 	mockFunc: func(mock mockFields, arg args) {
		// 		mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
		// 		mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(gomock.Any(), arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(mockGetEmailTemplate, nil)
		// 		mock.rEmail.EXPECT().SendEmailWithMailgun(atomicSessCtx, mockGetEmailTemplate.SenderEmail.String, mockGetEmailTemplate.SenderName.String, arg.params.Recipient, mockGetEmailTemplate.Cc.String, mockGetEmailTemplate.EmailSubject.String, mockGetEmailTemplate.EmailTemplate.String, arg.params.Params).Return(nil)
		// 		mockCreateSendLogEmailSuccess.Recipient = arg.params.Recipient
		// 		mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLogEmailSuccess).Return(assert.AnError)
		// 		atomicSess.EXPECT().Rollback(gomock.Any()).Times(1)
		// 	},
		// 	want:    &assert.AnError,
		// 	wantErr: true,
		// },
		// {
		// 	name: "Success to send Email and success to create send log",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		params: contract.Delivery{
		// 			Channel:   "email",
		// 			Service:   "academy",
		// 			Purpose:   "submission",
		// 			Recipient: "test@gmail.com",
		// 		},
		// 	},
		// 	mockFunc: func(mock mockFields, arg args) {
		// 		mock.atomicSession.EXPECT().BeginSession(gomock.Any()).Return(atomicSessCtx, nil).Times(1)
		// 		mock.rTemplate.EXPECT().GetTemplateByDeliveryPayload(gomock.Any(), arg.params.Channel, arg.params.Service, arg.params.Purpose).Return(mockGetEmailTemplate, nil)
		// 		mock.rEmail.EXPECT().SendEmailWithMailgun(atomicSessCtx, mockGetEmailTemplate.SenderEmail.String, mockGetEmailTemplate.SenderName.String, arg.params.Recipient, mockGetEmailTemplate.Cc.String, mockGetEmailTemplate.EmailSubject.String, mockGetEmailTemplate.EmailTemplate.String, arg.params.Params).Return(nil)
		// 		mockCreateSendLogEmailSuccess.Recipient = arg.params.Recipient
		// 		mock.rSendLog.EXPECT().CreateSendLog(atomicSessCtx, mockCreateSendLogEmailSuccess).Return(nil)
		// 		atomicSess.EXPECT().Commit(gomock.Any()).Times(1)
		// 	},
		// 	want:    &assert.AnError,
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			a := auth.NewAuthService(atomicMock, rUser, rUserPreference)
			err := a.SignUp(tt.args.ctx, tt.args.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("auth.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// assert.Equal(t, tt.wantErr, true)
		})
	}
}
