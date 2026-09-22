package kibana

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	kibanaeckv1alpha1 "github.com/xco-sk/eck-custom-resources/apis/kibana.eck/v1alpha1"
	"github.com/xco-sk/eck-custom-resources/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HandleDeleteResponse maps a delete call onto a reconcile result. A 404 counts as
// success: the object is already gone, and failing here would keep the finalizer on
// the Kubernetes object forever. Any other non-success status is reported, so that
// the finalizer holds until the object is really deleted from Kibana.
func HandleDeleteResponse(response *http.Response, err error) (ctrl.Result, error) {
	if err != nil {
		return utils.GetRequeueResult(), err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 && response.StatusCode != http.StatusNotFound {
		resBody, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return utils.GetRequeueResult(), readErr
		}
		return utils.GetRequeueResult(), fmt.Errorf("Non-success (%d) response: %s, ", response.StatusCode, string(resBody))
	}

	return ctrl.Result{}, nil
}

func InjectId(objectJson string, id string) (*string, error) {
	var body map[string]interface{}
	err := json.NewDecoder(strings.NewReader(objectJson)).Decode(&body)
	if err != nil {
		return nil, err
	}

	body["id"] = id

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sBody := string(marshalledBody)
	return &sBody, nil
}

func GetTargetInstance(cli client.Client, ctx context.Context, namespace string, targetName string, kibanaInstance *kibanaeckv1alpha1.KibanaInstance) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: targetName}, kibanaInstance); err != nil {
		return err
	}
	return nil
}
