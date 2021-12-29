package domain

import (
	"sync"

	"github.com/ory/ladon"
)

type Manager struct {
	sync.RWMutex //FIXME:是否需要
	// TODO: 两阶段提交，实现cache与DB的分布式事务
	Cache, DB ladon.Manager
}

func NewPolicyManager(cache, db ladon.Manager) *Manager {
	return &Manager{
		Cache: cache,
		DB:    db,
	}
}

// Create persists the policy.
func (m *Manager) Create(policy ladon.Policy) error {
	m.Cache.Create(policy)
	return nil
}

// Update updates an existing policy.
func (m *Manager) Update(policy ladon.Policy) error {
	m.Cache.Update(policy)
	return nil
}

// Get retrieves a policy.
func (m *Manager) Get(id string) (ladon.Policy, error) {
	m.Cache.Get(id)
	return nil, nil
}

// Delete removes a policy.
func (m *Manager) Delete(id string) error {
	m.Cache.Delete(id)

	return nil
}

// GetAll retrieves all policies.
func (m *Manager) GetAll(limit, offset int64) (ladon.Policies, error) {
	m.Cache.GetAll(limit,offset)

	return nil, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *Manager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	m.Cache.FindRequestCandidates(r)

	return nil, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *Manager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	m.Cache.FindPoliciesForSubject(subject)

	return nil, nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *Manager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	m.Cache.FindPoliciesForResource(resource)

	return nil, nil
}
