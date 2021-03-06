package domain

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
)

var (
	ErrPolicyExists  = errors.New("Policy exists")
	ErrBadConversion = errors.New("Could not convert policy from redis")
	ErrNotFound      = errors.New("Not found")
)

const (
	prefixPolicy   = "policy"
	prefixResource = "resource"
	prefixSubject  = "subject"
)

// Just returns strings.Join(vals, "_") for creating redis keys
func prefixKey(vals ...string) string {
	return strings.Join(vals, "_")
}

// RedisManager is a redis implementation of Manager to store policies persistently.
type RedisManager struct {
	db        *redis.Client
	ctx       context.Context
	keyPrefix string
}

// NewRedisManager initializes a new RedisManager with no policies
func NewRedisManager(ctx context.Context, db *redis.Client, keyPrefix string) *RedisManager {
	if keyPrefix == "" {
		keyPrefix = "ladon"
	}

	return &RedisManager{
		db:        db,
		keyPrefix: keyPrefix,
		ctx:       ctx,
	}
}

// Create a new policy in Redis. It will create a single key for the policy itself,
// and for each subject and resource the policy will also exist in a hashmap.
func (m *RedisManager) Create(policy ladon.Policy) error {
	// Make sure that the key doesn't already exist
	key := prefixKey(m.keyPrefix, prefixPolicy, policy.GetID())
	if err := m.db.Get(m.ctx, key).Err(); err == nil {
		return ErrPolicyExists
	}

	p, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	// Set the policy key
	cmd := m.db.Set(m.ctx, key, p, 0)

	if err := cmd.Err(); err != nil {
		return err
	}

	// Put this policy in the hashmap for each resource
	for _, v := range policy.GetResources() {
		hmkey := prefixKey(m.keyPrefix, prefixResource, v)
		field := policy.GetID()
		if err := m.db.HMSet(m.ctx, hmkey, map[string]interface{}{
			field: p,
		}).Err(); err != nil {
			return err
		}
	}

	// Put this policy in the hashmap for each subject
	for _, v := range policy.GetSubjects() {
		hmkey := prefixKey(m.keyPrefix, prefixSubject, v)
		field := policy.GetID()
		if err := m.db.HMSet(m.ctx, hmkey, map[string]interface{}{
			field: p,
		}).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetAll retrieves all policies. (Equivelant of db.keys + db.Mget)
func (m *RedisManager) GetAll(limit int64, offset int64) (ladon.Policies, error) {
	key := prefixKey(m.keyPrefix, prefixPolicy, "*")
	keyscmd := m.db.Keys(m.ctx, key)
	if err := keyscmd.Err(); err != nil {
		return nil, err
	}

	keys, err := keyscmd.Result()
	if err != nil {
		return nil, err
	}

	mgetcmd := m.db.MGet(m.ctx, keys...)
	if err := mgetcmd.Err(); err != nil {
		return nil, err
	}

	values := mgetcmd.Val()

	policies := make(ladon.Policies, len(values))
	for i, v := range values {
		p := &ladon.DefaultPolicy{}
		b := []byte(v.(string))
		// if !ok {
		// 	return nil, errors.Wrapf(ErrBadConversion, "value %+v is not a byte array", v)
		// }
		if err := json.Unmarshal(b, p); err != nil {
			return nil, errors.Wrap(ErrBadConversion, err.Error())
		}
		policies[i] = p
	}

	if offset+limit > int64(len(policies)) {
		limit = int64(len(policies))
		offset = 0
	}

	return policies[offset:limit], nil
}

// Get retrieves a policy.
func (m *RedisManager) Get(id string) (ladon.Policy, error) {
	var (
		key    = prefixKey(m.keyPrefix, prefixPolicy, id)
		cmd    = m.db.Get(m.ctx, key)
		policy = &ladon.DefaultPolicy{}
	)

	if err := cmd.Err(); err != nil {
		return nil, ErrNotFound
	}
	b, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, policy); err != nil {
		return nil, errors.Wrap(ErrBadConversion, err.Error())
	}
	return policy, nil
}

// Delete removes a policy.
func (m *RedisManager) Delete(id string) error {
	key := prefixKey(m.keyPrefix, prefixPolicy, id)
	getCmd := m.db.Get(m.ctx, key)
	if err := getCmd.Err(); err != nil {
		return ErrNotFound
	}
	policy := &ladon.DefaultPolicy{}
	res, err := getCmd.Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(res), policy); err != nil {
		return errors.Wrap(ErrBadConversion, err.Error())
	}

	if err := m.db.Del(m.ctx, key).Err(); err != nil {
		return err
	}

	// Put this policy in the hashmap for each resource
	for _, v := range policy.GetResources() {
		hmkey := prefixKey(m.keyPrefix, prefixResource, v)
		field := policy.GetID()
		if err := m.db.HDel(m.ctx, hmkey, field).Err(); err != nil {
			return err
		}
	}

	// Put this policy in the hashmap for each subject
	for _, v := range policy.GetSubjects() {
		hmkey := prefixKey(m.keyPrefix, prefixSubject, v)
		field := policy.GetID()
		if err := m.db.HDel(m.ctx, hmkey, field).Err(); err != nil {
			return err
		}
	}

	return nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *RedisManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	policies := ladon.Policies{}

	var (
		rKey    = prefixKey(m.keyPrefix, prefixResource, resource)
		rGetCmd = m.db.HGetAll(m.ctx, rKey)
	)
	if err := rGetCmd.Err(); err != nil {
		return nil, err
	}

	rPolicies, err := rGetCmd.Result()
	if err != nil {
		return nil, err
	}

	for _, v := range rPolicies {
		p := &ladon.DefaultPolicy{}
		b := []byte(v)
		if err := json.Unmarshal(b, p); err != nil {
			return nil, errors.Wrap(ErrBadConversion, err.Error())
		}
		policies = append(policies, p)
	}

	return policies, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *RedisManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	policies := ladon.Policies{}

	var (
		sKey    = prefixKey(m.keyPrefix, prefixSubject, subject)
		sGetCmd = m.db.HGetAll(m.ctx, sKey)
	)
	if err := sGetCmd.Err(); err != nil {
		return nil, err
	}

	sPolicies, err := sGetCmd.Result()
	if err != nil {
		return nil, err
	}

	for _, v := range sPolicies {
		p := &ladon.DefaultPolicy{}
		b := []byte(v)
		if err := json.Unmarshal(b, p); err != nil {
			return nil, errors.Wrap(ErrBadConversion, err.Error())
		}
		policies = append(policies, p)
	}

	return policies, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *RedisManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	policies := ladon.Policies{}
	var (
		rKey    = prefixKey(m.keyPrefix, prefixResource, r.Resource)
		sKey    = prefixKey(m.keyPrefix, prefixSubject, r.Subject)
		rGetCmd = m.db.HGetAll(m.ctx, rKey)
		sGetCmd = m.db.HGetAll(m.ctx, sKey)
	)
	if err := rGetCmd.Err(); err != nil {
		return nil, err
	}
	if err := sGetCmd.Err(); err != nil {
		return nil, err
	}

	rPolicies, err := rGetCmd.Result()
	if err != nil {
		return nil, err
	}
	sPolicies, err := sGetCmd.Result()
	if err != nil {
		return nil, err
	}

	for _, v := range rPolicies {
		p := &ladon.DefaultPolicy{}
		b := []byte(v)
		// if !ok {
		// 	return nil, errors.Wrapf(ErrBadConversion, "value %+v is not a byte array", v)
		// }
		if err := json.Unmarshal(b, p); err != nil {
			return nil, errors.Wrap(ErrBadConversion, err.Error())
		}
		policies = append(policies, p)
	}

	for _, v := range sPolicies {
		p := &ladon.DefaultPolicy{}
		b := []byte(v)
		// if !ok {
		// 	return nil, errors.Wrapf(ErrBadConversion, "value %+v is not a byte array", v)
		// }
		if err := json.Unmarshal(b, p); err != nil {
			return nil, errors.Wrap(ErrBadConversion, err.Error())
		}
		policies = append(policies, p)
	}

	return policies, nil
}

func (m *RedisManager) Update(policy ladon.Policy) error {
	// Make sure that the key doesn't already exist
	key := prefixKey(m.keyPrefix, prefixPolicy, policy.GetID())
	if err := m.db.Get(m.ctx, key).Err(); err != nil {
		return ErrNotFound
	}

	p, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	// Set the policy key
	cmd := m.db.Set(m.ctx, key, p, 0)

	if err := cmd.Err(); err != nil {
		return err
	}

	// Put this policy in the hashmap for each resource
	for _, v := range policy.GetResources() {
		hmkey := prefixKey(m.keyPrefix, prefixResource, v)
		field := policy.GetID()
		if err := m.db.HMSet(m.ctx, hmkey, map[string]interface{}{
			field: p,
		}).Err(); err != nil {
			return err
		}
	}

	// Put this policy in the hashmap for each subject
	for _, v := range policy.GetSubjects() {
		hmkey := prefixKey(m.keyPrefix, prefixSubject, v)
		field := policy.GetID()
		if err := m.db.HMSet(m.ctx, hmkey, map[string]interface{}{
			field: p,
		}).Err(); err != nil {
			return err
		}
	}

	return nil
}
