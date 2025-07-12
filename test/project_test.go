package test_test

import (
	"fmt"
	"github.com/hyteinc/dari"
)

type Project struct {
	TenantID  string `json:"tenantID" dynamodbav:"tenantID,omitempty"`
	ProjectID string `json:"projectID" dynamodbav:"projectID,omitempty"`

	Name string `json:"name" dynamodbav:"name,omitempty"`

	Version int `json:"version" dynamodbav:"version,omitempty"`
}

func NewProject(tenantId, projectId string) *Project {
	return &Project{
		TenantID:  tenantId,
		ProjectID: projectId,
	}
}

func (p Project) Keys() dari.KeySet {
	return dari.KeySet{
		dari.TableKey: {
			fmt.Sprintf("tenant#%s", p.TenantID),
			fmt.Sprintf("projectID#%s", p.ProjectID),
		},
		skPkIndexName: {
			fmt.Sprintf("projectID#%s", p.ProjectID),
			fmt.Sprintf("tenant#%s", p.TenantID),
		},
	}
}

func (p Project) VersionValue() int {
	return p.Version
}
