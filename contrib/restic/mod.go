// Extensions to the types vendored from restic.

package restic

func (s *Snapshot) GetId() *ID {
	return s.id
}
