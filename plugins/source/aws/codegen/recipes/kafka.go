package recipes

import (
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	"github.com/cloudquery/plugin-sdk/codegen"
	"github.com/cloudquery/plugin-sdk/schema"
)

func KafkaResources() []*Resource {
	resources := []*Resource{
		{
			SubService:          "clusters",
			Struct:              &types.Cluster{},
			Description:         "https://docs.aws.amazon.com/MSK/2.0/APIReference/v2-clusters-clusterarn.html#v2-clusters-clusterarn-properties",
			SkipFields:          []string{"ClusterArn"},
			PreResourceResolver: "getCluster",
			ExtraColumns: append(defaultAccountColumns, []codegen.ColumnDefinition{{
				Name:     "arn",
				Type:     schema.TypeString,
				Options:  schema.ColumnCreationOptions{PrimaryKey: true},
				Resolver: `schema.PathResolver("ClusterArn")`,
			},
			}...),
			Relations: []string{
				"Nodes()",
				"ClusterOperations()",
			},
			ShouldGenerateResolverAndMockTest: true,
			ResolverAndMockTestTemplate:       "list_and_describe_resources_1",
		},
		{
			SubService:  "nodes",
			Struct:      &types.NodeInfo{},
			Description: "https://docs.aws.amazon.com/msk/1.0/apireference/clusters-clusterarn-nodes.html#ListNodes",
			SkipFields:  []string{"NodeARN"},
			ExtraColumns: append(defaultAccountColumns, []codegen.ColumnDefinition{
				{
					Name:     "arn",
					Type:     schema.TypeString,
					Options:  schema.ColumnCreationOptions{PrimaryKey: true},
					Resolver: `schema.PathResolver("NodeARN")`,
				},
				{
					Name:     "cluster_arn",
					Type:     schema.TypeString,
					Resolver: `schema.ParentColumnResolver("arn")`,
				},
			}...),
			ShouldGenerateResolverAndMockTest: true,
			ResolverAndMockTestTemplate:       "list_resources_1",
			CustomListInput:                   `getListNodesInput(parent)`,
		},
		{
			SubService:  "cluster_operations",
			Struct:      &types.ClusterOperationInfo{},
			Description: "https://docs.aws.amazon.com/msk/1.0/apireference/clusters-clusterarn-operations.html",
			SkipFields:  []string{"OperationArn", "ClusterArn"},
			ExtraColumns: append(defaultAccountColumns, []codegen.ColumnDefinition{
				{
					Name:     "arn",
					Type:     schema.TypeString,
					Options:  schema.ColumnCreationOptions{PrimaryKey: true},
					Resolver: `schema.PathResolver("OperationArn")`,
				},
				{
					Name:     "cluster_arn",
					Type:     schema.TypeString,
					Resolver: `schema.ParentColumnResolver("arn")`,
				},
			}...),
			ShouldGenerateResolverAndMockTest: true,
			ResolverAndMockTestTemplate:       "list_resources_1",
			CustomListInput:                   `getListClusterOperationsInput(parent)`,
		},
		{
			SubService:  "configurations",
			Struct:      &types.Configuration{},
			Description: "https://docs.aws.amazon.com/msk/1.0/apireference/clusters-clusterarn-configuration.html",
			SkipFields:  []string{"Arn"},
			ExtraColumns: append(defaultAccountColumns, []codegen.ColumnDefinition{
				{
					Name:     "arn",
					Type:     schema.TypeString,
					Options:  schema.ColumnCreationOptions{PrimaryKey: true},
					Resolver: `schema.PathResolver("Arn")`,
				},
			}...),
			ShouldGenerateResolverAndMockTest: true,
			ResolverAndMockTestTemplate:       "list_resources_1",
		},
	}

	// set default values
	for _, r := range resources {
		r.Service = "kafka"
		r.Multiplex = `client.ServiceAccountRegionMultiplexer("kafka")`

		// Parameters for auto-generating the resolver and mock-test.
		// Only used when `ShouldGenerateResolverAndMockTest = true`
		r.Client = &kafka.Client{}
	}

	return resources
}
