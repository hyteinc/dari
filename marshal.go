package dari

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func marshalWithKeys(t *Table, item Keys) (map[string]types.AttributeValue, error) {
	putMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	pkValue, skValue, err := t.keysFor(item)
	if err != nil {
		return nil, err
	}

	putMap[t.pk] = pkValue
	if t.sk != "" && skValue.Value != "" {
		putMap[t.sk] = skValue
	}

	return putMap, nil
}