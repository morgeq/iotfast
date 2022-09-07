/*
* @desc:中间件处理
* @url:www.ddsiot.cn
* @Author: dwx<dwxdyx@qq.com>
* @Date:   2022/5/17 9:11
 */

package service

import "github.com/gogf/gf/v2/net/ghttp"

type IMiddleware interface {
	MiddlewareCORS(r *ghttp.Request)
}

type middlewareImpl struct{}

var middleService = middlewareImpl{}

func Middleware() IMiddleware {
	return &middleService
}

func (s *middlewareImpl) MiddlewareCORS(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	// you can set options
	//corsOptions.AllowDomain = []string{"goframe.org", "baidu.com"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}
