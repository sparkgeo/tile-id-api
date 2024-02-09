#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { IacStack } from '../lib/iac-stack';

const app = new cdk.App();
new IacStack(app, IacStack.stack_name, {
  env: {
    region: "us-west-2",
  }
});