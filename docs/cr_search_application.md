# Search Application (searchapplications.es.eck.github.com)

Representation of a Search application - a named set of indices plus a search
template, queried through a single simplified endpoint instead of the full
Query DSL.

## Lifecycle

No special lifecycle is applied - when the application is deleted from K8s, it is
also deleted from ES. Create and Update are done using the same
`PUT /_application/search_application/{name}` API.
See [Create or update search application API](https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-search-application-put)
in official documentation.

The indices listed in the body must exist before the application is created.
Declare them under `spec.dependencies.indices` and the operator will requeue
until they are present, rather than failing the reconcile.

## Fields

| Key                                    | Type   | Description                                                                                                        |
|----------------------------------------|--------|----------------------------------------------------------------------------------------------------------------------|
| `metadata.name`                        | string | Name of the Search application                                                                                      |
| `spec.targetInstance.name`             | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this SearchApplication will be deployed to |
| `spec.body`                            | string | Search application definition - same you would use when creating a search application using ES REST API              |
| `spec.dependencies.indices`            | list   | List of indices that have to be present in ES cluster before search application is created / updated                 |
| `spec.dependencies.indexTemplates`     | list   | List of index templates that have to be present in ES cluster before search application is created / updated         |
| `spec.dependencies.componentTemplates` | list   | List of component templates that have to be present in ES cluster before search application is created / updated     |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: SearchApplication
metadata:
  name: searchapplication-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  dependencies:
    indices:
      - index-sample
  body: |
    {
      "indices": ["index-sample"],
      "template": {
        "script": {
          "source": {
            "query": {
              "query_string": {
                "query": "{{query_string}}"
              }
            }
          },
          "params": {
            "query_string": "*"
          }
        }
      }
    }
```
