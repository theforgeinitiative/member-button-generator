package api

import (
	"github.com/theforgeinitiative/member-button-generator/dymo"
	"github.com/theforgeinitiative/member-button-generator/sfdc"
)

type Handlers struct {
	SFClient   *sfdc.Client
	DymoClient *dymo.Client
}
