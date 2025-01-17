package datasource

import (
	"context"
	"errors"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceConnectorsRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	connectors, err := c.ReadConnectors(ctx)
	if err != nil && !errors.Is(err, client.ErrGraphqlResultIsEmpty) {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(attr.Connectors, convertConnectorsToTerraform(connectors)); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId("all-connectors")

	return nil
}

func Connectors() *schema.Resource {
	return &schema.Resource{
		Description: "Connectors provide connectivity to Remote Networks. For more information, see Twingate's [documentation](https://docs.twingate.com/docs/understanding-access-nodes).",
		ReadContext: datasourceConnectorsRead,
		Schema: map[string]*schema.Schema{
			attr.Connectors: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Connectors",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						attr.ID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Connector.",
						},
						attr.Name: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the Connector.",
						},
						attr.RemoteNetworkID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Remote Network attached to the Connector.",
						},
						attr.StatusUpdatesEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether status notifications are enabled for the Connector.",
						},
					},
				},
			},
		},
	}
}
