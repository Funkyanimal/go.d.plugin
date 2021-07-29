package heliumvalidator

import "github.com/netdata/go.d.plugin/agent/module"

var charts = module.Charts{
  // response from JSONRPC request: command "block_height"[height]"
	{
		ID:    "blockheight",
		Title: "Validator Height",
		Units: "height",
		Fam:   "Validator",
		Ctx:   "heliumvalidator.blockheight",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "block_height", Name: "Height"},
		},
	},
	{
		ID:    "blockage",
		Title: "Block Age",
		Units: "age",
		Fam:   "Validator",
		Ctx:   "heliumvalidator.blockage",
		Dims: module.Dims{
			{ID: "block_age", Name: "Block Age"},
		},
	},
}
