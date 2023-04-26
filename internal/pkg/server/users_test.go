package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"homework-week-5/internal/pkg/repository"
	mock_repository "homework-week-5/internal/pkg/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/openpgp/errors"
)

func Test_CreateUsers(t *testing.T) {
	var (
		ctx = context.Background()
		w   = httptest.NewRecorder()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		//arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_repository.NewMockUsersRepo(ctrl)
		s := Server{userRepo: repo}

		req, err := http.NewRequest(http.MethodPost, "users", bytes.NewReader([]byte(`{"name": "Ivan", "surname": "Ivanov", "age": 25}`)))
		require.NoError(t, err)

		repo.EXPECT().Add(ctx, &repository.User{Name: "Ivan", Surname: "Ivanov", Age: 25}).Return(int64(1), nil)

		//act
		status := s.CreateUsers(w, ctx, repo, req)

		//assert
		require.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name   string
			url    string
			body   string
			isReq  bool
			error  error
			status int
		}{
			{
				name:   "fail_1",
				url:    "users",
				body:   `{"name": //, "surname": "Ivanov", "age": 25}`,
				isReq:  false,
				error:  nil,
				status: http.StatusInternalServerError,
			},
			{
				name:   "fail_2",
				url:    "users",
				body:   `{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab", "surname": "Ivanov", "age": 25}`,
				isReq:  true,
				error:  errors.InvalidArgumentError("Invalid argument"),
				status: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: repo}

				req, err := http.NewRequest(http.MethodPost, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				if tc.isReq {
					var unmarshalled usersRequest
					if err := json.Unmarshal([]byte(tc.body), &unmarshalled); err != nil {
					}
					repo.EXPECT().Add(ctx, &repository.User{Name: unmarshalled.Name, Surname: unmarshalled.Surname, Age: unmarshalled.Age}).Return(int64(1), tc.error)
				}
				//act
				status := s.CreateUsers(w, ctx, repo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
}

func Test_UpdateUsers(t *testing.T) {
	var (
		ctx = context.Background()
		w   = httptest.NewRecorder()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		//arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_repository.NewMockUsersRepo(ctrl)
		s := Server{userRepo: repo}

		req, err := http.NewRequest(http.MethodPut, "users?id=1", bytes.NewReader([]byte(`{"name": "Ivan", "surname": "Ivanov", "age": 25}`)))
		require.NoError(t, err)

		repo.EXPECT().Update(ctx, &repository.User{ID: 1, Name: "Ivan", Surname: "Ivanov", Age: 25}).Return(true, nil)

		//act

		status := s.UpdateUsers(w, ctx, repo, req)

		//assert
		require.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name   string
			url    string
			body   string
			isReq  bool
			error  error
			status int
		}{
			{
				name:   "fail_1",
				url:    "users?id=1",
				body:   `{"name": \\//, "surname": "Ivanov", "age": 25}`,
				isReq:  false,
				error:  nil,
				status: http.StatusInternalServerError,
			},
			{
				name:   "fail_2",
				url:    "users",
				body:   `{"name": "Ivan", "surname": "Ivanov", "age": 25}`,
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_3",
				url:    "users?id=Ivan",
				body:   `{"name": "Ivan", "surname": "Ivanov", "age": 25}`,
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_4",
				url:    "users?id=1",
				body:   `{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab", "surname": "Ivanov", "age": 25}`,
				isReq:  true,
				error:  errors.InvalidArgumentError("Invalid argument"),
				status: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: repo}

				req, err := http.NewRequest(http.MethodPut, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				if tc.isReq {
					var unmarshalled usersRequest
					if err := json.Unmarshal([]byte(tc.body), &unmarshalled); err != nil {
					}
					repo.EXPECT().Update(ctx, &repository.User{ID: 1, Name: unmarshalled.Name, Surname: unmarshalled.Surname, Age: unmarshalled.Age}).Return(false, tc.error)
				}

				//act
				status := s.UpdateUsers(w, ctx, repo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
}

func Test_DeleteUsers(t *testing.T) {
	var (
		ctx = context.Background()
		w   = httptest.NewRecorder()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		//arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock_repository.NewMockUsersRepo(ctrl)
		s := Server{userRepo: repo}

		req, err := http.NewRequest(http.MethodDelete, "users?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		repo.EXPECT().Delete(ctx, 1).Return(true, nil)

		//act

		status := s.DeleteUsers(w, ctx, repo, req)

		//assert
		require.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name   string
			url    string
			body   string
			isReq  bool
			error  error
			status int
		}{
			{
				name:   "fail_1",
				url:    "users",
				body:   "",
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_2",
				url:    "users?id=Ivan",
				body:   "",
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_3",
				url:    "users?id=-1",
				body:   "",
				isReq:  true,
				error:  errors.InvalidArgumentError("Invalid argument"),
				status: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: repo}

				req, err := http.NewRequest(http.MethodDelete, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				if tc.isReq {
					var unmarshalled usersRequest
					if err := json.Unmarshal([]byte(tc.body), &unmarshalled); err != nil {
					}
					repo.EXPECT().Delete(ctx, -1).Return(false, tc.error)
				}

				//act
				status := s.DeleteUsers(w, ctx, repo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
}

func Test_GetUsers(t *testing.T) {
	var (
		ctx = context.Background()
		w   = httptest.NewRecorder()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name   string
			url    string
			body   string
			isAll  bool
			status int
		}{
			{
				name:   "success_1",
				url:    "users?id=1",
				body:   "",
				isAll:  false,
				status: http.StatusOK,
			},
			{
				name:   "success_2",
				url:    "users?id=all",
				body:   "",
				isAll:  true,
				status: http.StatusOK,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: repo}

				req, err := http.NewRequest(http.MethodGet, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				if tc.isAll {
					repo.EXPECT().List(ctx).Return([]*repository.User{}, nil)
				} else {
					repo.EXPECT().GetById(ctx, 1).Return(&repository.User{}, nil)
				}

				//act
				status := s.GetUsers(w, ctx, repo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name   string
			url    string
			body   string
			isReq  bool
			error  error
			status int
		}{
			{
				name:   "fail_1",
				url:    "users",
				body:   "",
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_2",
				url:    "users?id=Ivan",
				body:   "",
				isReq:  false,
				error:  nil,
				status: http.StatusBadRequest,
			},
			{
				name:   "fail_3",
				url:    "users?id=-1",
				body:   "",
				isReq:  true,
				error:  errors.InvalidArgumentError("Invalid argument"),
				status: http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: repo}

				req, err := http.NewRequest(http.MethodGet, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				if tc.isReq {
					var unmarshalled usersRequest
					if err := json.Unmarshal([]byte(tc.body), &unmarshalled); err != nil {
					}
					repo.EXPECT().GetById(ctx, -1).Return(&repository.User{}, tc.error)
				}

				//act
				status := s.GetUsers(w, ctx, repo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
}
