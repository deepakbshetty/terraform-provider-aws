// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package waf

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	awstypes "github.com/aws/aws-sdk-go-v2/service/waf/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_waf_size_constraint_set", name="Size Constraint Set")
func resourceSizeConstraintSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSizeConstraintSetCreate,
		ReadWithoutTimeout:   resourceSizeConstraintSetRead,
		UpdateWithoutTimeout: resourceSizeConstraintSetUpdate,
		DeleteWithoutTimeout: resourceSizeConstraintSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size_constraints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"field_to_match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Type:     schema.TypeString,
										Optional: true,
									},
									names.AttrType: {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"text_transformation": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSizeConstraintSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFClient(ctx)

	name := d.Get(names.AttrName).(string)
	output, err := NewRetryer(conn).RetryWithToken(ctx, func(token *string) (interface{}, error) {
		input := &waf.CreateSizeConstraintSetInput{
			ChangeToken: token,
			Name:        aws.String(name),
		}

		return conn.CreateSizeConstraintSet(ctx, input)
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating WAF Size Constraint Set (%s): %s", name, err)
	}

	d.SetId(aws.ToString(output.(*waf.CreateSizeConstraintSetOutput).SizeConstraintSet.SizeConstraintSetId))

	return append(diags, resourceSizeConstraintSetUpdate(ctx, d, meta)...)
}

func resourceSizeConstraintSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFClient(ctx)

	sizeConstraintSet, err := findSizeConstraintSetByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] WAF Size Constraint Set (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading WAF Size Constraint Set (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "waf",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  "sizeconstraintset/" + d.Id(),
	}
	d.Set(names.AttrARN, arn.String())
	d.Set(names.AttrName, sizeConstraintSet.Name)
	if err := d.Set("size_constraints", flattenSizeConstraints(sizeConstraintSet.SizeConstraints)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting size_constraints: %s", err)
	}

	return diags
}

func resourceSizeConstraintSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFClient(ctx)

	if d.HasChange("size_constraints") {
		o, n := d.GetChange("size_constraints")
		oldConstraints, newConstraints := o.(*schema.Set).List(), n.(*schema.Set).List()
		if err := updateSizeConstraintSet(ctx, conn, d.Id(), oldConstraints, newConstraints); err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	return append(diags, resourceSizeConstraintSetRead(ctx, d, meta)...)
}

func resourceSizeConstraintSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFClient(ctx)

	if oldConstraints := d.Get("size_constraints").(*schema.Set).List(); len(oldConstraints) > 0 {
		if err := updateSizeConstraintSet(ctx, conn, d.Id(), oldConstraints, []interface{}{}); err != nil && !errs.IsA[*awstypes.WAFNonexistentItemException](err) && !errs.IsA[*awstypes.WAFNonexistentContainerException](err) {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	log.Printf("[INFO] Deleting WAF Size Constraint Set: %s", d.Id())
	_, err := NewRetryer(conn).RetryWithToken(ctx, func(token *string) (interface{}, error) {
		input := &waf.DeleteSizeConstraintSetInput{
			ChangeToken:         token,
			SizeConstraintSetId: aws.String(d.Id()),
		}

		return conn.DeleteSizeConstraintSet(ctx, input)
	})

	if errs.IsA[*awstypes.WAFNonexistentItemException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting WAF Size Constraint Set (%s): %s", d.Id(), err)
	}

	return diags
}

func findSizeConstraintSetByID(ctx context.Context, conn *waf.Client, id string) (*awstypes.SizeConstraintSet, error) {
	input := &waf.GetSizeConstraintSetInput{
		SizeConstraintSetId: aws.String(id),
	}

	output, err := conn.GetSizeConstraintSet(ctx, input)

	if errs.IsA[*awstypes.WAFNonexistentItemException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.SizeConstraintSet == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.SizeConstraintSet, nil
}

func updateSizeConstraintSet(ctx context.Context, conn *waf.Client, id string, oldS, newS []interface{}) error {
	_, err := NewRetryer(conn).RetryWithToken(ctx, func(token *string) (interface{}, error) {
		input := &waf.UpdateSizeConstraintSetInput{
			ChangeToken:         token,
			SizeConstraintSetId: aws.String(id),
			Updates:             diffSizeConstraints(oldS, newS),
		}

		return conn.UpdateSizeConstraintSet(ctx, input)
	})

	if err != nil {
		return fmt.Errorf("updating WAF Size Constraint Set (%s): %w", id, err)
	}

	return nil
}

func diffSizeConstraints(oldS, newS []interface{}) []awstypes.SizeConstraintSetUpdate {
	updates := make([]awstypes.SizeConstraintSetUpdate, 0)

	for _, os := range oldS {
		constraint := os.(map[string]interface{})

		if idx, contains := sliceContainsMap(newS, constraint); contains {
			newS = append(newS[:idx], newS[idx+1:]...)
			continue
		}

		updates = append(updates, awstypes.SizeConstraintSetUpdate{
			Action: awstypes.ChangeActionDelete,
			SizeConstraint: &awstypes.SizeConstraint{
				FieldToMatch:       expandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: awstypes.ComparisonOperator(constraint["comparison_operator"].(string)),
				Size:               int64(constraint["size"].(int)),
				TextTransformation: awstypes.TextTransformation(constraint["text_transformation"].(string)),
			},
		})
	}

	for _, ns := range newS {
		constraint := ns.(map[string]interface{})

		updates = append(updates, awstypes.SizeConstraintSetUpdate{
			Action: awstypes.ChangeActionInsert,
			SizeConstraint: &awstypes.SizeConstraint{
				FieldToMatch:       expandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: awstypes.ComparisonOperator(constraint["comparison_operator"].(string)),
				Size:               int64(constraint["size"].(int)),
				TextTransformation: awstypes.TextTransformation(constraint["text_transformation"].(string)),
			},
		})
	}
	return updates
}

func flattenSizeConstraints(sc []awstypes.SizeConstraint) []interface{} {
	out := make([]interface{}, len(sc))
	for i, c := range sc {
		m := make(map[string]interface{})
		m["comparison_operator"] = c.ComparisonOperator
		if c.FieldToMatch != nil {
			m["field_to_match"] = flattenFieldToMatch(c.FieldToMatch)
		}
		m["size"] = c.Size
		m["text_transformation"] = c.TextTransformation
		out[i] = m
	}
	return out
}
