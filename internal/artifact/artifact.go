package artifact

// nolint: gosec
import (
	"sync"

	"github.com/apex/log"
)

// Type defines the type of an artifact.
type Type int

const (
	// DockerImage is a published Docker image.
	DockerImage Type = iota
)

func (t Type) String() string {
	return "Docker Image"
}

// Artifact represents an artifact and its relevant info.
type Artifact struct {
	Name string
	Type Type
}

// Artifacts is a list of artifacts.
type Artifacts struct {
	items []*Artifact
	lock  *sync.Mutex
}

// New return a new list of artifacts.
func New() Artifacts {
	return Artifacts{
		items: []*Artifact{},
		lock:  &sync.Mutex{},
	}
}

// List return the actual list of artifacts.
func (artifacts Artifacts) List() []*Artifact {
	return artifacts.items
}

// Add safely adds a new artifact to an artifact list.
func (artifacts *Artifacts) Add(a *Artifact) {
	artifacts.lock.Lock()
	defer artifacts.lock.Unlock()
	log.WithFields(log.Fields{
		"name": a.Name,
		"type": a.Type,
	}).Debug("added new artifact")
	artifacts.items = append(artifacts.items, a)
}
