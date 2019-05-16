// nolint:lll
// Generates the onetimeadapter adapter's resource yaml. It contains the adapter's configuration, name, supported template
// names (metric in this case), and whether it is session or no-session based.
//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -a mixer/adapter/onetimeadapter/config/config.proto -x "-s=false -n onetimeadapter -t authorization"

package onetimeadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/adapter/onetimeadapter/config"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/authorization"
	"istio.io/pkg/log"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// OnetimeAdapter supports metric template.
	OnetimeAdapter struct {
		listener net.Listener
		server   *grpc.Server
	}
)

var _ authorization.HandleAuthorizationServiceServer = &OnetimeAdapter{}

func (s *OnetimeAdapter) verifyToken(user string, token string, url string) bool {
	reqBody, err := json.Marshal(map[string]string{
		"token": token,
	})

	if err != nil {
		log.Errorf("%v", err)
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Errorf("%v", err)
		return false
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-user", user)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}

	if resp.StatusCode != 200 {
		log.Errorf("Status Code: %d [%s]", resp.StatusCode, resp.Status)
		return false
	}

	log.Infof("body: %s", string(body))

	// Messed up somewhere along the line and response is getting double encoded...ah well double decode
	var firstPass string
	err = json.Unmarshal(body, &firstPass)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}

	var authResp map[string]interface{}

	err = json.Unmarshal([]byte(firstPass), &authResp)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}

	log.Infof("Resp: %v", authResp)
	return authResp["valid"].(bool)
}

// HandleAuthorization handle authorization requests
func (s *OnetimeAdapter) HandleAuthorization(ctx context.Context, r *authorization.HandleAuthorizationRequest) (*v1beta1.CheckResult, error) {

	log.Infof("received request %v\n", *r)

	cfg := &config.Params{}

	if r.AdapterConfig != nil {
		if err := cfg.Unmarshal(r.AdapterConfig.Value); err != nil {
			log.Errorf("error unmarshalling adapter config: %v", err)
			return nil, err
		}
	}
	log.Infof("PdpUrl: %s", cfg.PdpUrl)

	decodeValue := func(in interface{}) interface{} {
		switch t := in.(type) {
		case *policy.Value_StringValue:
			return t.StringValue
		case *policy.Value_Int64Value:
			return t.Int64Value
		case *policy.Value_DoubleValue:
			return t.DoubleValue
		default:
			return fmt.Sprintf("%v", in)
		}
	}

	decodeValueMap := func(in map[string]*policy.Value) map[string]interface{} {
		out := make(map[string]interface{}, len(in))
		for k, v := range in {
			out[k] = decodeValue(v.GetValue())
		}
		return out
	}

	props := decodeValueMap(r.Instance.Subject.Properties)
	user := r.Instance.Subject.User
	token := props["custom_token_header"]
	log.Infof("User[%s] Token[%s]", user, token)

	duration, _ := time.ParseDuration("1s")
	usecount := int32(1)

	if cfg.PdpUrl == "fake://deny" || token == "deny" {
		log.Infof("Fake Rejection!")
		return &v1beta1.CheckResult{
			Status:        status.WithPermissionDenied("Unauthorized..."),
			ValidDuration: duration,
			ValidUseCount: usecount,
		}, nil
	} else if cfg.PdpUrl == "fake://allow" || token == "allow" {
		return &v1beta1.CheckResult{
			Status:        status.OK,
			ValidDuration: duration,
			ValidUseCount: usecount,
		}, nil
	} else {
		authorized := s.verifyToken(user, token.(string), cfg.PdpUrl)
		if authorized {
			log.Infof("Authorization successful")
			return &v1beta1.CheckResult{
				Status:        status.OK,
				ValidDuration: duration,
				ValidUseCount: usecount,
			}, nil
		} else {
			log.Infof("Authorization failed")
			return &v1beta1.CheckResult{
				Status:        status.WithPermissionDenied("Access denied."),
				ValidDuration: duration,
				ValidUseCount: usecount,
			}, nil
		}
	}
}

// Addr returns the listening address of the server
func (s *OnetimeAdapter) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *OnetimeAdapter) Run(shutdown chan error) {
	shutdown <- s.server.Serve(s.listener)
}

// Close gracefully shuts down the server; used for testing
func (s *OnetimeAdapter) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// NewOnetimeAdapter creates a new IBP adapter that listens at provided port.
func NewOnetimeAdapter(addr string) (Server, error) {
	if addr == "" {
		addr = "0"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}
	s := &OnetimeAdapter{
		listener: listener,
	}
	log.Infof("listening on \"%v\"\n", s.Addr())
	s.server = grpc.NewServer()
	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)
	return s, nil
}
