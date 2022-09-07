/*
* @desc:字典数据
* @url:www.ddsiot.cn
* @Author: dwx<dwxdyx@qq.com>
* @Date:   2022/5/18 11:55
* @Ref: 参考于https://gitee.com/tiger1103/gfast/
 */

package service

import (
	"context"

	"github.com/morgeq/iotfast/api/v1/system"
	commonService "github.com/morgeq/iotfast/internal/app/common/service"
	"github.com/morgeq/iotfast/internal/app/system/consts"
	systemConsts "github.com/morgeq/iotfast/internal/app/system/consts"
	"github.com/morgeq/iotfast/internal/app/system/dao"
	"github.com/morgeq/iotfast/internal/app/system/model"
	"github.com/morgeq/iotfast/internal/app/system/model/do"
	"github.com/morgeq/iotfast/library/libErr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type IDictData interface {
	GetDictWithDataByType(ctx context.Context, req *system.GetDictReq) (dict *system.GetDictRes, err error)
	List(ctx context.Context, req *system.DictDataSearchReq) (res *system.DictDataSearchRes, err error)
	Add(ctx context.Context, req *system.DictDataAddReq, userId uint64) (err error)
	Get(ctx context.Context, dictCode uint) (res *system.DictDataGetRes, err error)
	Edit(ctx context.Context, req *system.DictDataEditReq, userId uint64) (err error)
	Delete(ctx context.Context, ids []int) (err error)
}

type dictDataImpl struct {
}

var dictData = dictDataImpl{}

func DictData() IDictData {
	return &dictData
}

// GetDictWithDataByType 通过字典键类型获取选项
func (s dictDataImpl) GetDictWithDataByType(ctx context.Context, req *system.GetDictReq) (dict *system.GetDictRes,
	err error) {
	cache := commonService.Cache()
	cacheKey := consts.CacheSysDict + "_" + req.DictType
	//从缓存获取
	iDict := cache.GetOrSetFuncLock(ctx, cacheKey, func(ctx context.Context) (value interface{}, err error) {
		err = g.Try(func() {
			//从数据库获取
			dict = &system.GetDictRes{}
			//获取类型数据
			err = dao.SysDictType.Ctx(ctx).Where(dao.SysDictType.Columns().DictType, req.DictType).
				Where(dao.SysDictType.Columns().Status, 1).Fields(model.DictTypeRes{}).Scan(&dict.Info)
			libErr.ErrPrint(ctx, err, "获取字典类型失败")
			err = dao.SysDictData.Ctx(ctx).Fields(model.DictDataRes{}).
				Where(dao.SysDictData.Columns().DictType, req.DictType).
				Order(dao.SysDictData.Columns().DictSort + " asc," +
					dao.SysDictData.Columns().DictCode + " asc").
				Scan(&dict.Values)
			libErr.ErrPrint(ctx, err, "获取字典数据失败")
		})
		value = dict
		return
	}, 0, consts.CacheSysDictTag)
	if iDict != nil {
		err = gconv.Struct(iDict, &dict)
		if err != nil {
			return
		}
	}
	//设置给定的默认值
	for _, v := range dict.Values {
		if req.DefaultValue != "" {
			if gstr.Equal(req.DefaultValue, v.DictValue) {
				v.IsDefault = 1
			} else {
				v.IsDefault = 0
			}
		}
	}
	return
}

// List 获取字典数据
func (s dictDataImpl) List(ctx context.Context, req *system.DictDataSearchReq) (res *system.DictDataSearchRes, err error) {
	res = new(system.DictDataSearchRes)
	err = g.Try(func() {
		m := dao.SysDictData.Ctx(ctx)
		if req != nil {
			if req.DictLabel != "" {
				m = m.Where(dao.SysDictData.Columns().DictLabel+" like ?", "%"+req.DictLabel+"%")
			}
			if req.Status != "" {
				m = m.Where(dao.SysDictData.Columns().Status+" = ", gconv.Int(req.Status))
			}
			if req.DictType != "" {
				m = m.Where(dao.SysDictData.Columns().DictType+" = ?", req.DictType)
			}
			res.Total, err = m.Count()
			libErr.ErrPrint(ctx, err, "获取字典数据失败")
			if req.PageNum == 0 {
				req.PageNum = 1
			}
			res.CurrentPage = req.PageNum
		}
		if req.PageSize == 0 {
			req.PageSize = systemConsts.PageSize
		}
		err = m.Page(req.PageNum, req.PageSize).Order(dao.SysDictData.Columns().DictSort + " asc," +
			dao.SysDictData.Columns().DictCode + " asc").Scan(&res.List)
		libErr.ErrPrint(ctx, err, "获取字典数据失败")
	})
	return
}

func (s *dictDataImpl) Add(ctx context.Context, req *system.DictDataAddReq, userId uint64) (err error) {
	err = g.Try(func() {
		_, err = dao.SysDictData.Ctx(ctx).Insert(do.SysDictData{
			DictSort:  req.DictSort,
			DictLabel: req.DictLabel,
			DictValue: req.DictValue,
			DictType:  req.DictType,
			CssClass:  req.CssClass,
			ListClass: req.ListClass,
			IsDefault: req.IsDefault,
			Status:    req.Status,
			CreateBy:  userId,
			Remark:    req.Remark,
		})
		libErr.ErrPrint(ctx, err, "添加字典数据失败")
		//清除缓存
		commonService.Cache().RemoveByTag(ctx, consts.CacheSysDictTag)
	})
	return
}

// Get 获取字典数据
func (s *dictDataImpl) Get(ctx context.Context, dictCode uint) (res *system.DictDataGetRes, err error) {
	res = new(system.DictDataGetRes)
	err = g.Try(func() {
		err = dao.SysDictData.Ctx(ctx).WherePri(dictCode).Scan(&res.Dict)
		libErr.ErrPrint(ctx, err, "获取字典数据失败")
	})
	return
}

// Edit 修改字典数据
func (s *dictDataImpl) Edit(ctx context.Context, req *system.DictDataEditReq, userId uint64) (err error) {
	err = g.Try(func() {
		_, err = dao.SysDictData.Ctx(ctx).WherePri(req.DictCode).Update(do.SysDictData{
			DictSort:  req.DictSort,
			DictLabel: req.DictLabel,
			DictValue: req.DictValue,
			DictType:  req.DictType,
			CssClass:  req.CssClass,
			ListClass: req.ListClass,
			IsDefault: req.IsDefault,
			Status:    req.Status,
			UpdateBy:  userId,
			Remark:    req.Remark,
		})
		libErr.ErrPrint(ctx, err, "修改字典数据失败")
		//清除缓存
		commonService.Cache().RemoveByTag(ctx, consts.CacheSysDictTag)
	})
	return
}

// Delete 删除字典数据
func (s *dictDataImpl) Delete(ctx context.Context, ids []int) (err error) {
	err = g.Try(func() {
		_, err = dao.SysDictData.Ctx(ctx).Where(dao.SysDictData.Columns().DictCode+" in(?)", ids).Delete()
		libErr.ErrPrint(ctx, err, "删除字典数据失败")
		//清除缓存
		commonService.Cache().RemoveByTag(ctx, consts.CacheSysDictTag)
	})
	return
}
