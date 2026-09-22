package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteEsqlView(esClient *elasticsearch.Client, esqlViewName string) (ctrl.Result, error) {
	res, err := esClient.EsqlDeleteView(esqlViewName)
	return HandleDeleteResponse(err, res)
}

func UpsertEsqlView(esClient *elasticsearch.Client, esqlView v1alpha1.EsqlView) (ctrl.Result, error) {
	res, err := esClient.EsqlPutView(esqlView.Name, strings.NewReader(esqlView.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
