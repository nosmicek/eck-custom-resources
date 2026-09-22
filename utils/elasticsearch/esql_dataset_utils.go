package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteEsqlDataset(esClient *elasticsearch.Client, esqlDatasetName string) (ctrl.Result, error) {
	res, err := esClient.EsqlDeleteDataset([]string{esqlDatasetName})
	return HandleDeleteResponse(err, res)
}

func UpsertEsqlDataset(esClient *elasticsearch.Client, esqlDataset v1alpha1.EsqlDataset) (ctrl.Result, error) {
	res, err := esClient.EsqlPutDataset(esqlDataset.Name, strings.NewReader(esqlDataset.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
