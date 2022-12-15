package resource

import (
	"context"
	"errors"
	"log"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	fieldID        = "id"
	fieldName      = "name"
	fieldServiceID = "service_account_id"
	fieldToken     = "token"
)

func ServiceKey() *schema.Resource {
	return &schema.Resource{
		Description:   "A Service Key authorizes access to all Resources assigned to a Service Account.",
		CreateContext: serviceKeyCreate,
		ReadContext:   serviceKeyRead,
		DeleteContext: serviceKeyDelete,
		UpdateContext: serviceKeyUpdate,

		Schema: map[string]*schema.Schema{
			fieldServiceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the Service Account",
			},
			// optional
			fieldName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Service Key",
			},
			// computed
			fieldID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Autogenerated ID of the Service Key",
			},
			fieldToken: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Autogenerated Token of the Service Key",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func serviceKeyCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	serviceKey, err := client.CreateServiceKey(ctx, &model.ServiceKey{
		Service: resourceData.Get(fieldServiceID).(string),
		Name:    resourceData.Get(fieldName).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Service key %s created with id %v", serviceKey.Name, serviceKey.ID)

	if err := resourceData.Set(fieldToken, serviceKey.Token); err != nil {
		return diag.FromErr(err)
	}

	return serviceKeyReadHelper(ctx, resourceData, serviceKey, nil, meta)
}

func serviceKeyUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	serviceKey, err := client.UpdateServiceKey(ctx,
		&model.ServiceKey{
			ID:   resourceData.Id(),
			Name: resourceData.Get(fieldName).(string),
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updated service key id %v", serviceKey.ID)

	return serviceKeyReadHelper(ctx, resourceData, serviceKey, err, meta)
}

func serviceKeyDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	serviceKey, err := client.ReadServiceKey(ctx, resourceData.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if serviceKey.IsActive() {
		err := client.RevokeServiceKey(ctx, resourceData.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = client.DeleteServiceKey(ctx, resourceData.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted service key id %s", resourceData.Id())

	return nil
}

func serviceKeyRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	serviceKey, err := client.ReadServiceKey(ctx, resourceData.Id())

	return serviceKeyReadHelper(ctx, resourceData, serviceKey, err, meta)
}

func serviceKeyReadHelper(ctx context.Context, resourceData *schema.ResourceData, serviceKey *model.ServiceKey, err error, meta interface{}) diag.Diagnostics {
	if err != nil {
		if errors.Is(err, client.ErrGraphqlResultIsEmpty) {
			// clear state
			resourceData.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	if !serviceKey.IsActive() {
		return reCreateServiceKey(ctx, resourceData, meta)
	}

	if err := resourceData.Set(fieldName, serviceKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(fieldServiceID, serviceKey.Service); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(serviceKey.ID)

	return nil
}

func reCreateServiceKey(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	err := client.DeleteServiceKey(ctx, resourceData.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return serviceKeyCreate(ctx, resourceData, meta)
}