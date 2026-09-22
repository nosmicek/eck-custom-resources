# Synonyms Set (synonymssets.es.eck.github.com)

Representation of a Synonyms set - a group of synonym rules stored in an
internal system index and referenced from `synonym_graph` or `synonym` token
filters by set name.

## Lifecycle

No special lifecycle is applied - when the set is deleted from K8s, it is also
deleted from ES. Create and Update are done using the same
`PUT /_synonyms/{id}` API.
See [Create or update synonyms set API](https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-synonyms-put-synonym)
in official documentation.

Updating a synonyms set reloads every analyzer that references it, so indices
using the set are briefly reloaded. Elasticsearch refuses to delete a set that is
still referenced by an index analyzer.

## Fields

| Key                        | Type   | Description                                                                                                  |
|----------------------------|--------|----------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Name of the Synonyms set                                                                                      |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this SynonymsSet will be deployed to |
| `spec.body`                | string | Synonyms set definition - same you would use when creating a synonyms set using ES REST API                   |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: SynonymsSet
metadata:
  name: synonymsset-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "synonyms_set": [
        {
          "id": "rule-1",
          "synonyms": "hello, hi, howdy"
        },
        {
          "id": "rule-2",
          "synonyms": "goodbye, bye, farewell"
        }
      ]
    }
```
