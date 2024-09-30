package dsu

import "fmt"

type LFTPHeader struct {
	Version       string
	ContentLength int
	TotalLength   int
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
		fmt.Sprint(l.ContentLength),
		fmt.Sprint(l.TotalLength),
		fmt.Sprint(l.StartOffset),
		fmt.Sprint(l.EndOffset),
		fmt.Sprint(l.ContentID),
	)
	fmt.Println(*l)
	return headerString
}
