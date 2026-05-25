package model

type CompatibilityReport struct {
	breaks []BreakingChange
}

func (r *CompatibilityReport) Append(b BreakingChange) {
	r.breaks = append(r.breaks, b)
}

func (r *CompatibilityReport) Breaks() []BreakingChange {
	if r.breaks == nil {
		return []BreakingChange{}
	}
	return r.breaks
}

type CompareInput struct {
	Consumer     Resource
	Provider     Resource
	UploaderRole UploaderRole
}

func NewMissingProviderBreak(consumer Resource, reason BreakingReason) BreakingChange {
	change := BreakingChange{
		UploaderRole: UploaderConsumer,
		Resource:     NewBrokenResource(consumer),
		Reason:       reason,
	}

	switch reason {
	case ReasonProviderNotSpecified:
		// counterpart is nil — the consumer did not name any provider
	case ReasonProviderNotFound, ReasonProviderResourceNotFound:
		// Name-only counterpart sourced from the consumer's declaration; owner
		// is unknown because no stored provider contract exists (or was found)
		// to source it from. Revisit if a ContractRepository.LoadProviderByName
		// is introduced for ReasonProviderResourceNotFound.
		change.Counterpart = &CounterpartInfo{
			Role: UploaderProvider,
			Name: consumer.Provider,
		}
	}

	return change
}

func Compare(in CompareInput) []BreakingChange {
	switch in.Consumer.Kind {
	case RestRequest:
		return compareRequest(in)
	default:
		return compareResponse(in)
	}
}

func counterpartFor(in CompareInput) *CounterpartInfo {
	counterpartRole := UploaderProvider
	storedInfo := in.Provider.Participant
	if in.UploaderRole == UploaderProvider {
		counterpartRole = UploaderConsumer
		storedInfo = in.Consumer.Participant
	}

	if storedInfo == nil {
		return nil
	}

	return &CounterpartInfo{
		Role: counterpartRole,
		Name: storedInfo.Name,
	}
}

func compareResponse(in CompareInput) []BreakingChange {
	brokenResource := NewBrokenResource(in.Consumer)
	counterpart := counterpartFor(in)

	var breaks []BreakingChange
	for consumerPropertyPath, consumerProperty := range in.Consumer.Properties {
		providerProperty, propertyExists := in.Provider.Properties[consumerPropertyPath]

		if !propertyExists {
			breaks = append(breaks, BreakingChange{
				UploaderRole: in.UploaderRole,
				Counterpart:  counterpart,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonMissingInProvider,
			})

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, BreakingChange{
				UploaderRole: in.UploaderRole,
				Counterpart:  counterpart,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonTypeMismatch,
				ConsumerType: consumerProperty.Type,
				ProviderType: providerProperty.Type,
			})

			continue
		}
		if !consumerProperty.Optional && providerProperty.Optional {
			breaks = append(breaks, BreakingChange{
				UploaderRole: in.UploaderRole,
				Counterpart:  counterpart,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonOptionalInProviderRequiredInConsumer,
			})
		}
	}
	return breaks
}

func compareRequest(in CompareInput) []BreakingChange {
	brokenResource := NewBrokenResource(in.Consumer)
	counterpart := counterpartFor(in)

	var breaks []BreakingChange
	for providerPropertyPath, providerProperty := range in.Provider.Properties {
		consumerProperty, propertyExists := in.Consumer.Properties[providerPropertyPath]
		if !propertyExists {
			if !providerProperty.Optional {
				breaks = append(breaks, BreakingChange{
					UploaderRole: in.UploaderRole,
					Counterpart:  counterpart,
					Resource:     brokenResource,
					Property:     providerPropertyPath,
					Reason:       ReasonMissingInConsumer,
				})
			}

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, BreakingChange{
				UploaderRole: in.UploaderRole,
				Counterpart:  counterpart,
				Resource:     brokenResource,
				Property:     providerPropertyPath,
				Reason:       ReasonTypeMismatch,
				ConsumerType: consumerProperty.Type,
				ProviderType: providerProperty.Type,
			})
			continue
		}

		if !providerProperty.Optional && consumerProperty.Optional {
			breaks = append(breaks, BreakingChange{
				UploaderRole: in.UploaderRole,
				Counterpart:  counterpart,
				Resource:     brokenResource,
				Property:     providerPropertyPath,
				Reason:       ReasonOptionalInConsumerRequiredInProvider,
			})
		}
	}

	return breaks
}
