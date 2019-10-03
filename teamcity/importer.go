package teamcity

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)


func subresourceImporter(readFunc schema.ReadFunc) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		var buildConfigID, resourceID string

		if n, err := fmt.Sscanf(d.Id(), "%s %s", &buildConfigID, &resourceID); err != nil {
			return nil, fmt.Errorf("invalid import ID '%s' - %v", d.Id(), err)
		} else if n != 2 {
			return nil, fmt.Errorf("invalid import ID '%s' - Unrecognized format", d.Id())
		}

		d.SetId(resourceID)

		if err := d.Set("build_config_id", buildConfigID); err != nil {
			return nil, err
		}

		if err := readFunc(d, meta); err != nil {
			return nil, err
		}

		return []*schema.ResourceData{d}, nil
	}
}
