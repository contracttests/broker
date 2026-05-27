package can_i_deploy

import "github.com/contracttesting/broker/server/internal/model"

func pairDeployable(asker, counterpart *model.Contract, counterpartName string) bool {
	if !askerConsumesAreCompatible(asker, counterpart, counterpartName) {
		return false
	}
	return counterpartConsumesAreCompatible(asker, counterpart)
}

func askerConsumesAreCompatible(asker, counterpart *model.Contract, counterpartName string) bool {
	for _, consumed := range asker.Resources {
		if consumed.Direction != model.Consumes || consumed.Provider != counterpartName {
			continue
		}
		provided, exists := counterpart.Resources[consumed.ProviderHash()]
		if !exists {
			continue
		}
		if hasBreaks(consumed, provided, model.UploaderConsumer) {
			return false
		}
	}
	return true
}

func counterpartConsumesAreCompatible(asker, counterpart *model.Contract) bool {
	askerName := asker.Participant.Name
	for _, provided := range asker.Resources {
		if provided.Direction != model.Provides {
			continue
		}
		if counterpartHasBreakingConsumption(counterpart, provided, askerName) {
			return false
		}
	}
	return true
}

func counterpartHasBreakingConsumption(counterpart *model.Contract, provided model.Resource, askerName string) bool {
	for _, consumed := range counterpart.Resources {
		if consumed.Direction != model.Consumes {
			continue
		}
		if consumed.Provider != askerName {
			continue
		}
		if consumed.ProviderHash() != provided.ProviderHash() {
			continue
		}
		if hasBreaks(consumed, provided, model.UploaderProvider) {
			return true
		}
	}
	return false
}

func hasBreaks(consumer, provider model.Resource, role model.UploaderRole) bool {
	return len(model.Compare(model.CompareInput{
		Consumer:     consumer,
		Provider:     provider,
		UploaderRole: role,
	})) > 0
}
