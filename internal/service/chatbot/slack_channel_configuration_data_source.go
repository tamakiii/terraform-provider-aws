// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chatbot

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/chatbot"
	awstypes "github.com/aws/aws-sdk-go-v2/service/chatbot/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for datasource registration to the Provider. DO NOT EDIT.
// @FrameworkDataSource("aws_chatbot_slack_channel_configuration", name="Slack Channel Configuration")
func newDataSourceSlackChannelConfiguration(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceSlackChannelConfiguration{}, nil
}

const (
	DSNameSlackChannelConfiguration = "Slack Channel Configuration Data Source"
)

type dataSourceSlackChannelConfiguration struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceSlackChannelConfiguration) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for managing an AWS Chatbot Slack Channel Configuration.",
		Attributes: map[string]schema.Attribute{
			"chat_configuration_arn": schema.StringAttribute{
				Description: "ARN of the Slack channel configuration.",
				CustomType:  fwtypes.ARNType,
				Required:    true,
			},
			"configuration_name": schema.StringAttribute{
				Description: "Name of the Slack channel configuration.",
				Computed:    true,
			},
			"iam_role_arn": schema.StringAttribute{
				Description: "ARN of the IAM role that defines the permissions for AWS Chatbot.",
				Computed:    true,
			},
			"logging_level": schema.StringAttribute{
				Description: "Specifies the logging level for this configuration.",
				Computed:    true,
			},
			"slack_channel_id": schema.StringAttribute{
				Description: "ID of the Slack channel.",
				Computed:    true,
			},
			"slack_channel_name": schema.StringAttribute{
				Description: "Name of the Slack channel.",
				Computed:    true,
			},
			"slack_team_id": schema.StringAttribute{
				Description: "ID of the Slack workspace authorized with AWS Chatbot.",
				Computed:    true,
			},
			"slack_team_name": schema.StringAttribute{
				Description: "Name of the Slack workspace.",
				Computed:    true,
			},
			"sns_topic_arns": schema.SetAttribute{
				Description: "ARNs of the SNS topics that deliver notifications to AWS Chatbot.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "State of the configuration.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Map of tags assigned to the resource.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"user_authorization_required": schema.BoolAttribute{
				Description: "Enables use of a user role requirement in your chat configuration.",
				Computed:    true,
			},
		},
	}
}

func (d *dataSourceSlackChannelConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().ChatbotClient(ctx)

	var data dataSourceSlackChannelConfigurationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := findSlackChannelConfigurationByArn(ctx, conn, data.ChatConfigurationArn.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Chatbot, create.ErrActionReading, DSNameSlackChannelConfiguration, data.ChatConfigurationArn.ValueString(), err),
			err.Error(),
		)
		return
	}

	// Using a field name prefix allows mapping fields such as `SlackChannelConfigurationId` to `ID`
	resp.Diagnostics.Append(flex.Flatten(ctx, out, &data, flex.WithFieldNamePrefix("SlackChannelConfiguration"))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set tags if present
	if len(out.Tags) > 0 {
		data.Tags = flex.FlattenFrameworkStringValueMap(ctx, keyValueTags(ctx, out.Tags).Map())
	} else {
		data.Tags = types.MapNull(types.StringType)
	}

	// Set SNS topic ARNs if present
	if len(out.SnsTopicArns) > 0 {
		data.SnsTopicArns = flex.FlattenFrameworkStringValueSetLegacy(ctx, out.SnsTopicArns)
	} else {
		data.SnsTopicArns = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findSlackChannelConfigurationByArn(ctx context.Context, conn *chatbot.Client, chatConfigurationArn string) (*awstypes.SlackChannelConfiguration, error) {
	input := &chatbot.DescribeSlackChannelConfigurationsInput{}

	pages := chatbot.NewDescribeSlackChannelConfigurationsPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, configuration := range page.SlackChannelConfigurations {
			if aws.ToString(configuration.ChatConfigurationArn) == chatConfigurationArn {
				return &configuration, nil
			}
		}
	}

	return nil, tfresource.NewEmptyResultError(input)
}

type dataSourceSlackChannelConfigurationModel struct {
	ChatConfigurationArn      fwtypes.ARN  `tfsdk:"chat_configuration_arn"`
	ConfigurationName         types.String `tfsdk:"configuration_name"`
	IamRoleArn                types.String `tfsdk:"iam_role_arn"`
	LoggingLevel              types.String `tfsdk:"logging_level"`
	SlackChannelId            types.String `tfsdk:"slack_channel_id"`
	SlackChannelName          types.String `tfsdk:"slack_channel_name"`
	SlackTeamId               types.String `tfsdk:"slack_team_id"`
	SlackTeamName             types.String `tfsdk:"slack_team_name"`
	SnsTopicArns              types.Set    `tfsdk:"sns_topic_arns"`
	State                     types.String `tfsdk:"state"`
	Tags                      types.Map    `tfsdk:"tags"`
	UserAuthorizationRequired types.Bool   `tfsdk:"user_authorization_required"`
}
