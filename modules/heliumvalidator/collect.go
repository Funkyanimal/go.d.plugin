package heliumvalidator

import (
        "bytes"
	      "encoding/json"
        "fmt"
	      "io"
	      "io/ioutil"
        "net/http"

	      "github.com/netdata/go.d.plugin/pkg/stm"
	      "github.com/netdata/go.d.plugin/pkg/web"
)

const (
	     jsonRPCVersion = "2.0"
  
	     methodBlockHeight = "block_height"
	     methodBlockAge    = "block_age"
	     methodInConsensus = "info_in_consensus"
)

var infoRequests = rpcRequests{
	      {JSONRPC: jsonRPCVersion, ID: 1, Method: methodBlockHeight},
	      {JSONRPC: jsonRPCVersion, ID: 2, Method: methodBlockAge},
	      {JSONRPC: jsonRPCVersion, ID: 3, Method: methodInConsensus},
}

func (e *Heliumvalidator) collect() (map[string]int64, error) {
	      responses, err := e.scrapeHeliumvalidator(infoRequests)
	      if err != nil {
	            	return nil, err
	      }

	      info, err := e.collectInfoResponse(infoRequests, responses)
      	      if err != nil {
		            return nil, err
	      }

	      return stm.ToMap(info), nil
}

func (e *Heliumvalidator) collectInfoResponse(requests rpcRequests, responses rpcResponses) (*heliumvalidatorInfo, error) {
	var info heliumvalidatorInfo
	for _, req := range requests {
		resp := responses.getByID(req.ID)
		if resp == nil {
			e.Warningf("method '%s' (id %d) not in responses", req.Method, req.ID)
			continue
		}

		if resp.Error != nil {
			e.Warningf("server returned an error on method '%s': %v", req.Method, resp.Error)
			continue
		}

		var err error
		switch req.Method {
		case methodBlockHeight:
			info.Blockchain, err = parseBlockchainInfo(resp.Result)
		case methodBlockAge:
			info.MemPool, err = parseMemPoolInfo(resp.Result)
		case methodInConsensus:
			info.Network, err = parseNetworkInfo(resp.Result)
		}
		if err != nil {
			return nil, fmt.Errorf("parse '%s' method result: %v", req.Method, err)
		}
	}

	return &info, nil
}

func parseBlockchainInfo(result []byte) (*blockchainInfo, error) {
	var m blockchainInfo
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func parseMemPoolInfo(result []byte) (*memPoolInfo, error) {
	var m memPoolInfo
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func parseNetworkInfo(result []byte) (*networkInfo, error) {
	var m networkInfo
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func parseTXOutSetInfo(result []byte) (*txOutSetInfo, error) {
	var m txOutSetInfo
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func parseMemoryInfo(result []byte) (*memoryInfo, error) {
	var m memoryInfo
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (e *Energid) scrapeEnergid(requests rpcRequests) (rpcResponses, error) {
	req, _ := web.NewHTTPRequest(e.Request)
	req.Method = http.MethodPost
	req.Header.Set("Content-Type", "application/json")
	body, _ := json.Marshal(requests)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	var resp rpcResponses
	if err := e.doOKDecode(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *Energid) doOKDecode(req *http.Request, in interface{}) error {
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on HTTP request '%s': %v", req.URL, err)
	}
	defer closeBody(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("'%s' returned HTTP status code: %d", req.URL, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(in); err != nil {
		return fmt.Errorf("error on decoding response from '%s': %v", req.URL, err)
	}

	return nil
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
