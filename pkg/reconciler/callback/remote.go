package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kyma-incubator/reconciler/pkg/reconciler"
	"go.uber.org/zap"
)

type RemoteCallbackHandler struct {
	logger      *zap.SugaredLogger
	callbackURL string
}

func NewRemoteCallbackHandler(callbackURL string, logger *zap.SugaredLogger) (Handler, error) {
	//validate URL
	if callbackURL != "" { //empty URLs are allowed (used in some test cases)
		if _, err := url.ParseRequestURI(callbackURL); err != nil {
			return nil, err
		}
	}

	//return new remote callback
	return &RemoteCallbackHandler{
		logger:      logger,
		callbackURL: callbackURL,
	}, nil
}

func (cb *RemoteCallbackHandler) Callback(msg *reconciler.CallbackMessage) error {
	if cb.callbackURL == "" { //test cases often don't provide a callback URL
		cb.logger.Warn("Empty callback-URL provided: remote callback not executed")
		return nil
	}

	requestBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(cb.callbackURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		cb.logger.Errorf("Callback request failed: %s", err)
		return err
	}
	//dump request for debugging purposes
	dumpResp, dumpErr := httputil.DumpResponse(resp, true)
	if dumpErr == nil {
		cb.logger.Debugf("HTTP response dump: %s", string(dumpResp))
	} else {
		cb.logger.Debugf("Failed to generate HTTP response dump: %s", dumpErr)
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Callback request ('%s')  failed with '%d' HTTP response code",
			msg,
			resp.StatusCode)
		cb.logger.Info(msg)
		return fmt.Errorf(msg)
	}

	return nil
}
