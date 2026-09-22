# ES|QL Data Source (esqldatasources.es.eck.github.com)

Representation of an ES|QL data source - a named, type-specific configuration
describing how to reach external data for ES|QL data federation. Datasets
reference a data source to read from it. `s3` is currently the only supported
type.

## Lifecycle

No special lifecycle is applied - when the data source is deleted from K8s, it is
also deleted from ES. Create and Update are done using the same
`PUT /_query/data_source/{name}` API.
See [Create or update an ES|QL data source](https://www.elastic.co/docs/api/doc/elasticsearch/v9/operation/operation-esql-put-data-source)
in official documentation.

A data source must exist before any [ES|QL Dataset](cr_esql_dataset.md) that
references it, and Elasticsearch refuses to delete one that datasets still
reference.

Credentials belong in `settings` and are sent to Elasticsearch verbatim. Prefer a
keyless authentication method where the deployment supports one, so that secrets
do not have to be written into the resource body.

## Fields

| Key                        | Type   | Description                                                                                                     |
|----------------------------|--------|-------------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Name of the ES|QL data source                                                                                    |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this EsqlDataSource will be deployed to |
| `spec.body`                | string | Data source definition - same you would use when creating a data source using ES REST API                        |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: EsqlDataSource
metadata:
  name: prod-s3-logs
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "type": "s3",
      "description": "Production S3 logs bucket",
      "settings": {
        "region": "us-east-1"
      }
    }
```
