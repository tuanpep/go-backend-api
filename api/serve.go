package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openAPIFile embed.FS

// SwaggerUIHTML is the HTML template for displaying OpenAPI documentation with testing capability
const SwaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Go Backend API - Interactive Documentation</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
	<style>
		html {
			box-sizing: border-box;
			overflow: -moz-scrollbars-vertical;
			overflow-y: scroll;
		}
		*, *:before, *:after {
			box-sizing: inherit;
		}
		body {
			margin: 0;
			background: #fafafa;
		}
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js"></script>
	<script>
		window.onload = function() {
			const ui = SwaggerUIBundle({
				url: "/openapi.yaml",
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				plugins: [
					SwaggerUIBundle.plugins.DownloadUrl
				],
				layout: "StandaloneLayout",
				validatorUrl: null,
				tryItOutEnabled: true
			});
		};
	</script>
</body>
</html>`

// ServeOpenAPISpec serves the OpenAPI specification
func ServeOpenAPISpec(c *gin.Context) {
	// Read the embedded OpenAPI spec file
	spec, err := openAPIFile.ReadFile("openapi.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load OpenAPI specification",
		})
		return
	}

	// Determine content type based on Accept header or query parameter
	format := c.Query("format")
	accept := c.GetHeader("Accept")

	if format == "yaml" || format == "yml" || accept == "application/yaml" || accept == "text/yaml" {
		c.Data(http.StatusOK, "application/yaml", spec)
		return
	}

	// Default to JSON (would need conversion, but for now serve YAML)
	// In production, you might want to convert YAML to JSON
	c.Data(http.StatusOK, "application/yaml", spec)
}

// ServeOpenAPIDocs serves the HTML documentation page with Swagger UI (with Try It Out feature)
func ServeOpenAPIDocs(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, SwaggerUIHTML)
}

// GetOpenAPIFileSystem returns the filesystem for the OpenAPI spec
func GetOpenAPIFileSystem() fs.FS {
	return openAPIFile
}
