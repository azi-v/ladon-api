package domain

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/azi-v/ladon-api/DB"
	"github.com/ory/ladon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBManager struct {
	DB  *mongo.Client
	Ctx context.Context
}

func NewPolicyMongoDBManager(ctx context.Context, db *mongo.Client) *MongoDBManager {
	return &MongoDBManager{
		DB:  db,
		Ctx: ctx,
	}
}

// Create persists the policy.
func (m *MongoDBManager) Create(policy ladon.Policy) error {
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")
	filter := bson.M{"policy_id": policy.GetID()}

	count, err := coll.CountDocuments(m.Ctx, filter)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New("the policy policy_id has existed!")
	}

	p, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies").InsertOne(m.Ctx, p)
	if err != nil {
		return err
	}
	for _, v := range policy.GetResources() {
		doc := bson.D{
			{Key: "policy_id", Value: policy.GetID()},
			{Key: "resource", Value: v},
		}
		coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-resources")
		if _, err := coll.InsertOne(m.Ctx, doc); err != nil {
			// TODO: 这里的事务问题
			return err
		}
	}

	for _, v := range policy.GetSubjects() {
		doc := bson.D{
			{Key: "policy_id", Value: policy.GetID()},
			{Key: "subject", Value: v},
		}
		coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-subjects")
		if _, err := coll.InsertOne(m.Ctx, doc); err != nil {
			// TODO: 这里的事务问题
			return err
		}
	}

	return nil
}

// Update updates an existing policy.
func (m *MongoDBManager) Update(policy ladon.Policy) error {
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")
	filter := bson.M{"policy_id": policy.GetID()}

	count, err := coll.CountDocuments(m.Ctx, filter)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.New("the policy policy_id has existed!")
	}

	// TODO: 这里考虑直接执行mongo shell的原生语句
	// m.DB.Database("").RunCommand(m.Ctx, "")
	p, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	filter = bson.M{"policy_id": policy.GetID()}
	_, err = coll.UpdateOne(m.Ctx, filter, p)
	if err != nil {
		return err
	}

	for _, v := range policy.GetResources() {
		filter := bson.D{
			{Key: "policy_id", Value: policy.GetID()},
		}
		doc := bson.D{{Key: "$set", Value: bson.D{{Key: "resource", Value: v}}}}
		coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-resources")
		if _, err := coll.UpdateOne(m.Ctx, filter, doc); err != nil {
			// TODO: 这里的事务问题
			return err
		}
	}

	for _, v := range policy.GetSubjects() {
		filter := bson.D{
			{Key: "policy_id", Value: policy.GetID()},
		}
		doc := bson.D{{Key: "$set", Value: bson.D{{Key: "subject", Value: v}}}}

		coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-subjects")
		if _, err := coll.UpdateOne(m.Ctx, filter, doc); err != nil {
			// TODO: 这里的事务问题
			return err
		}
	}

	return nil
}

// Get retrieves a policy.
func (m *MongoDBManager) Get(id string) (ladon.Policy, error) {
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")
	filter := bson.M{"policy_id": id} //TODO: policy_id增加唯一索引

	pol := &ladon.DefaultPolicy{}
	err := coll.FindOne(m.Ctx, filter).Decode(pol)
	return pol, err
}

// Delete removes a policy.
func (m *MongoDBManager) Delete(id string) error {
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")
	filter := bson.M{"policy_id": id}

	pol := &ladon.DefaultPolicy{}
	err := coll.FindOne(m.Ctx, filter).Decode(pol)
	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(m.Ctx, filter)
	if err != nil {
		return err
	}

	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-resources")
	if _, err := coll.DeleteOne(m.Ctx, filter); err != nil {
		// TODO: 这里的事务问题
		return err
	}

	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-subjects")
	if _, err := coll.DeleteOne(m.Ctx, filter); err != nil {
		// TODO: 这里的事务问题
		return err
	}

	return nil
}

// GetAll retrieves all policies.
func (m *MongoDBManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")
	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := coll.Find(m.Ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var pols []*ladon.DefaultPolicy
	if err := cursor.All(m.Ctx, &pols); err != nil {
		return nil, err
	}

	res := make(ladon.Policies, len(pols))
	for i, v := range pols {
		res[i] = v
	}

	return res, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *MongoDBManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	policies := ladon.Policies{}

	filter := bson.D{
		{Key: "subject", Value: r.Subject},
	}
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-subject")
	cursor, err := coll.Find(m.Ctx, filter)
	if err != nil {
		// TODO: 这里的事务问题
		return nil, err
	}

	allSubject := make([]DB.CollSubject, 0)
	err = cursor.All(m.Ctx, &allSubject)
	if err != nil {
		return nil, err
	}

	for _, v := range allSubject {
		filter := bson.M{"policy_id": v.PolicyID}
		p := &ladon.DefaultPolicy{}

		err := coll.FindOne(m.Ctx, filter).Decode(p)
		if err != nil {
			return nil, err
		}

		policies = append(policies, p)
	}

	filter = bson.D{
		{Key: "resource", Value: r.Resource},
	}
	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-resources")
	cursor, err = coll.Find(m.Ctx, filter)
	if err != nil {
		// TODO: 这里的事务问题
		return nil, err
	}

	allResource := make([]DB.CollResources, 0)
	err = cursor.All(m.Ctx, &allResource)
	if err != nil {
		return nil, err
	}

	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")

	for _, v := range allResource {
		filter := bson.M{"policy_id": v.PolicyID}
		p := &ladon.DefaultPolicy{}

		err := coll.FindOne(m.Ctx, filter).Decode(p)
		if err != nil {
			return nil, err
		}

		policies = append(policies, p)
	}

	return policies, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *MongoDBManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	filter := bson.D{
		{Key: "subject", Value: subject},
	}
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-subject")
	cursor, err := coll.Find(m.Ctx, filter)
	if err != nil {
		// TODO: 这里的事务问题
		return nil, err
	}

	allSubject := make([]DB.CollSubject, 0)
	err = cursor.All(m.Ctx, &allSubject)
	if err != nil {
		return nil, err
	}

	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")

	policies := ladon.Policies{}
	for _, v := range allSubject {
		filter := bson.M{"policy_id": v.PolicyID}
		p := &ladon.DefaultPolicy{}

		err := coll.FindOne(m.Ctx, filter).Decode(p)
		if err != nil {
			return nil, err
		}

		policies = append(policies, p)
	}

	return policies, nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *MongoDBManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {

	filter := bson.D{
		{Key: "resource", Value: resource},
	}
	coll := m.DB.Database("ymt-usercenter").Collection("ymt-ladon-resources")
	cursor, err := coll.Find(m.Ctx, filter)
	if err != nil {
		// TODO: 这里的事务问题
		return nil, err
	}

	allResource := make([]DB.CollResources, 0)
	err = cursor.All(m.Ctx, &allResource)
	if err != nil {
		return nil, err
	}

	coll = m.DB.Database("ymt-usercenter").Collection("ymt-ladon-policies")

	policies := ladon.Policies{}
	for _, v := range allResource {
		filter := bson.M{"policy_id": v.PolicyID}
		p := &ladon.DefaultPolicy{}

		err := coll.FindOne(m.Ctx, filter).Decode(p)
		if err != nil {
			return nil, err
		}

		policies = append(policies, p)
	}

	return policies, nil
}
