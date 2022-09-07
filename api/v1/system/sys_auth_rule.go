/*
* @desc:菜单api
* @url:www.ddsiot.cn
* @Author: dwx<dwxdyx@qq.com>
* @Date:   2022/5/18 10:27
* @Ref: 参考于https://gitee.com/tiger1103/gfast/
 */

package system

import (
	"github.com/morgeq/iotfast/internal/app/system/model"
	"github.com/morgeq/iotfast/internal/app/system/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type RuleSearchReq struct {
	g.Meta        `path:"/menu/list" tags:"菜单管理" method:"get" summary:"菜单列表"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
	Title         string `p:"menuName" `
	Component     string `p:"component"`
}

type RuleListRes struct {
	g.Meta `mime:"application/json"`
	Rules  []*model.SysAuthRuleTreeRes `json:"rules"`
}

type RuleAddReq struct {
	g.Meta        `path:"/menu/add" tags:"菜单管理" method:"post" summary:"添加菜单"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
	MenuType      uint   `p:"menuType"  v:"min:0|max:2#菜单类型最小值为:min|菜单类型最大值为:max"`
	Pid           uint   `p:"parentId"  v:"min:0"`
	Name          string `p:"name" v:"required#请填写规则名称"`
	Title         string `p:"menuName" v:"required|length:1,100#请填写标题|标题长度在:min到:max位"`
	Icon          string `p:"icon"`
	Weigh         int    `p:"menuSort" `
	Condition     string `p:"condition" `
	Remark        string `p:"remark" `
	IsHide        uint   `p:"isHide"`
	Path          string `p:"path"`
	Redirect      string `p:"redirect"` // 路由重定向
	Roles         []uint `p:"roles"`    // 角色ids
	Component     string `p:"component" v:"required-if:menuType,1#组件路径不能为空"`
	IsLink        uint   `p:"isLink"`
	IsIframe      uint   `p:"isIframe"`
	IsCached      uint   `p:"isKeepAlive"`
	IsAffix       uint   `p:"isAffix"`
	LinkUrl       string `p:"linkUrl"`
}

type RuleAddRes struct {
}

type RuleGetParamsReq struct {
	g.Meta        `path:"/menu/getParams" tags:"菜单管理" method:"get" summary:"获取添加、编辑菜单相关参数"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
}

type RuleGetParamsRes struct {
	g.Meta `mime:"application/json"`
	Roles  []*entity.SysRole           `json:"roles"`
	Menus  []*model.SysAuthRuleInfoRes `json:"menus"`
}

type RuleInfoReq struct {
	g.Meta        `path:"/menu/get" tags:"菜单管理" method:"get" summary:"获取菜单信息"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
	Id            uint   `p:"required#菜单id必须"`
}

type RuleInfoRes struct {
	g.Meta  `mime:"application/json"`
	Rule    *entity.SysAuthRule `json:"rule"`
	RoleIds []uint              `json:"roleIds"`
}

type RuleUpdateReq struct {
	g.Meta        `path:"/menu/update" tags:"菜单管理" method:"put" summary:"修改菜单"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
	Id            uint   `p:"id" v:"required#id必须"`
	MenuType      uint   `p:"menuType"  v:"min:0|max:2#菜单类型最小值为:min|菜单类型最大值为:max"`
	Pid           uint   `p:"parentId"  v:"min:0"`
	Name          string `p:"name" v:"required#请填写规则名称"`
	Title         string `p:"menuName" v:"required|length:1,100#请填写标题|标题长度在:min到:max位"`
	Icon          string `p:"icon"`
	Weigh         int    `p:"menuSort" `
	Condition     string `p:"condition" `
	Remark        string `p:"remark" `
	IsHide        uint   `p:"isHide"`
	Path          string `p:"path"`
	Redirect      string `p:"redirect"` // 路由重定向
	Roles         []uint `p:"roles"`    // 角色ids
	Component     string `p:"component" v:"required-if:menuType,1#组件路径不能为空"`
	IsLink        uint   `p:"isLink"`
	IsIframe      uint   `p:"isIframe"`
	IsCached      uint   `p:"isKeepAlive"`
	IsAffix       uint   `p:"isAffix"`
	LinkUrl       string `p:"linkUrl"`
}

type RuleUpdateRes struct {
}

type RuleDeleteReq struct {
	g.Meta        `path:"/menu/delete" tags:"菜单管理" method:"delete" summary:"删除菜单"`
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
	Ids           []int  `p:"ids" v:"required#菜单id必须"`
}

type RuleDeleteRes struct {
}
