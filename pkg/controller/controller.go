/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"github.com/trustbloc/edge-core/pkg/log"

	"github.com/trustbloc/edge-agent/pkg/controller/command"
	credentialclientcmd "github.com/trustbloc/edge-agent/pkg/controller/command/credentialclient"
	didclientcmd "github.com/trustbloc/edge-agent/pkg/controller/command/didclient"
	presentationclientcmd "github.com/trustbloc/edge-agent/pkg/controller/command/presentationclient"
	"github.com/trustbloc/edge-agent/pkg/controller/command/sdscomm"
)

var logger = log.New("edge-agent-didclient-controller")

type allOpts struct {
	blocDomain    string
	agentUsername string // TODO get username from the actual registration process instead of a cmd line arg #266
	sdsServerURL  string
}

// Opt represents a controller option.
type Opt func(opts *allOpts)

// WithBlocDomain is an option allowing for the trustbloc domain to be set.
func WithBlocDomain(blocDomain string) Opt {
	return func(opts *allOpts) {
		opts.blocDomain = blocDomain
	}
}

// TODO get username from the actual registration process instead of a cmd line arg #266
// WithAgentUsername is an option allowing for a username to be set.
func WithAgentUsername(agentUsername string) Opt {
	return func(opts *allOpts) {
		opts.agentUsername = agentUsername
	}
}

// WithSDSServerURL is an option allowing for the SDS server URL to be set.
func WithSDSServerURL(sdsServerURL string) Opt {
	return func(opts *allOpts) {
		opts.sdsServerURL = sdsServerURL
	}
}

// GetCommandHandlers returns all command handlers provided by controller.
func GetCommandHandlers(opts ...Opt) ([]command.Handler, error) {
	cmdOpts := &allOpts{}
	// Apply options
	for _, opt := range opts {
		opt(cmdOpts)
	}

	sdsComm := sdscomm.New(cmdOpts.sdsServerURL, cmdOpts.agentUsername)

	// did client command operation
	didClientCmd := didclientcmd.New(cmdOpts.blocDomain, sdsComm)

	// credential client command operation
	credentialClientCmd := credentialclientcmd.New(sdsComm)

	// presentation client command operation
	presentationClientCmd := presentationclientcmd.New(sdsComm)

	var allHandlers []command.Handler
	allHandlers = append(allHandlers, didClientCmd.GetHandlers()...)
	allHandlers = append(allHandlers, credentialClientCmd.GetHandlers()...)
	allHandlers = append(allHandlers, presentationClientCmd.GetHandlers()...)

	return allHandlers, nil
}
