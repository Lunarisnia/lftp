package dsu

import "bufio"

type ServerMemo map[string]*bufio.Writer

type ClientMemo map[string]*bufio.Reader
