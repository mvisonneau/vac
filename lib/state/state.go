package state

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/mvisonneau/vac/lib/client"
	"github.com/pkg/errors"
)

// State ..
type State struct {
	Current struct {
		Engine string `json:"engine,omitempty"`
		Role   string `json:"role,omitempty"`
	} `json:"current,omitempty"`
	AWSCredentials map[string]map[string]*client.AWSCredentials `json:"creds,omitempty"`
}

// Read ..
func Read(filePath string) (*State, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return &State{}, nil
	}

	s := &State{}
	byteValue, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &State{}, errors.Wrapf(err, "opening state file '%s'", filePath)
	}
	json.Unmarshal(byteValue, s)
	return s, nil
}

// Write ..
func Write(c *State, filePath string) error {
	marshalledContent, err := json.MarshalIndent(*c, "", " ")
	if err != nil {
		return errors.Wrap(err, "marshalling state into json")
	}
	return ioutil.WriteFile(filePath, marshalledContent, 0644)
}

// SetCurrentEngine ..
func (s *State) SetCurrentEngine(engine string) {
	s.Current.Engine = engine
}

// SetCurrentRole ..
func (s *State) SetCurrentRole(role string) {
	s.Current.Role = role
}

// SetCurrentAWSCredentials ..
func (s *State) SetCurrentAWSCredentials(creds *client.AWSCredentials) {
	s.SetAWSCredentials(s.Current.Engine, s.Current.Role, creds)
}

// SetAWSCredentials ..
func (s *State) SetAWSCredentials(engine, role string, creds *client.AWSCredentials) {
	if s.AWSCredentials == nil {
		s.AWSCredentials = make(map[string]map[string]*client.AWSCredentials)
	}
	if _, ok := s.AWSCredentials[engine]; !ok {
		s.AWSCredentials[engine] = make(map[string]*client.AWSCredentials)
	}
	s.AWSCredentials[engine][role] = creds
}

// GetCurrentAWSCredentials ..
func (s *State) GetCurrentAWSCredentials() *client.AWSCredentials {
	return s.GetAWSCredentials(s.Current.Engine, s.Current.Role)
}

// GetAWSCredentials ..
func (s *State) GetAWSCredentials(engine, role string) *client.AWSCredentials {
	if s.AWSCredentials != nil {
		if _, ok := s.AWSCredentials[engine]; ok {
			if creds, ok := s.AWSCredentials[engine][role]; ok {
				return creds
			}
		}
	}
	return nil
}
