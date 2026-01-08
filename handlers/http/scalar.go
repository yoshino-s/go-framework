package http

import (
	"html/template"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

const scalarHttp = `<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>

  <body>
    <div id="app"></div>

    <!-- Load the Script -->
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

    <!-- Initialize the Scalar API Reference -->
    <script>
      Scalar.createApiReference('#app', {
        // The URL of the OpenAPI/Swagger document
        url: {{.OpenAPIURL}},
      })
    </script>
  </body>
</html>
`

func (h *Handler) Scalar(prefix string, apiJSON []byte) {
	prefix = strings.TrimRight(prefix, "/")

	// create a template with name
	index, _ := template.New("scalar_index.html").Parse(scalarHttp)

	group := h.Group(prefix)

	group.GET("/doc.json", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = c.Response().Writer.Write(apiJSON)
		return nil
	})

	indexFunc := func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		return index.Execute(c.Response().Writer, map[string]interface{}{
			"OpenAPIURL": path.Join(prefix, "doc.json"),
		})
	}

	group.GET("", indexFunc)

	group.GET("/", indexFunc)
}
