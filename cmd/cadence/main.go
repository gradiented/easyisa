package cadence

import (
	"time"

	common "github.com/gradiented/easyisa/cmd/cadence/common"
	workflows "github.com/gradiented/easyisa/cmd/cadence/workflows"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
)

var h common.CadenceHelper

func Start() {
	h.SetupServiceConfig()
	startWorkers(&h)
}

func startWorkers(h *common.CadenceHelper) {
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, "easyisa-tasks", workerOptions)
}

func StartGreetWorkflow() {
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "greetings_",
		TaskList:                        "easyisa-tasks",
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}
	h.StartWorkflow(workflowOptions, workflows.HelloCadenceWorkflow)
}
