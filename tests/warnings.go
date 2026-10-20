package tests

import (
	"fmt"
	"strings"
	"sync"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type warningRecorder struct {
	mutex    sync.Mutex
	messages []string
}

func (recorder *warningRecorder) HandleWarningHeader(code int, agent, message string) {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	recorder.messages = append(recorder.messages, message)
}

func (recorder *warningRecorder) check(expected []string) error {
	if expected == nil {
		return nil
	}
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	for _, substring := range expected {
		found := false
		for _, message := range recorder.messages {
			if strings.Contains(message, substring) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected warning containing %q, got %q", substring, recorder.messages)
		}
	}
	for _, message := range recorder.messages {
		matched := false
		for _, substring := range expected {
			if strings.Contains(message, substring) {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("unexpected warning %q, expected substrings %q", message, expected)
		}
	}
	return nil
}

func clientForWarnings(configuration *rest.Config, baseClient client.Client, expected []string) (client.Client, *warningRecorder, error) {
	if expected == nil {
		return baseClient, nil, nil
	}
	recorder := &warningRecorder{}
	configuration = rest.CopyConfig(configuration)
	configuration.WarningHandler = recorder
	configuration.WarningHandlerWithContext = nil
	requestClient, err := client.New(configuration, client.Options{
		Scheme: baseClient.Scheme(),
		Mapper: baseClient.RESTMapper(),
	})
	return requestClient, recorder, err
}
