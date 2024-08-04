package aws_methods

import (
	"Codex-Backend/api/types"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func CreateTable(svc *dynamodb.DynamoDB) {
	tableName := "Movies"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)

	// TODO: Update this to session.NewSession()
	svc_a := applicationautoscaling.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	// Auto-scaling - ReadCapacityUnits
	_, err = svc_a.RegisterScalableTarget(&applicationautoscaling.RegisterScalableTargetInput{
		ServiceNamespace:  aws.String("dynamodb"),
		ResourceId:        aws.String(fmt.Sprintf("table/%s", tableName)),
		ScalableDimension: aws.String("dynamodb:table:ReadCapacityUnits"),
		MinCapacity:       aws.Int64(1),
		MaxCapacity:       aws.Int64(10),
	})

	if err != nil {
		log.Fatalf("Got error registering scalable target for read capacity: %s", err)
	}
	_, err = svc_a.PutScalingPolicy(&applicationautoscaling.PutScalingPolicyInput{
		PolicyName:        aws.String("ReadAutoScalingPolicy"),
		ServiceNamespace:  aws.String("dynamodb"),
		ResourceId:        aws.String(fmt.Sprintf("table/%s", tableName)),
		ScalableDimension: aws.String("dynamodb:table:ReadCapacityUnits"),
		PolicyType:        aws.String("TargetTrackingScaling"),
		TargetTrackingScalingPolicyConfiguration: &applicationautoscaling.TargetTrackingScalingPolicyConfiguration{
			TargetValue:      aws.Float64(70.0), // Target 70% utilization
			ScaleInCooldown:  aws.Int64(60),
			ScaleOutCooldown: aws.Int64(60),
			PredefinedMetricSpecification: &applicationautoscaling.PredefinedMetricSpecification{
				PredefinedMetricType: aws.String("DynamoDBReadCapacityUtilization"),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error putting scaling policy for read capacity: %s", err)
	}

	// Auto-scaling - WriteCapacityUnits
	_, err = svc_a.RegisterScalableTarget(&applicationautoscaling.RegisterScalableTargetInput{
		ServiceNamespace:  aws.String("dynamodb"),
		ResourceId:        aws.String(fmt.Sprintf("table/%s", tableName)),
		ScalableDimension: aws.String("dynamodb:table:WriteCapacityUnits"),
		MinCapacity:       aws.Int64(1),
		MaxCapacity:       aws.Int64(10),
	})
	if err != nil {
		log.Fatalf("Got error registering scalable target for write capacity: %s", err)
	}

	_, err = svc_a.PutScalingPolicy(&applicationautoscaling.PutScalingPolicyInput{
		PolicyName:        aws.String("WriteAutoScalingPolicy"),
		ServiceNamespace:  aws.String("dynamodb"),
		ResourceId:        aws.String(fmt.Sprintf("table/%s", tableName)),
		ScalableDimension: aws.String("dynamodb:table:WriteCapacityUnits"),
		PolicyType:        aws.String("TargetTrackingScaling"),
		TargetTrackingScalingPolicyConfiguration: &applicationautoscaling.TargetTrackingScalingPolicyConfiguration{
			TargetValue:      aws.Float64(70.0), // Target 70% utilization
			ScaleInCooldown:  aws.Int64(60),
			ScaleOutCooldown: aws.Int64(60),
			PredefinedMetricSpecification: &applicationautoscaling.PredefinedMetricSpecification{
				PredefinedMetricType: aws.String("DynamoDBWriteCapacityUtilization"),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error putting scaling policy for write capacity: %s", err)
	}

	fmt.Println("Configured autoscaling for the table", tableName)
}

func CreateNovel(svc *dynamodb.DynamoDB, novel types.Novel) error {
	tableName := "Novels"

	av, err := dynamodbattribute.MarshalMap(novel)
	if err != nil {
		return errors.New(fmt.Sprintf("Got error marshalling new movie item: %v", err))
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(novel.Title),
			},
			"Author": {
				S: aws.String(novel.Author),
			},
			"Novel": {
				M: av,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return errors.New(fmt.Sprintf("Got error calling PutItem: %v", err))
	}

	return errors.New(fmt.Sprintf("Successfully added '%s' to table '%s'", novel.Title, tableName))
}

func CreateChapter() {
	// TODO
}
