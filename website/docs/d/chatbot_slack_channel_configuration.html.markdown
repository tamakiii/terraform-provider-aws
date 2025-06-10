---
subcategory: "Chatbot"
layout: "aws"
page_title: "AWS: aws_chatbot_slack_channel_configuration"
description: |-
  Terraform data source for managing an AWS Chatbot Slack Channel Configuration.
---

# Data Source: aws_chatbot_slack_channel_configuration

Terraform data source for managing an AWS Chatbot Slack Channel Configuration.

## Example Usage

### Basic Usage

```terraform
data "aws_chatbot_slack_channel_configuration" "example" {
  chat_configuration_arn = "arn:aws:chatbot::123456789012:chat-configuration/slack-channel/example"
}
```

## Argument Reference

The following arguments are required:

* `chat_configuration_arn` - (Required) ARN of the Slack channel configuration.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `configuration_name` - Name of the Slack channel configuration.
* `iam_role_arn` - ARN of the IAM role that defines the permissions for AWS Chatbot.
* `logging_level` - Logging level for this configuration.
* `slack_channel_id` - ID of the Slack channel.
* `slack_channel_name` - Name of the Slack channel.
* `slack_team_id` - ID of the Slack workspace authorized with AWS Chatbot.
* `slack_team_name` - Name of the Slack workspace.
* `sns_topic_arns` - ARNs of the SNS topics that deliver notifications to AWS Chatbot.
* `state` - State of the configuration.
* `tags` - Map of tags assigned to the resource.
* `user_authorization_required` - Enables use of a user role requirement in your chat configuration.
