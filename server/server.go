package server

import (
	"log"
	"net/rpc"
	"strings"
	"sync"

	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-netrpc/internal/stdio"
	"github.com/diamondburned/cchat-netrpc/netrpc"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

var services struct {
	sync.RWMutex

	names    []cchat.ID
	services map[cchat.ID]*serviceMap
}

// serviceMap contains the cchat service and the associated entity registry.
type serviceMap struct {
	service  cchat.Service
	entities map[netrpc.EntityType]map[cchat.ID]cchat.Identifier
}

// newServiceMap creates a new serviceMap.
func newServiceMap(service cchat.Service) *serviceMap {
	svcMap := serviceMap{service: service}
	for i := range svcMap.entities {
		svcMap.entities[i] = make(map[cchat.ID]cchat.Identifier)
	}
	return &svcMap
}

// Register registers a cchat.Service into the global RPC server. This function
// must be called before Serve is called; it is best called in init(). The
// service's ID must not contain any slashes, or the function will panic.
func Register(service cchat.Service) {
	id := service.ID()
	if strings.Contains(id, "/") {
		log.Panicf("service ID %q contains a slash", id)
	}

	services.Lock()
	defer services.Unlock()

	services.names = append(services.names, id)
	services.services[id] = newServiceMap(service)
}

// ServiceIDs returns known service IDs.
func ServiceIDs() []cchat.ID {
	services.RLock()
	defer services.RUnlock()

	return append([]cchat.ID(nil), services.names...)
}

// Entity looks up the entity registry. Nil is returned if the ID is not found.
func Entity(id netrpc.EntityID) cchat.Identifier {
	services.RLock()
	defer services.RUnlock()

	svcMap, ok := services.services[id.Service]
	if !ok {
		return nil
	}

	if id.Type == netrpc.ServiceEntity {
		return svcMap.service
	}

	entMap, ok := svcMap.entities[id.Type]
	if !ok {
		return nil
	}

	return entMap[id.ID]
}

// FreeEntity frees the entity from the global registry. If Type is a Service,
// then the function does nothing.
func FreeEntity(id netrpc.EntityID) {
	if id.Type == netrpc.ServiceEntity {
		return
	}

	services.Lock()
	defer services.Unlock()

	svcMap, ok := services.services[id.Service]
	if !ok {
		return
	}

	delete(svcMap.entities[id.Type], id.ID)
}

// PutEntity puts the given Identifier entity into the global map. It does
// nothing if ider is of type cchat.Service.
func PutEntity(svc cchat.Service, ider cchat.Identifier) {
	typ := netrpc.QueryEntityType(ider)
	if typ == netrpc.ServiceEntity {
		return
	}

	id := ider.ID()
	serviceID := svc.ID()

	services.Lock()
	defer services.Unlock()

	svcMap, ok := services.services[serviceID]
	if !ok {
		return
	}

	svcMap.entities[typ][id] = ider
}

// Serve serves the RPC over stdio using JSONRPC2.
func Serve() {
	io, err := stdio.Take()
	if err != nil {
		log.Panicln("failed to take stdio:", err)
	}

	codec := jsonrpc2.NewServerCodec(io, nil)
	rpc.ServeCodec(codec)
}

// ServeWithCodec serves the giveen cchat service with the given RPC server
// codec over stdio.
func ServeWithCodec(codec rpc.ServerCodec) {
	rpc.ServeCodec(codec)
}
