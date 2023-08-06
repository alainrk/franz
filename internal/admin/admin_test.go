package admin

import (
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type MockConsumer struct {
	Metadata *kafka.Metadata
	Err      error
}

func (m *MockConsumer) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	return m.Metadata, m.Err
}

func TestDefaultKafkaAdminer_GetTopics(t *testing.T) {
	testCases := []struct {
		name             string
		includeInternals bool
		metadata         *kafka.Metadata
		expectedTopics   []string
		expectedError    error
	}{
		{
			name:             "Include internals",
			includeInternals: true,
			metadata: &kafka.Metadata{
				Topics: map[string]kafka.TopicMetadata{
					"topic1":  {},
					"_topic2": {},
					"topic3":  {},
				},
			},
			expectedTopics: []string{"topic1", "_topic2", "topic3"},
			expectedError:  nil,
		},
		{
			name:             "Exclude internals",
			includeInternals: false,
			metadata: &kafka.Metadata{
				Topics: map[string]kafka.TopicMetadata{
					"topic1":  {},
					"_topic2": {},
					"topic3":  {},
				},
			},
			expectedTopics: []string{"topic1", "topic3"},
			expectedError:  nil,
		},
		{
			name:             "Error getting metadata",
			includeInternals: true,
			metadata:         nil,
			expectedTopics:   nil,
			expectedError:    fmt.Errorf("simulated error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConsumer := &MockConsumer{
				Metadata: tc.metadata,
				Err:      tc.expectedError,
			}

			adminer := &DefaultKafkaAdminer{Client: mockConsumer}

			topics, err := adminer.GetTopics(tc.includeInternals)

			// Improve checking error (ATM upstream raises wrapped err)
			if tc.expectedError != nil && err == nil {
				t.Errorf("Expected error, but got: %v", err)
			}

			if !slicesHaveSameElements(topics, tc.expectedTopics) {
				t.Errorf("Expected topics: %v, but got: %v", tc.expectedTopics, topics)
			}
		})
	}
}

func slicesHaveSameElements(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for _, val := range slice1 {
		found := false
		for _, val2 := range slice2 {
			if val == val2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
