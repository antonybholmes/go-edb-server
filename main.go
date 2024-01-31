package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/antonybholmes/go-loctogene"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// func loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		// Log incoming request
// 		log.Printf("Incoming Request: %s %s", c.Request().Method, c.Request().URL.String())

// 		// Call next handler
// 		if err := next(c); err != nil {
// 			c.Error(err)
// 		}

// 		// Log outgoing response
// 		log.Printf("Outgoing Response: %d %s", c.Response().Status, http.StatusText(c.Response().Status))

// 		return nil
// 	}
// }

const DEFAULT_ASSEMBLY = "grch38"
const DEFAULT_LEVEL = 1
const DEFAULT_CHR = "chr3"
const DEFAULT_START = 187728170
const DEFAULT_END = 187752257
const DEFAULT_CLOSEST_N = 10

type ParsedLocation struct {
	Loc      *loctogene.Location
	Level    int
	DB       *sql.DB
	Assembly string
}

func parseLocation(c echo.Context) (*ParsedLocation, error) {
	assembly := DEFAULT_ASSEMBLY
	chr := DEFAULT_CHR
	start := DEFAULT_START
	end := DEFAULT_END
	level := DEFAULT_LEVEL

	var v string
	var err error

	v = c.QueryParam("assembly")

	if v != "" {
		assembly = v
	} else {
		c.Logger().Warn(fmt.Sprintf("assembly was not, using default %s", DEFAULT_ASSEMBLY))
	}

	v = c.QueryParam("chr")

	if v != "" {
		chr = v
	} else {
		c.Logger().Warn("chr was not, using default...")
	}

	v = c.QueryParam("start")

	if v != "" {
		start, err = strconv.Atoi(v)

		if err != nil {
			c.Logger().Warn(fmt.Sprintf("%s is an invalid start, using default %d...", v, DEFAULT_START))
			start = DEFAULT_START
		}
	} else {
		c.Logger().Warn(fmt.Sprintf("start was not set, using default %d...", DEFAULT_START))
	}

	v = c.QueryParam("end")

	if v != "" {
		end, err = strconv.Atoi(v)

		if err != nil {
			c.Logger().Warn(fmt.Sprintf("%s is an invalid end, using default %d...", v, DEFAULT_END))
			end = DEFAULT_END
		}
	} else {
		c.Logger().Warn(fmt.Sprintf("end was not set, using default %d...", DEFAULT_END))
	}

	v = c.QueryParam("level")

	if v != "" {
		level = loctogene.GetLevel(v)
	} else {
		c.Logger().Warn(fmt.Sprintf("level was not set, using default %d...", DEFAULT_LEVEL))
	}

	loc := loctogene.Location{Chr: chr, Start: start, End: end}

	db, err := loctogene.GetDB(fmt.Sprintf("data/modules/loctogene/%s.db", assembly))

	if err != nil {
		return nil, err
	}

	return &ParsedLocation{Loc: &loc, Assembly: assembly, DB: db, Level: level}, nil
}

func main() {
	//zerolog.SetGlobalLevel(zerolog.DebugLevel)

	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(loggerMiddleware)
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/genes/within", func(c echo.Context) error {
		loc, err := parseLocation(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
		}

		genes, err := loctogene.GetGenesWithin(loc.DB, loc.Loc, loc.Level)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
		}

		return c.JSON(http.StatusOK, genes)
	})

	// e.POST("/genes/within", func(c echo.Context) error {
	// 	jsonBody := make(map[string]interface{})

	// 	json.NewDecoder(c.Request().Body).Decode(&jsonBody)

	// 	chr := DEFAULT_CHR
	// 	var ok bool

	// 	_, ok = jsonBody["chr"]

	// 	if ok {
	// 		c.Logger().Info("chr set through body.")
	// 		chr = jsonBody["chr"].(string)
	// 	} else {
	// 		c.Logger().Warn("chr not set through body, using default.")
	// 	}

	// 	start := DEFAULT_START

	// 	_, ok = jsonBody["start"]

	// 	if ok {
	// 		c.Logger().Info("start set through body.")
	// 		start = jsonBody["start"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("start not set, using default %d.", DEFAULT_START))
	// 	}

	// 	end := DEFAULT_END

	// 	_, ok = jsonBody["end"]

	// 	if ok {
	// 		c.Logger().Info("end set through body.")
	// 		end = jsonBody["end"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("end not set, using default %d.", DEFAULT_END))
	// 	}

	// 	assembly := DEFAULT_ASSEMBLY

	// 	_, ok = jsonBody["assembly"]

	// 	if ok {
	// 		c.Logger().Info("assembly set through body.")
	// 		assembly = jsonBody["assembly"].(string)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("assembly not set, using default %s.", DEFAULT_ASSEMBLY))
	// 	}

	// 	level := DEFAULT_LEVEL

	// 	_, ok = jsonBody["level"]

	// 	if ok {
	// 		c.Logger().Info("level set through body.")
	// 		level = jsonBody["level"].(int)
	// 	} else {
	// 		c.Logger().Warn(fmt.Sprintf("level not set, using default %d.", DEFAULT_LEVEL))
	// 	}

	// 	c.Logger().Info(fmt.Sprintf("loc: %s:%d-%d on %s", chr, start, end, assembly))

	// 	db, err := loctogene.GetDB(fmt.Sprintf("data/modules/loctogene/%s.db", assembly))

	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "DB error"})
	// 	}

	// 	loc := loctogene.Location{Chr: chr, Start: start, End: end}

	// 	genes, err := loctogene.GetGenesWithin(db, &loc, level)

	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
	// 	}

	// 	return c.JSON(http.StatusOK, genes)
	// })

	e.GET("/genes/closest", func(c echo.Context) error {

		loc, err := parseLocation(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
		}

		n := DEFAULT_CLOSEST_N

		v := c.QueryParam("n")

		if v != "" {
			n, err = strconv.Atoi(v)

			if err != nil {
				c.Logger().Warn(fmt.Sprintf("%s is an invalid, using default n=%d...", v, DEFAULT_CLOSEST_N))
				n = DEFAULT_CLOSEST_N
			}
		} else {
			c.Logger().Warn(fmt.Sprintf("n was not set, using default %d...", DEFAULT_CLOSEST_N))
		}

		genes, err := loctogene.GetClosestGenes(loc.DB, loc.Loc, n, loc.Level)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, struct{ Status string }{Status: "Error"})
		}

		return c.JSON(http.StatusOK, genes)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
