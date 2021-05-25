package netrpc

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/diamondburned/cchat"
)

//go:generate go run ./cmd/entity_types entities.go

// EntityID describes the global ID for any type. It is service-scoped, meaning
// that entity IDs are required to be unique to a single service.
type EntityID struct {
	ID      cchat.ID   `json:"id"` // ignored if .Type == ServiceEntity
	Type    EntityType `json:"t"`
	Service cchat.ID   `json:"s"`
}

// GetID gets the EntityID of the given value.
func GetID(svc cchat.Service, v cchat.Identifier) (EntityID, bool) {
	typ := QueryEntityType(v)
	if typ == "" {
		return EntityID{}, false
	}

	svcID := svc.ID()
	if strings.Contains(svcID, "/") {
		return EntityID{}, false
	}

	return EntityID{
		ID:      v.ID(),
		Type:    typ,
		Service: svc.ID(),
	}, true
}

// ParseID parses the given ID into an EntityID.
func ParseID(id string) (EntityID, bool) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 {
		return EntityID{}, false
	}

	entID := EntityID{
		ID:      parts[2],
		Type:    EntityType(parts[1]),
		Service: parts[0],
	}

	if !entID.Type.IsValid() {
		return EntityID{}, false
	}

	return entID, true
}

// MarshalJSON marshals EntityID into a JSON object. It preallocates the
// returned byte slice.
func (id EntityID) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Grow(len(id.ID) + len(id.Type) + len(id.Service) + 64)

	if err := json.NewEncoder(&buf).Encode(id); err != nil {
		log.Panicln("failed to encode EntityID:", err)
	}

	return buf.Bytes(), nil
}

// EntityType is a string type for each entity.
type EntityType string
