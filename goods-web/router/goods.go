package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/goods"
	"mxshop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                            //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)          //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                      //商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                               //获取商品库存
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
	//服务注册和发现
}
