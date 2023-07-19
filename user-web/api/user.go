package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	ip := "127.0.0.1"
	port := 50051
	//拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 链接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败 ")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	//claims, _ := ctx.Get("claims")
	//currentUser := claims.(*models.CustomClaims)
	//zap.S().Infof("访问用户: %d", currentUser.ID)
	////生成grpc的client并调用接口
	//
	//pn := ctx.DefaultQuery("pn", "0")
	//pnInt, _ := strconv.Atoi(pn)
	//pSize := ctx.DefaultQuery("psize", "10")
	//pSizeInt, _ := strconv.Atoi(pSize)
	//
	//rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
	//	Pn:    uint32(pnInt),
	//	PSize: uint32(pSizeInt),
	//})
	//if err != nil {
	//	zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
	//	HandleGrpcErrorToHttp(err, ctx)
	//	return
	//}
	//
	//reMap := gin.H{
	//	"total": rsp.Total,
	//}
	//result := make([]interface{}, 0)
	//for _, value := range rsp.Data {
	//	user := reponse.UserResponse{
	//		Id:       value.Id,
	//		NickName: value.NickName,
	//		//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
	//		Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
	//		Gender:   value.Gender,
	//		Mobile:   value.Mobile,
	//	}
	//	result = append(result, user)
	//}
	//
	//reMap["data"] = result
	//ctx.JSON(http.StatusOK, reMap)
}
