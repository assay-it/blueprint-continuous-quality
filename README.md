<p align="center">
  <img src="./blueprint.gif" height="240" />
  <h3 align="center">Blueprint: Continuous Proof of the Quality</h3>
  <p align="center"><strong>A simple strategy on continuous quality automation with CI/CD system.</strong></p>

  <p align="center">
  </p>
</p>

Quality assurance of serverless applications is more complex than doing it for other runtimes. Engineering teams spend twice as much time maintaining testing environments and mocks of cloud dependencies instead of building a loyal relationship with their customers, assay.it has you covered.

https://assay.it is a Software as a Service for developers to perform formal proofs of quality using type safe Behavior as a Code. It automates validation of cause-and-effect in loosely coupled topologies such as serverless applications, microservices and other systems that rely on interface syntaxes and its behaviors. It emphasizes deployment and quality assessment as a key feature along the development pipelines. Continuous proofs of the quality helps to eliminate defects at earlier phases of the feature lifecycle. It impacts on engineering teams philosophy and commitments, ensuring that your microservice(s) are always in a release-ready state.

This example is inspired by the blog post [Everything is Continuous: A Simple Strategy on Software Delivery](https://assay.it/case%20study/continuous%20deployment/2020/07/01/everything-is-continuos) and provides the reference implementation of the eponymous [workflow](https://assay.it/doc/case-study/everything-is-continuous).


### Continuous Deployment

Continuous Deployment defines a right philosophy and commitment that ensures the always ready state of your code. It also implements pipelines to deploy every commit straight to sandbox with the following promotion to production.Let's us emphasizes deployment and quality assessment as a key feature along the development process. Continuous proofs of the quality helps to eliminate defects at earlier phases of the feature lifecycle. This example application implements a reference workflow to emphasis role of [https://assay.it](https://assay.it) as part of quality check of each code change.

Let's take a look:
 
1. The `main` branch of your project is always the latest deployable snapshot of a microservice. CI/CD have to automate the `main` snapshot deployments every time when a new feature is merged. 

2. The feature integration into `main` is implemented through pull request (no exceptions whatsoever).

3. GitHub Actions (or other CI/CD) executes automated pull request deployment to the sandbox environment every time new changes are proposed (each commit). The sandbox environment is a disposable deployment dedicated only for pull request validation.

4. Once the sandbox environment is ready, the integration and behavioral testing phase begins.

5. https://assay.it make a formal proofs of the quality using Behavior as a Code paradigm.

6. Everything is alright, the pull request is merged to `main` branch. 

Exactly same quality assessment technique is applied to the `main` branch and each release - Behavior as a Code runs the quality assessment for each deployment.


## Getting started

1. **Sign up for [assay.it](https://assay.it)** with your GitHub developer account. Initially, the service requires only access to your public profile, public repositories and access to commit status of connected repositories. Later, you can enable quality assessments of private repositories.

2. **Fork [assay-it/blueprint-continuous-quality](https://github.com/assay-it/blueprint-continuous-quality)** to your own GitHub account and then add to the assay.it workspace. The example implements advanced quality assessment suites using [category pattern](https://assay.it/doc/core/category) to connect cause-and-effect (Given/When/Then) with the networking concepts (Input/Process/Output). Just write [pure functional code](https://assay.it/doc/core) instead of clicking through UI or maintaining endless XML, YAML or JSON documents.

3. Either create a new **AWS Account** or obtain an access to existing one. You shall be able to create/delete AWS resources (e.g. IAM roles, Lambda functions, API Gateway, etc). AWS account is required to reveal the whole essence of this example. Ultimately, you are deploying a fictional serverless application and validating its deployment. `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY` are required to enable deployment by GitHub Actions. Store these credentials to secret key vault at your fork settings (Your Fork > Settings > Secrets).

4. **Allow GitHub Action** to run quality assessment with assay.it on your behalf, an access key is requires. Go to your profile settings at assay.it and generate a new personal access key. Store this key to secret key vault at your fork settings (Your Fork > Settings > Secrets) under the name `ASSAY_SECRET_KEY`.

5. **Open a pull request** against your fork with some change. The change triggers GitHub Action [pipeline](.github/workflow/check.yml) that builds this example projects and deploys it to AWS account using Infrastructure as a Code. Upon the deployment, GitHub Action passes quality assessment job of newly deployed application to assay.it. The service runs the quality check and updates the status of the pull request with results.


**Let's have a look on the content of repository**:
* [News Feed](newsfeed.go) is a fictional serverless application that implement REST API using Golang;
* [Infrastructure as a Code](cloud) uses TypeScript and AWS CDK to declare and deploy a cloud runtime environment for the application;
* [GitHub Actions Workflow](.github/workflows/check.yml) orchestrates a typical CI/CD activities - building, deploying and quality assurance for each pull request (commit);
* [Behavior as a Code suites](suite) declares a desired behavior (contract) of the api, assay.it uses these code to confirm the quality.

```
/
├── .github/workflows
|   └── check.yml      // GitHub Action, the CI/CD workflow
|
├── cloud
|   ├── index.ts       // Infrastructure as a Code
|   └── ...
|
├── suite              // Behavior as a Code
|   ├── elementary.go  // shows (unit) contract testing of api
|   ├── minimal.go     // a minimal BaC suite
|   ├── recursive.go   // demonstrates ability of recursive scenarios
|   └── scenario.go    // shows behavioral scenario, chain of requests to model consumer 
|
├── .assay.json        // the configuration file, lists scenarios to execute
├── newsfeed.go        // fictional News Feed serverless application
```

## Further Reading

Please continue to [the core](https://assay.it/doc/core) sections for details about Behavior as a Code development. Study the annotated examples given here, they demonstrate a typical design patterns.


## License

[![See LICENSE](https://img.shields.io/github/license/assay-it/example.assay.it.svg?style=for-the-badge)](LICENSE)

