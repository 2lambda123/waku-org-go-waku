//go:build gowaku_rln
// +build gowaku_rln

package waku

import (
	"crypto/ecdsa"
	"errors"

	"github.com/status-im/go-rln/rln"
	"github.com/status-im/go-waku/waku/v2/node"
	"go.uber.org/zap"
)

var loadedCredentialsFromFile bool = false

func checkForRLN(logger *zap.Logger, options Options, nodeOpts *[]node.WakuNodeOption) {
	if options.RLNRelay.Enable {
		if !options.Relay.Enable {
			failOnErr(errors.New("relay not available"), "Could not enable RLN Relay")
		}
		if !options.RLNRelay.Dynamic {
			*nodeOpts = append(*nodeOpts, node.WithStaticRLNRelay(options.RLNRelay.PubsubTopic, options.RLNRelay.ContentTopic, rln.MembershipIndex(options.RLNRelay.MembershipIndex), nil))
		} else {

			var ethPrivKey *ecdsa.PrivateKey
			if options.RLNRelay.ETHPrivateKey != nil {
				ethPrivKey = options.RLNRelay.ETHPrivateKey
			}

			loaded, idKey, idCommitment, membershipIndex, err := getMembershipCredentials(logger, options)
			failOnErr(err, "Invalid membership credentials")

			loadedCredentialsFromFile = loaded

			*nodeOpts = append(*nodeOpts, node.WithDynamicRLNRelay(
				options.RLNRelay.PubsubTopic,
				options.RLNRelay.ContentTopic,
				membershipIndex,
				idKey,
				idCommitment,
				nil,
				options.RLNRelay.ETHClientAddress,
				ethPrivKey,
				options.RLNRelay.MembershipContractAddress,
				nil,
			))
		}
	}
}

func onStartRLN(wakuNode *node.WakuNode, options Options) {
	if options.RLNRelay.Enable && options.RLNRelay.Dynamic && !loadedCredentialsFromFile && options.RLNRelay.CredentialsPath != "" {
		err := writeRLNMembershipCredentialsToFile(wakuNode.RLNRelay().MembershipKeyPair(), wakuNode.RLNRelay().MembershipIndex(), options.RLNRelay.CredentialsPath, []byte(options.KeyPasswd), options.Overwrite)
		failOnErr(err, "Could not write membership credentials file")
	}
}