// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chatbot_test

import (
	"testing"

	"github.com/YakDriver/regexache"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccChatbotSlackChannelConfigurationDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_chatbot_slack_channel_configuration.test"

	// The slack workspace must be created via the AWS Console. It cannot be created via APIs or Terraform.
	// Once it is created, export the workspace details in the env variables for this test
	teamID := acctest.SkipIfEnvVarNotSet(t, envSlackTeamID)
	channelID := acctest.SkipIfEnvVarNotSet(t, envSlackChannelID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ChatbotServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSlackChannelConfigurationDataSourceConfig_basic(rName, channelID, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "chat_configuration_arn"),
					resource.TestCheckResourceAttr(dataSourceName, "configuration_name", "TestDatasourceChatbotAndCodestar-Test"),
					resource.TestCheckResourceAttrSet(dataSourceName, names.AttrIAMRoleARN),
					resource.TestCheckResourceAttr(dataSourceName, "logging_level", "NONE"),
					resource.TestCheckResourceAttr(dataSourceName, "slack_channel_id", channelID),
					resource.TestCheckResourceAttr(dataSourceName, "slack_channel_name", "test"),
					resource.TestCheckResourceAttr(dataSourceName, "slack_team_id", teamID),
					resource.TestCheckResourceAttr(dataSourceName, "slack_team_name", "tamakiii"),
					resource.TestCheckResourceAttr(dataSourceName, "user_authorization_required", "false"),
					resource.TestCheckResourceAttr(dataSourceName, names.AttrState, "ENABLED"),
					resource.TestCheckResourceAttr(dataSourceName, "sns_topic_arns.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Service", "TestDatasourceChatbotAndCodestar"),
					// Use pattern matching for account-specific ARNs
					acctest.MatchResourceAttrGlobalARN(ctx, dataSourceName, "chat_configuration_arn", "chatbot", regexache.MustCompile(`chat-configuration/slack-channel/TestDatasourceChatbotAndCodestar-Test$`)),
					acctest.MatchResourceAttrGlobalARN(ctx, dataSourceName, names.AttrIAMRoleARN, "iam", regexache.MustCompile(`role/TestDatasourceChatbotAndCodestar-ChatBot$`)),
				),
			},
		},
	})
}

func TestAccChatbotSlackChannelConfigurationDataSource_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_chatbot_slack_channel_configuration.test"

	teamID := acctest.SkipIfEnvVarNotSet(t, envSlackTeamID)
	channelID := acctest.SkipIfEnvVarNotSet(t, envSlackChannelID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ChatbotServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSlackChannelConfigurationDataSourceConfig_tags(rName, channelID, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "chat_configuration_arn"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.Service", "TestDatasourceChatbotAndCodestar"),
					acctest.MatchResourceAttrGlobalARN(ctx, dataSourceName, "chat_configuration_arn", "chatbot", regexache.MustCompile(`chat-configuration/slack-channel/TestDatasourceChatbotAndCodestar-Test$`)),
				),
			},
		},
	})
}

func TestAccChatbotSlackChannelConfigurationDataSource_notFound(t *testing.T) {
	ctx := acctest.Context(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ChatbotServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSlackChannelConfigurationDataSourceConfig_notFound(),
				ExpectError: regexache.MustCompile(`empty result`),
			},
		},
	})
}

func testAccSlackChannelConfigurationDataSourceConfig_basic(rName, channelID, teamID string) string {
	return `
data "aws_caller_identity" "current" {}

data "aws_chatbot_slack_channel_configuration" "test" {
  chat_configuration_arn = "arn:aws:chatbot::${data.aws_caller_identity.current.account_id}:chat-configuration/slack-channel/TestDatasourceChatbotAndCodestar-Test"
}
`
}

func testAccSlackChannelConfigurationDataSourceConfig_tags(rName, channelID, teamID string) string {
	return `
data "aws_caller_identity" "current" {}

data "aws_chatbot_slack_channel_configuration" "test" {
  chat_configuration_arn = "arn:aws:chatbot::${data.aws_caller_identity.current.account_id}:chat-configuration/slack-channel/TestDatasourceChatbotAndCodestar-Test"
}
`
}

func testAccSlackChannelConfigurationDataSourceConfig_notFound() string {
	return `
data "aws_caller_identity" "current" {}

data "aws_chatbot_slack_channel_configuration" "test" {
  chat_configuration_arn = "arn:aws:chatbot::${data.aws_caller_identity.current.account_id}:chat-configuration/slack-channel/does-not-exist"
}
`
}
