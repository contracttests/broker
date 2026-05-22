package upload_contract

import (
	"context"
	"errors"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
)

type CompatibilityChecker struct {
	repo *repository.ContractRepository
}

func NewCompatibilityChecker(repo *repository.ContractRepository) *CompatibilityChecker {
	return &CompatibilityChecker{repo: repo}
}

func (c *CompatibilityChecker) Run(
	ctx context.Context,
	uploaded *model.Contract,
) *model.CompatibilityReport {
	report := &model.CompatibilityReport{}

	for _, resource := range uploaded.Resources {
		switch resource.Direction {
		case model.Consumes:
			c.checkConsumer(ctx, resource, report)
		case model.Provides:
			c.checkProvider(ctx, resource, report)
		}
	}

	return report
}

func (c *CompatibilityChecker) checkConsumer(
	ctx context.Context,
	consumer model.Resource,
	report *model.CompatibilityReport,
) {
	if consumer.Provider == "" {
		report.Append(model.NewMissingProviderBreak(
			consumer,
			model.ReasonProviderNotSpecified,
		))

		return
	}

	provider, err := c.repo.LoadProviderResource(ctx, consumer)

	if errors.Is(err, repository.ErrProviderResourceNotFound) {
		report.Append(model.NewMissingProviderBreak(
			consumer,
			model.ReasonProviderResourceNotFound,
		))

		return
	}

	for _, breakingChange := range model.Compare(consumer, provider) {
		report.Append(breakingChange)
	}
}

func (c *CompatibilityChecker) checkProvider(
	ctx context.Context,
	provider model.Resource,
	report *model.CompatibilityReport,
) {
	consumers := c.repo.FindConsumersOfProvider(ctx, provider)

	for _, consumer := range consumers {
		for _, breakingChange := range model.Compare(consumer, provider) {
			report.Append(breakingChange)
		}
	}
}
