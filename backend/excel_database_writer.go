package backend

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func createExcelFile() error {
	file := excelize.NewFile()

	columns := []string{"URL", "Ergebnis", "Notiz"}

	for i, column := range columns {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), column)
	}

	err := file.SaveAs("leads.xlsx")
	if err != nil {
		return err
	}

	return nil
}

func WriteBusinessUrlToExcelDatabase(websiteUrl string) error {
	filePath := "leads.xlsx"

	url := "http://" + websiteUrl

	_, err := os.Stat(filePath)
	if err != nil {
		err = createExcelFile()
		if err != nil {
			return err
		}
	}
	
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return err
	}

	rowIndex := strconv.Itoa(len(rows) + 1)

	err = file.SetCellValue("Sheet1", "A" + rowIndex, url)
	if err != nil {
		return err
	}
	
	err = file.SetCellHyperLink("Sheet1", "A" + rowIndex, url, "External")
	if err != nil {
		return err
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}
