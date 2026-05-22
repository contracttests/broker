package model

type CompatibilityReport struct {
	breaks []BreakingChange
}

func (r *CompatibilityReport) HasBreaks() bool {
	return len(r.breaks) > 0
}

func (r *CompatibilityReport) Append(b BreakingChange) {
	r.breaks = append(r.breaks, b)
}

func NewMissingProviderBreak(consumer Resource, reason BreakingReason) BreakingChange {
	return BreakingChange{
		ContractInfo: *consumer.ContractInfo,
		Resource:     NewBrokenResource(consumer),
		Reason:       reason,
	}
}

func Compare(consumer, provider Resource) []BreakingChange {
	switch consumer.Kind {
	case RestRequest:
		return compareRequest(consumer, provider)
	default:
		return compareResponse(consumer, provider)
	}
}

func compareResponse(consumer, provider Resource) []BreakingChange {
	brokenResource := NewBrokenResource(consumer)

	// Initialize the contract info to an empty struct.
	// When the consumer is uploaded, the contract info is not set.
	consumerContractInfo := ContractInfo{}
	if consumer.ContractInfo != nil {
		consumerContractInfo = *consumer.ContractInfo
	}

	var breaks []BreakingChange
	for consumerPropertyPath, consumerProperty := range consumer.Properties {
		providerProperty, propertyExists := provider.Properties[consumerPropertyPath]

		if !propertyExists {
			breaks = append(breaks, BreakingChange{
				ContractInfo: consumerContractInfo,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonMissingInProvider,
			})

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, BreakingChange{
				ContractInfo: consumerContractInfo,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonTypeMismatch,
				ExpectedType: consumerProperty.Type,
				ActualType:   providerProperty.Type,
			})

			continue
		}
		if !consumerProperty.Optional && providerProperty.Optional {
			breaks = append(breaks, BreakingChange{
				ContractInfo: consumerContractInfo,
				Resource:     brokenResource,
				Property:     consumerPropertyPath,
				Reason:       ReasonOptionalInProviderRequiredInConsumer,
			})
		}
	}
	return breaks
}

func compareRequest(consumer, provider Resource) []BreakingChange {
	brokenResource := NewBrokenResource(consumer)

	// Initialize the contract info to an empty struct.
	// When the consumer is uploaded, the contract info is not set.
	consumerContractInfo := ContractInfo{}
	if consumer.ContractInfo != nil {
		consumerContractInfo = *consumer.ContractInfo
	}

	var breaks []BreakingChange
	for providerPropertyPath, providerProperty := range provider.Properties {
		consumerProperty, propertyExists := consumer.Properties[providerPropertyPath]
		if !propertyExists {
			if !providerProperty.Optional {
				breaks = append(breaks, BreakingChange{
					ContractInfo: consumerContractInfo,
					Resource:     brokenResource,
					Property:     providerPropertyPath,
					Reason:       ReasonMissingInConsumer,
				})
			}

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, BreakingChange{
				ContractInfo: consumerContractInfo,
				Resource:     brokenResource,
				Property:     providerPropertyPath,
				Reason:       ReasonTypeMismatch,
				ExpectedType: consumerProperty.Type,
				ActualType:   providerProperty.Type,
			})
			continue
		}

		if !providerProperty.Optional && consumerProperty.Optional {
			breaks = append(breaks, BreakingChange{
				ContractInfo: consumerContractInfo,
				Resource:     brokenResource,
				Property:     providerPropertyPath,
				Reason:       ReasonOptionalInConsumerRequiredInProvider,
			})
		}
	}

	return breaks
}
