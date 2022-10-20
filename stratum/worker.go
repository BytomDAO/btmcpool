package superstratum

import (
	"fmt"
	"strings"
)

type Worker struct {
	fullName string // fullName = account + "." + name. unique id for worker
	account  string // account name
	name     string // worker name
	version  string // client type & version
}

func NewWorker(loginWorkerPair, version string) (*Worker, error) {
	s := strings.SplitN(loginWorkerPair, ".", 2)
	if len(s) == 0 {
		return nil, fmt.Errorf("invalid name")
	} else if len(s) == 1 {
		// default worker
		s = append(s, "0")
	}

	return &Worker{
		fullName: loginWorkerPair,
		account:  s[0],
		name:     s[1],
		version:  version,
	}, nil
}

func (w *Worker) GetId() string {
	return w.fullName
}

func (w *Worker) GetWorker() (string, string) {
	return w.account, w.name
}

func (w *Worker) GetFullName() string {
	return w.fullName
}
