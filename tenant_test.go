package dari

import "fmt"

type Tenant struct {
	TenantId string `json:"tenantId" dynamodbav:"tenantId,omitempty"`
	Name     string `json:"name" dynamodbav:"name,omitempty"`

	Version int `json:"version" dynamodbav:"version,omitempty"`
}

func NewTenant(tenantId string) *Tenant {
	return &Tenant{
		TenantId: tenantId,
	}
}

func (t Tenant) Keys() KeySet {
	return KeySet{
		TableKey: {
			"tenant",
			fmt.Sprintf("tenantId#%s", t.TenantId),
		},
	}
}

func (t Tenant) VersionValue() int {
	return t.Version
}
