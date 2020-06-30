import * as cdk from '@aws-cdk/core'
import * as iam from '@aws-cdk/aws-iam'
import * as lambda from '@aws-cdk/aws-lambda'
import * as api from '@aws-cdk/aws-apigateway'
import * as logs from '@aws-cdk/aws-logs'
import * as pure from 'aws-cdk-pure'
import * as path from 'path'

//
//
const app = new cdk.App()
const stack = new cdk.Stack(app, 'example', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  }
})

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
      code: new lambda.AssetCode('', { bundling: gocc() }),
      handler: 'main',
      runtime: lambda.Runtime.GO_1_X,
      logRetention: logs.RetentionDays.FIVE_DAYS,
      functionName: `${cdk.Aws.STACK_NAME}-func`,
      role,
    })
    return (pure.iaac(lambda.Function)(Lambda))
  }
)

const gocc = (): cdk.BundlingOptions => {
  const gopath = process.env.GOPATH || '/go'
  const fnpath = path.join(__dirname, '..').split(gopath).join('')
  console.log('gopath => ', gopath)
  console.log('fnpath => ', fnpath)

  return {
    image: cdk.BundlingDockerImage.fromAsset(`${gopath}${fnpath}`),
    command: ["go", "build", "-o", `${cdk.AssetStaging.BUNDLING_OUTPUT_DIR}/main`],
    user: 'root',
    volumes: [
      {
        containerPath: '/go/src',
        hostPath: `${gopath}/src`,
      },
    ],
    workingDirectory: `/go${fnpath}`,
  }
}

//
//
const Gateway = (): api.RestApiProps => ({
  deploy: true,
  deployOptions: {
    stageName: 'api',
  },
  endpointTypes: [api.EndpointType.REGIONAL],
  failOnWarnings: true,  
})

pure.join(stack,
  pure.use({
    rest: pure.iaac(api.RestApi)(Gateway),
    func: pure.wrap(api.LambdaIntegration)(func),
  }).effect(eff => {
    const seq = eff.rest.root.addResource('news')
    const els = seq.addResource('{any+}')
    seq.addMethod('ANY', eff.func)
    els.addMethod('ANY', eff.func)
  })
)

app.synth()