// ==========================================================================
// 物联网快速开发自动生成model代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2022-06-18 13:58:41
// 生成路径: github.com/morgeq/iotfast/internal/app/system/model/entity/sys_job.go
// 生成人：dwx
// ==========================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// SysJob is the golang structure for table sys_job.
type SysJob struct {    
         JobId       int64         `orm:"job_id,primary" json:"jobId"`    // 任务ID    
         JobName       string         `orm:"job_name,primary" json:"jobName"`    // 任务名称    
         JobParams    string         `orm:"job_params" json:"jobParams"`    // 参数    
         JobGroup       string         `orm:"job_group,primary" json:"jobGroup"`    // 任务组名    
         InvokeTarget    string         `orm:"invoke_target" json:"invokeTarget"`    // 调用目标字符串    
         CronExpression    string         `orm:"cron_expression" json:"cronExpression"`    // cron执行表达式    
         MisfirePolicy    int         `orm:"misfire_policy" json:"misfirePolicy"`    // 计划执行策略（1多次执行 2执行一次）    
         Concurrent    int         `orm:"concurrent" json:"concurrent"`    // 是否并发执行（0允许 1禁止）    
         Status    int         `orm:"status" json:"status"`    // 状态（0正常 1暂停）    
         CreateBy    uint64         `orm:"create_by" json:"createBy"`    // 创建者    
         UpdateBy    uint64         `orm:"update_by" json:"updateBy"`    // 更新者    
         Remark    string         `orm:"remark" json:"remark"`    // 备注信息    
         CreatedAt    *gtime.Time         `orm:"created_at" json:"createdAt"`    // 创建时间    
         UpdatedAt    *gtime.Time         `orm:"updated_at" json:"updatedAt"`    // 更新时间    
         DeletedAt    *gtime.Time         `orm:"deleted_at" json:"deletedAt"`    // 删除时间    
}