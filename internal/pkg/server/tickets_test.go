// test cover = 51.1% - дальше доделывать не стал
package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"homework-week-5/internal/pkg/repository"
	mock_repository "homework-week-5/internal/pkg/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/openpgp/errors"
)

func Test_GetTickets(t *testing.T) {
	var (
		ctx = context.Background()
		w   = httptest.NewRecorder()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name    string
			url     string
			body    string
			reqCode int
			status  int
		}{
			{
				name:    "success_1",
				url:     "tickets?id=all",
				body:    "",
				reqCode: 0,
				status:  http.StatusOK,
			},
			{
				name:    "success_2",
				url:     "tickets?id=1",
				body:    "",
				reqCode: 1,
				status:  http.StatusOK,
			},
			{
				name:    "success_3",
				url:     "tickets?user_id=1",
				body:    "",
				reqCode: 2,
				status:  http.StatusOK,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ticketRepo := mock_repository.NewMockTicketRepo(ctrl)
				userRepo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: userRepo, ticketRepo: ticketRepo}

				req, err := http.NewRequest(http.MethodPost, tc.url, bytes.NewReader([]byte(tc.body)))
				require.NoError(t, err)

				switch tc.reqCode {
				case 0:
					ticketRepo.EXPECT().List(ctx).Return([]*repository.Ticket{}, nil)
				case 1:
					ticketRepo.EXPECT().GetById(ctx, 1).Return(&repository.Ticket{}, nil)
				case 2:
					userRepo.EXPECT().GetById(ctx, 1).Return(&repository.User{}, nil)
					ticketRepo.EXPECT().GetByUserId(ctx, 1).Return([]*repository.Ticket{}, nil)
				}

				//act
				status := s.GetTickets(w, ctx, userRepo, ticketRepo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			name      string
			url       string
			isNeedReq bool
			reqCode   int
			error     error
			status    int
		}{
			{
				name:      "fail_1",
				url:       "tickets?user_id=Ivan",
				isNeedReq: false,
				reqCode:   2,
				error:     nil,
				status:    http.StatusBadRequest,
			},
			{
				name:      "fail_2",
				url:       "tickets?user_id=-1",
				isNeedReq: true,
				reqCode:   2,
				error:     errors.InvalidArgumentError("Invalid argument"),
				status:    http.StatusBadRequest,
			},
			{
				name:      "fail_3",
				url:       "tickets",
				isNeedReq: false,
				reqCode:   1,
				error:     nil,
				status:    http.StatusBadRequest,
			},
			{
				name:      "fail_4",
				url:       "tickets?id=Ivan",
				isNeedReq: false,
				reqCode:   1,
				error:     nil,
				status:    http.StatusBadRequest,
			},
			{
				name:      "fail_5",
				url:       "tickets?id=-1",
				isNeedReq: true,
				reqCode:   1,
				error:     errors.InvalidArgumentError("Invalid argument"),
				status:    http.StatusBadRequest,
			},
		}

		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				//arrange
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ticketRepo := mock_repository.NewMockTicketRepo(ctrl)
				userRepo := mock_repository.NewMockUsersRepo(ctrl)
				s := Server{userRepo: userRepo, ticketRepo: ticketRepo}

				req, err := http.NewRequest(http.MethodPost, tc.url, bytes.NewReader([]byte{}))
				require.NoError(t, err)

				if tc.isNeedReq {
					switch tc.reqCode {
					case 0:
						ticketRepo.EXPECT().List(ctx).Return([]*repository.Ticket{}, tc.error)
					case 1:
						ticketRepo.EXPECT().GetById(ctx, -1).Return(&repository.Ticket{}, tc.error)
					case 2:
						userRepo.EXPECT().GetById(ctx, -1).Return(&repository.User{}, tc.error)
						ticketRepo.EXPECT().GetByUserId(ctx, -1).Return([]*repository.Ticket{}, tc.error)
					}
				}

				//act
				status := s.GetTickets(w, ctx, userRepo, ticketRepo, req)

				//assert
				require.Equal(t, tc.status, status)
			})
		}
	})
}
