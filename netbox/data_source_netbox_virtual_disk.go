package netbox

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/virtualization"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetboxVirtualDisk() *schema.Resource {
    return &schema.Resource{
        Read: dataSourceNetboxVirtualDiskRead,
		Description: `:meta:subcategory:Virtualization:From the [official documentation](https://docs.netbox.dev/en/stable/models/virtualization/virtualdisk/):
		> A virtual disk is used to model discrete virtual hard disks assigned to virtual machines.`,
        Schema: map[string]*schema.Schema {
            "vm_id": {
                Type: schema.TypeInt,
                Optional: true,
            },
            "filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
                        "name": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "description": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "disk_size_gb": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                      //"tag_ids": {
                      //    Type:     schema.TypeList,
                      //    Computed: true,
                      //    Elem: &schema.Schema{
                      //        Type: schema.TypeInt,
                      //    },
                      //},
						"custom_fields": {
							Type:     schema.TypeMap,
							Computed: true,
						},
                        "vm_id": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
					},
				},
			},
        },
    }
}

func dataSourceNetboxVirtualDiskRead(d *schema.ResourceData, m interface{}) error {
    api := m.(*client.NetBoxAPI)

    params := virtualization.NewVirtualizationVirtualDisksListParams()

    if vmID, ok := d.Get("vm_id").(int); ok && vmID != 0 {
        params.VirtualMachineIDn = strToPtr(strconv.FormatInt(int64(vmID), 10))
    }

    if filter, ok := d.GetOk("filter"); ok {
        var filterParams = filter.(*schema.Set)
        for _, f := range filterParams.List() {
            k := f.(map[string]interface{})["name"]
            v := f.(map[string]interface{})["value"]

            vString := v.(string)
            switch k {
            case "name":
                params.Name = &vString
            case "tag":
                params.Tag = []string{vString}
            default:
                return fmt.Errorf("'%', is not a supported filter parameter", k)
            }
        }
    }

    params.Limit = getOptionalInt(d, "limit")

    res, err := api.Virtualization.VirtualizationVirtualDisksList(params, nil)
    if err != nil {
        return nil
    }

    if *res.GetPayload().Count == int64(0) {
        return errors.New(fmt.Sprintf("%+v", res.GetPayload().Results))
    }

    var s []map[string]interface{}
    for _, v := range res.GetPayload().Results {
        fmt.Printf("Debugging value: %+v", v)
        var mapping = make(map[string]interface{})
        mapping["id"] = v.ID

        if v.Name != nil {
            mapping["name"] = *v.Name
        }

        if v.Description != "" {
            mapping["description"] = v.Description
        }

        if v.Size != nil {
            mapping["disk_size_gb"] = *v.Size
        }

      //if v.Tags != nil {
      //    var tags []int64
      //    for _, t := range v.Tags {
      //        tags = append(tags, t.ID)
      //    }

      //    mapping["tag_Ids"] = tags
      //}

        if v.CustomFields != nil {
            mapping["custom_fields"] = v.CustomFields
        }

        mapping["vm_id"] = v.VirtualMachine.ID

        fmt.Printf("%+v", mapping)

        s = append(s, mapping)
    }

    d.SetId(id.UniqueId())
    return d.Set("disks", s)
}
