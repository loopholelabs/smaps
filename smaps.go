// SPDX-License-Identifier: Apache-2.0

package smaps

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Identifier string

const (
	UnknownIdentifier Identifier = "unknown"
)

type Smap struct {
	Identifier     Identifier
	Address        string
	Permissions    string
	Offset         string
	Dev            string
	Inode          string
	IsPath         bool
	Size           int64
	KernelPageSize int64
	MMUPageSize    int64
	Rss            int64
	Pss            int64
	PssDirty       int64
	SharedClean    int64
	SharedDirty    int64
	PrivateClean   int64
	PrivateDirty   int64
	Referenced     int64
	Swap           int64
}

type Smaps map[Identifier][]Smap

func parseSize(sizeStr string) int64 {
	parts := strings.Fields(sizeStr)
	if len(parts) != 2 {
		return 0
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	unit := strings.ToLower(parts[1])
	switch unit {
	case "kb":
		return int64(value * 1024)
	case "mb":
		return int64(value * 1024 * 1024)
	case "gb":
		return int64(value * 1024 * 1024 * 1024)
	default:
		return 0
	}
}

func parse(reader io.Reader, smaps Smaps) {
	scanner := bufio.NewScanner(reader)
	var line string
	var fields []string
	var prefixIndex int
	for scanner.Scan() {
		line = scanner.Text()
		fields = strings.Fields(line)
		if len(fields) >= 5 {
			smap := Smap{
				Address:     fields[0],
				Permissions: fields[1],
				Offset:      fields[2],
				Dev:         fields[3],
				Inode:       fields[4],
			}
			if len(fields) > 5 {
				identifier := strings.Join(fields[5:], " ")
				smap.IsPath = !strings.HasPrefix(identifier, "[") || !strings.HasSuffix(identifier, "]")
				if smap.IsPath {
					smap.Identifier = Identifier(identifier)
				} else {
					smap.Identifier = Identifier(identifier[1 : len(identifier)-1])
				}
			} else {
				smap.IsPath = false
				smap.Identifier = UnknownIdentifier
			}

			for scanner.Scan() {
				line = scanner.Text()
				if strings.HasPrefix(line, "VmFlags:") {
					break
				}
				prefixIndex = strings.Index(line, ":")
				switch line[:prefixIndex] {
				case "Size":
					smap.Size = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "KernelPageSize":
					smap.KernelPageSize = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "MMUPageSize":
					smap.MMUPageSize = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Rss":
					smap.Rss = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Pss":
					smap.Pss = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Pss_Dirty":
					smap.PssDirty = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Shared_Clean":
					smap.SharedClean = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Shared_Dirty":
					smap.SharedDirty = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Private_Clean":
					smap.PrivateClean = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Private_Dirty":
					smap.PrivateDirty = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Referenced":
					smap.Referenced = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				case "Swap":
					smap.Swap = parseSize(strings.TrimSpace(line[prefixIndex+1:]))
				}
			}
			if smapArray, ok := smaps[smap.Identifier]; ok {
				smaps[smap.Identifier] = append(smapArray, smap)
			} else {
				smaps[smap.Identifier] = []Smap{smap}
			}
		}
	}
}

func Parse(file *os.File) (Smaps, error) {
	smaps := make(Smaps)
	_, err := file.Seek(0, 0)
	if err == nil {
		parse(file, smaps)
	}
	return smaps, err
}
