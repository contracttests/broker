package compatibility_checker

import (
	"context"
	"errors"
	"fmt"

	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
)

type BrokenResource struct {
	Kind       model.ResourceKind
	Provider   string
	Endpoint   string
	Method     string
	StatusCode string
}

func NewBrokenResource(resource model.Resource) BrokenResource {
	return BrokenResource{
		Kind:       resource.Kind,
		Provider:   resource.Provider,
		Endpoint:   resource.Endpoint,
		Method:     resource.Method,
		StatusCode: resource.StatusCode,
	}
}

type BreakingReason string

const (
	ReasonProviderNotFound                     BreakingReason = "provider_not_found"
	ReasonProviderResourceNotFound             BreakingReason = "provider_resource_not_found"
	ReasonMissingInProvider                    BreakingReason = "missing_in_provider"
	ReasonMissingInConsumer                    BreakingReason = "missing_in_consumer"
	ReasonTypeMismatch                         BreakingReason = "type_mismatch"
	ReasonOptionalInProviderRequiredInConsumer BreakingReason = "optional_in_provider_required_in_consumer"
	ReasonOptionalInConsumerRequiredInProvider BreakingReason = "optional_in_consumer_required_in_provider"
)

type BreakingChange struct {
	LeftResource  *model.Resource `json:"left_resource"`
	RightResource *model.Resource `json:"right_resource"`
	Reason        BreakingReason  `json:"reason"`
	Property      string          `json:"property"`
	HumanReadable string          `json:"human_readable"`
}

func NewBreakingChange(
	leftResource *model.Resource,
	rightResource *model.Resource,
	reason BreakingReason,
	property string,
) BreakingChange {
	breakingChange := BreakingChange{
		LeftResource:  leftResource,
		RightResource: rightResource,
		Reason:        reason,
		Property:      property,
	}

	breakingChange.humanReadable()

	return breakingChange
}

func (b *BreakingChange) humanReadable() {
	switch b.Reason {
	case ReasonProviderNotFound:
		b.HumanReadable = fmt.Sprintf(
			"Provider %s not found",
			b.LeftResource.ParticipantName(),
		)

	case ReasonProviderResourceNotFound:
		b.HumanReadable = fmt.Sprintf(
			"Provider resource %s not found",
			b.LeftResource.HumanReadable(),
		)

	case ReasonMissingInProvider:
		b.HumanReadable = fmt.Sprintf(
			"Property %s is missing in provider %s",
			b.Property,
			b.RightResource.ParticipantName(),
		)

	case ReasonMissingInConsumer:
		b.HumanReadable = fmt.Sprintf(
			"Property %s is missing in consumer %s",
			b.Property,
			b.LeftResource.HumanReadable(),
		)

	case ReasonTypeMismatch:
		consumerType := b.LeftResource.Properties[b.Property].Type
		providerType := b.RightResource.Properties[b.Property].Type

		b.HumanReadable = fmt.Sprintf(
			"Property %s type mismatch, provider %s expects %s but consumer %s expects %s",
			b.Property,
			b.RightResource.ParticipantName(),
			providerType,
			b.LeftResource.ParticipantName(),
			consumerType,
		)

	case ReasonOptionalInProviderRequiredInConsumer:
		b.HumanReadable = fmt.Sprintf(
			"Property %s is optional in provider %s but required in consumer %s",
			b.Property,
			b.RightResource.ParticipantName(),
			b.LeftResource.ParticipantName(),
		)

	case ReasonOptionalInConsumerRequiredInProvider:
		b.HumanReadable = fmt.Sprintf(
			"Property %s is optional in consumer %s but required in provider %s",
			b.Property,
			b.LeftResource.ParticipantName(),
			b.RightResource.ParticipantName(),
		)

	default:
		b.HumanReadable = "Unknown reason"
	}
}

type CompatibilityChecker struct {
	repository *repository.ContractRepository
}

type CompatibilityReport struct {
	Breaks []BreakingChange `json:"breaks"`
}

func (r *CompatibilityReport) Append(b BreakingChange) {
	r.Breaks = append(r.Breaks, b)
}

func NewCompatibilityChecker(repository *repository.ContractRepository) *CompatibilityChecker {
	return &CompatibilityChecker{
		repository: repository,
	}
}

func (c *CompatibilityChecker) Check(
	ctx context.Context,
	uploaded *model.Contract,
	environment *model.Environment,
) *CompatibilityReport {
	report := &CompatibilityReport{}

	for _, resource := range uploaded.Resources {
		switch resource.Direction {
		case model.Consumes:
			c.checkConsumer(ctx, resource, environment, report)
		case model.Provides:
			c.checkProvider(ctx, resource, environment, report)
		}
	}

	return report
}

func (c *CompatibilityChecker) checkConsumer(
	ctx context.Context,
	consumer model.Resource,
	environment *model.Environment,
	report *CompatibilityReport,
) {
	provider, err := c.repository.LoadProviderResourceOfConsumerAndEnvironment(ctx, consumer, environment)

	if errors.Is(err, repository.ErrProviderResourceNotFound) {
		report.Append(BreakingChange{
			LeftResource:  &consumer,
			RightResource: nil,
			Reason:        ReasonProviderResourceNotFound,
		})

		return
	}

	for _, breakingChange := range checkResources(&consumer, &provider) {
		report.Append(breakingChange)
	}
}

func (c *CompatibilityChecker) checkProvider(
	ctx context.Context,
	provider model.Resource,
	environment *model.Environment,
	report *CompatibilityReport,
) {
	consumers := c.repository.FindConsumersOfProviderAndEnvironment(ctx, provider, environment)

	for _, consumer := range consumers {
		for _, breakingChange := range checkResources(&consumer, &provider) {
			report.Append(breakingChange)
		}
	}
}

func checkResources(leftResource *model.Resource, rightResource *model.Resource) []BreakingChange {
	switch leftResource.Kind {
	case model.RestRequest:
		return checkRequestResource(leftResource, rightResource)
	default:
		return checkResponseResource(leftResource, rightResource)
	}
}

func checkResponseResource(leftResource *model.Resource, rightResource *model.Resource) []BreakingChange {
	var breaks []BreakingChange
	for consumerPropertyPath, consumerProperty := range leftResource.Properties {
		providerProperty, propertyExists := rightResource.Properties[consumerPropertyPath]

		if !propertyExists {
			breaks = append(breaks, NewBreakingChange(
				leftResource,
				rightResource,
				ReasonMissingInProvider,
				consumerPropertyPath,
			))

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, NewBreakingChange(
				leftResource,
				rightResource,
				ReasonTypeMismatch,
				consumerPropertyPath,
			))

			continue
		}

		if !consumerProperty.Optional && providerProperty.Optional {
			breaks = append(breaks, NewBreakingChange(
				leftResource,
				rightResource,
				ReasonOptionalInProviderRequiredInConsumer,
				consumerPropertyPath,
			))
		}
	}

	return breaks
}

func checkRequestResource(leftResource *model.Resource, rightResource *model.Resource) []BreakingChange {
	var breaks []BreakingChange
	for providerPropertyPath, providerProperty := range rightResource.Properties {
		consumerProperty, propertyExists := leftResource.Properties[providerPropertyPath]
		if !propertyExists {
			if !providerProperty.Optional {
				breaks = append(breaks, NewBreakingChange(
					leftResource,
					rightResource,
					ReasonMissingInConsumer,
					providerPropertyPath,
				))
			}

			continue
		}

		if consumerProperty.Type != providerProperty.Type {
			breaks = append(breaks, NewBreakingChange(
				leftResource,
				rightResource,
				ReasonTypeMismatch,
				providerPropertyPath,
			))

			continue
		}

		if !providerProperty.Optional && consumerProperty.Optional {
			breaks = append(breaks, NewBreakingChange(
				leftResource,
				rightResource,
				ReasonOptionalInConsumerRequiredInProvider,
				providerPropertyPath,
			))
		}
	}

	return breaks
}
