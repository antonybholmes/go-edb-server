package utilroutes

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"

	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-sys"
	"github.com/labstack/echo/v4"
)

type XlsxReq struct {
	Sheet    string `json:"sheet"`
	Headers  int    `json:"headers"`
	Indexes  int    `json:"indexes"`
	SkipRows int    `json:"skipRows"`
	Xlsx     string `json:"b64xlsx"`
}

type XlsxResp struct {
	Table *sys.Table `json:"table"`
}

type XlsxSheetsResp struct {
	Sheets []string `json:"sheets"`
}

func makeXlsxReader(data string) (*bytes.Reader, error) {
	xlsxb, err := b64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(xlsxb)

	return reader, nil
}

func XlsxSheetsRoute(c echo.Context) error {

	var req XlsxReq

	err := c.Bind(&req)

	if err != nil {
		return err
	}

	//log.Debug().Msgf("m:%s", req.Xlsx)

	reader, err := makeXlsxReader(req.Xlsx)

	if err != nil {
		return err
	}

	sheets, err := sys.XlsxSheetNames(reader)

	if err != nil {
		return err
	}

	resp := XlsxSheetsResp{Sheets: sheets}

	return routes.MakeDataPrettyResp(c, "", resp)
}

func XlsxToRoute(c echo.Context) error {

	format := c.Param("format")

	if format != "json" {
		return routes.ErrorReq(fmt.Errorf("unsupported format: %s", format))
	}

	var req XlsxReq

	err := c.Bind(&req)

	if err != nil {
		return routes.ErrorReq(err)
	}

	reader, err := makeXlsxReader(req.Xlsx)

	if err != nil {
		return routes.ErrorReq(err)
	}

	table, err := sys.XlsxToJson(reader, req.Sheet, req.Indexes, req.Headers, req.SkipRows)

	if err != nil {
		return routes.ErrorReq(err)
	}

	resp := XlsxResp{Table: table}

	return routes.MakeDataPrettyResp(c, "", resp)
}
