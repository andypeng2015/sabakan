package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/cybozu-go/log"
	sabakan "github.com/cybozu-go/sabakan/v2"
	"github.com/cybozu-go/sabakan/v2/gql"
	"github.com/cybozu-go/sabakan/v2/gql/graph/generated"
	"github.com/cybozu-go/sabakan/v2/gql/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *bMCResolver) BmcType(ctx context.Context, obj *sabakan.MachineBMC) (string, error) {
	return obj.Type, nil
}

func (r *bMCResolver) Ipv4(ctx context.Context, obj *sabakan.MachineBMC) (*gql.IPAddress, error) {
	return &gql.IPAddress{IP: net.ParseIP(obj.IPv4)}, nil
}

func (r *machineSpecResolver) Labels(ctx context.Context, obj *sabakan.MachineSpec) ([]*model.Label, error) {
	if len(obj.Labels) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(obj.Labels))
	for k := range obj.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	labels := make([]*model.Label, 0, len(obj.Labels))
	for _, k := range keys {
		labels = append(labels, &model.Label{Name: k, Value: obj.Labels[k]})
	}
	return labels, nil
}

func (r *machineSpecResolver) Rack(ctx context.Context, obj *sabakan.MachineSpec) (int, error) {
	return int(obj.Rack), nil
}

func (r *machineSpecResolver) IndexInRack(ctx context.Context, obj *sabakan.MachineSpec) (int, error) {
	return int(obj.IndexInRack), nil
}

func (r *machineSpecResolver) Ipv4(ctx context.Context, obj *sabakan.MachineSpec) ([]*gql.IPAddress, error) {
	addresses := make([]*gql.IPAddress, len(obj.IPv4))
	for i, a := range obj.IPv4 {
		addresses[i] = &gql.IPAddress{IP: net.ParseIP(a)}
	}
	return addresses, nil
}

func (r *machineSpecResolver) RegisterDate(ctx context.Context, obj *sabakan.MachineSpec) (*gql.DateTime, error) {
	t := gql.DateTime(obj.RegisterDate)
	return &t, nil
}

func (r *machineSpecResolver) RetireDate(ctx context.Context, obj *sabakan.MachineSpec) (*gql.DateTime, error) {
	t := gql.DateTime(obj.RetireDate)
	return &t, nil
}

func (r *machineStatusResolver) Timestamp(ctx context.Context, obj *sabakan.MachineStatus) (*gql.DateTime, error) {
	t := gql.DateTime(obj.Timestamp)
	return &t, nil
}

func (r *mutationResolver) SetMachineState(ctx context.Context, serial string, state sabakan.MachineState) (*sabakan.MachineStatus, error) {
	now := time.Now()

	log.Info("SetMachineState is called", map[string]interface{}{
		"serial": serial,
		"state":  state,
	})

	err := r.Model.Machine.SetState(ctx, serial, state)
	if err != nil {
		switch err {
		case sabakan.ErrNotFound:
			return &sabakan.MachineStatus{}, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"serial": serial,
					"type":   gql.ErrMachineNotFound,
				},
			}
		case sabakan.ErrEncryptionKeyExists:
			return &sabakan.MachineStatus{}, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"serial": serial,
					"type":   gql.ErrEncryptionKeyExists,
				},
			}
		default:
			var from, to string
			_, err2 := fmt.Sscanf(err.Error(), sabakan.SetStateErrorFormat, &from, &to)
			if err2 != nil {
				return &sabakan.MachineStatus{}, &gqlerror.Error{
					Message: err.Error(),
					Extensions: map[string]interface{}{
						"serial": serial,
						"type":   gql.ErrInternalServerError,
					},
				}
			}
			return &sabakan.MachineStatus{}, &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"serial": serial,
					"type":   gql.ErrInvalidStateTransition,
				},
			}
		}
	}

	machine, err := r.Model.Machine.Get(ctx, serial)
	if err != nil {
		return &sabakan.MachineStatus{}, err
	}
	machine.Status.Duration = now.Sub(machine.Status.Timestamp).Seconds()
	return &machine.Status, nil
}

func (r *nICConfigResolver) Address(ctx context.Context, obj *sabakan.NICConfig) (*gql.IPAddress, error) {
	return &gql.IPAddress{IP: net.ParseIP(obj.Address)}, nil
}

func (r *nICConfigResolver) Netmask(ctx context.Context, obj *sabakan.NICConfig) (*gql.IPAddress, error) {
	return &gql.IPAddress{IP: net.ParseIP(obj.Netmask)}, nil
}

func (r *nICConfigResolver) Gateway(ctx context.Context, obj *sabakan.NICConfig) (*gql.IPAddress, error) {
	return &gql.IPAddress{IP: net.ParseIP(obj.Gateway)}, nil
}

func (r *queryResolver) Machine(ctx context.Context, serial string) (*sabakan.Machine, error) {
	now := time.Now()

	log.Info("Machine is called", map[string]interface{}{
		"serial": serial,
	})

	machine, err := r.Model.Machine.Get(ctx, serial)
	if err != nil {
		return &sabakan.Machine{}, err
	}
	machine.Status.Duration = now.Sub(machine.Status.Timestamp).Seconds()
	return machine, nil
}

func (r *queryResolver) SearchMachines(ctx context.Context, having *model.MachineParams, notHaving *model.MachineParams) ([]*sabakan.Machine, error) {
	now := time.Now()

	log.Info("SearchMachines is called", map[string]interface{}{
		"having":    having,
		"nothaving": notHaving,
	})

	machines, err := r.Model.Machine.Query(ctx, sabakan.Query{})
	if err != nil {
		return nil, err
	}
	var filtered []*sabakan.Machine
	for _, m := range machines {
		m.Status.Duration = now.Sub(m.Status.Timestamp).Seconds()
		if gql.MatchMachine(m, having, notHaving, now) {
			filtered = append(filtered, m)
		}
	}
	return filtered, nil
}

// BMC returns generated.BMCResolver implementation.
func (r *Resolver) BMC() generated.BMCResolver { return &bMCResolver{r} }

// MachineSpec returns generated.MachineSpecResolver implementation.
func (r *Resolver) MachineSpec() generated.MachineSpecResolver { return &machineSpecResolver{r} }

// MachineStatus returns generated.MachineStatusResolver implementation.
func (r *Resolver) MachineStatus() generated.MachineStatusResolver { return &machineStatusResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// NICConfig returns generated.NICConfigResolver implementation.
func (r *Resolver) NICConfig() generated.NICConfigResolver { return &nICConfigResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type bMCResolver struct{ *Resolver }
type machineSpecResolver struct{ *Resolver }
type machineStatusResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type nICConfigResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }