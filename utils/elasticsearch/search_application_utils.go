package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteSearchApplication(esClient *elasticsearch.Client, searchApplicationName string) (ctrl.Result, error) {
	res, err := esClient.SearchApplicationDelete(searchApplicationName)
	return HandleDeleteResponse(err, res)
}

func UpsertSearchApplication(esClient *elasticsearch.Client, searchApplication v1alpha1.SearchApplication) (ctrl.Result, error) {
	res, err := esClient.SearchApplicationPut(searchApplication.Name, strings.NewReader(searchApplication.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
