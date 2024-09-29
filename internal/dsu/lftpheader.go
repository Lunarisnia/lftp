package dsu

import "fmt"

type LFTPHeader struct {
	Version       string
	ContentLength uint
	TotalLength   uint
	StartOffset   int
	EndOffset     int
	ContentID     string
}

func (l *LFTPHeader) ConstructString() string {
	headerString := "LFTP"
	headerString = fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%s|",
		headerString,
		l.Version,
		string(l.ContentLength),
		string(l.TotalLength),
		string(l.StartOffset),
		string(l.EndOffset),
		string(l.ContentID),
	)
	return headerString
}
