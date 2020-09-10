import * as cdk from '@aws-cdk/core'
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

// ----------------------------------------------------------------------------
//
// Common Stack
//
// ----------------------------------------------------------------------------
const stack = new cdk.Stack(app, `example-${vsn}`, { ...config })

const Lambda = (): lambda.FunctionProps => ({
  code: hoc.common.AssetCodeGo(path.join(__dirname, '..')),
  handler: 'main',
  runtime: lambda.Runtime.GO_1_X,
  logRetention: logs.RetentionDays.FIVE_DAYS,
  functionName: `${cdk.Aws.STACK_NAME}-api`,
})

const Gateway = (): api.RestApiProps => ({
  deploy: true,
  deployOptions: { stageName: 'api' },
  endpointTypes: [api.EndpointType.REGIONAL],
  failOnWarnings: true,
  defaultCorsPreflightOptions: {
    allowOrigins: api.Cors.ALL_ORIGINS,
    maxAge: cdk.Duration.minutes(10),
  }
})

const rest = pure.iaac(api.RestApi)(Gateway)
const func = pure.wrap(api.LambdaIntegration)(pure.iaac(lambda.Function)(Lambda))
  
pure.join(stack,
  pure.use({ rest, func })
  .effect(
    eff => {
      const seq = eff.rest.root.addResource('news')
      seq.addMethod('ANY', eff.func)

      const els = seq.addResource('{any+}')
      els.addMethod('ANY', eff.func)
    }
  )
  .effect(
    eff =>
      new cdk.CfnOutput(stack, 'Foo', { value: eff.rest.url })
  )
)



app.synth()