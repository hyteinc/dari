package dari

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func ListWithPrefix[T any](ctx context.Context, t Queryable, k Keys) ([]T, error) {
	if t.Sk() == "" {
		return []T{}, ErrNotSupportNoSk
	}

	pkValue, skValue, err := t.keysFor(k)
	if err != nil {
		return []T{}, err
	}

	queryInput := dynamodb.QueryInput{
		TableName:              aws.String(t.TableName()),
		KeyConditionExpression: aws.String("#pk = :pk AND begins_with(#sk, :sk)"),
		ExpressionAttributeNames: map[string]string{
			"#pk": t.Pk(),
			"#sk": t.Sk(),
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": pkValue,
			":sk": skValue,
		},
		ScanIndexForward: aws.Bool(false),
	}

	if string(t.IndexName()) != "" {
		queryInput.IndexName = aws.String(string(t.IndexName()))
	}

	// queryInputJson, err := json.Marshal(queryInput)
	// if err == nil {
	// 	fmt.Println(string(queryInputJson))
	// }

	results, err := t.Client().Query(ctx, &queryInput)
	if err != nil {
		return nil, fmt.Errorf("list with prefix error: %w", err)
	}

	var items []T
	for _, item := range results.Items {
		var tItem T

		err = attributevalue.UnmarshalMap(item, &tItem)
		if err != nil {
			continue
		}

		items = append(items, tItem)
	}

	return items, nil
}
