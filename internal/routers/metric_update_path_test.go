package routers

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
)

func TestRegisterMetricUpdatePathRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRouter := NewMockRouter(ctrl)
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {}
	mockRouter.EXPECT().AddHandler("POST", "/update/:type/:name/:value", gomock.Any())
	RegisterMetricUpdatePathRouter(mockRouter, handler)
}
