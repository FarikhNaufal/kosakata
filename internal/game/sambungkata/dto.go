package sambungkata

type TodayWordDTO struct {
	ID    string `json:"id"`
	Start string `json:"start"`
}

type WrongWordDTO struct {
	Clue     []string       `json:"clue"`
	Length   int            `json:"length"`
	PrevWord *string        `json:"prev_word"`
	Details  []WordPosition `json:"details"`
	Progress float64        `json:"progress"`
}
