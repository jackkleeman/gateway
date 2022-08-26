package message

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/envoyproxy/gateway/internal/ir"
	xdstypes "github.com/envoyproxy/gateway/internal/xds/types"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	xdscachetypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resourcev3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func TestProviderResources(t *testing.T) {
	resources := new(ProviderResources)
	gc1 := &gwapiv1b1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gc1",
		},
	}
	gw1 := &gwapiv1b1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test1",
			Namespace: "test",
		},
	}
	r1 := &gwapiv1b1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route1",
			Namespace: "test",
		},
	}

	// Check init state
	assert.Nil(t, resources.GetGatewayClasses())
	assert.Nil(t, resources.GetGateways())
	assert.Nil(t, resources.GetHTTPRoutes())

	// Add resources
	resources.GatewayClasses.Store("test-gc1", gc1)

	gw1Key := types.NamespacedName{
		Namespace: gw1.GetNamespace(),
		Name:      gw1.GetName(),
	}
	resources.Gateways.Store(gw1Key, gw1)

	r1Key := types.NamespacedName{
		Namespace: r1.GetNamespace(),
		Name:      r1.GetName(),
	}
	resources.HTTPRoutes.Store(r1Key, r1)

	// Test
	gcs := resources.GetGatewayClasses()
	assert.Equal(t, len(gcs), 1)

	gws := resources.GetGateways()
	assert.Equal(t, len(gws), 1)

	hrs := resources.GetHTTPRoutes()
	assert.Equal(t, len(hrs), 1)

	// Add more resources
	gc2 := &gwapiv1b1.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-gc2",
		},
	}
	gw2 := &gwapiv1b1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test2",
			Namespace: "test",
		},
	}
	r2 := &gwapiv1b1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route2",
			Namespace: "test",
		},
	}

	resources.GatewayClasses.Store("test-gc2", gc2)
	gw2Key := types.NamespacedName{
		Namespace: gw2.GetNamespace(),
		Name:      gw2.GetName(),
	}
	resources.Gateways.Store(gw2Key, gw2)

	r2Key := types.NamespacedName{
		Namespace: r2.GetNamespace(),
		Name:      r2.GetName(),
	}
	resources.HTTPRoutes.Store(r2Key, r2)

	// Test contents

	gcs = resources.GetGatewayClasses()
	assert.ElementsMatch(t, gcs, []*gwapiv1b1.GatewayClass{gc1, gc2})

	gws = resources.GetGateways()
	assert.ElementsMatch(t, gws, []*gwapiv1b1.Gateway{gw1, gw2})

	hrs = resources.GetHTTPRoutes()
	assert.ElementsMatch(t, hrs, []*gwapiv1b1.HTTPRoute{r1, r2})
}

func TestXdsIR(t *testing.T) {
	msg := new(XdsIR)
	assert.Nil(t, msg.Get())
	in := &ir.Xds{
		HTTP: []*ir.HTTPListener{{Name: "test"}},
	}
	msg.Store("xds-ir", in)
	assert.Equal(t, msg.Get(), in)
}

func TestInfraIR(t *testing.T) {
	msg := new(InfraIR)
	assert.Nil(t, msg.Get())
	in := &ir.Infra{
		Proxy: &ir.ProxyInfra{Name: "test"},
	}
	msg.Store("infra-ir", in)
	assert.Equal(t, msg.Get(), in)
}

func TestXds(t *testing.T) {
	msg := new(Xds)
	assert.Nil(t, msg.Get())
	in := &xdstypes.ResourceVersionTable{
		XdsResources: xdstypes.XdsResources{
			resourcev3.ListenerType: []xdscachetypes.Resource{&listenerv3.Listener{Name: "test"}},
		},
	}
	msg.Store("xds", in)
	diff := cmp.Diff(in, msg.Get(), protocmp.Transform())
	require.Empty(t, diff)
}