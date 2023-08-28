package api

import (
	"github.com/labstack/echo/v4"
	"github.com/theforgeinitiative/member-button-generator/pdf"
	"github.com/theforgeinitiative/member-button-generator/sfdc"
)

func (h *Handlers) GenerateButtonPDF(c echo.Context) error {
	var buttons []sfdc.MemberButton
	err := c.Bind(&buttons)
	if err != nil {
		return err
	}

	// extract names
	var names []string
	for _, b := range buttons {
		names = append(names, b.Name)
	}

	doc, err := pdf.RenderButtons(names)
	if err != nil {
		return err
	}
	return c.Blob(200, "application/pdf", doc)
}

func (h *Handlers) PrintLabels(c echo.Context) error {
	var buttons []sfdc.MemberButton
	err := c.Bind(&buttons)
	if err != nil {
		return err
	}

	err = h.DymoClient.PrintLabels(buttons)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
