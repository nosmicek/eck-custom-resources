# Query Ruleset (queryrulesets.es.eck.github.com)

Representation of a Query ruleset - an ordered set of rules that pin, exclude or
otherwise reorder documents for queries matching given criteria, applied through
the `rule` query.

## Lifecycle

No special lifecycle is applied - when the ruleset is deleted from K8s, it is
also deleted from ES. Create and Update are done using the same
`PUT /_query_rules/{ruleset_id}` API.
See [Create or update query ruleset API](https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-query-rules-put-ruleset)
in official documentation.

The whole ruleset is replaced on every update - individual rules are not merged.
Rules are evaluated in the order they appear in `spec.body`.

## Fields

| Key                        | Type   | Description                                                                                                   |
|----------------------------|--------|-----------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Name of the Query ruleset                                                                                      |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this QueryRuleset will be deployed to |
| `spec.body`                | string | Query ruleset definition - same you would use when creating a query ruleset using ES REST API                  |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: QueryRuleset
metadata:
  name: queryruleset-sample
spec:
  targetInstance:
    name: elasticsearch-quickstart
  body: |
    {
      "rules": [
        {
          "rule_id": "pin-eck",
          "type": "pinned",
          "criteria": [
            {
              "type": "contains",
              "metadata": "query_string",
              "values": ["eck", "operator"]
            }
          ],
          "actions": {
            "ids": ["doc-1", "doc-2"]
          }
        }
      ]
    }
```
