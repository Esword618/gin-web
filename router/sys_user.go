package router

import (
	v1 "gin-web/api/v1"
	"github.com/piupuer/go-helper/router"
)

func InitUserRouter(r *router.Router) {
	router1 := r.Casbin("/user")
	router2 := r.CasbinAndIdempotence("/user")
	router1.POST("/info", v1.GetUserInfo)
	router1.GET("/list", v1.FindUser)
	router1.GET("/list/:ids", v1.FindUserByIds)
	router2.POST("/create", v1.CreateUser)
	router1.PATCH("/update/:id", v1.UpdateUserById)
	router1.DELETE("/delete/batch", v1.BatchDeleteUserByIds)
}
