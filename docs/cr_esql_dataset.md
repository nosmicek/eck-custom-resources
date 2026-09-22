# ES|QL Dataset (esqldatasets.es.eck.github.com)

Representation of an ES|QL dataset - a named reference to external data, read
through an [ES|QL Data Source](cr_esql_data_source.md). Dataset names participate
in the index namespace alongside indices, aliases and views, so they must follow
index naming rules.

## Lifecycle

No special lifecycle is applied - when the dataset is deleted from K8s, it is
also deleted from ES. Create and Update are done using the same
`PUT /_query/dataset/{name}` API.
See [Create or update an ES|QL dataset](https://www.elastic.co/docs/api/doc/elasticsearch/v9/operation/operation-esql-put-dataset)
in official documentation.

The data source named in `data_source` must already exist. There is no dependency
declaration for data sources, so a dataset created before its data source will
fail its reconcile and be retried until the data source appears.

## Fields

| Key                        | Type   | Description                                                                                                  |
|----------------------------|--------|----------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Name of the ES|QL dataset                                                                                     |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this EsqlDataset will be deployed to |
| `spec.body`                | string | Dataset definition - same you would use when creating a dataset using ES REST API                             |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: EsqlDataset
metadata:
  name: access-logs
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "data_source": "prod-s3-logs",
      "resource": "s3://logs-bucket/access/**/*.parquet",
      "description": "Production access logs",
      "settings": {
        "partition_detection": "hive"
      }
    }
```
