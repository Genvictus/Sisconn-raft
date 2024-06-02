package raft

// util
func (r *RaftNode) GetVotedFor() string {
	return r.votedFor
}

func (r *RaftNode) SetCurrentTerm(term uint64) {
	r.currentTerm = term
}
