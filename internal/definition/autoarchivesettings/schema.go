// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package autoarchivesettings

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoarch "github.com/signalfx/signalfx-go/automated-archival"
	"go.uber.org/multierr"
)

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"creator": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the creator of the automated archival setting",
		},
		"last_updated_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of user who last updated the automated archival setting",
		},
		"created": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Timestamp of when the automated archival setting was created",
		},
		"last_updated": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Timestamp of when the automated archival setting was last updated",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Version of the automated archival setting",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Whether the automated archival is enabled for this organization or not",
		},
		"lookback_period": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "This tracks if a metric was unused in the past N number of days (N one of 30, 45, or 60). We’ll archive a metric if it wasn’t used in the lookback period. The value here uses ISO 8061 duration format. Examples - 'P30D', 'P45D', 'P60D'",
		},
		"grace_period": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Grace period is an org level setting that applies to the newly created metrics. This allows customers to protect newly added metrics that users haven’t had the time to use in charts and detectors from being automatically archived The value here uses ISO 8061 duration format. Examples - 'P0D', 'P15D', 'P30D', 'P45D', 'P60D'",
		},
		"ruleset_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Org limit for the number of rulesets that can be created",
		},
	}
}

func decodeTerraform(data *schema.ResourceData) (*autoarch.AutomatedArchivalSettings, error) {
	settings := &autoarch.AutomatedArchivalSettings{
		Enabled:        data.Get("enabled").(bool),
		LookbackPeriod: data.Get("lookback_period").(string),
		GracePeriod:    data.Get("grace_period").(string),
	}
	if creatorStr, ok := data.GetOk("creator"); ok {
		if s, ok := creatorStr.(string); ok {
			settings.Creator = &s
		}
	}
	if lastUpdatedByStr, ok := data.GetOk("last_updated_by"); ok {
		if s, ok := lastUpdatedByStr.(string); ok {
			settings.LastUpdatedBy = &s
		}
	}
	if createdStr, ok := data.GetOk("created"); ok {
		if i, ok := createdStr.(int); ok {
			i64 := int64(i)
			settings.Created = &i64
		}
	}
	if lastUpdatedStr, ok := data.GetOk("last_updated"); ok {
		if i, ok := lastUpdatedStr.(int); ok {
			i64 := int64(i)
			settings.LastUpdated = &i64
		}
	}
	if versionStr, ok := data.GetOk("version"); ok {
		version, err := strconv.ParseInt(versionStr.(string), 10, 64)
		if err == nil {
			settings.Version = version
		}
	}
	if rulesetLimit, ok := data.GetOk("ruleset_limit"); ok {
		rulesetLimit := rulesetLimit.(int)
		if rulesetLimit > math.MaxInt32 || rulesetLimit < math.MinInt32 {
			return nil, fmt.Errorf("ruleset_limit %d is out of range", rulesetLimit)
		}
		settings.RulesetLimit = autoarch.PtrInt32(int32(rulesetLimit))
	}

	return settings, nil
}

func encodeTerraform(settings *autoarch.AutomatedArchivalSettings, data *schema.ResourceData) error {
	errs := multierr.Combine(
		data.Set("enabled", settings.Enabled),
		data.Set("lookback_period", settings.LookbackPeriod),
		data.Set("grace_period", settings.GracePeriod),
		data.Set("version", strconv.FormatInt(settings.Version, 10)),
		data.Set("ruleset_limit", settings.RulesetLimit),
	)
	if settings.Creator != nil {
		errs = multierr.Append(errs, data.Set("creator", *settings.Creator))
	}
	if settings.LastUpdatedBy != nil {
		errs = multierr.Append(errs, data.Set("last_updated_by", *settings.LastUpdatedBy))
	}
	if settings.Created != nil {
		errs = multierr.Append(errs, data.Set("created", *settings.Created))
	}
	if settings.LastUpdated != nil {
		errs = multierr.Append(errs, data.Set("last_updated", *settings.LastUpdated))
	}

	return errs
}
