package dnsimple

import (
	"context"
	"fmt"
)

// Collaborator represents a Collaborator in DNSimple.
type Collaborator struct {
	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	UserEmail  string `json:"user_email,omitempty"`
	Invitation bool   `json:"invitation,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	AcceptedAt string `json:"accepted_at,omitempty"`
}

func collaboratorsPath(accountID, domainIdentifier string) (string, error) {
	basePath, err := domainPath(accountID, domainIdentifier)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/collaborators", basePath), nil
}

func collaboratorPath(accountID, domainIdentifier string, collaboratorID int64) (string, error) {
	basePath, err := collaboratorsPath(accountID, domainIdentifier)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/%v", basePath, collaboratorID), nil
}

// CollaboratorAttributes represents Collaborator attributes for AddCollaborator operation.
type CollaboratorAttributes struct {
	Email string `json:"email,omitempty"`
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
func (s *DomainsService) ListCollaborators(ctx context.Context, accountID, domainIdentifier string, options *ListOptions) (*CollaboratorsResponse, error) {
	path, err := collaboratorsPath(accountID, domainIdentifier)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	collaboratorsResponse := &CollaboratorsResponse{}

	path, err = addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(ctx, path, collaboratorsResponse)
	if err != nil {
		return collaboratorsResponse, err
	}

	collaboratorsResponse.HTTPResponse = resp
	return collaboratorsResponse, nil
}

// AddCollaborator adds a new collaborator to the domain in the account.
//
// See https://developer.dnsimple.com/v2/domains/collaborators#add
func (s *DomainsService) AddCollaborator(ctx context.Context, accountID string, domainIdentifier string, attributes CollaboratorAttributes) (*CollaboratorResponse, error) {
	path, err := collaboratorsPath(accountID, domainIdentifier)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	collaboratorResponse := &CollaboratorResponse{}

	resp, err := s.client.post(ctx, path, attributes, collaboratorResponse)
	if err != nil {
		return nil, err
	}

	collaboratorResponse.HTTPResponse = resp
	return collaboratorResponse, nil
}

// RemoveCollaborator PERMANENTLY deletes a domain from the account.
//
// See https://developer.dnsimple.com/v2/domains/collaborators#remove
func (s *DomainsService) RemoveCollaborator(ctx context.Context, accountID string, domainIdentifier string, collaboratorID int64) (*CollaboratorResponse, error) {
	path, err := collaboratorPath(accountID, domainIdentifier, collaboratorID)
	if err != nil {
		return nil, err
	}

	path = versioned(path)

	collaboratorResponse := &CollaboratorResponse{}

	resp, err := s.client.delete(ctx, path, nil, nil)
	if err != nil {
		return nil, err
	}

	collaboratorResponse.HTTPResponse = resp
	return collaboratorResponse, nil
}
