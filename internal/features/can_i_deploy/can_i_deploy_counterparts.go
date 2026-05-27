package can_i_deploy

import (
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
)

type counterpartParticipant struct {
	id      int64
	name    string
	version string
}

func counterpartFromDeployment(deployment model.Deployment) counterpartParticipant {
	return counterpartParticipant{
		id:      deployment.Participant.ID,
		name:    deployment.Participant.Name,
		version: deployment.Version,
	}
}

func counterpartFromConsumer(consumer repository.CurrentConsumerInEnv) counterpartParticipant {
	return counterpartParticipant{
		id:      consumer.ParticipantID,
		name:    consumer.ParticipantName,
		version: consumer.Version,
	}
}

type counterpartSet struct {
	pairs       map[int64]counterpartParticipant
	strictFalse map[int64]struct{}
	ghostSeen   bool
}

func newCounterpartSet() *counterpartSet {
	return &counterpartSet{
		pairs:       map[int64]counterpartParticipant{},
		strictFalse: map[int64]struct{}{},
	}
}

func (s *counterpartSet) addPair(pair counterpartParticipant) {
	s.pairs[pair.id] = pair
}

func (s *counterpartSet) markStrictFalse(counterpartID int64) {
	s.strictFalse[counterpartID] = struct{}{}
}

func (s *counterpartSet) markGhost() {
	s.ghostSeen = true
}

func (s *counterpartSet) hasGhost() bool {
	return s.ghostSeen
}

func (s *counterpartSet) pairCount() int {
	return len(s.pairs)
}

func (s *counterpartSet) dropStrictFalsesAlreadyPaired() {
	for id := range s.pairs {
		delete(s.strictFalse, id)
	}
}

func (s *counterpartSet) eachPair(visit func(pair counterpartParticipant)) {
	for _, pair := range s.pairs {
		visit(pair)
	}
}

func (s *counterpartSet) eachStrictFalse(visit func(counterpartID int64)) {
	for id := range s.strictFalse {
		visit(id)
	}
}
