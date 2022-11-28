// Code generated by codegen; DO NOT EDIT.

package incidents

import (
	"github.com/cloudquery/plugin-sdk/schema"
)

func Crowdscore() *schema.Table {
	return &schema.Table{
		Name:     "crowdstrike_incidents_crowdscore",
		Resolver: fetchCrowdscore,
		Columns: []schema.Column{
			{
				Name:     "errors",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Errors"),
			},
			{
				Name:     "meta",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Meta"),
			},
			{
				Name:     "resources",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Resources"),
			},
		},
	}
}
