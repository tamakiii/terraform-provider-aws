// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chatbot_test

import (
	"fmt"
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
	resourceName := "aws_chatbot_slack_channel_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Chatbot)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ChatbotServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSlackChannelConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSlackChannelConfigurationDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "chat_configuration_arn", resourceName, "chat_configuration_arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "configuration_name", resourceName, "configuration_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "iam_role_arn", resourceName, "iam_role_arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logging_level", resourceName, "logging_level"),
					resource.TestCheckResourceAttrPair(dataSourceName, "slack_channel_id", resourceName, "slack_channel_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "slack_channel_name", resourceName, "slack_channel_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "slack_team_id", resourceName, "slack_team_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "slack_team_name", resourceName, "slack_team_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_authorization_required", resourceName, "user_authorization_required"),
					acctest.MatchResourceAttrRegionalARN(ctx, dataSourceName, "chat_configuration_arn", "chatbot", regexache.MustCompile(`chat-configuration/slack-channel/.+$`)),
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
			acctest.PreCheckPartitionHasService(t, names.Chatbot)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ChatbotServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSlackChannelConfigurationDataSourceConfig_notFound(),
				ExpectError: regexache.MustCompile(`missing`),
			},
		},
	})
}

func testAccSlackChannelConfigurationDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccSlackChannelConfigurationConfig_basic(rName, "C07EZ1NHXEZ", "T07EA7JMZ"), fmt.Sprintf(`
data "aws_chatbot_slack_channel_configuration" "test" {
  chat_configuration_arn = aws_chatbot_slack_channel_configuration.test.chat_configuration_arn
}
`))
}

func testAccSlackChannelConfigurationDataSourceConfig_notFound() string {
	return `
data "aws_chatbot_slack_channel_configuration" "test" {
  chat_configuration_arn = "arn:aws:chatbot::123456789012:chat-configuration/slack-channel/non-existent"
}
`
}
