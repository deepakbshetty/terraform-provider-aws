// Code generated by "internal/generate/listpages/main.go -AWSSDKVersion=2 -ListOps=DescribeACLs,DescribeClusters,DescribeParameterGroups,DescribeSnapshots,DescribeSubnetGroups,DescribeUsers"; DO NOT EDIT.

package memorydb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"
)

func describeACLsPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeACLsInput, fn func(*memorydb.DescribeACLsOutput, bool) bool) error {
	for {
		output, err := conn.DescribeACLs(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeClustersPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeClustersInput, fn func(*memorydb.DescribeClustersOutput, bool) bool) error {
	for {
		output, err := conn.DescribeClusters(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeParameterGroupsPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeParameterGroupsInput, fn func(*memorydb.DescribeParameterGroupsOutput, bool) bool) error {
	for {
		output, err := conn.DescribeParameterGroups(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeSnapshotsPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeSnapshotsInput, fn func(*memorydb.DescribeSnapshotsOutput, bool) bool) error {
	for {
		output, err := conn.DescribeSnapshots(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeSubnetGroupsPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeSubnetGroupsInput, fn func(*memorydb.DescribeSubnetGroupsOutput, bool) bool) error {
	for {
		output, err := conn.DescribeSubnetGroups(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeUsersPages(ctx context.Context, conn *memorydb.Client, input *memorydb.DescribeUsersInput, fn func(*memorydb.DescribeUsersOutput, bool) bool) error {
	for {
		output, err := conn.DescribeUsers(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
