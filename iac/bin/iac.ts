#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { IacStack } from '../lib/iac-stack';

const app = new cdk.App();

new IacStack(app, IacStack.stack_name, {
  env: {
    // region: app.node.tryGetContext("aws_region"),
    // account: app.node.tryGetContext("aws_account"),
    account: process.env.CDK_DEFAULT_ACCOUNT, 
    region: process.env.CDK_DEFAULT_REGION 
  }
});
