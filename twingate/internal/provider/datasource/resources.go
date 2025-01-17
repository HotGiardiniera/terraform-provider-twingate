package datasource

import (
	"context"
	"errors"
	"fmt"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceResourcesRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	resourceName := resourceData.Get(attr.Name).(string)

	resources, err := c.ReadResourcesByName(ctx, resourceName)
	if err != nil && !errors.Is(err, client.ErrGraphqlResultIsEmpty) {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(attr.Resources, convertResourcesToTerraform(resources)); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId("query resources by name: " + resourceName)

	return nil
}

func Resources() *schema.Resource { //nolint:funlen
	portsResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			attr.Policy: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Whether to allow or deny all ports, or restrict protocol access within certain port ranges: Can be `%s` (only listed ports are allowed), `%s`, or `%s`", model.PolicyRestricted, model.PolicyAllowAll, model.PolicyDenyAll),
			},
			attr.Ports: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of port ranges between 1 and 65535 inclusive, in the format `100-200` for a range, or `8080` for a single port",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	return &schema.Resource{
		Description: "Resources in Twingate represent servers on the private network that clients can connect to. Resources can be defined by IP, CIDR range, FQDN, or DNS zone. For more information, see the Twingate [documentation](https://docs.twingate.com/docs/resources-and-access-nodes).",
		ReadContext: datasourceResourcesRead,
		Schema: map[string]*schema.Schema{
			attr.Name: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Resource",
			},
			// computed
			attr.Resources: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						attr.ID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the Resource",
						},
						attr.Name: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Resource",
						},
						attr.Address: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Resource's IP/CIDR or FQDN/DNS zone",
						},
						attr.RemoteNetworkID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote Network ID where the Resource lives",
						},
						attr.Protocols: {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Restrict access to certain protocols and ports. By default or when this argument is not defined, there is no restriction, and all protocols and ports are allowed.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									attr.AllowIcmp: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to allow ICMP (ping) traffic",
									},
									attr.TCP: {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &portsResource,
									},
									attr.UDP: {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &portsResource,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
