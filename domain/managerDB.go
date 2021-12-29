package domain

import (
	"sync"

	"github.com/ory/ladon"
)

type DBManager struct {
	sync.RWMutex //FIXME:是否需要
	// TODO: 两阶段提交，实现cache与DB的分布式事务
}

func NewPolicyDBManager() *DBManager {
	return &DBManager{}
}

// Create persists the policy.
func (m *DBManager) Create(policy ladon.Policy) error {
	return nil
}

// Update updates an existing policy.
func (m *DBManager) Update(policy ladon.Policy) error {
	return nil
}

// Get retrieves a policy.
func (m *DBManager) Get(id string) (ladon.Policy, error) {
	return nil, nil
}

// Delete removes a policy.
func (m *DBManager) Delete(id string) error {
	return nil
}

// GetAll retrieves all policies.
func (m *DBManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	return nil,nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *DBManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *DBManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *DBManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return nil, nil
}
