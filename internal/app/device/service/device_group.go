// ==========================================================================
// 物联网快速开发自动生成业务逻辑层相关代码，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2022-06-21 22:06:17
// 生成路径: github.com/morgeq/iotfast/internal/app/service/device_group.go
// 生成人：dwx
// ==========================================================================

package service

import (
	"context"

	"github.com/morgeq/iotfast/api/v1/device"
	"github.com/morgeq/iotfast/internal/app/device/dao"
	"github.com/morgeq/iotfast/internal/app/device/model/entity"
	systemConsts "github.com/morgeq/iotfast/internal/app/system/consts"
	"github.com/morgeq/iotfast/library/libErr"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

//type deviceGroup struct {
//}
//var DeviceGroup = new(deviceGroup)
type IDeviceGroup interface {
	List(ctx context.Context, req *device.DeviceGroupSearchReq) (total, page int, list []*entity.DeviceGroup, err error)
	Get(ctx context.Context, id int) (info *entity.DeviceGroup, err error)
	Add(ctx context.Context, req *device.DeviceGroupAddReq) (err error)
	Edit(ctx context.Context, req *device.DeviceGroupEditReq) error
	DeleteByIds(ctx context.Context, ids []int) (err error)
}
type deviceGroupImpl struct {
}

var deviceGroupService = deviceGroupImpl{}

func DeviceGroup() IDeviceGroup {
	return &deviceGroupService
}

// List 获取任务列表
func (s *deviceGroupImpl) List(ctx context.Context, req *device.DeviceGroupSearchReq) (total, page int, list []*entity.DeviceGroup, err error) {
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	if req.PageSize == 0 {
		req.PageSize = systemConsts.PageSize
	}
	m := dao.DeviceGroup.Ctx(ctx)
	if req.Name != "" {
		m = m.Where(dao.DeviceGroup.Columns().Name+" like ?", "%"+req.Name+"%")
	}
	if req.BeginTime != "" {
		m = m.Where(dao.DeviceGroup.Columns().CreatedAt+" >= ", req.BeginTime)
	}
	if req.EndTime != "" {
		m = m.Where(dao.DeviceGroup.Columns().CreatedAt+" <= ", req.EndTime)
	}
	if req.Id != 0 {
		m = m.Where(dao.DeviceGroup.Columns().Id+" = ?", req.Id)
	}

	err = g.Try(func() {
		total, err = m.Count()
		libErr.ErrPrint(ctx, err, "获取DeviceGroup列表失败")
		if err != nil {
			g.Log().Error(ctx, err)
			err = gerror.New("获取总行数失败")
			return
		}
		order := "id asc"
		if req.OrderBy != "" {
			order = req.OrderBy
		}
		err = m.Page(page, req.PageSize).Order(order).Scan(&list)
		if err != nil {
			g.Log().Error(ctx, err)
			err = gerror.New("获取数据失败")
		}
	})
	return
}

// Get 通过id获取
func (s *deviceGroupImpl) Get(ctx context.Context, id int) (info *entity.DeviceGroup, err error) {
	if id == 0 {
		err = gerror.New("参数错误")
		return
	}
	err = dao.DeviceGroup.Ctx(ctx).Where(dao.DeviceGroup.Columns().Id, id).Scan(&info)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if info == nil || err != nil {
		err = gerror.New("获取信息失败")
	}
	return
}

// Add 添加
func (s *deviceGroupImpl) Add(ctx context.Context, req *device.DeviceGroupAddReq) (err error) {
	_, err = dao.DeviceGroup.Ctx(ctx).Insert(req)
	return
}

// Edit 修改
func (s *deviceGroupImpl) Edit(ctx context.Context, req *device.DeviceGroupEditReq) error {
	_, err := dao.DeviceGroup.Ctx(ctx).FieldsEx(dao.DeviceGroup.Columns().Id, dao.DeviceGroup.Columns().CreatedAt).Where(dao.DeviceGroup.Columns().Id, req.Id).
		Update(req)
	return err
}

// DeleteByIds 删除
func (s *deviceGroupImpl) DeleteByIds(ctx context.Context, ids []int) (err error) {
	if len(ids) == 0 {
		err = gerror.New("参数错误")
		return
	}
	_, err = dao.DeviceGroup.Ctx(ctx).Delete(dao.DeviceGroup.Columns().Id+" in (?)", ids)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("删除失败")
	}
	return
}
