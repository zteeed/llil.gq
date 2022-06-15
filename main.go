/*
 * URL Shortener API
 *
 * This is an URL Shortener API
 *
 * API version: 1.0.0
 * Contact: aurelien@duboc.xyz
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"os"
)

func main() {
	app := App{}
	app.Initialize(
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_BASE_URL"))

	defer app.DB.Close()
	app.Run("0.0.0.0:8888")
}
