/*
* @desc:后台路由
* @url:www.ddsiot.cn
* @Author: dwx
* @Date:   2022/5/18 17:34
 */

package router

import (
	commonService "github.com/morgeq/iotfast/internal/app/common/service"
	"github.com/morgeq/iotfast/internal/app/system/service"
	controller "github.com/morgeq/iotfast/server/http/controller/system"

	"github.com/gogf/gf/v2/net/ghttp"
)

func BindSystemController(group *ghttp.RouterGroup) {
	group.Group("/system", func(group *ghttp.RouterGroup) {
		group.Middleware(commonService.Middleware().MiddlewareCORS)
		// 系统初始化
		group.Bind(
			controller.DbInit,
		)
		group.Bind(
			//登录
			controller.Login,
		)
		//登录验证拦截
		service.GfToken().Middleware(group)
		//context拦截器
		group.Middleware(service.Middleware().Ctx, service.Middleware().Auth)
		group.Bind(
			controller.User,
			controller.Menu,
			controller.Role,
			controller.Dept,
			controller.Post,
			controller.DictType,
			controller.DictData,
			controller.Config,
			controller.Monitor,
			controller.LoginLog,
			controller.Gen,
			controller.SysJob,
		)
	})
}
