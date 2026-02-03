package dari

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go"
)

func Put[T Keys](ctx context.Context, t *Table, item *T) error {
	if item == nil {
		return fmt.Errorf("nil item")
	}

	// Reuse the tx builder to get identical Item + condition/version logic.
	txItem, err := buildPutTxItem(t, any(item).(Keys), 0)
	if err != nil {
		return err
	}

	put := txItem.Put
	input := &dynamodb.PutItemInput{
		TableName:                 put.TableName,
		Item:                      put.Item,
		ConditionExpression:       put.ConditionExpression,
		ExpressionAttributeNames:  put.ExpressionAttributeNames,
		ExpressionAttributeValues: put.ExpressionAttributeValues,
	}

	_, err = t.client.PutItem(ctx, input)
	if err == nil {
		return nil
	}

	var oe smithy.APIError
	if errors.As(err, &oe) && oe.ErrorCode() == "ConditionalCheckFailedException" {
		return ErrAlreadyExists
	}
	return err
}
