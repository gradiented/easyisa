package cadence

import (
	"context"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func init() {
	activity.Register(GreetingActivity)
	workflow.Register(HelloCadenceWorkflow)
}

func GreetingActivity(ctx context.Context) (string, error) {
	activity.GetLogger(ctx).Info("GreetingActivity called.")
	return "Oh hi. :)", nil
}

func HelloCadenceWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		TaskList:               "easyisa-tasks",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteActivity(ctx, GreetingActivity)
	var result string
	if err := future.Get(ctx, &result); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Done", zap.String("result", result))
	return nil
}
