package restic

func (s *Snapshot) GetId() *ID {
	return s.id
}
