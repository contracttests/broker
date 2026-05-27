package can_i_deploy

import (
	"fmt"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/gofiber/fiber/v3"
)

func (h *CanIDeployHandler) evaluate(ctx fiber.Ctx, askerID int64, version string, environmentID int64, askerName string) bool {
	askerContract := h.mustLoadContract(ctx, askerName, version)
	txStart := h.clockRepository.DBNow(ctx.Context())

	counterparts := h.resolveCounterparts(ctx, askerID, askerContract, environmentID)
	pairResults := h.evaluatePairs(ctx, askerContract, counterparts)
	h.persistMatrix(ctx, askerID, version, counterparts, pairResults)

	if counterparts.hasGhost() {
		return false
	}
	return !h.compatibilityMatrixRepository.AnyFailureSince(ctx.Context(), askerID, version, txStart)
}

func (h *CanIDeployHandler) resolveCounterparts(ctx fiber.Ctx, askerID int64, askerContract *model.Contract, environmentID int64) *counterpartSet {
	set := newCounterpartSet()
	deployedProviders := h.indexDeploymentsByName(ctx, askerID, environmentID)

	for _, resource := range askerContract.Resources {
		if resource.Direction != model.Consumes {
			continue
		}
		h.classifyConsumedProvider(ctx, resource.Provider, askerID, deployedProviders, set)
	}

	for _, resource := range askerContract.Resources {
		if resource.Direction != model.Provides {
			continue
		}
		consumers := h.contractRepository.FindCurrentConsumersOfProviderInEnv(ctx.Context(), resource.ProviderHash(), environmentID)
		for _, consumer := range consumers {
			if consumer.ParticipantID == askerID {
				continue
			}
			set.addPair(counterpartFromConsumer(consumer))
		}
	}

	set.dropStrictFalsesAlreadyPaired()
	return set
}

func (h *CanIDeployHandler) classifyConsumedProvider(
	ctx fiber.Ctx,
	providerName string,
	askerID int64,
	deployedProviders map[string]model.Deployment,
	set *counterpartSet,
) {
	if deployed, ok := deployedProviders[providerName]; ok {
		set.addPair(counterpartFromDeployment(deployed))
		return
	}
	participant, found := h.participantRepository.FindByName(ctx.Context(), providerName)
	if !found {
		set.markGhost()
		return
	}
	if participant.ID == askerID {
		return
	}
	set.markStrictFalse(participant.ID)
}

func (h *CanIDeployHandler) indexDeploymentsByName(ctx fiber.Ctx, askerID, environmentID int64) map[string]model.Deployment {
	deployments := h.deploymentRepository.ListCurrentDeploymentsInEnv(ctx.Context(), environmentID)
	index := make(map[string]model.Deployment, len(deployments))
	for _, deployment := range deployments {
		if deployment.Participant.ID == askerID {
			continue
		}
		index[deployment.Participant.Name] = deployment
	}
	return index
}

func (h *CanIDeployHandler) evaluatePairs(ctx fiber.Ctx, askerContract *model.Contract, set *counterpartSet) map[int64]bool {
	results := make(map[int64]bool, set.pairCount())
	set.eachPair(func(pair counterpartParticipant) {
		counterpartContract := h.mustLoadContract(ctx, pair.name, pair.version)
		results[pair.id] = pairDeployable(askerContract, counterpartContract, pair.name)
	})
	return results
}

func (h *CanIDeployHandler) persistMatrix(ctx fiber.Ctx, askerID int64, version string, set *counterpartSet, pairResults map[int64]bool) {
	tx := h.compatibilityMatrixRepository.BeginTx(ctx.Context())
	defer tx.Rollback(ctx.Context())

	inserted := 0
	set.eachStrictFalse(func(counterpartID int64) {
		h.compatibilityMatrixRepository.Insert(ctx.Context(), tx, model.NewStrictFalseRow(askerID, version, counterpartID))
		inserted++
	})
	set.eachPair(func(pair counterpartParticipant) {
		h.compatibilityMatrixRepository.Insert(ctx.Context(), tx, model.NewPairCheckedRow(askerID, version, pair.id, pair.version, pairResults[pair.id]))
		inserted++
	})
	if inserted == 0 {
		h.compatibilityMatrixRepository.Insert(ctx.Context(), tx, model.NewVacuousTrueRow(askerID, version))
	}

	if err := tx.Commit(ctx.Context()); err != nil {
		panic(fmt.Errorf("error committing compatibility_matrix transaction: %w", err))
	}
}

func (h *CanIDeployHandler) mustLoadContract(ctx fiber.Ctx, name, version string) *model.Contract {
	contract, found := h.contractRepository.LoadContractByNameAndVersion(ctx.Context(), name, version)
	if !found {
		panic(fmt.Errorf("contract not found: %s@%s", name, version))
	}
	return contract
}
