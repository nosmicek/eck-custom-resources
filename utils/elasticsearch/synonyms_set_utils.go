package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteSynonymsSet(esClient *elasticsearch.Client, synonymsSetId string) (ctrl.Result, error) {
	res, err := esClient.SynonymsDeleteSynonym(synonymsSetId)
	return HandleDeleteResponse(err, res)
}

func UpsertSynonymsSet(esClient *elasticsearch.Client, synonymsSet v1alpha1.SynonymsSet) (ctrl.Result, error) {
	res, err := esClient.SynonymsPutSynonym(synonymsSet.Name, strings.NewReader(synonymsSet.Spec.Body))

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
