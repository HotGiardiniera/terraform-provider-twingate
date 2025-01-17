package resource

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ErrAllowedToChangeOnlyManualUsers = fmt.Errorf("only users of type %s may be modified", model.UserTypeManual)

func User() *schema.Resource { //nolint:funlen
	return &schema.Resource{
		Description:   "Users provides different levels of write capabilities across the Twingate Admin Console. For more information, see Twingate's [documentation](https://www.twingate.com/docs/users).",
		CreateContext: userCreate,
		ReadContext:   userRead,
		DeleteContext: userDelete,
		UpdateContext: userUpdate,
		Schema: map[string]*schema.Schema{
			attr.Email: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The User's email address",
			},
			// optional
			attr.FirstName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The User's first name",
			},
			attr.LastName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The User's last name",
			},
			attr.SendInvite: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to send an email invitation to the User. True by default.",
			},
			attr.IsActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether the User is active or not. Inactive users will be not able to sign in.",
			},
			attr.Role: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  fmt.Sprintf("Determines the User's role. Either %s.", utils.DocList(model.UserRoles)),
				ValidateFunc: validation.StringInSlice(model.UserRoles, false),
			},
			// computed
			attr.Type: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("Indicates the User's type. Either %s.", utils.DocList(model.UserTypes)),
			},
			attr.ID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Autogenerated ID of the User, encoded in base64.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func userCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	user, err := client.CreateUser(ctx, convertUser(resourceData))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] User %s created with id %v", user.Email, user.ID)

	return resourceUserReadHelper(resourceData, user, nil)
}

func userUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	err := isAllowedToChangeUser(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.UpdateUser(ctx, convertUserUpdate(resourceData))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updated user id %v", user.ID)

	return resourceUserReadHelper(resourceData, user, err)
}

func userDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	err := isAllowedToChangeUser(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteUser(ctx, resourceData.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted user id %s", resourceData.Id())

	return nil
}

func userRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	user, err := c.ReadUser(ctx, resourceData.Id())

	return resourceUserReadHelper(resourceData, user, err)
}

func resourceUserReadHelper(resourceData *schema.ResourceData, user *model.User, err error) diag.Diagnostics {
	if err != nil {
		if errors.Is(err, client.ErrGraphqlResultIsEmpty) {
			// clear state
			resourceData.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	resourceData.SetId(user.ID)

	if err := resourceData.Set(attr.Email, user.Email); err != nil {
		return ErrAttributeSet(err, attr.Email)
	}

	if err := resourceData.Set(attr.FirstName, user.FirstName); err != nil {
		return ErrAttributeSet(err, attr.FirstName)
	}

	if err := resourceData.Set(attr.LastName, user.LastName); err != nil {
		return ErrAttributeSet(err, attr.LastName)
	}

	if err := resourceData.Set(attr.Role, user.Role); err != nil {
		return ErrAttributeSet(err, attr.Role)
	}

	if err := resourceData.Set(attr.Type, user.Type); err != nil {
		return ErrAttributeSet(err, attr.Type)
	}

	if err := resourceData.Set(attr.IsActive, user.IsActive); err != nil {
		return ErrAttributeSet(err, attr.IsActive)
	}

	return nil
}

func isAllowedToChangeUser(data *schema.ResourceData) error {
	userType := data.Get(attr.Type).(string)
	if userType != model.UserTypeManual {
		return ErrAllowedToChangeOnlyManualUsers
	}

	return nil
}

func convertUser(data *schema.ResourceData) *model.User {
	return &model.User{
		ID:         data.Id(),
		Email:      data.Get(attr.Email).(string),
		FirstName:  data.Get(attr.FirstName).(string),
		LastName:   data.Get(attr.LastName).(string),
		SendInvite: convertSendInviteFlag(data),
		Role:       withDefaultValue(data.Get(attr.Role).(string), model.UserRoleMember),
		Type:       data.Get(attr.Type).(string),
		IsActive:   convertIsActiveFlag(data),
	}
}

func convertUserUpdate(data *schema.ResourceData) *model.UserUpdate {
	req := &model.UserUpdate{ID: data.Id()}

	if data.HasChange(attr.FirstName) {
		req.FirstName = stringPtr(data.Get(attr.FirstName).(string))
	}

	if data.HasChange(attr.LastName) {
		req.LastName = stringPtr(data.Get(attr.LastName).(string))
	}

	if data.HasChange(attr.Role) {
		req.Role = stringPtr(data.Get(attr.Role).(string))
	}

	if data.HasChange(attr.IsActive) {
		req.IsActive = boolPtr(convertIsActiveFlag(data))
	}

	return req
}

func convertSendInviteFlag(data *schema.ResourceData) bool {
	return getBooleanFlag(data, attr.SendInvite, true)
}

func convertIsActiveFlag(data *schema.ResourceData) bool {
	return getBooleanFlag(data, attr.IsActive, true)
}
