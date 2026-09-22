# ES|QL View (esqlviews.es.eck.github.com)

Representation of an ES|QL view - a named, saved ES|QL query that can be queried
like an index. View names participate in the index namespace alongside indices,
aliases and datasets.

## Lifecycle

No special lifecycle is applied - when the view is deleted from K8s, it is also
deleted from ES. Create and Update are done using the same
`PUT /_query/view/{name}` API.
See [Create or update an ES|QL view](https://www.elastic.co/docs/api/doc/elasticsearch/v9/operation/operation-esql-put-view)
in official documentation.

The indices, datasets or views the query reads from must exist when the view is
created. Creating a view requires the `create_view` index privilege.

## Fields

| Key                        | Type   | Description                                                                                               |
|----------------------------|--------|-------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Name of the ES|QL view                                                                                     |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this EsqlView will be deployed to |
| `spec.body`                | string | View definition - same you would use when creating a view using ES REST API                                |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: EsqlView
metadata:
  name: esqlview-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "query": "FROM index-sample | STATS count() BY category"
    }
```
