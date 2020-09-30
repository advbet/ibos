package ibos

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Amfii/ftp"
)

const dialTimeout = 10 * time.Second
const retrTimeout = 1 * time.Minute

// Data wraps feed events and error to be returned via streamed channel
type Data struct {
	Data     []byte
	Filename string
	Error    error
}

// Client is iBOS feed client for feed consumption with FTP-pull
// delivery method.
//
// For FTP-pull client usage examples see demo applications in bin directory.
type Client struct {
	username string
	password string
	hostname string
	baseDir  string
}

// NewFTPClient creates a new instance of FTP-pull iBOS client. Parameter `baseURL`
// must be a ftp protocol URL and include username, password and path to racing
// documents. Port is optional, if unspecified 21 will be used.
// Example:
//
//   ftp://user:pass@azftp.phumelela.com
func NewFTPClient(baseURL string) (*Client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "ftp" {
		return nil, errors.New("only ftp URL scheme is supported")
	}

	if strings.IndexByte(url.Host, ':') == -1 {
		url.Host = url.Host + ":21"
	}
	password, _ := url.User.Password()
	client := &Client{
		username: url.User.Username(),
		password: password,
		hostname: url.Host,
		baseDir:  strings.TrimSuffix(url.Path, "/"),
	}
	return client, nil
}

// retrv returns raw file data on the FTP server
func (c *Client) retrv(conn *ftp.ServerConn, fileName string) ([]byte, error) {
	resp, err := conn.Retr(fmt.Sprintf("%s/%s", c.baseDir, fileName))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Stream starts a goroutine for continuous horses documents delivery. Stream can
// be stopped by stopping `ctx` context. Parameter `lastFile` is used to skip
// restreaming already processed documents. Polling will be performed every
// `interval` time duration. Recommended value for `interval` is one minute.
func (c *Client) Stream(ctx context.Context, lastFile string, interval time.Duration) <-chan Data {
	ch := make(chan Data)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		lastFile = c.streamPoll(ch, lastFile)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastFile = c.streamPoll(ch, lastFile)
			}
		}
	}()
	return ch
}

func (c *Client) streamPoll(ch chan<- Data, lastFile string) string {
	files, err := c.List()
	if err != nil {
		ch <- Data{Error: err}
		return lastFile
	}
	missing := missingFiles(lastFile, files)
	if len(missing) == 0 {
		return lastFile
	}

	conn, err := ftp.DialTimeouts(c.hostname, dialTimeout, retrTimeout)
	if err != nil {
		ch <- Data{Error: err}
		return lastFile
	}
	defer conn.Quit()

	if err := conn.Login(c.username, c.password); err != nil {
		ch <- Data{Error: err}
		return lastFile
	}

	for _, fileName := range missing {
		data, err := c.retrv(conn, fileName)
		if err != nil {
			ch <- Data{Error: err}
			return lastFile
		}

		ch <- Data{Data: data, Filename: fileName}
		lastFile = fileName
	}
	return missing[len(missing)-1]
}

// List returns a list of documents available on the FTP server sorted in
// file creation order.
func (c *Client) List() ([]string, error) {
	conn, err := ftp.DialTimeouts(c.hostname, dialTimeout, retrTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	if err = conn.Login(c.username, c.password); err != nil {
		return nil, err
	}

	items, err := conn.List(c.baseDir)
	if err != nil {
		return nil, err
	}

	// Sort files in chronological order by server create/modify time, as
	// primary key and file name as secondary key.
	sort.Slice(items, func(i, j int) bool {
		if items[i].Time.Equal(items[j].Time) {
			return items[i].Name < items[j].Name
		}
		return items[i].Time.Before(items[j].Time)
	})

	var files []string
	for _, item := range items {
		if item.Type != ftp.EntryTypeFile {
			continue
		}
		files = append(files, item.Name)
	}
	return files, nil
}
