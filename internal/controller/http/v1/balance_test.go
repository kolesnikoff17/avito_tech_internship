package v1

import (
	"balance_api/internal/entity"
	ucmock "balance_api/internal/mocks/usecase"
	"balance_api/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	h := gin.New()
	uc := ucmock.NewBalance(t)
	l, _ := logger.New("debug")
	NewRouter(h, uc, l)

	uc.On("GetByID", ctx, 1).Return(entity.Balance{ID: 1, Amount: "200"}, nil)
	uc.On("GetByID", ctx, 2).Return(entity.Balance{}, entity.ErrNoID)
	uc.On("GetByID", ctx, 3).Return(entity.Balance{}, errors.New("aboba"))

	req := "/v1/user"

	type testCases struct {
		name    string
		query   string
		expCode int
		resp    interface{}
	}

	cases := []testCases{{
		name:    "valid",
		query:   "?id=1",
		expCode: http.StatusOK,
		resp:    entity.Balance{ID: 1, Amount: "200"},
	}, {
		name:    "empty query",
		query:   "",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "wrong id",
		query:   "?id=-1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "text id",
		query:   "?id=a",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "no such id",
		query:   "?id=2",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "No such id"},
	}, {
		name:    "db error",
		query:   "?id=3",
		expCode: http.StatusInternalServerError,
		resp:    response{Msg: "Database error"},
	},
	}

	for _, tc := range cases {
		r, _ := http.NewRequest(http.MethodGet, req+tc.query, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		require.Equal(t, tc.expCode, w.Code)
		b, _ := json.Marshal(tc.resp)
		require.Equal(t, string(b), w.Body.String())
	}
}

func TestIncrease(t *testing.T) {
	ctx := context.Background()
	h := gin.New()
	uc := ucmock.NewBalance(t)
	l, _ := logger.New("debug")
	NewRouter(h, uc, l)

	req := "/v1/user"

	uc.On("Increase", ctx, entity.Balance{ID: 1, Amount: "200"}).Return(nil)
	uc.On("Increase", ctx, entity.Balance{ID: 2, Amount: "200"}).Return(errors.New("aboba"))

	type testCases struct {
		name    string
		body    userPostRequest
		expCode int
		resp    interface{}
	}

	cases := []testCases{{
		name:    "valid",
		body:    userPostRequest{ID: 1, Amount: "200"},
		expCode: http.StatusOK,
		resp:    struct{}{},
	}, {
		name:    "wrong id",
		body:    userPostRequest{ID: -1, Amount: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request body format"},
	}, {
		name:    "negative money",
		body:    userPostRequest{ID: 2, Amount: "-20"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid money format"},
	}, {
		name:    "wrong money format",
		body:    userPostRequest{ID: 2, Amount: "1.2.3.4"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid money format"},
	}, {
		name:    "db err",
		body:    userPostRequest{ID: 2, Amount: "200"},
		expCode: http.StatusInternalServerError,
		resp:    response{Msg: "Database error"},
	},
	}

	for _, tc := range cases {
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(tc.body)
		r, _ := http.NewRequest(http.MethodPost, req, &buf)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		require.Equal(t, tc.expCode, w.Code)
		b, _ := json.Marshal(tc.resp)
		require.Equal(t, string(b), w.Body.String())
	}
}

func TestOrder(t *testing.T) {
	ctx := context.Background()
	h := gin.New()
	uc := ucmock.NewBalance(t)
	l, _ := logger.New("debug")
	NewRouter(h, uc, l)

	req := "/v1/order"

	uc.On("CreateOrder", ctx, entity.Order{ID: 1, ServiceID: 2, UserID: 1, Sum: "200"}).
		Return(nil)
	uc.On("ChangeOrderStatus", ctx, entity.Order{ID: 1, ServiceID: 2, UserID: 1, Sum: "200", StatusID: 2}).
		Return(nil)
	uc.On("CreateOrder", ctx, entity.Order{ID: 1, ServiceID: 10, UserID: 1, Sum: "200"}).
		Return(entity.ErrNoService)
	uc.On("CreateOrder", ctx, entity.Order{ID: 1, ServiceID: 2, UserID: 10, Sum: "200"}).
		Return(entity.ErrNoID)
	uc.On("CreateOrder", ctx, entity.Order{ID: 10, ServiceID: 2, UserID: 1, Sum: "200"}).
		Return(entity.ErrOrderExists)
	uc.On("CreateOrder", ctx, entity.Order{ID: 1, ServiceID: 2, UserID: 1, Sum: "1000"}).
		Return(entity.ErrNotEnoughMoney)
	uc.On("ChangeOrderStatus", ctx, entity.Order{ID: 2, ServiceID: 2, UserID: 1, Sum: "200", StatusID: 2}).
		Return(entity.ErrOrderNoExists)
	uc.On("ChangeOrderStatus", ctx, entity.Order{ID: 3, ServiceID: 2, UserID: 1, Sum: "200", StatusID: 2}).
		Return(entity.ErrOrderMismatch)
	uc.On("ChangeOrderStatus", ctx, entity.Order{ID: 4, ServiceID: 2, UserID: 1, Sum: "200", StatusID: 2}).
		Return(entity.ErrCantChangeStatus)
	uc.On("ChangeOrderStatus", ctx, entity.Order{ID: 5, ServiceID: 2, UserID: 1, Sum: "200", StatusID: 3}).
		Return(errors.New("aboba"))

	type testCases struct {
		name    string
		body    orderPostRequest
		expCode int
		resp    interface{}
	}

	cases := []testCases{{
		name:    "valid create",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusOK,
		resp:    struct{}{},
	}, {
		name:    "valid change",
		body:    orderPostRequest{Action: "approve", ID: 1, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusOK,
		resp:    struct{}{},
	}, {
		name:    "wrong id",
		body:    orderPostRequest{Action: "create", ID: -1, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request body format"},
	}, {
		name:    "wrong service id",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: -1, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request body format"},
	}, {
		name:    "invalid user id",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: -1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request body format"},
	}, {
		name:    "wrong money format",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: 1, Sum: "aboba"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid money format"},
	}, {
		name:    "negative money",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: 1, Sum: "-2"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid money format"},
	}, {
		name:    "wrong action",
		body:    orderPostRequest{Action: "aboba", ID: 1, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid order action"},
	}, {
		name:    "no service",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 10, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "No such service"},
	}, {
		name:    "no id",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: 10, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "No such id"},
	}, {
		name:    "not enough money",
		body:    orderPostRequest{Action: "create", ID: 1, ServiceID: 2, UserID: 1, Sum: "1000"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Not enough money"},
	}, {
		name:    "order exists",
		body:    orderPostRequest{Action: "create", ID: 10, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Order already exists"},
	}, {
		name:    "order not exists",
		body:    orderPostRequest{Action: "approve", ID: 2, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Order not exists"},
	}, {
		name:    "mismatch",
		body:    orderPostRequest{Action: "approve", ID: 3, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Wrong order data"},
	}, {
		name:    "cant change status",
		body:    orderPostRequest{Action: "approve", ID: 4, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Order already approved/canceled"},
	}, {
		name:    "db error",
		body:    orderPostRequest{Action: "cancel", ID: 5, ServiceID: 2, UserID: 1, Sum: "200"},
		expCode: http.StatusInternalServerError,
		resp:    response{Msg: "Database error"},
	},
	}

	for _, tc := range cases {
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(tc.body)
		r, _ := http.NewRequest(http.MethodPost, req, &buf)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		require.Equal(t, tc.expCode, w.Code)
		b, _ := json.Marshal(tc.resp)
		require.Equal(t, string(b), w.Body.String())
	}
}

func TestHistory(t *testing.T) {
	ctx := context.Background()
	h := gin.New()
	uc := ucmock.NewBalance(t)
	l, _ := logger.New("debug")
	NewRouter(h, uc, l)

	req := "/v1/history"

	uc.On("GetHistory", ctx,
		entity.History{UserID: 1, Limit: 10, OrderBy: "sum", Desc: true, Page: 1}).
		Return(entity.History{Orders: []entity.Order{{
			Sum:         "200",
			ServiceName: "aboba",
			Status:      "approved",
			Time:        entity.MyTime{Time: time.Unix(10, 0)},
		}, {
			Sum:         "1",
			ServiceName: "aboba2",
			Status:      "canceled",
			Time:        entity.MyTime{Time: time.Unix(10, 0)},
		},
		}, Limit: 10, OrderBy: "sum", Desc: true, Page: 1}, nil)

	uc.On("GetHistory", ctx,
		entity.History{UserID: 2, OrderBy: "date"}).Return(entity.History{}, entity.ErrNoID)

	uc.On("GetHistory", ctx,
		entity.History{UserID: 3, Limit: 100, Page: 99, OrderBy: "date"}).Return(entity.History{}, entity.ErrEmptyPage)

	uc.On("GetHistory", ctx,
		entity.History{UserID: 4, OrderBy: "date"}).Return(entity.History{}, errors.New("aboba"))

	type testCases struct {
		name    string
		query   string
		expCode int
		resp    interface{}
	}

	cases := []testCases{{
		name:    "valid",
		query:   "?id=1&limit=10&page=1&order_by=sum&desc=true",
		expCode: http.StatusOK,
		resp: entity.History{Orders: []entity.Order{{
			Sum:         "200",
			ServiceName: "aboba",
			Status:      "approved",
			Time:        entity.MyTime{Time: time.Unix(10, 0)},
		}, {
			Sum:         "1",
			ServiceName: "aboba2",
			Status:      "canceled",
			Time:        entity.MyTime{Time: time.Unix(10, 0)},
		},
		}},
	}, {
		name:    "wrong id",
		query:   "?id=-1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "wrong limit",
		query:   "?id=1&limit=-1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "wrong page",
		query:   "?id=1&page=-1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "wrong desc param",
		query:   "?id=1&desc=-1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "page zero and non zero limit",
		query:   "?id=1&limit=10",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Limit and page should be both zero or non zero"},
	}, {
		name:    "limit zero and non zero page",
		query:   "?id=1&page=10",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Limit and page should be both zero or non zero"},
	}, {
		name:    "wrong order by param",
		query:   "?id=1&order_by=1",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Wrong \"order by\" value"},
	}, {
		name:    "no id",
		query:   "?id=2&order_by=date",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "No such id"},
	}, {
		name:    "empty page",
		query:   "?id=3&limit=100&page=99",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "The page is empty"},
	}, {
		name:    "db error",
		query:   "?id=4",
		expCode: http.StatusInternalServerError,
		resp:    response{Msg: "Database error"},
	},
	}

	for _, tc := range cases {
		r, _ := http.NewRequest(http.MethodGet, req+tc.query, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		require.Equal(t, tc.expCode, w.Code)
		b, _ := json.Marshal(tc.resp)
		require.Equal(t, string(b), w.Body.String())
	}
}

func TestReport(t *testing.T) {
	ctx := context.Background()
	h := gin.New()
	uc := ucmock.NewBalance(t)
	l, _ := logger.New("debug")
	NewRouter(h, uc, l)

	req := "/v1/report"

	uc.On("UpdateReport", ctx, 2022, 10).Return("2022-10.csv", nil)
	uc.On("UpdateReport", ctx, 2000, 1).Return("", entity.ErrEmptyReport)
	uc.On("UpdateReport", ctx, 2000, 2).Return("", errors.New("aboba"))

	type testCases struct {
		name    string
		query   string
		expCode int
		resp    interface{}
	}

	cases := []testCases{{
		name:    "valid",
		query:   "?year=2022&month=10",
		expCode: http.StatusOK,
		resp: struct {
			Link string `json:"link"`
		}{Link: "/v1/reports/2022-10.csv"}, // should be localhost:8080/v1/reports/2022-10.csv
	}, {
		name:    "wrong year",
		query:   "?year=-1&month=10",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "wrong month",
		query:   "?month=-1&year=2000",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Invalid request query"},
	}, {
		name:    "empty report",
		query:   "?month=1&year=2000",
		expCode: http.StatusBadRequest,
		resp:    response{Msg: "Report is empty"},
	}, {
		name:    "db error",
		query:   "?month=2&year=2000",
		expCode: http.StatusInternalServerError,
		resp:    response{Msg: "Database error"},
	},
	}

	for _, tc := range cases {
		r, _ := http.NewRequest(http.MethodGet, req+tc.query, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		require.Equal(t, tc.expCode, w.Code)
		b, _ := json.Marshal(tc.resp)
		require.Equal(t, string(b), w.Body.String())
	}
}
