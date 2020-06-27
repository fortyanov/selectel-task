package vscale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"selectel-task/log"
)

var cfg *config

func Init() (err error) {
	if cfg, err = initConfig(); err != nil {
		return err
	}
	return nil
}

type createScaletReqBody struct {
	MakeFrom string `json:"make_from"`
	RPlan    string `json:"rplan"`
	DoStart  bool   `json:"do_start"`
	Name     string `json:"name"`
	//Keys []string	`json:"keys"`
	Password string `json:"password"`
	Location string `json:"location"`
}

type ScaletInfo struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	CtId     int    `json:"ctid"`
}

type VscaleClient struct {
	Scalets []*ScaletInfo
	client  *http.Client
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewVscaleClient(timeout time.Duration) (vc *VscaleClient) {
	vc = &VscaleClient{
		Scalets: make([]*ScaletInfo, 0),
		client: &http.Client{
			Timeout: timeout,
		},
	}

	return vc
}

func (vc *VscaleClient) CreateScalets(count int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer vc.wg.Wait()

	for i := 0; i < count; i++ {
		vc.wg.Add(1)
		reqBody := &createScaletReqBody{
			MakeFrom: "ubuntu_20.04_64_001_master",
			RPlan:    "medium",
			DoStart:  true,
			Name:     fmt.Sprintf("MyTest_%d", i),
			Password: "MyPassword1!",
			Location: "spb0",
		}

		select {
		case <-ctx.Done():
			vc.deleteCreatedScalets()
			return
		default:
			log.Info(fmt.Sprintf("Creating scalet %s", reqBody.Name))
			go vc.createScalet(ctx, cancel, reqBody)
		}
	}
}

func (vc *VscaleClient) createScalet(ctx context.Context, cancel context.CancelFunc, structReqBody *createScaletReqBody) {
	var err error

	defer vc.wg.Done()

	var reqBody []byte
	reqBody, err = json.Marshal(structReqBody)

	var request *http.Request
	if request, err = http.NewRequestWithContext(ctx, "POST", "https://api.vscale.io/v1/scalets", bytes.NewBuffer(reqBody)); err != nil {
		log.Error(err)
		cancel()
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
	if resp, err = vc.client.Do(request); err != nil {
		log.Error(err)
		cancel()
		return
	}
	defer resp.Body.Close()

	var respBody []byte
	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error(err)
		cancel()
		return
	}

	structRespBody := new(ScaletInfo)
	if err = json.Unmarshal(respBody, structRespBody); err != nil {
		log.Error(err)
		cancel()
		return
	}
	log.Info(fmt.Sprintf("Created scalet [%d] %s", structRespBody.CtId, structRespBody.Name))

	vc.mutex.Lock()
	vc.Scalets = append(vc.Scalets, structRespBody)
	vc.mutex.Unlock()
}

func (vc *VscaleClient) deleteCreatedScalets() {
	for _, scalet := range vc.Scalets {
		vc.wg.Add(1)

		go vc.deleteScalet(scalet.CtId)
	}
}

func (vc *VscaleClient) deleteScalet(ctid int) {
	var err error

	defer vc.wg.Done()

	var request *http.Request
	if request, err = http.NewRequest("DELETE", fmt.Sprintf("https://api.vscale.io/v1/scalets/%d", ctid), bytes.NewBuffer(nil)); err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
	if resp, err = vc.client.Do(request); err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	var respBody []byte
	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error(err)
		return
	}

	structRespBody := new(ScaletInfo)
	if err = json.Unmarshal(respBody, structRespBody); err != nil {
		log.Error(err)
		return
	}
	log.Info(fmt.Sprintf("Deleted scalet [%d] %s", structRespBody.CtId, structRespBody.Name))
}

func (vc *VscaleClient) DeleteAllScalets() {
	var err error

	defer vc.wg.Wait()

	var request *http.Request
	if request, err = http.NewRequest("GET", "https://api.vscale.io/v1/scalets", bytes.NewBuffer(nil)); err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
	if resp, err = vc.client.Do(request); err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	var respBody []byte
	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(respBody, &vc.Scalets); err != nil {
		log.Error(err)
		return
	}

	for _, scalet := range vc.Scalets {
		vc.wg.Add(1)
		log.Info(fmt.Sprintf("Deleting scalet [%d] %s", scalet.CtId, scalet.Hostname))
		go vc.deleteScalet(scalet.CtId)
	}
}
