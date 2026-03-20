package index

type Posting struct {
	DocID     int
	Frequency int
	Positions []int
}

type PostingList []Posting
