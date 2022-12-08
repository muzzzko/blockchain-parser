package error

import (
	"errors"
	"fmt"
)

var (
	DomainErr  = errors.New("domain error")
	TimeoutErr = fmt.Errorf("http timeout: %w", DomainErr)
	HttpErr    = fmt.Errorf("http error: %w", DomainErr)

	SubscriberNotFound = fmt.Errorf("subscriber not found: %w", DomainErr)
	BlockNotFound      = fmt.Errorf("block not found: %w", DomainErr)
	NoBlockForParsing  = fmt.Errorf("no block for parsing: %w", DomainErr)
	UnknownBlockStatus = fmt.Errorf("unknown block status: %w", DomainErr)
)
