package services

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

// ECSFaceRecognitionDispatcher 프로덕션용: job 삽입 후 ECS task 트리거
type ECSFaceRecognitionDispatcher struct {
	jobStore       store.FaceRecognitionJobStore
	ecsClient      *ecs.Client
	clusterArn     string
	taskDefArn     string
	subnets        []string
	securityGroups []string
}

func NewECSFaceRecognitionDispatcher(
	jobStore store.FaceRecognitionJobStore,
	ecsClient *ecs.Client,
	clusterArn, taskDefArn string,
	subnets, securityGroups []string,
) *ECSFaceRecognitionDispatcher {
	return &ECSFaceRecognitionDispatcher{
		jobStore:       jobStore,
		ecsClient:      ecsClient,
		clusterArn:     clusterArn,
		taskDefArn:     taskDefArn,
		subnets:        subnets,
		securityGroups: securityGroups,
	}
}

func (d *ECSFaceRecognitionDispatcher) Dispatch(ctx context.Context, params FaceRecognitionJobParams) error {
	// 1. job INSERT
	if _, err := d.jobStore.CreateFaceRecognitionJob(ctx, db.CreateFaceRecognitionJobParams{
		MediaItemID:    params.MediaItemID,
		FamilyID:       params.FamilyID,
		ViewStorageKey: params.ViewStorageKey,
	}); err != nil {
		return err
	}

	// 2. 이미 실행 중인 task가 있으면 RunTask 스킵
	listOut, err := d.ecsClient.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       aws.String(d.clusterArn),
		Family:        aws.String(d.taskDefArn),
		DesiredStatus: types.DesiredStatusRunning,
	})
	if err != nil {
		log.Printf("ECS ListTasks 실패 (RunTask 진행): %v", err)
	} else if len(listOut.TaskArns) > 0 {
		return nil
	}

	// 3. ECS RunTask (실패해도 job은 DB에 있으므로 log만)
	_, err = d.ecsClient.RunTask(ctx, &ecs.RunTaskInput{
		Cluster:        aws.String(d.clusterArn),
		TaskDefinition: aws.String(d.taskDefArn),
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        d.subnets,
				SecurityGroups: d.securityGroups,
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
	})
	if err != nil {
		log.Printf("ECS RunTask 실패 (job은 생성됨, 다음 트리거에서 재시도): %v", err)
	}

	return nil
}
