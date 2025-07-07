package collab

// EditMsg carries the details of a spreadsheet cell edit
// made by a client, for broadcast to other collaborators.
type EditMsg struct {
	Row  int    `json:"row"`
	Col  int    `json:"col"`
	Data string `json:"data"`
}

// BroadCastMsg represents a message to be broadcasted to sheet collaborators.
// It has a SheetID field to identify the sheet whose clients should receive the
// message and an Edit field, which is an EditMsg, holding the edit that will be
// broadcasted.
type BroadCastMsg struct {
	SheetID string
	Edit    EditMsg
}
