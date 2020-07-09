package adapter

// 说明：参考 https://github.com/casbin/casbin/blob/master/persist/file-adapter/adapter.go

import (
	"context"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/nilorg/naas/pkg/proto"
	"google.golang.org/grpc"
)

// Adapter is the file adapter for Casbin.
// It can load policy from file or save policy to file.
type Adapter struct {
	ctx            context.Context
	client         proto.CasbinAdapterClient
	ResourceID     string
	ResourceSecret string
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(ctx context.Context, clientConn grpc.ClientConnInterface) *Adapter {
	if ctx != nil {
		ctx = context.Background()
	}
	return &Adapter{
		ctx:    ctx,
		client: proto.NewCasbinAdapterClient(clientConn),
	}
}

// loadPolicyLine loads a text line as a policy rule to model.
func loadPolicyLine(line string, model model.Model) {
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}
	tokens := strings.Split(line, ",")
	for i := 0; i < len(tokens); i++ {
		tokens[i] = strings.TrimSpace(tokens[i])
	}

	key := tokens[0]
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])
}

// LoadPolicy loads all policy rules from the storage.
func (a *Adapter) LoadPolicy(model model.Model) error {
	return a.loadPolicyGrpc(model, loadPolicyLine)
}

func (a *Adapter) loadPolicyGrpc(model model.Model, handler func(string, model.Model)) (err error) {
	var resp *proto.LoadPolicyResponse
	resp, err = a.client.LoadPolicy(a.ctx, &proto.LoadPolicyRequest{
		ResourceId:     a.ResourceID,
		ResourceSecret: a.ResourceSecret,
	})
	if err != nil {
		return
	}
	for _, policy := range resp.Policys {
		line := strings.TrimSpace(policy)
		handler(line, model)
	}
	return
}

// SavePolicy saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	return errors.New("not implemented")
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// AddPolicies adds policy rules to the storage.
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicies removes policy rules from the storage.
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
