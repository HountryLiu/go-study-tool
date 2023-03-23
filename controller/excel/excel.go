package excel

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/HountryLiu/go-study-tool/config"
	"github.com/HountryLiu/go-study-tool/model"
	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags Excel操作
// @Summary 数据导出
// @Description 数据导出
// @accept application/json
// @Produce application/json
// @Param cur_page query int false "当前第几页，默认1"
// @Param page_size query int false "一页显示数据量，默认20"
// @Param file_type query string false "导出文件类型(csv,xlxs)"
// @Success 200 {array} byte "文件内容字节流"
// @Failure 500 {object} object{no=int,data=string,errors=string}
// @Router /api/excel/export [get]
func Export(ctx *gin.Context) {
	cur_page, _ := strconv.Atoi(ctx.DefaultQuery("cur_page", "1"))
	page_size, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	//默认csv
	file_type := ctx.DefaultQuery("file_type", CSV)

	var excel_data model.ExcelData
	list, err := excel_data.List(cur_page, page_size)
	if err != nil {
		utils.Error(ctx, utils.DatabaseSelectError, err)
		return
	}
	file_name := ""
	file_path := ""
	if file_type == CSV {
		file_name, file_path, err = exportByCsv(list["data"].([]model.ExcelData))
		if err != nil {
			utils.Error(ctx, utils.InternalServerError, err)
			return
		}
	} else if file_type == XLSX {
		file_name, file_path, err = exportByXlsx(list["data"].([]model.ExcelData))
		if err != nil {
			utils.Error(ctx, utils.InternalServerError, err)
			return
		}
	} else {
		utils.Error(ctx, utils.UnSupportFileType, nil)
		return
	}

	defer func() {
		err := utils.RmDir(file_path) //下载后，删除文件
		if err != nil {
			fmt.Println("remove excel file failed", err)
		}
	}()

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file_name))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream") //设置下载文件格式，流式下载
	ctx.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.File(file_path) //直接返回文件
}

func exportByCsv(data []model.ExcelData) (file_name, file_path string, err error) {
	//创建csv文件
	file_name = fmt.Sprintf("人员信息-%s.csv", utils.GetNowTimeString())
	file_path = config.GetExcelPath() + file_name
	header := []string{"ID", "创建时间", "最后更新时间", "姓名", "性别", "年龄", "身份证号", "余额"}

	err = utils.MkdirAll(config.GetExcelPath())
	if err != nil {
		return
	}
	csvFile, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return
	}
	defer csvFile.Close()
	//开始写入内容
	//写入UTF-8 BOM,此处如果不写入就会导致写入的汉字乱码
	csvFile.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(csvFile)
	wStr.Write(header)

	for _, v := range data {
		wStr.Write([]string{
			strconv.Itoa(int(v.ID)),
			v.CreatedAt.String(),
			v.UpdatedAt.String(),
			v.Name,
			v.Sex,
			strconv.Itoa(v.Age),
			v.IDCard,
			strconv.FormatFloat(v.Money, 'f', 2, 64),
		})
	}
	wStr.Flush() //写入文件
	return
}

func exportByXlsx(data []model.ExcelData) (file_name, file_path string, err error) {
	file_name = fmt.Sprintf("人员信息-%s.xlsx", utils.GetNowTimeString())
	file_path = config.GetExcelPath() + file_name
	header := []string{"ID", "创建时间", "最后更新时间", "姓名", "性别", "年龄", "身份证号", "余额"}

	err = utils.MkdirAll(config.GetExcelPath())
	if err != nil {
		return
	}

	sheet_name := "Sheet1"
	f := excelize.NewFile()
	_ = f.SetSheetRow(sheet_name, "A1", &header)
	for k, v := range data {
		row := []interface{}{
			strconv.Itoa(int(v.ID)),
			v.CreatedAt.String(),
			v.UpdatedAt.String(),
			v.Name,
			v.Sex,
			strconv.Itoa(v.Age),
			v.IDCard,
			strconv.FormatFloat(v.Money, 'f', 2, 64),
		}
		err = f.SetSheetRow(sheet_name, "A"+strconv.Itoa(k+2), &row)
		if err != nil {
			return
		}
	}
	if err = f.SaveAs(file_path); err != nil {
		return
	}
	return
}

// @Tags Excel操作
// @Summary 数据导入
// @Description 数据导入
// @accept application/json
// @Produce application/json
// @Param file formData file true "file"
// @Failure 200 {object} object{no=int,data=string,msg=string}
// @Router /api/excel/import [post]
func Import(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
		return
	}
	var excel_datas []model.ExcelData
	if filepath.Ext(header.Filename)[1:] == CSV {
		excel_datas, err = importByCsv(file)
		if err != nil {
			utils.Error(ctx, utils.FileSystemError, err)
			return
		}
	} else if filepath.Ext(header.Filename)[1:] == XLSX {
		excel_datas, err = importByXlsx(file)
		if err != nil {
			utils.Error(ctx, utils.FileSystemError, err)
			return
		}
	} else if filepath.Ext(header.Filename)[1:] == XLS {
		excel_datas, err = importByXls(file)
		if err != nil {
			utils.Error(ctx, utils.FileSystemError, err)
			return
		}
	} else {
		utils.Error(ctx, utils.UnSupportFileType, nil)
		return
	}

	//批量创建数据
	if err = model.DB().Create(&excel_datas).Error; err != nil {
		utils.Error(ctx, utils.DatabaseInsertError, err)
		return
	}

	utils.Success(ctx)

}
func importByCsv(r io.Reader) (excel_datas []model.ExcelData, err error) {
	file := csv.NewReader(r)

	rows, err := file.ReadAll()
	if err != nil {
		return
	}

	for iRow, row := range rows {
		//排除表头 第一行
		if iRow > 0 {
			// "ID", "创建时间", "最后更新时间", "姓名", "性别", "年龄", "身份证号", "余额"
			age, _ := strconv.Atoi(row[5])
			money, _ := strconv.ParseFloat(row[7], 64)
			excel_datas = append(excel_datas, model.ExcelData{
				Name:   row[3],
				Sex:    row[4],
				Age:    age,
				IDCard: row[6],
				Money:  money,
			})
		}
	}
	return
}
func importByXlsx(r io.Reader) (excel_datas []model.ExcelData, err error) {
	sheet_name := "Sheet1"
	file, err := excelize.OpenReader(r)
	if err != nil {
		return
	}
	rows, err := file.GetRows(sheet_name)
	if err != nil {
		return
	}
	for iRow, row := range rows {
		//排除表头 第一行
		if iRow > 0 {
			// "ID", "创建时间", "最后更新时间", "姓名", "性别", "年龄", "身份证号", "余额"
			age, _ := strconv.Atoi(row[5])
			money, _ := strconv.ParseFloat(row[7], 64)
			excel_datas = append(excel_datas, model.ExcelData{
				Name:   row[3],
				Sex:    row[4],
				Age:    age,
				IDCard: row[6],
				Money:  money,
			})
		}
	}
	return
}

func importByXls(r io.Reader) (excel_datas []model.ExcelData, err error) {
	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, r)
	if err != nil {
		return
	}
	seeker := bytes.NewReader(buffer.Bytes())
	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	xlFile, err := xls.OpenReader(seeker, "utf-8")
	if err != nil {
		return
	}
	sheet := xlFile.GetSheet(0)
	if sheet.MaxRow != 0 {
		for i := 0; i <= int(sheet.MaxRow); i++ {
			//排除表头 第一行
			if i <= 0 {
				continue
			}
			row := sheet.Row(i)
			// "ID", "创建时间", "最后更新时间", "姓名", "性别", "年龄", "身份证号", "余额"
			age, _ := strconv.Atoi(row.Col(5))
			money, _ := strconv.ParseFloat(row.Col(7), 64)
			excel_datas = append(excel_datas, model.ExcelData{
				Name:   row.Col(3),
				Sex:    row.Col(4),
				Age:    age,
				IDCard: row.Col(6),
				Money:  money,
			})
		}
	} else {
		err = errors.New("没有内容")
		return
	}
	return
}
