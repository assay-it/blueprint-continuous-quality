import * as cdk from '@aws-cdk/core'
import * as iam from '@aws-cdk/aws-iam'
import * as lambda from '@aws-cdk/aws-lambda'
import * as api from '@aws-cdk/aws-apigateway'
import * as logs from '@aws-cdk/aws-logs'
import * as pure from 'aws-cdk-pure'
import * as hoc from 'aws-cdk-pure-hoc'
import * as path from 'path'

//
//
const app = new cdk.App()
const config = {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  }
}
const vsn: string = app.node.tryGetContext('vsn') || 'latest'
const domain: string = app.node.tryGetContext('domain') || 'example.com'

// ----------------------------------------------------------------------------
//
// Common Stack
//
// ----------------------------------------------------------------------------
const common = new cdk.Stack(app, 'example-common', { ...config })

const tlsCertificate = pure.join(common,
  hoc.common.HostedZone(domain).flatMap(
    zone => hoc.common.Certificate(`*.${domain}`, zone)
  ),
).certificateArn


// ----------------------------------------------------------------------------
//
// Common Stack
//
// ----------------------------------------------------------------------------
const stack = new cdk.Stack(app, `example-${vsn}`, { ...config })

//
//
const Role = (): iam.RoleProps => ({
  assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
  managedPolicies: [
    iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
  ],
})

const func = pure.iaac(iam.Role)(Role).flatMap(
  role => {
    const Lambda = (): lambda.FunctionProps => ({
      code: hoc.common.AssetCodeGo(path.join(__dirname, '..')),
      handler: 'main',
      runtime: lambda.Runtime.GO_1_X,
      logRetention: logs.RetentionDays.FIVE_DAYS,
      functionName: `${cdk.Aws.STACK_NAME}-func`,
      role,
    })
    return (pure.iaac(lambda.Function)(Lambda))
  }
)

//
//
const restapi = hoc.gateway.Api({
  domain,
  subdomain: vsn,
  tlsCertificate
})


pure.join(stack,
  pure.use({
    restapi,
    func: pure.wrap(api.LambdaIntegration)(func),
  }).effect(eff => {
    const seq = eff.restapi.root.addResource('news')
    const els = seq.addResource('{any+}')
    seq.addMethod('ANY', eff.func)
    els.addMethod('ANY', eff.func)
  })
)

app.synth()