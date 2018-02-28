package interfaces

type IElections interface {
}

type IElectionAdapter interface {
	Execute(IMsg) IMsg
	GetDBHeight() int
	GetMinute() int
	GetElecting() int
}

type IElectionMsg interface {
	ElectionProcess(IState, IElections)
	String() string
}

type IElectionsFactory interface {
	// Messages
	NewAddLeaderInternal(Name string, dbheight uint32, serverID IHash) IMsg
	NewAddAuditInternal(name string, dbheight uint32, serverID IHash) IMsg
	NewRemoveLeaderInternal(name string, dbheight uint32, serverID IHash) IMsg
	NewRemoveAuditInternal(name string, dbheight uint32, serverID IHash) IMsg
	NewEomSigInternal(name string, dbheight uint32, minute uint32, height uint32, serverID IHash) IMsg
	NewDBSigSigInternal(name string, dbheight uint32, minute uint32, height uint32, serverID IHash) IMsg

	//
	NewElectionAdapter(el IElections) IElectionAdapter
}