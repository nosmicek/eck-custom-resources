package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteInferenceEndpoint(esClient *elasticsearch.Client, inferenceEndpointId string, taskType string) (ctrl.Result, error) {
	res, err := esClient.InferenceDelete(inferenceEndpointId,
		esClient.InferenceDelete.WithTaskType(taskType),
	)
	return HandleDeleteResponse(err, res)
}

func UpsertInferenceEndpoint(esClient *elasticsearch.Client, inferenceEndpoint v1alpha1.InferenceEndpoint) (ctrl.Result, error) {
	res, err := esClient.InferencePut(
		strings.NewReader(inferenceEndpoint.Spec.Body),
		inferenceEndpoint.Name,
		esClient.InferencePut.WithTaskType(inferenceEndpoint.Spec.TaskType),
	)

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
