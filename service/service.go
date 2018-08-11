package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrecronje/hashgraph/node"
	"github.com/sirupsen/logrus"
)

type Service struct {
	bindAddress string
	node        *node.Node
	logger      *logrus.Logger
}

func NewService(bindAddress string, node *node.Node, logger *logrus.Logger) *Service {
	service := Service{
		bindAddress: bindAddress,
		node:        node,
		logger:      logger,
	}

	return &service
}

func (s *Service) Serve() {
	s.logger.WithField("bind_address", s.bindAddress).Debug("Service serving")
	http.HandleFunc("/stats", s.GetStats)
	http.HandleFunc("/participants/", s.GetParticipants)
	http.HandleFunc("/event/", s.GetEvent)
	http.HandleFunc("/lasteventfrom/", s.GetLastEventFrom)
	http.HandleFunc("/events/", s.GetKnownEvents)
	http.HandleFunc("/consensusevents/", s.GetConsensusEvents)
	http.HandleFunc("/round/", s.GetRound)
	http.HandleFunc("/lastround/", s.GetLastRound)
	http.HandleFunc("/roundwitnesses/", s.GetRoundWitnesses)
	http.HandleFunc("/roundevents/", s.GetRoundEvents)
	http.HandleFunc("/root/", s.GetRoot)
	http.HandleFunc("/block/", s.GetBlock)
	err := http.ListenAndServe(s.bindAddress, nil)
	if err != nil {
		s.logger.WithField("error", err).Error("Service failed")
	}
}

func (s *Service) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := s.node.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Service) GetParticipants(w http.ResponseWriter, r *http.Request) {
	participants, err := s.node.GetParticipants()
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing participants parameter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

func (s *Service) GetEvent(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/event/"):]
	event, err := s.node.GetEvent(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %d", event)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}


func (s *Service) GetLastEventFrom(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/lasteventfrom/"):]
	event, _, err := s.node.GetLastEventFrom(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %d", event)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Service) GetKnownEvents(w http.ResponseWriter, r *http.Request) {
	knownEvents := s.node.GetKnownEvents()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(knownEvents)
}

func (s *Service) GetConsensusEvents(w http.ResponseWriter, r *http.Request) {
	consensusEvents := s.node.GetConsensusEvents()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consensusEvents)
}

func (s *Service) GetRound(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/round/"):]
	roundIndex, err := strconv.Atoi(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	round, err := s.node.GetRound(roundIndex)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving round %d", roundIndex)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(round)
}

func (s *Service) GetLastRound(w http.ResponseWriter, r *http.Request) {
	lastRound := s.node.GetLastRound()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastRound)
}


func (s *Service) GetRoundWitnesses(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundwitnesses/"):]
	roundWitnessesIndex, err := strconv.Atoi(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundWitnessesIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundWitnesses := s.node.GetRoundWitnesses(roundWitnessesIndex)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roundWitnesses)
}

func (s *Service) GetRoundEvents(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundevents/"):]
	roundEventsIndex, err := strconv.Atoi(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundEventsIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundEvent := s.node.GetRoundEvents(roundEventsIndex)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roundEvent)
}

func (s *Service) GetRoot(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/root/"):]
	root, err := s.node.GetRoot(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving root %d", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(root)
}

func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/block/"):]
	blockIndex, err := strconv.Atoi(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing block_index parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	block, err := s.node.GetBlock(blockIndex)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving block %d", blockIndex)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}
