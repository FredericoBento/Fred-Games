package services

type StatusChecker interface {
	IsActive() bool
	IsInactive() bool
	SetActive()
	SetInactive()
	HasStartedOnce() bool
}

type Status struct {
	value          string
	hasStartedOnce bool
}

const (
	statusInactive = "inactive"
	statusActive   = "active"
)

func NewStatus() *Status {
	return &Status{
		value:          statusInactive,
		hasStartedOnce: false,
	}
}

func (s *Status) IsActive() bool {
	return s.value == statusActive
}

func (s *Status) IsInactive() bool {
	return s.value == statusInactive
}

func (s *Status) SetActive() {
	s.value = statusActive
	s.hasStartedOnce = true
}

func (s *Status) SetInactive() {
	s.value = statusInactive
}

func (s *Status) HasStartedOnce() bool {
	return s.hasStartedOnce
}
