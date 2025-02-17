package psql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/al3ksus/messengerusers/internal/domain/models"
)

var (
	TestUsername       = "user1"
	TestPass           = []byte("qwerty")
	TestUserId   int64 = 1
	EmptyUserId  int64 = 0
)

var TestUser = models.User{
	Id:           TestUserId,
	Username:     TestUsername,
	PasswordHash: TestPass,
	IsActive:     true,
}

func TestRepository_SaveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rep := New(db)

	type args struct {
		ctx      context.Context
		username string
		password []byte
	}
	type mockBehavior func(ctx context.Context, username string, pass []byte)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int64
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(ctx context.Context, username string, pass []byte) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(TestUserId)
				mock.ExpectQuery("INSERT INTO users").WithArgs(username, pass).WillReturnRows(rows)

			},
			want: TestUserId,
		},
		{
			name: "Error",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
				password: TestPass,
			},
			mockBehavior: func(ctx context.Context, username string, pass []byte) {
				mock.ExpectQuery("INSERT INTO users").WithArgs(username, pass).WillReturnError(errors.New(""))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args.ctx, tt.args.username, tt.args.password)

			got, err := rep.SaveUser(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.SaveUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rep := New(db)

	type args struct {
		ctx      context.Context
		username string
	}
	type mockBehavior func(ctx context.Context, username string)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         models.User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
			},
			mockBehavior: func(ctx context.Context, username string) {
				rows := sqlmock.
					NewRows([]string{"id", "username", "pass_hash", "is_active"}).
					AddRow(TestUser.Id, TestUser.Username, TestUser.PasswordHash, TestUser.IsActive)

				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(username).WillReturnRows(rows)
			},
			want: TestUser,
		},
		{
			name: "Error",
			args: args{
				ctx:      context.Background(),
				username: TestUsername,
			},
			mockBehavior: func(ctx context.Context, username string) {
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(username).WillReturnError(errors.New(""))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args.ctx, tt.args.username)

			got, err := rep.GetUser(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_SetInactive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rep := New(db)

	type args struct {
		ctx    context.Context
		userId int64
	}
	type mockBehavior func(ctx context.Context, userId int64)
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(ctx context.Context, userId int64) {
				mock.ExpectBegin()

				rows := sqlmock.
					NewRows([]string{"is_active"}).
					AddRow(true)

				mock.ExpectQuery("SELECT is_active FROM users").WithArgs(userId).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET is_active = FALSE").WithArgs(userId).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "ErrorSelect",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(ctx context.Context, userId int64) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT is_active FROM users").WithArgs(userId).WillReturnError(errors.New(""))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "ErrorUpdate",
			args: args{
				ctx:    context.Background(),
				userId: TestUserId,
			},
			mockBehavior: func(ctx context.Context, userId int64) {
				mock.ExpectBegin()

				rows := sqlmock.
					NewRows([]string{"is_active"}).
					AddRow(true)

				mock.ExpectQuery("SELECT is_active FROM users").WithArgs(userId).WillReturnRows(rows)
				mock.ExpectExec("UPDATE users SET is_active = FALSE").WithArgs(userId).WillReturnError(errors.New(""))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args.ctx, tt.args.userId)

			err := rep.SetInactive(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.SetInactive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
