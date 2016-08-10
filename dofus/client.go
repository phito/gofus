package dofus

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

// A Client represents a Dofus client instance
type Client struct {
	process              *os.Process
	fingerprint, payload []byte
}

// RunClient runs and patches a new Dofus client
// executable
func RunClient(executable, fingerprint, payload string) (*Client, error) {
	var client Client
	var err error

	client.loadPayload(fingerprint, payload)

	client.process, err = os.StartProcess(executable, []string{"--lang=fr"}, &os.ProcAttr{})

	if err != nil {
		return nil, err
	}

	if err = client.patch(); err != nil {
		return nil, err
	}

	return &client, err
}

func (client *Client) loadPayload(fingerprint, payload string) (err error) {
	client.fingerprint, err = ioutil.ReadFile(fingerprint)
	client.payload, err = ioutil.ReadFile(payload)
	return
}

func (client *Client) patch() (err error) {
	time.Sleep(time.Second)
	pid := client.process.Pid

	// attaching to the tracee
	if err = syscall.PtraceAttach(pid); err != nil {
		return
	}

	// map the process' memory
	mmap, err := client.mapMemory()
	if err != nil {
		return
	}

	// open the process' mem file
	file, err := os.OpenFile("/proc/"+strconv.Itoa(pid)+"/mem", os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	// defer it to insure that it will be closed no matter what
	defer file.Close()

	for _, region := range mmap {
		size := region.end - region.start

		if size > 5000000 {
			buffer := make([]byte, region.end-region.start)

			_, e := file.ReadAt(buffer, region.start)
			if e != nil {
				println("unable to read region:", e.Error())
			} else {

				offset, e := find(buffer, client.fingerprint)
				if e == nil {
					position := region.start + int64(offset)
					_, err = file.WriteAt(client.payload, position)
					if err != nil {
						return
					}
				}
			}
		}
	}

	// detaching the tracee
	if err = syscall.PtraceDetach(pid); err != nil {
		return
	}
	return
}

func find(data []byte, fingerprint []byte) (int, error) {
	j := 0
	for i := 0; i < len(data); i++ {
		if data[i] == fingerprint[j] {
			j++
		} else {
			j = 0
		}
		if j == len(fingerprint) {
			return i - (j - 1), nil
		}
	}
	return 0, errors.New("fingerprint not found")
}

type region struct {
	start, end int64
}

func (client *Client) mapMemory() (regions []region, err error) {
	pid := client.process.Pid

	// compile the regex used to parse the maps file
	r := regexp.MustCompile("([0-9A-Fa-f]+)-([0-9A-Fa-f]+) ([-r])")

	file, err := os.Open("/proc/" + strconv.Itoa(pid) + "/maps")
	if err != nil {
		return
	}

	// defer it to insure that it will be closed no matter what
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			matches := r.FindAllStringSubmatch(line, -1)[0]
			// if we have the right to read the region
			if matches[3] == "r" {
				start, e := strconv.ParseInt("0x"+matches[1], 0, 64)
				if e != nil {
					return regions, e
				}
				end, e := strconv.ParseInt("0x"+matches[2], 0, 64)
				if e != nil {
					return regions, e
				}
				regions = append(regions, region{start: start, end: end})
			}
		}
	}

	err = scanner.Err()
	return
}
