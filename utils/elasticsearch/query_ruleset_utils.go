package elasticsearch

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xco-sk/eck-custom-resources/apis/es.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

func DeleteQueryRuleset(esClient *elasticsearch.Client, queryRulesetId string) (ctrl.Result, error) {
	res, err := esClient.QueryRulesDeleteRuleset(queryRulesetId)
	return HandleDeleteResponse(err, res)
}

func UpsertQueryRuleset(esClient *elasticsearch.Client, queryRuleset v1alpha1.QueryRuleset) (ctrl.Result, error) {
	res, err := esClient.QueryRulesPutRuleset(strings.NewReader(queryRuleset.Spec.Body), queryRuleset.Name)

	if err != nil || res.IsError() {
		return utils.GetRequeueResult(), GetClientErrorOrResponseError(err, res)
	}

	return ctrl.Result{}, nil
}
