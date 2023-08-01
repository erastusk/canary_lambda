# Canary Lambda Function

## Description

A lambda function that monitors ALB or DynamoDB AWS endpoints and switches traffic between two regions. When using Private-facing ALBs, Route53 do not have access to any resources within private subnets making it difficult
to use weighted routing policies with healthchecks.

https://aws.amazon.com/blogs/networking-and-content-delivery/performing-route-53-health-checks-on-private-resources-in-a-vpc-with-aws-lambda-and-amazon-cloudwatch/


