import { Stack, StackProps, RemovalPolicy } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Cluster, ContainerImage, LogDriver } from "aws-cdk-lib/aws-ecs";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { ApplicationLoadBalancedFargateService } from "aws-cdk-lib/aws-ecs-patterns";
import { LogGroup } from "aws-cdk-lib/aws-logs";
import { ARecord, HostedZone, RecordTarget } from 'aws-cdk-lib/aws-route53';
import { DnsValidatedCertificate } from 'aws-cdk-lib/aws-certificatemanager';
import { ListenerAction } from 'aws-cdk-lib/aws-elasticloadbalancingv2';
import { LoadBalancerTarget } from "aws-cdk-lib/aws-route53-targets";

export class IacStack extends Stack {

  public static readonly stack_name: string = "tile-id-api";

  private readonly baseDomain = "sparkgeo.dev";
  private readonly subDomain = `tile-id.${this.baseDomain}`;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const hostedZone = HostedZone.fromLookup(
      this,
      "hostedZone", {
        domainName: this.baseDomain,
      }
    ); 

    const service = new ApplicationLoadBalancedFargateService(this, "ALBFargateSvc", {
      cluster: new Cluster(this, "cluster", {
        vpc: new Vpc(this, "vpc", {
          maxAzs: 2,
        }),
        containerInsights: true,
      }),
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
    service.loadBalancer.addListener("https-listener", {
      port: 443,
      certificates: [
        new DnsValidatedCertificate(this, "certificate", {
          domainName: this.subDomain,
          hostedZone: hostedZone,
        }),
      ],
      defaultAction: ListenerAction.forward([service.targetGroup]),
    })
    
    new ARecord(this, "domain-alias", {
      recordName: this.subDomain,
      target: RecordTarget.fromAlias(new LoadBalancerTarget(service.loadBalancer)),
      zone: hostedZone,
    });
  }
}
