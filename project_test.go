package dari

import "fmt"

type Project struct {
	TenantId  string `json:"tenantId" dynamodbav:"tenantId,omitempty"`
	ProjectId string `json:"projectId" dynamodbav:"projectId,omitempty"`

	Name string `json:"name" dynamodbav:"name,omitempty"`

	Version int `json:"version" dynamodbav:"version,omitempty"`
}

func NewProject(tenantId, projectId string) *Project {
	return &Project{
		TenantId:  tenantId,
		ProjectId: projectId,
	}
}

func (p Project) Keys() KeySet {
	return KeySet{
		TableKey: {
			fmt.Sprintf("tenant#%s", p.TenantId),
			fmt.Sprintf("projectId#%s", p.ProjectId),
		},
		skPkIndexName: {
			fmt.Sprintf("projectId#%s", p.ProjectId),
			fmt.Sprintf("tenant#%s", p.TenantId),
		},
	}
}

func (p Project) VersionValue() int {
	return p.Version
}
