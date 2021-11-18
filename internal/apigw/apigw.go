package apigw

import (
	"net/http"

	"github.com/morzhanov/async-api/api/order"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/async-api/internal/rest"
	"go.uber.org/zap"
)

type controller struct {
	rest.BaseController
	client Client
}

type Controller interface {
	Listen(port string)
}

func (c *controller) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.BaseController.Logger().Info("error in the REST handler", zap.Error(err))
}

func (c *controller) handleCreateOrder(ctx *gin.Context) {
	d := order.CreateOrderMessage{}
	if err := c.BaseController.ParseRestBody(ctx, &d); err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	err := c.client.CreateOrder(ctx, &d)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *controller) handleProcessOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.client.ProcessOrder(ctx, id)
	if err != nil {
		c.handleHttpErr(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *controller) Listen(port string) {
	c.BaseController.Listen(port)
}

func NewController(
	client Client,
	log *zap.Logger,
) Controller {
	bc := rest.NewBaseController(log)
	c := controller{BaseController: bc, client: client}
	r := bc.Router()
	r.POST("/order", bc.Handler(c.handleCreateOrder))
	r.PUT("/order/:id", bc.Handler(c.handleProcessOrder))
	return &c
}
