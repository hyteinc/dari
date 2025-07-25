package dari_test

import (
	"fmt"

	"github.com/hyteinc/dari"
)

type Tenant struct {
	TenantID string `json:"tenantID" dynamodbav:"tenantID,omitempty"`
	Name     string `json:"name" dynamodbav:"name,omitempty"`

	Version int `json:"version" dynamodbav:"version,omitempty"`
}

func NewTenant(tenantID string) *Tenant {
	return &Tenant{
		TenantID: tenantID,
	}
}

func (t Tenant) Keys() dari.KeySet {
	return dari.KeySet{
		dari.TableKey: {
			"tenant",
			fmt.Sprintf("tenantId#%s", t.TenantID),
		},
	}
}

func (t Tenant) VersionValue() int {
	return t.Version
}
