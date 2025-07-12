package dari

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Get[T any](ctx context.Context, t *Table, k Keys) (*T, error) {
	pkValue, skValue, err := t.keysFor(k)
	if err != nil {
		return nil, err
	}

	key := map[string]types.AttributeValue{
		t.pk: pkValue,
	}

	if t.sk != "" && skValue.Value != "" {
		key[t.sk] = skValue
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(t.name),
		Key:       key,
	}

	out, err := t.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	} else if out.Item != nil {
		var tItem T
		err := attributevalue.UnmarshalMap(out.Item, &tItem)
		if err != nil {
			return nil, err
		}

		return &tItem, nil
	}

	return nil, ErrNotFound
}
