package controller

type Sender struct {
	jobs Jobs
}

func NewSender(jobs Jobs) *Sender {
	return &Sender{
		jobs: jobs,
	}
}

func (s *Sender) Send() error {
	return nil
}
