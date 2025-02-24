package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
  "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
  "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
  "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	awscdk.StackProps
}

func NewGoCdkStack(scope constructs.Construct, id string, props *GoCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)


  // create table here
  userTable := awsdynamodb.NewTable(stack, jsii.String("myUserTable"), &awsdynamodb.TableProps{
    PartitionKey: &awsdynamodb.Attribute{
      Name: jsii.String("username"),
      Type: awsdynamodb.AttributeType_STRING,
    },

    // this table name maps to const in database.go const table name
    TableName: jsii.String("userTable"),
  })

  blogTable := awsdynamodb.NewTable(stack, jsii.String("myBlogTable"), &awsdynamodb.TableProps{
    PartitionKey: &awsdynamodb.Attribute{
      Name: jsii.String("slug"),
      Type: awsdynamodb.AttributeType_STRING,
    },

    // this table name maps to const in database.go const table name
    TableName: jsii.String("blogsTable"),
  })



	// The code that defines your stack goes here

  myFunction := awslambda.NewFunction(stack, jsii.String("myLambdaFunction"), &awslambda.FunctionProps{
    //go run time, meaning the lambda function can run in go, it serverless architure to run a specific language as you can't install language on a server
    //AL means amazon linux
    Runtime: awslambda.Runtime_PROVIDED_AL2023(),
    //jsii compiles from go to typescript as cdk is built in typescript, options here is where the lambda code is from, it can be in s3 buckets
    Code: awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
    Handler: jsii.String("main"),
  })
  
  userTable.GrantReadWriteData(myFunction)
  blogTable.GrantReadWriteData(myFunction)

  api := awsapigateway.NewRestApi(stack, jsii.String("myAPIGateway"), &awsapigateway.RestApiProps{
    DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
      AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
      AllowMethods: jsii.Strings("POST", "GET", "PUT", "DELETE", "OPTIONS"),
      AllowOrigins: jsii.Strings("*"),
    },
    // need to enable cloudwatch logging for this to work
    // DeployOptions: &awsapigateway.StageOptions{
    //   LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
    // },
  })

  integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

  //define routes
  registerResource := api.Root().AddResource(jsii.String("register"), nil)
  registerResource.AddMethod(jsii.String("POST"), integration, nil)

  loginResource := api.Root().AddResource(jsii.String("login"), nil)
  loginResource.AddMethod(jsii.String("POST"), integration, nil)

  blogResource := api.Root().AddResource(jsii.String("blog"), nil)
  blogResource.AddMethod(jsii.String("POST"), integration, nil)

  blogWithSlugResource := blogResource.AddResource(jsii.String("{slug}"), nil)
  blogWithSlugResource.AddMethod(jsii.String("GET"), integration, nil)
  blogWithSlugResource.AddMethod(jsii.String("PUT"), integration, nil)
  blogWithSlugResource.AddMethod(jsii.String("DELETE"), integration, nil)

  blogsResource := api.Root().AddResource(jsii.String("blogs"), nil)
  blogsResource.AddMethod(jsii.String("GET"), integration, nil)

  protectedResource := api.Root().AddResource(jsii.String("protected"), nil)
  protectedResource.AddMethod(jsii.String("GET"), integration, nil)

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("GoCdkQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoCdkStack(app, "GoCdkStack", &GoCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
