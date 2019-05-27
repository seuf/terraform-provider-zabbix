package provider

import (
	"errors"
	"fmt"
	"log"

	"github.com/seuf/zabbix-1"

	"github.com/hashicorp/terraform/helper/schema"
)

var OPERATION_TYPE = map[string]zabbix.OperationType{
	"send_message":            0,
	"remote_command":          1,
	"add_host":                2,
	"remove_host":             3,
	"add_to_host_group":       4,
	"remove_from_host_group":  5,
	"link_to_template":        6,
	"unlink_from_template":    7,
	"enable_host":             8,
	"disable_host":            9,
	"set_host_inventory_mode": 10,
}

var EVAL_TYPES = map[string]zabbix.EVAL_TYPES


func resourceZabbixAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixActionCreate,
		Read:   resourceZabbixActionRead,
		Update: resourceZabbixActionUpdate,
		Delete: resourceZabbixActionDelete,
		Schema: map[string]*schema.Schema{
			"actionid": &schema.Schema{
				Type:	schema.TypeString,
				Required: false,
				ForceNew: false,
				Description: "(readonly) ID of the action",

			},
			"esc_period ": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				Description: "Default operation step duration. Must be greater than 60 seconds. Accepts seconds, time unit with suffix and user macro. "
			},
			"eventsource ": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Descrption: "Type of events that the action will handle. "
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"def_longdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"def_shortdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"r_longdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"r_shortdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"ack_longdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"ack_shortdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: false,
			},
			"status": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"pause_suppressed": &schema.Schema{
				Type: 	schema.TypeBool,
				Default: false,
				Optional: true,				
			}
			"operation": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     operationSchema,
				Required: true,
				ForceNew: true,
			},
			
		},
	}
}


var operationSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"operationid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"operationtype": &schema.Schema{
			Type: schema.TypeString,
			Required: true,
			ForceNew: true,
			Description: "Type of operation."
		},
		"actionid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"esc_period ": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: false,
			Description: "Default operation step duration. Must be greater than 60 seconds. Accepts seconds, time unit with suffix and user macro. "
		},
		"esc_step_from": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			ForceNew: false,
			Default: "1",
			Description: "Step to start escalation from."
		},
		"esc_step_to": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			ForceNew: false,
			Default: "1",
			Description: "Step to end escalation at."
		},
		"evaltype": &schema.Schema{
			Type:     schema.TypeString,
			Required: false,
			ForceNew: false,
			Default: "and_or",
			Description: "Operation condition evaluation method."
		},
		// "opcommand": &schema.Schema{
		// 	Type:     schema.TypeBool,
		// 	Required: false,
		// 	ForceNew: false,
		// 	Default: "and_or",
		// 	Description: "Object containing the data about the command run by the operation."
		// },
		"opgroup": &schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{
				"operationid": &schema.Schema{
					Type: schema.TypeString,
					Optional: true,
				},
				"groupid": &schema.Schema{
					Type: schema.TypeString,
					Required: false,
					Description: " Required for “add to host group” and “remove from host group” operations."
				},
			}
			Required: false,
			ForceNew: false,
			Description: "Host groups to add hosts to."
	}
}




func resourceZabbixActionCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	action := zabbix.Action{
		EscPeriod: 		d.Get("esc_perdiod").(string),
		EventSource: 	d.Get("eventsource").(string),
		Name:   		d.Get("name").(string),
	}
	if !d.Get("status").(bool) {
		action.Status = 1
	}
	typeId, ok := OPERATION_TYPE[d.Get("operationtype")]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s isnt valid interface type", interfaceType))
	}
	action.OperationType = typeId

	actions = zabbix.Actions{*action}

	err = api.ActionsCreate(actions)

	if err != nil {
		return err
	}

	log.Printf("Created action id is %s", actions[0].ActionId)

	d.Set("action_id", actions[0].ActionId)
	d.SetId(actions[0].ActionId)

	return nil
}

func resourceZabbixActionRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	log.Printf("Will read action with id %s", d.Get("actionid").(string))

	action, err := api.ActionGetById(d.Get("actionid").(string))

	if err != nil {
		return err
	}

	log.Printf("Action name is %s", action.Name)

	d.Set("action", action.Action)
	d.Set("name", action.Name)

	d.Set("status", action.Status == 0)

	return nil
}

func resourceZabbixActionUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	action := zabbix.Action{
		EscPeriod: 		d.Get("esc_perdiod").(string),
		EventSource: 	d.Get("eventsource").(string),
		Name:   		d.Get("name").(string),
	}
	if !d.Get("status").(bool) {
		action.Status = 1
	}
	typeId, ok := OPERATION_TYPE[d.Get("operationtype")]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s isnt valid interface type", interfaceType))
	}
	action.OperationType = typeId

	action.ActionId = d.Id()

	actions = zabbix.Actions{*action}

	err = api.ActionsUpdate(actions)

	if err != nil {
		return err
	}

	log.Printf("Updated action id is %s", hosts[0].ActionId)

	return nil
}

func resourceZabbixActionDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	return api.ActionsDeleteByIds([]string{d.Id()})
}
