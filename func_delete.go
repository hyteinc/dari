package dari

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Delete[T any](ctx context.Context, t *Table, k Keys) error {
	pkValue, skValue, err := t.keysFor(k)
	if err != nil {
		return err
	}

	key := map[string]types.AttributeValue{
		t.pk: pkValue,
	}

	if t.sk != "" && skValue.Value != "" {
		key[t.sk] = skValue
	}

	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(t.name),
		Key:          key,
		ReturnValues: types.ReturnValueAllOld,
	}

	out, err := t.client.DeleteItem(ctx, input)
	if err != nil {
		return err
	} else if out.Attributes != nil {
		return nil
	}

	return ErrNotFound
}
