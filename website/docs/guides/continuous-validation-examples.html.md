---
subcategory: ""
layout: "aws"
page_title: "Using Terraform Cloud's Continuous Validation feature with the AWS Provider"
description: |-
  Using Terraform Cloud's Continuous Validation feature with the AWS Provider
---

# Using Terraform Cloud's Continuous Validation feature with the AWS Provider

## Continuous Validation in Terraform Cloud

The Continuous Validation feature in Terraform Cloud (TFC) allows users to make assertions about their infrastructure between applied runs. This helps users to identify issues at the time they first appear and avoid situations where a change is only identified once it causes a customer-facing problem.

Users can add checks to their Terraform configuration using check blocks. Check blocks contain assertions that are defined with a custom condition expression and an error message. When the condition expression evaluates to true the check passes, but when the expression evaluates to false Terraform will show a warning message that includes the user-defined error message.

Custom conditions can be created using data from Terraform providers’ resources and data sources. Data can also be combined from multiple sources; for example, you can use checks to monitor expirable resources by comparing a resource’s expiration date attribute to the current time returned by Terraform’s built-in time functions.

Below, this guide shows examples of how data returned by the AWS provider can be used to define checks in your Terraform configuration.

## Example - Ensure your AWS account is within budget (aws_budgets_budget)

AWS Budgets allows you to track and take action on your AWS costs and usage. You can use AWS Budgets to monitor your aggregate utilization and coverage metrics for your Reserved Instances (RIs) or Savings Plans.

- You can use AWS Budgets to enable simple-to-complex cost and usage tracking. Some examples include:

- Setting a monthly cost budget with a fixed target amount to track all costs associated with your account.

- Setting a monthly cost budget with a variable target amount, with each subsequent month growing the budget target by 5 percent.

- Setting a monthly usage budget with a fixed usage amount and forecasted notifications to help ensure that you are staying within the service limits for a specific service.

- Setting a daily utilization or coverage budget to track your RI or Savings Plans.

The example below shows how a check block can be used to assert that you remain in compliance for the budgets that have been established.

```hcl
check "check_budget_exceeded" {
  data "aws_budgets_budget" "example" {
    name = aws_budgets_budget.example.name
  }

  assert {
    condition = !data.aws_budgets_budget.example.budget_exceeded
    error_message = format("AWS budget has been exceeded! Calculated spend: '%s' and budget limit: '%s'",
      data.aws_budgets_budget.example.calculated_spend[0].actual_spend[0].amount,
      data.aws_budgets_budget.example.budget_limit[0].amount
    )
  }
}
```

If the budget exceeds the set limit, the check block assertion will return a warning similar to the following:

```
│ Warning: Check block assertion failed
│ 
│   on main.tf line 43, in check "check_budget_exceeded":
│   43:     condition = !data.aws_budgets_budget.example.budget_exceeded
│     ├────────────────
│     │ data.aws_budgets_budget.example.budget_exceeded is true
│ 
│ AWS budget has been exceeded! Calculated spend: '1550.0' and budget limit: '1200.0'
```

## Example - Check GuardDuty for Threats (aws_guardduty_finding_ids)

Amazon GuardDuty is a threat detection service that continuously monitors for malicious activity and unauthorized behavior to protect your Amazon Web Services accounts, workloads, and data stored in Amazon S3. With the cloud, the collection and aggregation of account and network activities is simplified, but it can be time consuming for security teams to continuously analyze event log data for potential threats. With GuardDuty, you now have an intelligent and cost-effective option for continuous threat detection in Amazon Web Services Cloud.

The following example outlines how a check block can be utilized to assert that no threats have been identified from AWS GuardDuty.

```hcl
data "aws_guardduty_detector" "example" {}

check "check_guardduty_findings" {
  data "aws_guardduty_finding_ids" "example" {
    detector_id = data.aws_guardduty_detector.example.id
  }

  assert {
    condition = !data.aws_guardduty_finding_ids.example.has_findings
    error_message = format("AWS GuardDuty detector '%s' has %d open findings!",
      data.aws_guardduty_finding_ids.example.detector_id,
      length(data.aws_guardduty_finding_ids.example.finding_ids),
    )
  }
}
```

If findings are present, the check block assertion will return a warning similar to the following:

```
│ Warning: Check block assertion failed
│
│   on main.tf line 24, in check "check_guardduty_findings":
│   24:     condition = !data.aws_guardduty_finding_ids.example.has_findings
│     ├────────────────
│     │ data.aws_guardduty_finding_ids.example.has_findings is true
│
│ AWS GuardDuty detector 'abcdef123456' has 9 open findings!
```
