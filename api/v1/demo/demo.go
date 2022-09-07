/*
* @desc:demo
* @url:www.ddsiot.cn
* @Author: dwx
* @Date:   2022/5/2 15:25
 */

package demo

import "github.com/gogf/gf/v2/frame/g"

type DmReq struct {
	g.Meta `path:"/demo" tags:"Demo" method:"post" summary:"demo api"`
	Name   string `p:"name" v:"required#标签内容不能为空"`
}
type DmRes struct {
	g.Meta `mime:"text/html"`
	Name   string `json:"name"`
}
