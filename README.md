# Custom resources for ECK
[![docker-publish](https://github.com/nosmicek/eck-custom-resources/actions/workflows/docker-publish.yaml/badge.svg)](https://github.com/nosmicek/eck-custom-resources/actions/workflows/docker-publish.yaml)
[![helm-publish](https://github.com/nosmicek/eck-custom-resources/actions/workflows/helm-publish.yml/badge.svg)](https://github.com/nosmicek/eck-custom-resources/actions/workflows/helm-publish.yml)

Kubernetes operator that enables the installation of various resources for
Elasticsearch and Kibana.

Fork of [xco-sk/eck-custom-resources](https://github.com/xco-sk/eck-custom-resources)
with Elasticsearch 9.x support. Charts and images are published from this repository.

Currently supported resources: 
- For Elasticsearch:
  - [Elasticsearch Instance](docs/cr_elasticsearch_instance.md)
  - [Index](docs/cr_index.md)
  - [Index template](docs/cr_index_template.md)
  - [Index lifecycle policy](docs/cr_index_lifecycle_policy.md)
  - [Ingest pipeline](docs/cr_ingest_pipeline.md)
  - [Snapshot repository](docs/cr_snapshot_repo.md)
  - [Snapshot lifecycle policy](docs/cr_snapshot_lifecycle_policy.md)
  - [User](docs/cr_user.md)
  - [Role](docs/cr_role.md)
  - [API key](docs/cr_apikey.md)
  - [Component template](docs/cr_component_template.md)
  - [Inference endpoint](docs/cr_inference_endpoint.md)
  - [Synonyms set](docs/cr_synonyms_set.md)
  - [Query ruleset](docs/cr_query_ruleset.md)
  - [Search application](docs/cr_search_application.md)
  - [ES|QL view](docs/cr_esql_view.md)
  - [ES|QL dataset](docs/cr_esql_dataset.md)
  - [ES|QL data source](docs/cr_esql_data_source.md)
- For Kibana:
  - [Kibana Instance](docs/cr_kibana_instance.md)
  - [Space](docs/cr_space.md)
  - [Index pattern](docs/cr_index_pattern.md)
  - [Saved search](docs/cr_saved_search.md)
  - [Visualization](docs/cr_visualization.md)
  - [Lens](docs/cr_lens.md)
  - [Dashboard](docs/cr_dashboard.md)
  - [Data View](docs/cr_data_view.md)

## Installation

```shell
# Add eck-custom-resources helm repo
helm repo add eck-custom-resources https://nosmicek.github.io/eck-custom-resources/

# Install chart
helm install eck-cr eck-custom-resources/eck-custom-resources-operator
```
Configuration options are documented in [chart README file](charts/eck-custom-resources-operator/README.md)

## Resource deletion

Every managed resource carries a finalizer, so deleting the Kubernetes object also
deletes its counterpart in Elasticsearch or Kibana. If that call fails the object
stays in `Terminating` and is retried, rather than disappearing and leaving the
remote object behind.

Because a blocked finalizer also blocks deletion of the namespace holding it, the
operator releases the finalizer without attempting remote deletion when the call
could never succeed:

- the target `ElasticsearchInstance` / `KibanaInstance` no longer exists
- the referenced authentication or certificate secret no longer exists
- the target instance has `enabled: false`

A resource that is already absent from Elasticsearch or Kibana (HTTP 404) counts as
deleted, so re-deleting never wedges.

To keep a resource in Elasticsearch or Kibana while removing the Kubernetes object,
annotate it:

```yaml
metadata:
  annotations:
    eck.github.com/skip-remote-delete: "true"
```

Anything else - a cluster that is unreachable, credentials without delete
permission, a resource still referenced by another - keeps the object in
`Terminating` until it is resolved. To force it through, remove the finalizer:

```shell
kubectl patch <kind> <name> --type=merge -p '{"metadata":{"finalizers":[]}}'
```

## Upgrade guide

Helm installs CRDs only on first install and never upgrades them, so apply the
CRD bundle attached to each release before upgrading the chart. From 0.9.0 on,
the full set is published as a single `crds.yaml`.

### To 0.9.0
0.9.0 is the first chart published from this repository. A release installed
from the original `xco-sk` chart repository has to be pointed at this one first:
```shell
helm repo add eck-custom-resources https://nosmicek.github.io/eck-custom-resources/ --force-update
```
Seven new CRDs for Elasticsearch 9.x resources were introduced. Apply the full
CRD set, then upgrade the chart:
```shell
kubectl apply --server-side --force-conflicts \
  -f https://github.com/nosmicek/eck-custom-resources/releases/download/v0.9.0/crds.yaml
helm upgrade eck-cr eck-custom-resources/eck-custom-resources-operator
```
`--force-conflicts` is needed because the existing CRDs are owned by Helm's
field manager. The operator image moved to `nosmo/eck-custom-resources`; if
`image.repository` is overridden in your values, update it as well.

The ES|QL view, dataset and data source resources require Elasticsearch 9.2 or
newer; the remaining resources work on any Elasticsearch 9.x.

### From 0.7.0 to 0.7.1
Existing `ComponentTemplate` CRD was fixed. To apply the CRD, run:
```
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/v0.7.1/config/crd/bases/es.eck.github.com_componenttemplates.yaml
```

### From 0.6.0 to 0.7.0
There is a new `ComponentTemplate` CRD present. To apply the CRD, run:
```
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/v0.7.0/config/crd/bases/es.eck.github.com_componenttemplates.yaml
```

### From 0.5.0 to 0.6.0
The Elasticsearch API Key support was introduced. To apply the CRD, run:
```
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/v0.6.0/config/crd/bases/es.eck.github.com_elasticsearchapikeys.yaml
```

### From 0.4.1 to 0.5.0
The Multi-target support was introduced. This changes is backward compatible, but in order to make use of the multi-target support
apply the new CRDs manually:
```

kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_elasticsearchinstances.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_elasticsearchroles.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_elasticsearchusers.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_indexlifecyclepolicies.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_indextemplates.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_indices.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_ingestpipelines.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_snapshotlifecyclepolicies.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/es.eck.github.com_snapshotrepositories.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_kibanainstances.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_dashboards.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_indexpatterns.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_lens.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_savedsearches.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_spaces.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_visualizations.yaml
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.5.0/config/crd/bases/kibana.eck.github.com_dataviews.yaml
```

There are 2 new CRDs, `ElasticsearchInstance` and `KibanaInstance` that allows you to deploy the target configuration for
both Kibana and Elasticsearch. The rest of the CRDs were extended with optional `spec.targetInstance.name` field, that should reference
the `ElasticsearchInstance`/`KibanaInstance`. If `targetInstance` field is not present, the default operator configuration (`elasticsearch` and `kibana`
fields) is used.
This approach should ensure the backward compatibility with previously deployed CRDs.
See [samples](config/samples).

### From 0.3.2 to 0.4.1
There is new `DataView` CRD present. To apply the CRD, run:
```
kubectl apply --server-side -f https://raw.githubusercontent.com/nosmicek/eck-custom-resources/eck-custom-resources-operator-0.4.1/config/crd/bases/kibana.eck.github.com_dataviews.yaml
```


## Uninstallation
To uninstall the eck-cr from Kubernetes cluster, run:

```shell
helm uninstall eck-cr
```

This removes all resources related to eck-custom-resources operator. It won't remove the CRDs nor any deployed custom resource
(e.g. Index, Index Template ...), they will remain in K8s and also in Elasticsearch.

## Working with custom resources
After the operator is installed, you can deploy Elasticsearch/Kibana resources from the list above. The reconciler
will take care of propagating the change to Elasticsearch or Kibana, whether it is creation of new resource, deletion
or update. Definition of target Elasticsearch/Kibana is done using [Elasticsearch Instance](docs/cr_elasticsearch_instance.md) and 
[Kibana Instance](docs/cr_kibana_instance.md) resources. These are then referenced (by name) from other resources through `spec.targetInstance.name` field.

For detailed documentation for each resource, see [List of supported resources](docs/cr_list.md)

## Help and Troubleshooting
In case you need help or found a bug, please create an [Issue on Github](https://github.com/nosmicek/eck-custom-resources/issues).

## License
Licensed under the Apache License, Version 2.0; see [LICENSE.md](LICENSE.md)
