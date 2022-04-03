package integ

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Protocol struct {
	cmd      cmd
	settings Settings
	config   []byte
	states   map[string][]byte
	_w       io.Writer
	wMtx     sync.Mutex
}

func (i *Protocol) ActiveStreams() Streams {
	return i.settings.Streams
}

func (i *Protocol) encode(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return i.Write(append(b, '\n'))
}

func (i *Protocol) Write(b []byte) error {
	i.wMtx.Lock()
	defer i.wMtx.Unlock()
	_, err := i._w.Write(b)
	return err
}

func (i *Protocol) Load(stream string, config, state interface{}) error {
	if len(i.config) > 0 {
		if err := json.NewDecoder(bytes.NewReader(i.config)).Decode(config); err != nil {
			return err
		}
	} else if config != nil {
		return fmt.Errorf("expected config")
	}

	if v := i.states[stream]; len(v) == 0 {
		return nil
	} else if err := json.NewDecoder(bytes.NewReader(v)).Decode(state); err != nil {
		return err
	}
	return nil
}
