package common

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/bootstrap"
	"github.com/onflow/flow-go/model/flow"
)

// ReadFullPartnerNodeInfos reads partner node info and partner weight information from the specified paths and constructs
// a list of full bootstrap.NodeInfo for each partner node.
// Args:
// - log: the logger instance.
// - partnerWeightsPath: path to partner weights configuration file.
// - partnerNodeInfoDir: path to partner nodes configuration file.
// Returns:
// - []bootstrap.NodeInfo: the generated node info list.
// - error: if any error occurs. Any error returned from this function is irrecoverable.
func ReadFullPartnerNodeInfos(log zerolog.Logger, partnerWeightsPath, partnerNodeInfoDir string) []bootstrap.NodeInfo {
	partners := ReadPartnerNodes(log, partnerNodeInfoDir)
	log.Info().Msgf("read %d partner node configuration files", len(partners))

	weights, err := ReadPartnerWeights(partnerWeightsPath)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to read partner weights: %w", err))
	}

	var nodes []bootstrap.NodeInfo
	for _, partner := range partners {
		// validate every single partner node
		nodeID := ValidateNodeID(log, partner.NodeID)
		networkPubKey := ValidateNetworkPubKey(log, partner.NetworkPubKey)
		stakingPubKey := ValidateStakingPubKey(log, partner.StakingPubKey)
		weight, valid := ValidateWeight(weights[partner.NodeID])
		if !valid {
			log.Error().Msgf("weights: %v", weights)
			log.Fatal().Msgf("partner node id %x has no weight", nodeID)
		}
		if weight != flow.DefaultInitialWeight {
			log.Warn().Msgf("partner node (id=%x) has non-default weight (%d != %d)", partner.NodeID, weight, flow.DefaultInitialWeight)
		}

		node := bootstrap.NewPublicNodeInfo(
			nodeID,
			partner.Role,
			partner.Address,
			weight,
			networkPubKey.PublicKey,
			stakingPubKey.PublicKey,
		)
		nodes = append(nodes, node)
	}

	return nodes
}

// ReadPartnerWeights reads the partner weights configuration file and returns a list of PartnerWeights.
// Args:
// - partnerWeightsPath: path to partner weights configuration file.
// Returns:
// - PartnerWeights: the generated partner weights list.
// - error: if any error occurs. Any error returned from this function is irrecoverable.
func ReadPartnerWeights(partnerWeightsPath string) (PartnerWeights, error) {
	var weights PartnerWeights

	err := ReadJSON(partnerWeightsPath, &weights)
	if err != nil {
		return nil, fmt.Errorf("failed to read partner weights json: %w", err)
	}

	return weights, nil
}

// ReadPartnerNodes reads the partner node info from the configuration file and returns a list of []bootstrap.NodeInfoPub.
// Args:
// - partnerNodeInfoDir: path to partner nodes configuration file.
// Returns:
// - []bootstrap.NodeInfoPub: the generated partner node info list.
// - error: if any error occurs. Any error returned from this function is irrecoverable.
func ReadPartnerNodes(log zerolog.Logger, partnerNodeInfoDir string) []bootstrap.NodeInfoPub {
	var partners []bootstrap.NodeInfoPub
	files, err := FilesInDir(partnerNodeInfoDir)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read partner node infos")
	}
	for _, f := range files {
		// skip files that do not include node-infos
		if !strings.Contains(f, bootstrap.PathPartnerNodeInfoPrefix) {
			continue
		}

		// read file and append to partners
		var p bootstrap.NodeInfoPub
		err = ReadJSON(f, &p)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read node info")
		}
		partners = append(partners, p)
	}
	return partners
}

// ReadInternalNodeInfos returns a list of internal nodes after collecting weights
// from configuration files.
func ReadInternalNodeInfos(log zerolog.Logger, internalNodePrivInfoDir, internalWeightsConfig string) []bootstrap.NodeInfo {
	privInternals := ReadInternalNodes(log, internalNodePrivInfoDir)
	log.Info().Msgf("read %v internal private node-info files", len(privInternals))

	weights := internalWeightsByAddress(log, internalWeightsConfig)
	log.Info().Msgf("read %d weights for internal nodes", len(weights))

	var nodes []bootstrap.NodeInfo
	for _, internal := range privInternals {
		// check if address is valid format
		ValidateAddressFormat(log, internal.Address)

		// validate every single internal node
		nodeID := ValidateNodeID(log, internal.NodeID)
		weight, valid := ValidateWeight(weights[internal.Address])
		if !valid {
			log.Error().Msgf("weights: %v", weights)
			log.Fatal().Msgf("internal node %v has no weight. Did you forget to update the node address?", internal)
		}
		if weight != flow.DefaultInitialWeight {
			log.Warn().Msgf("internal node (id=%x) has non-default weight (%d != %d)", internal.NodeID, weight, flow.DefaultInitialWeight)
		}

		node := bootstrap.NewPrivateNodeInfo(
			nodeID,
			internal.Role,
			internal.Address,
			weight,
			internal.NetworkPrivKey,
			internal.StakingPrivKey,
		)

		nodes = append(nodes, node)
	}

	return nodes
}

// ReadInternalNodes reads our internal node private infos generated by
// `keygen` command and returns it
func ReadInternalNodes(log zerolog.Logger, internalNodePrivInfoDir string) []bootstrap.NodeInfoPriv {
	var internalPrivInfos []bootstrap.NodeInfoPriv

	// get files in internal priv node infos directory
	files, err := FilesInDir(internalNodePrivInfoDir)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read partner node infos")
	}

	// for each of the files
	for _, f := range files {
		// skip files that do not include node-infos
		if !strings.Contains(f, bootstrap.PathPrivNodeInfoPrefix) {
			continue
		}

		// read file and append to partners
		var p bootstrap.NodeInfoPriv
		err = ReadJSON(f, &p)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read json")
		}
		internalPrivInfos = append(internalPrivInfos, p)
	}

	return internalPrivInfos
}

// internalWeightsByAddress returns a mapping of node address by weight for internal nodes
func internalWeightsByAddress(log zerolog.Logger, config string) map[string]uint64 {
	// read json
	var configs []bootstrap.NodeConfig
	err := ReadJSON(config, &configs)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read json")
	}
	log.Info().Interface("config", configs).Msgf("read internal node configurations")

	weights := make(map[string]uint64)
	for _, config := range configs {
		if _, ok := weights[config.Address]; !ok {
			weights[config.Address] = config.Weight
		} else {
			log.Error().Msgf("duplicate internal node address %s", config.Address)
		}
	}

	return weights
}
