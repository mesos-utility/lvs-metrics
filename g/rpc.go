package g

import (
	"bytes"
	"encoding/json"
	"github.com/toolkits/net"
	"math"
	"net/http"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/open-falcon/common/model"
)

var (
	TransferClient *SingleConnRpcClient
	SendMetrics    func(metrics []*model.MetricValue)
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RpcServer string
	Timeout   time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) insureConn() {
	if this.rpcClient != nil {
		return
	}

	var err error
	var retry int = 1

	for {
		if this.rpcClient != nil {
			return
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServer, this.Timeout)
		if err == nil {
			return
		}

		glog.Warningf("dial %s fail: %v", this.RpcServer, err)

		if retry > 6 {
			retry = 1
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	this.insureConn()

	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		glog.Warningf("rpc call timeout %v => %v", this.rpcClient, this.RpcServer)
		this.close()
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}

// Set diff metrics send methods for Transfer:
//      PostToAgent     ->  http://127.0.0.1:1988/v1/push
//      SendToTransfer  ->  127.0.0.1:8433
func InitRpcClients() {
	if Config().Transfer.Enable {
		taddr := Config().Transfer.Addr
		if strings.HasPrefix(taddr, "http") {
			SendMetrics = PostToAgent
		} else {
			TransferClient = &SingleConnRpcClient{
				RpcServer: Config().Transfer.Addr,
				Timeout:   time.Duration(Config().Transfer.Timeout) * time.Millisecond,
			}
			SendMetrics = SendToTransfer
		}
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := Config().Debug

	if debug {
		glog.Infof("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	var resp model.TransferResponse
	err := TransferClient.Call("Transfer.Update", metrics, &resp)
	if err != nil {
		glog.Warningf("call Transfer.Update fail", err)
	}

	if debug {
		glog.Infof("<= %v", resp)
	}
}

func PostToAgent(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := Config().Debug

	if debug {
		glog.Infof("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	contentJson, err := json.Marshal(metrics)
	if err != nil {
		glog.Warningf("Error for PostToAgent json Marshal: %v", err)
		return
	}
	contentReader := bytes.NewReader(contentJson)
	req, err := http.NewRequest("POST", Config().Transfer.Addr, contentReader)
	if err != nil {
		glog.Warningf("Error for PostToAgent in NewRequest: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Warningf("Error for PostToAgent in http client Do: %v", err)
		return
	}
	defer resp.Body.Close()

	if debug {
		glog.Infof("<= %v", resp.Body)
	}
}
