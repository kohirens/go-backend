package backend

import "fmt"

type AccountNotFoundError struct {
	id string
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf(stderr.AccountNotFound, e.id)
}

type ProviderNotFound struct {
	name string
}

func (e *ProviderNotFound) Error() string {
	return fmt.Sprintf(stderr.ProviderNotFound, e.name)
}

// ReferralError Thrown when processing the request has failed in some way and
// there is somewhere else to send the client.
// NOTE: This is good for HTTP codes such as SeeOther, 302,301, etc. Consider
// this a treatable error, that you choose to log or not.
type ReferralError struct {
	Body        []byte
	Code        int
	ContentType string
	Location    string
	Log         bool
	Message     string
}

func (e *ReferralError) Error() string {
	return fmt.Sprintf(stderr.SeeOther, e.Location)
}

func NewReferralError(contentType, msg, loc string, code int, log bool) *ReferralError {
	return &ReferralError{
		Code:        code,
		ContentType: contentType,
		Message:     msg,
		Location:    loc,
		Log:         log,
	}
}

type ServiceNotFoundError struct {
	name string
}

func (e *ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service %v not found", e.name)
}
