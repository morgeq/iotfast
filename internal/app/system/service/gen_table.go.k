package service

import (
	"context"
	"iotfast/api/v1/system"
	"iotfast/internal/app/system/consts"
	"iotfast/internal/app/system/dao"
	"iotfast/internal/app/system/model"
	"iotfast/internal/app/system/model/entity"
	"iotfast/library/libErr"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

type IGenTable interface {
	GetDbTable(ctx context.Context, req *system.GenDbSearchReq) (res *system.GenTableSearchRes, err error)
	ImportGenTable(ctx context.Context, req *system.GenTableImportReq) (err error)
	TableList(ctx context.Context, req *system.GenTableSearchReq) (res *system.GenTableSearchRes, err error)
	//TableDelete(ctx context.Context, ids []int) (err error)
	ColumnList(ctx context.Context, req *system.GenTableColumnSearchReq) (res *system.GenTableColumnSearchRes, err error)
	Edit(ctx context.Context, req *system.GenTableEditReq) (err error)
	Delete(ctx context.Context, ids []int) error
	Preview(ctx context.Context, tableId int) (data g.MapStrStr, err error)
	BatchGenCode(ctx context.Context, tableIds []int, r *ghttp.Request) (err error)
}

var genTableService = genTableImpl{}

func GenTable() IGenTable {
	return &genTableService
}

type genTableImpl struct {
}

func (s *genTableImpl) GetDbTable(ctx context.Context, req *system.GenDbSearchReq) (res *system.GenTableSearchRes, err error) {
	db := g.DB()
	res = new(system.GenTableSearchRes)
	sql := " from information_schema.tables where table_schema = (select database())" +
		" and table_name NOT LIKE 'qrtz_%' AND table_name NOT LIKE 'gen_%' and table_name NOT IN (select table_name from gen_table) "
	if req != nil {
		if req.TableName != "" {
			sql += gdb.FormatSqlWithArgs(" and lower(table_name) like lower(?)", []interface{}{"%" + req.TableName + "%"})
		}

		if req.TableComment != "" {
			sql += gdb.FormatSqlWithArgs(" and lower(table_comment) like lower(?)", []interface{}{"%" + req.TableComment + "%"})
		}
	}

	countSql := "select count(1) " + sql
	res.Total, err = db.GetCount(ctx, countSql)
	if err != nil {
		libErr.ErrPrint(ctx, err, "?????????????????????")
		err = gerror.New("?????????????????????")
		return
	}

	sql = "table_name, table_comment, create_time, update_time " + sql

	if req.PageNum == 0 {
		req.PageNum = 1
	}

	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	page := (req.PageNum - 1) * req.PageSize
	sql += " order by create_time desc,table_name asc limit  " + gconv.String(page) + "," + gconv.String(req.PageSize)

	err = g.Try(func() {
		var result gdb.Result
		result, err = db.GetAll(ctx, "select "+sql)
		result.Structs(&res.List)
		libErr.ErrPrint(ctx, err, "???????????????????????????????????????")
	})

	return
}

//???????????????
func (s *genTableImpl) ImportGenTable(ctx context.Context, req *system.GenTableImportReq) (err error) {
	if req != nil && len(req.Tables) > 0 {
		tx, err := dao.GenTable.DB().Begin(ctx)
		if err != nil {
			return err
		}

		tableList, err := SelectDbTableListByNames(ctx, req.Tables)
		if err != nil {
			return err
		}
		operName := Context().GetLoginUser(ctx).UserName

		for _, table := range tableList {
			tableName := table.TableName
			InitTable(ctx, table, operName)
			result, err := tx.Insert(dao.GenTable.Table(), table)
			if err != nil {
				return err
			}

			tmpid, err := result.LastInsertId()

			if err != nil || tmpid <= 0 {
				tx.Rollback()
				return gerror.New("??????????????????")
			}

			table.TableId = tmpid

			// ???????????????
			genTableColumns, err := SelectDbTableColumnsByName(ctx, tableName)

			if err != nil || len(genTableColumns) <= 0 {
				tx.Rollback()
				return gerror.New("?????????????????????")
			}
			for _, column := range genTableColumns {
				InitColumnField(column, table)
				_, err = tx.Insert(dao.GenTableColumn.Table(), column)
				if err != nil {
					tx.Rollback()
					return gerror.New("?????????????????????")
				}
			}
		}
		return tx.Commit()
	} else {
		return gerror.New("????????????")
	}
}

func (s *genTableImpl) TableList(ctx context.Context, req *system.GenTableSearchReq) (res *system.GenTableSearchRes, err error) {
	res = new(system.GenTableSearchRes)
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	res.CurrentPage = req.PageNum
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	m := dao.GenTable.Ctx(ctx)

	if req.TableName != "" {
		m = m.Where("table_name like ?", "%"+req.TableName+"%")
	}
	if req.TableComment != "" {
		m = m.Where("table_comment like ?", "%"+req.TableComment+"%")
	}

	err = g.Try(func() {
		res.Total, err = m.Count()
		libErr.ErrPrint(ctx, err, "?????????????????????????????????")
		err = m.Page(req.PageNum, req.PageSize).Order("id asc").Scan(&res.List)
		libErr.ErrPrint(ctx, err, "???????????????????????????????????????")
	})
	return
}

// DeleteByIds ?????????
/*
func (s *genTableImpl) TableDelete(ctx context.Context, ids []int) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			_, err = dao.GenTable.Ctx(ctx).TX(tx).Where(dao.GenTable.Columns().TableId+" in(?)", ids).Delete()
			libErr.ErrPrint(ctx, err, "???????????????")
		})
		return err
	})
	return
}
*/

func (s *genTableImpl) ColumnList(ctx context.Context, req *system.GenTableColumnSearchReq) (res *system.GenTableColumnSearchRes, err error) {
	res = new(system.GenTableColumnSearchRes)
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	res.CurrentPage = req.PageNum
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}

	if req.TableId == 0 || req == nil {
		return nil, gerror.New("????????????")
	}
	res.Rows, err = SelectGenTableColumnListByTableId(ctx, req.TableId)
	if err != nil {
		return nil, err
	}

	var tableInfo entity.GenTable
	err = dao.GenTable.Ctx(ctx).WherePri(req.TableId).Scan(&tableInfo)
	res.Info = gconv.Map(tableInfo)
	//????????????????????????????????????
	if res.Info != nil && tableInfo.TplCategory == "tree" {
		options := gjson.New(tableInfo.Options)
		res.Info["tree_code"] = options.Get("tree_code")
		res.Info["tree_parent_code"] = options.Get("tree_parent_code")
		res.Info["tree_name"] = options.Get("tree_name")
	}

	return
}

//?????????????????????
func (s *genTableImpl) Edit(ctx context.Context, req *system.GenTableEditReq) (err error) {
	if req == nil {
		err = gerror.New("????????????")
		return
	}
	table := &entity.GenTable{}
	err = dao.GenTable.Ctx(ctx).WherePri(req.TableId).Scan(table)
	if err != nil || table == nil {
		err = gerror.New("???????????????")
		return
	}
	if req.TableName != "" {
		table.TableName = req.TableName
	}
	if req.TableComment != "" {
		table.TableComment = req.TableComment
	}
	if req.BusinessName != "" {
		table.BusinessName = req.BusinessName
	}
	if req.ClassName != "" {
		table.ClassName = req.ClassName
	}
	if req.FunctionAuthor != "" {
		table.FunctionAuthor = req.FunctionAuthor
	}
	if req.FunctionName != "" {
		table.FunctionName = req.FunctionName
	}
	if req.ModuleName != "" {
		table.ModuleName = req.ModuleName
	}
	if req.PackageName != "" {
		table.PackageName = req.PackageName
	}
	if req.Remark != "" {
		table.Remark = req.Remark
	}
	if req.TplCategory != "" {
		table.TplCategory = req.TplCategory
	}
	if req.Params != "" {
		table.Options = req.Params
	}
	table.UpdateTime = gtime.Now()
	table.UpdateBy = Context().GetLoginUser(ctx).UserName
	if req.TplCategory == "tree" {
		//????????????options
		options := g.Map{
			"tree_code":        req.TreeCode,
			"tree_parent_code": req.TreeParentCode,
			"tree_name":        req.TreeName,
		}
		table.Options = gconv.String(options)
	} else {
		table.Options = ""
	}

	var tx *gdb.TX
	tx, err = dao.GenTable.DB().Begin(ctx)
	if err != nil {
		return
	}
	_, err = tx.Save(dao.GenTable.Table(), table)
	if err != nil {
		tx.Rollback()
		return err
	}

	//???????????????
	if req.Columns != "" {
		var j *gjson.Json
		if j, err = gjson.DecodeToJson([]byte(req.Columns)); err != nil {
			tx.Rollback()
			return
		} else {
			var columnList []entity.GenTableColumn
			err = j.Scan(&columnList)
			if err == nil && columnList != nil && len(columnList) > 0 {
				for _, column := range columnList {
					if column.ColumnId > 0 {
						tmp := &entity.GenTableColumn{}
						_ = dao.GenTableColumn.Ctx(ctx).WherePri("column_id=?", column.ColumnId).Scan(&tmp)
						if tmp != nil {
							tmp.ColumnComment = column.ColumnComment
							tmp.GoType = column.GoType
							tmp.HtmlType = column.HtmlType
							tmp.QueryType = column.QueryType
							tmp.GoField = column.GoField
							tmp.DictType = column.DictType
							tmp.IsInsert = column.IsInsert
							tmp.IsEdit = column.IsEdit
							tmp.IsList = column.IsList
							tmp.IsQuery = column.IsQuery
							_, err = tx.Save(dao.GenTableColumn.Table(), tmp)
							if err != nil {
								tx.Rollback()
								return
							}
						}
					}
				}
			}
		}
	}
	tx.Commit()
	return
}

//????????????
func (s *genTableImpl) Delete(ctx context.Context, ids []int) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			_, err = dao.GenTable.Ctx(ctx).TX(tx).Where(dao.GenTable.Columns().TableId+" in(?)", ids).Delete()
			libErr.ErrPrint(ctx, err, "???????????????")

			_, err = dao.GenTableColumn.Ctx(ctx).TX(tx).Where(dao.GenTableColumn.Columns().TableId+" in(?)", ids).Delete()
			libErr.ErrPrint(ctx, err, "??????????????????????????????")
		})
		return err
	})
	return
}

func (s *genTableImpl) Preview(ctx context.Context, tableId int) (data g.MapStrStr, err error) {
	//g.Log().Printf(ctx, "tableId:%v", tableId)
	//err = g.Try(func() {
	data, _, err = s.genCode(ctx, tableId)
	//})
	//g.Log().Printf(ctx, "tableId:%v, data:%v", tableId, data)
	return
}

//?????????????????????
func (s *genTableImpl) BatchGenCode(ctx context.Context, tableIds []int, r *ghttp.Request) (err error) {

	if len(tableIds) < 1 {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		err = gerror.New("????????????????????????")
		return err
	}

	dataFilePath := g.Cfg().MustGet(ctx, "adminInfo.dataDir").String()
	dataFileRange := grand.S(10)
	//????????????
	for _, id := range tableIds {
		data, entity, err := s.genCode(ctx, id)
		if err != nil {
			libErr.ErrPrint(ctx, err, "??????????????????")
			err = gerror.New("??????????????????")
			return err
		}
		pathMap := s.getPath(entity)
		for key, val := range data {
			switch key {
			case "template/go/" + entity.BusinessName + "_controller.go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/"+pathMap["controller"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/go/" + entity.BusinessName + "_service.go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/"+pathMap["service"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/go/" + entity.BusinessName + "_model.go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/"+pathMap["gfgenModel"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/go/" + entity.BusinessName + "_entity.go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/"+pathMap["gfgenEntity"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/go/" + entity.BusinessName + ".go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/"+pathMap["model"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/go/" + entity.BusinessName + "_route.go.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/go/????????????.txt", val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/html/" + entity.BusinessName + "_api.js.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/vue/"+pathMap["api"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			case "template/html/" + entity.BusinessName + "_vue.js.vm":
				err = gfile.PutContents(dataFilePath+"/gen/"+dataFileRange+"/vue/"+pathMap["vue"], val)
				if err != nil {
					libErr.ErrPrint(ctx, err, "????????????????????????")
					return err
				}
			}
		}
	}
	//??????
	err = gcompress.ZipPathWriter(dataFilePath+"/gen/"+dataFileRange, r.Response.Writer)
	if err != nil {
		return err
	}
	//?????????????????????
	gfile.Remove(dataFilePath + "/gen/" + dataFileRange)
	//?????????????????????
	r.Response.Header().Set("Content-Length", gconv.String(r.Response.BufferLength()))
	r.Response.Header().Set("Content-Type", "application/force-download")
	r.Response.Header().Set("Accept-Ranges", "bytes")
	r.Response.Header().Set("Content-Disposition", "attachment; filename=IOTFast.zip")
	r.Response.Buffer()

	return nil
}

//???????????????????????????
func (s *genTableImpl) getPath(entity *model.GenTableEntityExtend) g.MapStrStr {
	controller := "app/controller/" + entity.ModuleName + "/" + entity.ClassName + ".go"
	service := "app/service/" + entity.ModuleName + "/" + entity.BusinessName + "_service/" + entity.ClassName + ".go"
	model := "app/model/" + entity.ModuleName + "/" + entity.BusinessName + "/" + entity.ClassName + ".go"
	gfgenModel := "app/model/" + entity.ModuleName + "/" + entity.BusinessName + "/" + entity.ClassName + "_model.go"
	gfgenEntity := "app/model/" + entity.ModuleName + "/" + entity.BusinessName + "/" + entity.ClassName + "_entity.go"
	vue := "views/" + entity.ModuleName + "/" + entity.BusinessName + "/index.vue"
	api := "api/" + entity.ModuleName + "/" + entity.BusinessName + ".js"
	return g.MapStrStr{
		"controller":  controller,
		"service":     service,
		"model":       model,
		"gfgenModel":  gfgenModel,
		"gfgenEntity": gfgenEntity,
		"vue":         vue,
		"api":         api,
	}
}

//??????????????????
func (s *genTableImpl) genCode(ctx context.Context, tableId int) (data g.MapStrStr, entity *model.GenTableEntityExtend, err error) {
	entity, err = SelectRecordById(ctx, tableId)
	if err != nil {
		libErr.ErrPrint(ctx, err, "?????????????????????")
		return
	}
	if entity == nil {
		libErr.ErrPrint(ctx, err, "?????????????????????")
		err = gerror.New("?????????????????????")
		return
	}
	SetPkColumn(entity, entity.Columns)
	controllerKey := "template/go/" + entity.BusinessName + "_controller.go.vm"
	controllerValue := ""
	serviceKey := "template/go/" + entity.BusinessName + "_service.go.vm"
	serviceValue := ""
	modelKey := "template/go/" + entity.BusinessName + ".go.vm"
	modelValue := ""
	routeKey := "template/go/" + entity.BusinessName + "_route.go.vm"
	routeValue := ""
	gfgenModelKey := "template/go/" + entity.BusinessName + "_model.go.vm"
	gfgenModelValue := ""
	gfgenEntityKey := "template/go/" + entity.BusinessName + "_entity.go.vm"
	gfgenEntityValue := ""
	apiJsKey := "template/html/" + entity.BusinessName + "_api.js.vm"
	apiJsValue := ""
	vueKey := "template/html/" + entity.BusinessName + "_vue.js.vm"
	vueValue := ""

	view := gview.New()
	view.BindFuncMap(g.Map{
		"UcFirst": func(str string) string {
			return gstr.UcFirst(str)
		},
		"add": func(a, b int) int {
			return a + b
		},
	})
	view.SetConfigWithMap(g.Map{
		"Paths":      []string{"resource"},
		"Delimiters": []string{"${", "}"},
	})
	//??????????????????
	var options g.Map
	if entity.TplCategory == "tree" {
		options = gjson.New(entity.Options).Map()
	}
	var tmpController string
	if tmpController, err = view.Parse(ctx, "template/go/"+entity.TplCategory+"/controller.template", g.Map{"table": entity}); err == nil {
		controllerValue = tmpController
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpService string
	if tmpService, err = view.Parse(ctx, "template/go/"+entity.TplCategory+"/service.template", g.Map{"table": entity, "options": options}); err == nil {
		serviceValue = tmpService
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpModel string
	if tmpModel, err = view.Parse(ctx, "template/go/"+entity.TplCategory+"/model.template", g.Map{"table": entity}); err == nil {
		modelValue = tmpModel
		modelValue, err = trimBreak(modelValue)
		libErr.ErrPrint(ctx, err, "???????????????????????????")
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpJs string
	if tmpJs, err = view.Parse(ctx, "template/html/js.template", g.Map{"table": entity}); err == nil {
		apiJsValue = tmpJs
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpVue string
	if tmpVue, err = view.Parse(ctx, "template/html/vue_"+entity.TplCategory+".template", g.Map{"table": entity, "options": options}); err == nil {
		vueValue = tmpVue
		vueValue, err = trimBreak(vueValue)
		libErr.ErrPrint(ctx, err, "???????????????????????????")
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}

	var tmpRouter string
	if tmpRouter, err = view.Parse(ctx, "template/go/common/route.template", g.Map{"table": entity}); err == nil {
		routeValue = tmpRouter
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpGfModel string
	if tmpGfModel, err = view.Parse(ctx, "template/go/common/model.template", g.Map{"table": entity}); err == nil {
		gfgenModelValue = tmpGfModel
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}
	var tmpGfEntity string
	if tmpGfEntity, err = view.Parse(ctx, "template/go/common/entity.template", g.Map{"table": entity}); err == nil {
		gfgenEntityValue = tmpGfEntity
	} else {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return
	}

	data = g.MapStrStr{
		routeKey:       routeValue,
		gfgenModelKey:  gfgenModelValue,
		gfgenEntityKey: gfgenEntityValue,
		modelKey:       modelValue,
		serviceKey:     serviceValue,
		controllerKey:  controllerValue,
		apiJsKey:       apiJsValue,
		vueKey:         vueValue,
	}
	return
}

//???????????????????????????
func GetDbType(columnType string) string {
	if strings.Index(columnType, "(") > 0 {
		return columnType[0:strings.Index(columnType, "(")]
	} else {
		return columnType
	}
}

//????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????? ?????????HELLO_WORLD->HelloWorld
func ConvertToCamelCase(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// ???????????????????????????????????????
		return strings.ToUpper(name[0:1]) + name[1:len(name)]
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
	}
	return result
}

////????????????????????????????????????????????????????????????,???????????????????????????????????????????????????????????????????????????????????????????????????????????? ?????????HELLO_WORLD->helloWorld
func ConvertToCamelCase1(name string) string {
	if name == "" {
		return ""
	} else if !strings.Contains(name, "_") {
		// ??????????????????????????????
		return name
	}
	var result string = ""
	camels := strings.Split(name, "_")
	for index := range camels {
		if camels[index] == "" {
			continue
		}
		camel := camels[index]
		if result == "" {
			result = strings.ToLower(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		} else {
			result = result + strings.ToUpper(camel[0:1]) + strings.ToLower(camel[1:len(camel)])
		}
	}
	return result
}

//??????????????????
func GetColumnLength(columnType string) int {
	start := strings.Index(columnType, "(")
	end := strings.Index(columnType, ")")
	result := ""
	if start >= 0 && end >= 0 {
		result = columnType[start+1 : end-1]
	}
	return gconv.Int(result)
}

//??????????????????
func SelectDbTableListByNames(ctx context.Context, tableNames []string) ([]*entity.GenTable, error) {
	db := g.DB()
	sql := "select * from information_schema.tables where table_name NOT LIKE 'qrtz_%' and table_name NOT LIKE 'gen_%' " +
		" and table_schema = (select database()) "
	if len(tableNames) > 0 {
		in := gstr.TrimRight(gstr.Repeat("?,", len(tableNames)), ",")
		sql += " and " + gdb.FormatSqlWithArgs("table_name in ("+in+")", gconv.SliceAny(tableNames))
	}
	var result []*entity.GenTable
	res, err := db.GetAll(ctx, sql)
	if err != nil {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return nil, gerror.New("????????????????????????")
	}
	err = res.Structs(&result)
	if err != nil {
		libErr.ErrPrint(ctx, err, "????????????????????????")
		return nil, gerror.New("????????????????????????")
	}
	return result, err
}

//????????????????????????
func InitColumnField(column *entity.GenTableColumn, table *entity.GenTable) {
	dataType := GetDbType(column.ColumnType)
	columnName := column.ColumnName
	column.TableId = table.TableId
	column.CreateBy = table.CreateBy
	column.CreateTime = gtime.Now()
	column.UpdateTime = column.CreateTime
	//???????????????
	column.GoField = ConvertToCamelCase(columnName)
	column.HtmlField = ConvertToCamelCase1(columnName)

	if IsStringObject(dataType) {
		//????????????????????????
		column.GoType = "string"
		columnLength := GetColumnLength(column.ColumnType)
		if columnLength >= 500 {
			column.HtmlType = "textarea"
		} else {
			column.HtmlType = "input"
		}
	} else if IsTimeObject(dataType) {
		//?????????????????????
		column.GoType = "Time"
		column.HtmlType = "datetime"
	} else if IsNumberObject(dataType) {
		//?????????????????????
		column.HtmlType = "input"
		t, _ := gregex.ReplaceString(`\(.+\)`, "", column.ColumnType)
		t = gstr.Split(gstr.Trim(t), " ")[0]
		t = gstr.ToLower(t)
		// ??????????????????
		switch t {
		case "float", "double", "decimal":
			column.GoType = "float64"
		case "bit", "int", "tinyint", "small_int", "smallint", "medium_int", "mediumint":
			if gstr.ContainsI(column.ColumnType, "unsigned") {
				column.GoType = "uint"
			} else {
				column.GoType = "int"
			}
		case "big_int", "bigint":
			if gstr.ContainsI(column.ColumnType, "unsigned") {
				column.GoType = "uint64"
			} else {
				column.GoType = "int64"
			}
		}
	}
	//????????????
	if columnName == "create_by" || columnName == "create_time" || columnName == "update_by" || columnName == "update_time" {
		column.IsRequired = "0"
		column.IsInsert = "0"
	} else {
		column.IsRequired = "0"
		column.IsInsert = "1"
		if strings.Index(columnName, "name") >= 0 || strings.Index(columnName, "status") >= 0 {
			column.IsRequired = "1"
		}
	}

	// ????????????
	if IsNotEdit(columnName) {
		if column.IsPk == "1" {
			column.IsEdit = "0"
		} else {
			column.IsEdit = "1"
		}
	} else {
		column.IsEdit = "0"
	}
	// ????????????
	if IsNotList(columnName) {
		column.IsList = "1"
	} else {
		column.IsList = "0"
	}
	// ????????????
	if IsNotQuery(columnName) {
		column.IsQuery = "1"
	} else {
		column.IsQuery = "0"
	}

	// ??????????????????
	if CheckNameColumn(columnName) {
		column.QueryType = "LIKE"
	} else {
		column.QueryType = "EQ"
	}

	// ???????????????????????????
	if CheckStatusColumn(columnName) {
		column.HtmlType = "radio"
	} else if CheckTypeColumn(columnName) || CheckSexColumn(columnName) {
		// ??????&???????????????????????????
		column.HtmlType = "select"
	}
}

//??????????????????3????????????sex
func CheckSexColumn(columnName string) bool {
	if len(columnName) >= 3 {
		end := len(columnName)
		start := end - 3

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "sex" {
			return true
		}
	}
	return false
}

//??????????????????4????????????type
func CheckTypeColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		if columnName[start:end] == "type" {
			return true
		}
	}
	return false
}

//??????????????????6????????????status
func CheckStatusColumn(columnName string) bool {
	if len(columnName) >= 6 {
		end := len(columnName)
		start := end - 6

		if start <= 0 {
			start = 0
		}
		tmp := columnName[start:end]

		if tmp == "status" {
			return true
		}
	}

	return false
}

//??????????????????4????????????name
func CheckNameColumn(columnName string) bool {
	if len(columnName) >= 4 {
		end := len(columnName)
		start := end - 4

		if start <= 0 {
			start = 0
		}

		tmp := columnName[start:end]

		if tmp == "name" {
			return true
		}
	}
	return false
}

//??????????????????
func InitTable(ctx context.Context, table *entity.GenTable, operName string) {
	table.ClassName = ConvertClassName(ctx, table.TableName)
	table.PackageName = g.Cfg().MustGet(ctx, "gen.packageName").String()
	table.ModuleName = g.Cfg().MustGet(ctx, "gen.moduleName").String()
	table.BusinessName = GetBusinessName(table.TableName)
	table.FunctionName = strings.ReplaceAll(table.TableComment, "???", "")
	table.FunctionAuthor = g.Cfg().MustGet(ctx, "gen.author").String()
	table.CreateBy = operName
	table.TplCategory = "crud"
	table.CreateTime = gtime.Now()
	table.UpdateTime = table.CreateTime
}

//?????????????????????
func ConvertClassName(ctx context.Context, tableName string) string {
	autoRemovePre := g.Cfg().MustGet(ctx, "gen.autoRemovePre").Bool()
	tablePrefix := g.Cfg().MustGet(ctx, "gen.tablePrefix").String()
	if autoRemovePre && tablePrefix != "" {
		searchList := strings.Split(tablePrefix, ",")
		for _, str := range searchList {
			if strings.HasPrefix(tableName, str) {
				tableName = strings.Replace(tableName, str, "", 1) //????????????????????????
			}
		}
	}
	return tableName
}

//???????????????
func GetBusinessName(tableName string) string {
	lastIndex := strings.LastIndex(tableName, "_")
	nameLength := len(tableName)
	businessName := tableName[lastIndex+1 : nameLength]
	return businessName
}

func GetTableInfoByTableId(ctx context.Context, tableId int64) (info *entity.GenTable, err error) {
	//return gen_table.GetInfoById(tableId)
	err = dao.GenTable.Ctx(ctx).WherePri(tableId).Scan(&info)
	return
}

func SelectRecordById(ctx context.Context, tableId int) (entityExtend *model.GenTableEntityExtend, err error) {
	var entity *entity.GenTable
	err = dao.GenTable.Ctx(ctx).WherePri(tableId).Scan(&entity)
	if err != nil {
		return
	}
	m := gconv.Map(entity)
	gconv.Struct(m, &entityExtend)

	//???????????????
	entityExtend.Columns, err = SelectGenTableColumnListByTableId(ctx, tableId)
	if err != nil {
		return
	}
	return
}

//?????????????????????
func trimBreak(str string) (s string, err error) {
	var b []byte
	if b, err = gregex.Replace("(([\\s\t]*)\r?\n){2,}", []byte("$2\n"), []byte(str)); err == nil {
		s = gconv.String(b)
	}
	return
}

//?????????????????????
func SetPkColumn(table *model.GenTableEntityExtend, columns []*entity.GenTableColumn) {
	for _, column := range columns {
		if column.IsPk == "1" {
			table.PkColumn = column
			break
		}
	}
	if table.PkColumn == nil {
		table.PkColumn = columns[0]
	}
}
