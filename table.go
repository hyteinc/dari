package dari

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableKey KeyName

type Table struct {
	name   string
	pk     string
	pkOrig string
	sk     string

	indexName KeyName

	version string

	client         *dynamodb.Client
	prefixPkValues string
}

func NewTable(client *dynamodb.Client, name, pk, sk, version string) *Table {
	return &Table{
		name:      name,
		pk:        pk,
		pkOrig:    pk,
		sk:        sk,
		indexName: TableKey,
		version:   version,
		client:    client,
	}
}

func (t *Table) PrefixPkValues(prefix string) {
	t.prefixPkValues = prefix
}

func (t *Table) keysFor(k Keys) (*types.AttributeValueMemberS, *types.AttributeValueMemberS, error) {
	keySet := k.Keys()
	key, ok := keySet[t.indexName]
	if !ok {
		return nil, nil, fmt.Errorf("not supported, key set[%s] does not exist on %t", t.indexName, k)
	}

	pkValue := key[0]
	skValue := key[1]

	if t.pkOrig == t.sk && t.prefixPkValues != "" {
		skValue = fmt.Sprintf("%s-%s", t.prefixPkValues, skValue)
	} else if t.prefixPkValues != "" {
		pkValue = fmt.Sprintf("%s-%s", t.prefixPkValues, pkValue)
	}

	return &types.AttributeValueMemberS{Value: pkValue}, &types.AttributeValueMemberS{Value: skValue}, nil
}

func (t *Table) Client() *dynamodb.Client {
	return t.client
}

func (t *Table) TableName() string {
	return t.name
}

func (t *Table) Pk() string {
	return t.pk
}

func (t *Table) Sk() string {
	return t.sk
}

func (t *Table) IndexName() KeyName {
	return t.indexName
}

func (t *Table) WithIndex(indexName KeyName, pk, sk string) Queryable {
	return &Table{
		name:           t.name,
		pk:             pk,
		pkOrig:         t.pkOrig,
		prefixPkValues: t.prefixPkValues,
		sk:             sk,
		version:        "",
		indexName:      indexName,
		client:         t.client,
	}
}

func (t *Table) WithSimilarTable(tableName string) *Table {
	return &Table{
		name:           tableName,
		pk:             t.pk,
		pkOrig:         t.pkOrig,
		prefixPkValues: t.prefixPkValues,
		sk:             t.sk,
		version:        t.version,
		indexName:      t.indexName,
		client:         t.client,
	}
}

