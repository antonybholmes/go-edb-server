package utilroutes

import (
	"github.com/antonybholmes/go-edb-server/routes"
	"github.com/antonybholmes/go-sys"
	"github.com/labstack/echo/v4"
)

type XlsxReq struct {
	Header   int    `json:"header"`
	IndexCol int    `json:"indexCol"`
	Xlsx     []byte `json:"xslx"`
}

type XlsxResp struct {
	Table *sys.Table `json:"table"`
}

func XlsxToTextRoute(c echo.Context) error {

	var req XlsxReq

	err := c.Bind(&req)

	if err != nil {
		return err
	}

	table, err := sys.XlsxToText(req.Xlsx, req.IndexCol, req.Header)

	if err != nil {
		return err
	}

	resp := XlsxResp{Table: table}

	return routes.MakeDataPrettyResp(c, "", resp)

}
