package v1

import (
	mw "balance_api/internal/controller/http/v1/middleware"
	"balance_api/internal/entity"
	"balance_api/internal/usecase"
	"balance_api/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

type balanceRouters struct {
	b usecase.Balance
	l logger.Interface
}

func newBalanceRoutes(handler *gin.RouterGroup, b usecase.Balance, l logger.Interface) {
	r := &balanceRouters{
		b: b,
		l: l,
	}

	handler.GET("/user", mw.ValidateQuery[userGetRequest](r.l), r.getByID)
	handler.POST("/user", mw.ValidateJSONBody[userPostRequest](r.l), r.increaseAmount)
	handler.POST("/order", mw.ValidateJSONBody[orderPostRequest](r.l), r.order)
	handler.GET("/history", mw.ValidateQuery[historyGetRequest](r.l), r.getHistory)
	handler.GET("/report", mw.ValidateQuery[reportGetRequest](r.l), r.createReport)
	handler.GET("/reports/:name", func(c *gin.Context) {
		name := c.Param("name")
		dir := r.b.GetReportDir()
		c.FileAttachment(dir+name, name)
	})
}

type userGetRequest struct {
	ID int `form:"id" binding:"required,gte=1"`
}

func (r *balanceRouters) getByID(c *gin.Context) {
	q := mw.GetQueryParams[userGetRequest](c)
	balance, err := r.b.GetByID(c.Request.Context(), q.ID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("err \"%s\" with request params: %v", err, q)
		errorResponse(c, http.StatusBadRequest, "No such id")
		return
	case err != nil:
		r.l.Error(err)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, balance)
}

type userPostRequest struct {
	ID     int    `json:"id" binding:"required,gte=1"`
	Amount string `json:"amount" binding:"required"`
}

func (r *balanceRouters) increaseAmount(c *gin.Context) {
	b := mw.GetJSONBody[userPostRequest](c)
	num, err := decimal.NewFromString(b.Amount)
	if err != nil || !num.IsPositive() {
		r.l.Infof("err \"%s\" with request params: %v", err, b)
		errorResponse(c, http.StatusBadRequest, "Invalid money format")
		return
	}
	err = r.b.Increase(c.Request.Context(), entity.Balance{ID: b.ID, Amount: b.Amount})
	if err != nil {
		r.l.Error(err)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}

type orderPostRequest struct {
	Action    string `json:"action" binding:"required"`
	ID        int    `json:"order_id" binding:"required,gte=1"`
	ServiceID int    `json:"service_id" binding:"required,gte=1"`
	UserID    int    `json:"user_id" binding:"required,gte=1"`
	Sum       string `json:"sum" binding:"required"`
}

func (r *balanceRouters) order(c *gin.Context) {
	b := mw.GetJSONBody[orderPostRequest](c)
	num, err := decimal.NewFromString(b.Sum)
	if err != nil || !num.IsPositive() {
		r.l.Infof("err \"%s\" with request params: %v", err, b)
		errorResponse(c, http.StatusBadRequest, "Invalid money format")
		return
	}
	switch b.Action {
	case "create":
		err = r.b.CreateOrder(c.Request.Context(),
			entity.Order{ID: b.ID, ServiceID: b.ServiceID, UserID: b.UserID, Sum: b.Sum})
	case "approve":
		err = r.b.ChangeOrderStatus(c.Request.Context(),
			entity.Order{ID: b.ID, ServiceID: b.ServiceID, UserID: b.UserID, Sum: b.Sum, StatusID: 2})
	case "cancel":
		err = r.b.ChangeOrderStatus(c.Request.Context(),
			entity.Order{ID: b.ID, ServiceID: b.ServiceID, UserID: b.UserID, Sum: b.Sum, StatusID: 3})
	default:
		r.l.Infof("err \"wrong order action\" with request params: %v", b)
		errorResponse(c, http.StatusBadRequest, "Invalid order action")
		return
	}
	errMsg := ""
	switch {
	case errors.Is(err, entity.ErrNoService):
		errMsg = "No such service"
	case errors.Is(err, entity.ErrNoID):
		errMsg = "No such id"
	case errors.Is(err, entity.ErrNotEnoughMoney):
		errMsg = "Not enough money"
	case errors.Is(err, entity.ErrOrderExists):
		errMsg = "Order already exists"
	case errors.Is(err, entity.ErrOrderNoExists):
		errMsg = "Order not exists"
	case errors.Is(err, entity.ErrOrderMismatch):
		errMsg = "Wrong order data"
	case errors.Is(err, entity.ErrCantChangeStatus):
		errMsg = "Order already approved/canceled"
	case err != nil:
		r.l.Error(err)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	if err != nil {
		r.l.Infof("err \"%s\" with request params: %v", err, b)
		errorResponse(c, http.StatusBadRequest, errMsg)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}

type historyGetRequest struct {
	ID      int    `form:"id" binding:"required,gte=1"`
	Limit   int    `form:"limit" binding:"omitempty,gte=0,lte=200"`
	Page    int    `form:"page" binding:"omitempty,gte=1"`
	Desc    bool   `form:"desc" binding:"omitempty"`
	OrderBy string `form:"order_by" binding:"omitempty"`
}

func (r *balanceRouters) getHistory(c *gin.Context) {
	q := mw.GetQueryParams[historyGetRequest](c)
	q, msg := setHistoryParams(q)
	if msg != "" {
		r.l.Infof("err \"%s\" with request params: %v", msg, q)
		errorResponse(c, http.StatusBadRequest, msg)
		return
	}
	h, err := r.b.GetHistory(c.Request.Context(),
		entity.History{UserID: q.ID, Limit: q.Limit, OrderBy: q.OrderBy, Desc: q.Desc, Page: q.Page})
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("err \"%s\" with request params: %v", err, q)
		errorResponse(c, http.StatusBadRequest, "No such id")
		return
	case errors.Is(err, entity.ErrEmptyPage):
		r.l.Infof("err \"%s\" with request params: %v", err, q)
		errorResponse(c, http.StatusBadRequest, "The page is empty")
		return
	case err != nil:
		r.l.Error(err)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, h)
}

func setHistoryParams(h historyGetRequest) (historyGetRequest, string) {
	if (h.Limit == 0 && h.Page != 0) || (h.Limit != 0 && h.Page == 0) {
		return historyGetRequest{}, "Limit and page should be both zero or non zero"
	}
	switch h.OrderBy {
	case "":
		h.OrderBy = "date"
	case "date":
	case "sum":
	default:
		return historyGetRequest{}, "Wrong order by param"
	}
	return h, ""
}

type reportGetRequest struct {
	Year  int `form:"year" binding:"required,gte=1900"`
	Month int `form:"month" binding:"required,gte=1,lte=12"`
}

func (r *balanceRouters) createReport(c *gin.Context) {
	q := mw.GetQueryParams[reportGetRequest](c)
	name, err := r.b.UpdateReport(c.Request.Context(), q.Year, q.Month)
	switch {
	case errors.Is(err, entity.ErrEmptyReport):
		r.l.Infof("err \"%s\" with request params: %v", err, q)
		errorResponse(c, http.StatusBadRequest, "Report is empty")
		return
	case err != nil:
		r.l.Error(err)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, struct {
		Link string `json:"link"`
	}{Link: c.Request.Host + c.Request.URL.Path + "s/" + name})
}
