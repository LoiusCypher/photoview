package resolvers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var rec1 map[string]time.Duration
var rec2 map[string]time.Duration

func recordTiming( caller string, time1 time.Duration, time2 time.Duration) {
	if rec1 == nil {
		rec1 = make(map[string]time.Duration)
		rec2 = make(map[string]time.Duration)
	}
	rec1[caller] = rec1[caller] + time1
	rec2[caller] = rec2[caller] + time2
}

func logTiming() {
	if rec1 == nil {
		log.Println("TRACE rec1 is uninitialized")
		return
	}
	for _caller, _ := range rec1 {
		log.Println("TRACE Diff", _caller, rec1[_caller]-rec2[_caller],
				"Sum1", rec1[_caller],
				"Sum2", rec2[_caller],
				"Diff %", 100*float64(rec1[_caller]-rec2[_caller])/float64(rec1[_caller]))
	}
}

func recordChange( path string, incMirror bool, incRotate bool, remove bool) {
	log.Printf("recordChange Path %s\n", path)

	filePath := "media-cache/rotate.txt"
 	// If the file doesn't exist, create it, or append to the file
	fr, err := os.OpenFile( filePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
        	log.Fatal(err)
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(fr)

	found := false
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		str := text[3:]
		//log.Printf("recordChange str: %s\n", str)
		var rot, mir int
		args, err := fmt.Sscanf(text, "%1d%1d ", &mir, &rot)
		//log.Printf("recordChange err: %t args: %d\n", err == nil, args)
		if err == nil && args == 2 && path == str {
			if !found {
				if incMirror && mir != 9 {
					mir = (mir + 1) % 2
				}
				if incRotate && rot != 9 {
					rot = (rot + 1) % 4
				}
				if remove {
					mir = 9
					rot = 9
				}
				log.Printf("recordChange mir: %d rot: %d\n", mir, rot)
				text = fmt.Sprintf("%1d%1d %s", mir, rot, str)
				lines = append(lines, text + "\n")
				found = true
			}
		} else {
			lines = append(lines, text + "\n")
		}
	}
	fr.Close()
	if found {
		fw, _ := os.OpenFile( filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
  		defer fw.Close()
		for _, line := range lines {
			if _, err := fw.Write([]byte(line)); err != nil {
				panic("a problem")
			}
		}
	} else {
		fa, _ := os.OpenFile( filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  		defer fa.Close()
		mir := 0
		if incMirror {
			mir = 1
		}
		rot := 0
		if incRotate {
			rot = 1
		}
		if remove {
			mir = 9
			rot = 9
		}
  		fa.WriteString( fmt.Sprintf("%1d%1d %s\n", mir, rot, path))
	}
}
