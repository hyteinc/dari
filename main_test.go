package dari_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/hyteinc/dari"
)

var table *dari.Table

var skPkIndexName dari.KeyName = "sk-pk-index"

var skPkIndex dari.Queryable

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("hyte"),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
		return
	}

	testUUID, err := uuid.NewV7()
	if err != nil {
		panic(err)
		return
	}

	dynamoDbClient := dynamodb.NewFromConfig(cfg)
	table = dari.NewTable(dynamoDbClient, "test-table", "pk", "sk", "version")
	table.PrefixPkValues(testUUID.String())

	skPkIndex = table.WithIndex(skPkIndexName, "sk", "pk")

	m.Run()
}

func TestPutGetTenant(t *testing.T) {
	frog := NewTenant("frog")
	frog.Name = "Frog Tenant"

	err := dari.Put[Tenant](t.Context(), table, frog)
	if err != nil {
		t.Errorf("failed to create tenant[%s]: %v", frog.TenantID, err)
		return
	}

	frogKey := NewTenant("frog")
	gotFrog, err := dari.Get[Tenant](t.Context(), table, frogKey)
	if err != nil {
		t.Errorf("failed to get tenant[%s]: %v", frogKey.TenantID, err)
		return
	}

	if gotFrog.TenantID != frogKey.TenantID {
		t.Errorf("failed to get correct tenant, got %s, want %s", gotFrog.TenantID, frogKey.TenantID)
	}

	toad := NewTenant("toad")
	toad.Name = "Toad Tenant"

	err = dari.Put[Tenant](t.Context(), table, toad)
	if err != nil {
		t.Errorf("failed to create tenant[%s]: %v", toad.TenantID, err)
		return
	}
}

func TestPutListProjects(t *testing.T) {
	project1 := NewProject("frog", "project1")
	project2 := NewProject("frog", "project2")
	project3 := NewProject("frog", "project3")

	project4 := NewProject("toad", "project4")
	project5 := NewProject("toad", "project5")

	_ = dari.Put[Project](t.Context(), table, project1)
	_ = dari.Put[Project](t.Context(), table, project2)
	_ = dari.Put[Project](t.Context(), table, project3)
	_ = dari.Put[Project](t.Context(), table, project4)
	_ = dari.Put[Project](t.Context(), table, project5)

	frogProject := NewProject("frog", "")
	frogProjects, err := dari.ListWithPrefix[Project](t.Context(), table, frogProject)
	if err != nil {
		t.Errorf("failed to list projects: %v", err)
		return
	}

	if len(frogProjects) != 3 {
		t.Errorf("got %d projects, want 3", len(frogProjects))
	}
}

func TestIndex(t *testing.T) {
	proj := NewProject("", "project1")
	_, err := dari.ListWithPrefix[Project](t.Context(), skPkIndex, proj)
	if err != nil {
		t.Errorf("failed to list projects: %v", err)
		return
	}
}

func TestUpdateVersion(t *testing.T) {
	frogKey := NewTenant("frog")
	gotFrog, err := dari.Get[Tenant](t.Context(), table, frogKey)
	if err != nil {
		t.Errorf("failed to get tenant[%s]: %v", frogKey.TenantID, err)
		return
	}

	gotFrog.Name = fmt.Sprintf("Frog Tenant [%d]", gotFrog.Version)
	err = dari.Put[Tenant](t.Context(), table, gotFrog)
	if err != nil {
		t.Errorf("failed to update tenant[%s]: %v", gotFrog.TenantID, err)
	}

	gotFrog.Name = "Bad Update"
	err = dari.Put[Tenant](t.Context(), table, gotFrog)
	if err == nil {
		t.Errorf("failed to prevent update due to version mismatch[%s]: %v", gotFrog.TenantID, err)
	}

}
