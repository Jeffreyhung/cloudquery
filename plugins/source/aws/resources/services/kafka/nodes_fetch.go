// Code generated by codegen; DO NOT EDIT.

package kafka

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
)

func fetchKafkaNodes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	var input kafka.ListNodesInput = getListNodesInput(parent)
	c := meta.(*client.Client)
	svc := c.Services().Kafka
	for {
		response, err := svc.ListNodes(ctx, &input)
		if err != nil {
			return err
		}
		res <- response.NodeInfoList
		if aws.ToString(response.NextToken) == "" {
			break
		}
		input.NextToken = response.NextToken
	}
	return nil
}
