package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteEsqlDataSource(esClient *elasticsearch.Client, esqlDataSourceName string) (ctrl.Result, error) {
	res, err := esClient.EsqlDeleteDataSource([]string{esqlDataSourceName})
	return HandleDeleteResponse(err, res)
}

func UpsertEsqlDataSource(esClient *elasticsearch.Client, esqlDataSource v1alpha1.EsqlDataSource) (ctrl.Result, error) {
	res, err := esClient.EsqlPutDataSource(esqlDataSource.Name, strings.NewReader(esqlDataSource.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
