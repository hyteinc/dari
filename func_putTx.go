package dari

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

func PutTx(ctx context.Context, t *Table, items ...Keys) error {
	switch n := len(items); {
	case n == 0:
		return nil
	case n > 25:
		return fmt.Errorf("too many items (%d): DynamoDB transactions support up to 25 actions", n)
	}

	txItems := make([]types.TransactWriteItem, 0, len(items))
	for i, item := range items {
		txItem, err := buildPutTxItem(t, item, i)
		if err != nil {
			return err
		}
		txItems = append(txItems, txItem)
	}

	_, err := t.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: txItems,
	})
	if err == nil {
		return nil
	}
	if isConditionalTxCancel(err) {
		return ErrAlreadyExists
	}
	return err
}

func buildPutTxItem(t *Table, item Keys, index int) (types.TransactWriteItem, error) {
	if item == nil {
		return types.TransactWriteItem{}, fmt.Errorf("nil item at index %d", index)
	}

	putMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return types.TransactWriteItem{}, err
	}

	pkValue, skValue, err := t.keysFor(item)
	if err != nil {
		return types.TransactWriteItem{}, err
	}

	putMap[t.pk] = pkValue
	if t.sk != "" && skValue.Value != "" {
		putMap[t.sk] = skValue
	}

	put := &types.Put{
		TableName: aws.String(t.name),
		Item:      putMap,
	}

	putApplyVersioning(t, item, put, putMap)

	return types.TransactWriteItem{Put: put}, nil
}

func putApplyVersioning(t *Table, item Keys, put *types.Put, putMap map[string]types.AttributeValue) {
	if t.version == "" {
		return
	}

	v, ok := any(item).(Version)
	if !ok {
		return
	}

	versionValue := v.VersionValue()

	if versionValue == 0 && t.sk != "" {
		put.ConditionExpression = aws.String(fmt.Sprintf("attribute_not_exists(%s)", t.sk))
	} else {
		put.ConditionExpression = aws.String("#v = :v")
		put.ExpressionAttributeNames = map[string]string{"#v": t.version}
		put.ExpressionAttributeValues = map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberN{Value: strconv.Itoa(versionValue)},
		}
	}

	putMap[t.version] = &types.AttributeValueMemberN{Value: strconv.Itoa(versionValue + 1)}
}

func isConditionalTxCancel(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "TransactionCanceledException" {
		return false
	}

	var tce *types.TransactionCanceledException
	if !errors.As(err, &tce) {
		return false
	}

	for _, r := range tce.CancellationReasons {
		if r.Code != nil && *r.Code == "ConditionalCheckFailed" {
			return true
		}
	}
	return false
}
