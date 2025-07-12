package dari

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

func Put[T Keys](ctx context.Context, t *Table, item *T) error {
	putMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	k, _ := any(item).(Keys)
	pkValue, skValue, err := t.keysFor(k)
	if err != nil {
		return err
	}

	putMap[t.pk] = pkValue
	if t.sk != "" && skValue.Value != "" {
		putMap[t.sk] = skValue
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(t.name),
		Item:      putMap,
	}

	if t.version != "" {
		version, versionOk := any(item).(Version)
		if versionOk {
			versionValue := version.VersionValue()
			if versionValue == 0 && t.sk != "" {
				input.ConditionExpression = aws.String(fmt.Sprintf("attribute_not_exists(%s)", t.sk))
			} else {
				input.ConditionExpression = aws.String("#v = :v")
				input.ExpressionAttributeNames = map[string]string{
					"#v": t.version,
				}
				input.ExpressionAttributeValues = map[string]types.AttributeValue{
					":v": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", versionValue)},
				}
			}

			putMap[t.version] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", versionValue+1)}
		}
	}

	_, err = t.client.PutItem(ctx, input)
	if err != nil {
		var oe smithy.APIError
		if errors.As(err, &oe) {
			if oe.ErrorCode() == "ConditionalCheckFailedException" {
				return ErrAlreadyExists
			}
		}

		return err
	}
	return nil
}
