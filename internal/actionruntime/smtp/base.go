package smtp

import (
	"os"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/gomail.v2"
)

func (s *Connector) getConnectionWithOptions(resourceOptions map[string]interface{}) (*gomail.Dialer, error) {
	if err := mapstructure.Decode(resourceOptions, &s.ResourceOpts); err != nil {
		return nil, err
	}

	smtpDialer := gomail.NewDialer(s.ResourceOpts.Host, s.ResourceOpts.Port, s.ResourceOpts.Username, s.ResourceOpts.Password)

	return smtpDialer, nil
}

func attachSizeLimiter(contentLength int64) bool {
	limitStr := os.Getenv("ILLA_S3_LIMIT")
	limit64, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit64 = 5
	}
	objectSize := contentLength / 1024
	objectSize /= 1024

	if objectSize > limit64 {
		return true
	}
	return false
}
