package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/theforgeinitiative/member-button-generator/api"
	"github.com/theforgeinitiative/member-button-generator/dymo"
	"github.com/theforgeinitiative/member-button-generator/sfdc"
)

var (
	sfURL        = os.Getenv("SF_URL")
	clientID     = os.Getenv("SF_CLIENT_ID")
	clientSecret = os.Getenv("SF_CLIENT_SECRET")
	printer      = os.Getenv("DYMO_PRINTER_NAME")
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	// setup SFDC connection
	sfClient, err := sfdc.NewClient(sfURL, clientID, clientSecret)
	if err != nil {
		e.Logger.Fatal("error making SFDC client", err)
	}

	// setup Dymo client
	dymoClient := dymo.NewClient()
	dymoClient.Printer = printer

	// create handler struct
	app := api.Handlers{
		DymoClient: &dymoClient,
		SFClient:   &sfClient,
	}

	// static assets
	e.File("/", "public/index.html")
	e.Static("/assets", "public/assets")

	// api routes
	e.GET("/api/members", app.GetMembers)
	e.POST("/api/members/:id/complete", app.CompleteMember)
	e.POST("/api/buttons/pdf", app.GenerateButtonPDF)
	e.POST("/api/labels/print", app.PrintLabels)

	e.Logger.Fatal(e.Start(":3000"))
}
