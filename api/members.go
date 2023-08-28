package api

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/theforgeinitiative/member-button-generator/sfdc"
)

func (h *Handlers) GetMembers(c echo.Context) error {
	barcodes := c.QueryParam("barcodes")

	var buttons []sfdc.MemberButton
	var err error
	if len(barcodes) > 0 {
		buttons, err = h.SFClient.ButtonsByBarcode(strings.Split(barcodes, ",")...)
	} else {
		buttons, err = h.SFClient.UnprintedButtons()
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve barcodes: %w", err)
	}

	return c.JSON(200, buttons)
}

func (h *Handlers) CompleteMember(c echo.Context) error {
	id := c.Param("id")

	err := h.SFClient.CompleteButton(id)
	if err != nil {
		return fmt.Errorf("failed to update completion date for contact %s: %w", id, err)
	}
	return c.NoContent(204)
}
