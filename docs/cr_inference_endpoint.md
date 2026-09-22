# Inference Endpoint (inferenceendpoints.es.eck.github.com)

Representation of an Inference endpoint, used to integrate Elasticsearch with
machine learning models served either by Elasticsearch itself (ELSER, E5) or by
an external provider such as OpenAI, Cohere, Anthropic, Azure AI Studio,
Amazon Bedrock, Google Vertex AI, HuggingFace, Mistral or Watsonx.

## Lifecycle

No special lifecycle is applied - when the endpoint is deleted from K8s, it is
also deleted from ES. Create and Update are done using the same
`PUT /_inference/{task_type}/{inference_id}` API.
See [Create inference API](https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-inference-put)
in official documentation.

`spec.taskType` forms part of the endpoint's path in Elasticsearch, so it cannot
be changed on an existing object - the CRD rejects the update. To move an
endpoint to a different task type, delete it and create it again.

Deleting an endpoint that is still referenced by a `semantic_text` field or by an
ingest pipeline will be refused by Elasticsearch; the object stays until the
reference is removed.

## Fields

| Key                        | Type   | Description                                                                                                        |
|----------------------------|--------|--------------------------------------------------------------------------------------------------------------------|
| `metadata.name`            | string | Inference ID of the endpoint                                                                                        |
| `spec.taskType`            | string | Inference task the endpoint serves. One of `sparse_embedding`, `text_embedding`, `rerank`, `completion`, `chat_completion`, `embedding`. Immutable |
| `spec.targetInstance.name` | string | Name of the [Elasticsearch Instance](cr_elasticsearch_instance.md) to which this InferenceEndpoint will be deployed to |
| `spec.body`                | string | Inference endpoint definition - same you would use when creating an inference endpoint using ES REST API             |

## Example

```yaml
apiVersion: es.eck.github.com/v1alpha1
kind: InferenceEndpoint
metadata:
  name: my-elser-model
spec:
  targetInstance:
    name: elasticsearch-quickstart
  taskType: sparse_embedding
  body: |
    {
      "service": "elser",
      "service_settings": {
        "num_allocations": 1,
        "num_threads": 1
      }
    }
```
