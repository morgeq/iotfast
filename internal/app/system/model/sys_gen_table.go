package model

import (
	"github.com/xiaodingding/iotfast/internal/app/system/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

//实体扩展
type SysGenTableEntityExtend struct {
	TableId        int64                       `orm:"table_id,primary" json:"table_id"`        // 编号
	TableName      string                      `orm:"table_name"       json:"table_name"`      // 表名称
	TableComment   string                      `orm:"table_comment"    json:"table_comment"`   // 表描述
	ClassName      string                      `orm:"class_name"       json:"class_name"`      // 实体类名称
	TplCategory    string                      `orm:"tpl_category"     json:"tpl_category"`    // 使用的模板（crud单表操作 tree树表操作）
	PackageName    string                      `orm:"package_name"     json:"package_name"`    // 生成包路径
	ModuleName     string                      `orm:"module_name"      json:"module_name"`     // 生成模块名
	BusinessName   string                      `orm:"business_name"    json:"business_name"`   // 生成业务名
	FunctionName   string                      `orm:"function_name"    json:"function_name"`   // 生成功能名
	FunctionAuthor string                      `orm:"function_author"  json:"function_author"` // 生成功能作者
	Options        string                      `orm:"options"          json:"options"`         // 其它生成选项
	CreateBy       string                      `orm:"create_by"        json:"create_by"`       // 创建者
	CreateTime     *gtime.Time                 `orm:"create_time"      json:"create_time"`     // 创建时间
	UpdateBy       string                      `orm:"update_by"        json:"update_by"`       // 更新者
	UpdateTime     *gtime.Time                 `orm:"update_time"      json:"update_time"`     // 更新时间
	Remark         string                      `orm:"remark"           json:"remark"`          // 备注
	TreeCode       string                      `json:"tree_code"`                              // 树编码字段
	TreeParentCode string                      `json:"tree_parent_code"`                       // 树父编码字段
	TreeName       string                      `json:"tree_name"`                              // 树名称字段
	Columns        []*entity.SysGenTableColumn `json:"columns"`                                // 表列信息
	PkColumn       *entity.SysGenTableColumn   `json:"pkColumn"`                               // 表列信息
}
