package domain

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/ladon"
)

type MySQLManager struct {
	Ctx context.Context
	DB  *sql.DB
}

func NewPolicyMySQLManager(ctx context.Context, db *sql.DB) *MySQLManager {

	return &MySQLManager{
		DB:  db,
		Ctx: ctx,
	}
}

// Create persists the policy.
func (mm *MySQLManager) Create(policy ladon.Policy) error {
	return nil
}

// Update updates an existing policy.
func (mm *MySQLManager) Update(policy ladon.Policy) error {
	return nil
}

// Get retrieves a policy.
func (mm *MySQLManager) Get(id string) (ladon.Policy, error) {
	return nil, nil
}

// Delete removes a policy.
func (mm *MySQLManager) Delete(id string) error {
	return nil
}

// GetAll retrieves all policies.
func (mm *MySQLManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	return nil, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (mm *MySQLManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (mm *MySQLManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return nil, nil

}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (mm *MySQLManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return nil, nil
}
