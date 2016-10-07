package dnsimple

import (
	"fmt"
	// "strconv"
)

// CollaboratorsService handles communication with the collaborator related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/domains/collaborators
type CollaboratorsService struct {
	client *Client
}

// Collaborator represents a Collaborator in DNSimple.
type Collaborator struct {
	ID         int    `json:"id,omitempty"`
	DomainID   int    `json:"domain_id,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
	UserEmail  string `json:"user_email,omitempty"`
	Invitation bool   `json:"invitation,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	AcceptedAt string `json:"accepted_at,omitempty"`
}

func collaboratorPath(accountID, domainIdentifier, collaboratorID string) string {
	path := fmt.Sprintf("%v/collaborators", domainPath(accountID, domainIdentifier))

	if collaboratorID != "" {
		path += fmt.Sprintf("/%v", collaboratorID)
	}
	return path
}

// CollaboratorResponse represents a response from an API method that returns a Collaborator struct.
type CollaboratorResponse struct {
	Response
	Data *Collaborator `json:"data"`
}

// CollaboratorsResponse represents a response from an API method that returns a collection of Collaborator struct.
type CollaboratorsResponse struct {
	Response
	Data []Collaborator `json:"data"`
}

// ListCollaborators list the collaborators for a domain.
//
// See https://developer.dnsimple.com/v2/domains/collaborators#list
func (s *CollaboratorsService) ListCollaborators(accountID, domainIdentifier string, options *ListOptions) (*CollaboratorsResponse, error) {
	path := versioned(collaboratorPath(accountID, domainIdentifier, ""))
	collaboratorsResponse := &CollaboratorsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, collaboratorsResponse)
	if err != nil {
		return collaboratorsResponse, err
	}

	collaboratorsResponse.HttpResponse = resp
	return collaboratorsResponse, nil
}
