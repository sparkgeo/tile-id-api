import { Stack, StackProps, RemovalPolicy, CfnOutput } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Cluster, ContainerImage, LogDriver } from "aws-cdk-lib/aws-ecs";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { ApplicationLoadBalancedFargateService } from "aws-cdk-lib/aws-ecs-patterns";
import { LogGroup } from "aws-cdk-lib/aws-logs";

export class IacStack extends Stack {

  public static readonly stack_name: string = "tile-id-api";

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const vpc = new Vpc(this, "vpc", {
      maxAzs: 2,
    });

    const cluster = new Cluster(this, "cluster", {
      vpc: vpc,
      containerInsights: true,
    });

    const service = new ApplicationLoadBalancedFargateService(this, "ALBFargateSvc", {
      cluster: cluster,
      memoryLimitMiB: 512,
      assignPublicIp: true,
      desiredCount: 1,
      taskImageOptions: {
        image: ContainerImage.fromAsset("../", {
          file: "dockerfile.api",
        }),
        containerPort: 8080,
        logDriver: LogDriver.awsLogs({
          streamPrefix: IacStack.stack_name,
          logGroup: new LogGroup(this, "apiLogs", {
            logGroupName: "apiLogGroup",
            removalPolicy: RemovalPolicy.DESTROY,
          }),
        }),
        enableLogging: true,
      },
    })
    service.targetGroup.configureHealthCheck({
      path: "/healthz",
    });
    service.service.autoScaleTaskCount({ maxCapacity: 10 }).scaleOnCpuUtilization("CpuScaling", {
      targetUtilizationPercent: 50,
    });

    new CfnOutput(this, "apiURL", {
      value: service.loadBalancer.loadBalancerDnsName,
    });
  }
}
