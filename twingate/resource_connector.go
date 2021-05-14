package twingate

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Remote networks are the logical containers that group Resources together.\n" +
			"Checkout the [twingate docs](https://docs.twingate.com/docs/remote-networks) for detailed information",
		CreateContext: resourceConnectorCreate,
		ReadContext:   resourceConnectorRead,
		DeleteContext: resourceConnectorDelete,

		Schema: map[string]*schema.Schema{
			// required
			"remote_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the remote network to attach the connector to",
			},
			// computed
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Autogenerated ID of the connector in encoded in base64",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Autogenerated name of the connector (can't be changed)",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	remoteNetworkID := d.Get("remote_network_id").(string)
	connector, err := client.createConnector(remoteNetworkID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(connector.ID)
	log.Printf("[INFO] Created conector %s", connector.Name)

	return resourceConnectorRead(ctx, d, m)
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	connectorID := d.Id()

	err := client.deleteConnector(connectorID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Destroyed connector id %s", d.Id())

	return diags
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var diags diag.Diagnostics

	connectorID := d.Id()
	connector, err := client.readConnector(connectorID)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", connector.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w ", err))
	}

	return diags
}