package dari

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Queryable interface {
	TableName() string
	Pk() string
	Sk() string
	IndexName() KeyName

	keysFor(k Keys) (*types.AttributeValueMemberS, *types.AttributeValueMemberS, error)

	Client() *dynamodb.Client
}
