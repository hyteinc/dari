package dari

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"testing"
)

//go:generate go test -coverprofile=coverage .
//go:generate go tool cover -html=coverage

var table *Table

var skPkIndexName KeyName = "sk-pk-index"

var skPkIndex Queryable

func TestMain(m *testing.M) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("hyte"),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return
	}

	testUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("unable to generate uuid, %v", err)
		return
	}

	dynamoDbClient := dynamodb.NewFromConfig(cfg)
	table = NewTable(dynamoDbClient, "test-table", "pk", "sk", "version")
	table.PrefixPkValues(testUuid.String())

	skPkIndex = table.WithIndex(skPkIndexName, "sk", "pk")

	m.Run()
}

func TestPutGetTenant(t *testing.T) {
	frog := NewTenant("frog")
	frog.Name = "Frog Tenant"

	err := Put[Tenant](t.Context(), table, frog)
	if err != nil {
		t.Errorf("failed to create tenant[%s]: %v", frog.TenantId, err)
		return
	}

	frogKey := NewTenant("frog")
	gotFrog, err := Get[Tenant](t.Context(), table, frogKey)
	if err != nil {
		t.Errorf("failed to get tenant[%s]: %v", frogKey.TenantId, err)
		return
	}

	if gotFrog.TenantId != frogKey.TenantId {
		t.Errorf("failed to get correct tenant, got %s, want %s", gotFrog.TenantId, frogKey.TenantId)
	}

	toad := NewTenant("toad")
	toad.Name = "Toad Tenant"

	err = Put[Tenant](t.Context(), table, toad)
	if err != nil {
		t.Errorf("failed to create tenant[%s]: %v", toad.TenantId, err)
		return
	}
}

func TestPutListProjects(t *testing.T) {
	project1 := NewProject("frog", "project1")
	project2 := NewProject("frog", "project2")
	project3 := NewProject("frog", "project3")

	project4 := NewProject("toad", "project4")
	project5 := NewProject("toad", "project5")

	_ = Put[Project](t.Context(), table, project1)
	_ = Put[Project](t.Context(), table, project2)
	_ = Put[Project](t.Context(), table, project3)
	_ = Put[Project](t.Context(), table, project4)
	_ = Put[Project](t.Context(), table, project5)

	frogProject := NewProject("frog", "")
	frogProjects, err := ListWithPrefix[Project](t.Context(), table, frogProject)
	if err != nil {
		t.Errorf("failed to list projects: %v", err)
		return
	}

	if len(frogProjects) != 3 {
		t.Errorf("got %d projects, want 3", len(frogProjects))
	}

	fmt.Println(frogProjects)
}

func TestIndex(t *testing.T) {

	proj := NewProject("", "project1")
	items, err := ListWithPrefix[Project](t.Context(), skPkIndex, proj)
	if err != nil {
		t.Errorf("failed to list projects: %v", err)
		return
	}

	fmt.Println(items)
}

func TestUpdateVersion(t *testing.T) {
	frogKey := NewTenant("frog")
	gotFrog, err := Get[Tenant](t.Context(), table, frogKey)
	if err != nil {
		t.Errorf("failed to get tenant[%s]: %v", frogKey.TenantId, err)
		return
	}

	gotFrog.Name = fmt.Sprintf("Frog Tenant [%d]", gotFrog.Version)
	err = Put[Tenant](t.Context(), table, gotFrog)
	if err != nil {
		t.Errorf("failed to update tenant[%s]: %v", gotFrog.TenantId, err)
	}

	gotFrog.Name = "Bad Update"
	err = Put[Tenant](t.Context(), table, gotFrog)
	if err == nil {
		t.Errorf("failed to prevent update due to version mismatch[%s]: %v", gotFrog.TenantId, err)
	}

}
