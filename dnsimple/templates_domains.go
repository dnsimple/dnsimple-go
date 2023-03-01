package dnsimple

import (
	"context"
	"errors"
	"fmt"
)

// ApplyTemplate applies a template to the given domain.
//
// See https://developer.dnsimple.com/v2/templates/domains/#applyTemplateToDomain
func (s *TemplatesService) ApplyTemplate(ctx context.Context, accountID string, templateIdentifier string, domainIdentifier string) (*TemplateResponse, error) {
	if templateIdentifier == "" {
		return nil, errors.New("empty path param")
	}

	domainPath, err := domainPath(accountID, domainIdentifier)
	if err != nil {
		return nil, err
	}

	path := versioned(fmt.Sprintf("%v/templates/%v", domainPath, templateIdentifier))

	templateResponse := &TemplateResponse{}

	resp, err := s.client.post(ctx, path, nil, nil)
	if err != nil {
		return nil, err
	}

	templateResponse.HTTPResponse = resp
	return templateResponse, nil
}
