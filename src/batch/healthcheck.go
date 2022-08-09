package batch

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Healthcheck interface {
	DoWhile() error
	checkLocalStack() error
}

type HealthcheckImpl struct {
	url string
}

func NewHealthcheck(url string) *HealthcheckImpl {
	return &HealthcheckImpl{url: url}
}

func (h *HealthcheckImpl) checkLocalStack() error {
	resp, err := http.Get(h.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != `{"status": "running"}` {
		return fmt.Errorf("localstack not ready")
	}

	return nil
}

func (h *HealthcheckImpl) DoWhile() error {
	timesTried := 0
	for {
		err := h.checkLocalStack()
		if err == nil {
			break
		} else {
			if timesTried > 10 {
				return fmt.Errorf("tried %v times and it didn't work. Check localstack", timesTried)
			}
			time.Sleep(time.Second * 3)
			timesTried++
		}
	}
	return nil
}
