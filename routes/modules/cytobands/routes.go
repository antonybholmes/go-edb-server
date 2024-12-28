package seqroutes

import (
	"github.com/antonybholmes/go-edb-server/routes"

	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/labstack/echo/v4"
)

func CytobandsRoute(c echo.Context) error {

	cytobands, _ := cytobandsdbcache.Cytobands(c.Param("assembly"), c.Param("chr"))

	// if err != nil {
	// 	return routes.ErrorReq(err)
	// }

	return routes.MakeDataPrettyResp(c, "", cytobands)
}
