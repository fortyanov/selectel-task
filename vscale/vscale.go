package vscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"selectel-task/log"
)

const API = "https://api.vscale.io/v1/scalets"

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

type Scalet struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	CtId     int    `json:"ctid"`
}

type Error struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Client struct {
	Scalets []*Scalet
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewClient() (vc *Client) {
	vc = &Client{
		Scalets: make([]*Scalet, 0),
	}

	return vc
}

func (vc *Client) CreateScalets(count int) error {
	isFailed := false
	for i := 0; i < count; i++ {
		i := i
		reqBody := &createScaletReqBody{
			MakeFrom: "ubuntu_20.04_64_001_master",
			RPlan:    "medium",
			DoStart:  true,
			Name:     fmt.Sprintf("MyTest_%d", i),
			Password: "MyPassword1!",
			Location: "spb0",
		}
		vc.wg.Add(1)
		go func() {
			log.Info(fmt.Sprintf("Creating scalet %s", reqBody.Name))
			err := vc.createScalet(reqBody)
			if err != nil && isFailed == false {
				isFailed = true
			}
		}()
	}

	vc.wg.Wait()

	if isFailed {
		vc.deleteCreatedScalets()
		return fmt.Errorf("Error on scalets creation")
	}

	log.Info("Successfully created all scalets")
	return nil
}

func (vc *Client) createScalet(structReqBody *createScaletReqBody) (err error) {
	defer vc.wg.Done()

	var reqBody []byte
	reqBody, err = json.Marshal(structReqBody)

	var request *http.Request
	if request, err = http.NewRequest("POST", API, bytes.NewBuffer(reqBody)); err != nil {
		log.Error(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
repeatCreate:
	client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.RequestTimeout),
	}
	if resp, err = client.Do(request); err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	var respBody []byte
	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error(err)
		return err
	}

	switch {
	case resp.StatusCode > 199 && resp.StatusCode < 300:
		scalet := new(Scalet)
		if err = json.Unmarshal(respBody, scalet); err != nil {
			log.Error(err)
			return err
		}
		log.Info(fmt.Sprintf("Created scalet [%d] %s", scalet.CtId, scalet.Name))
		vc.mutex.Lock()
		vc.Scalets = append(vc.Scalets, scalet)
		vc.mutex.Unlock()
		return nil
	case resp.StatusCode == 429:
		time.Sleep(time.Duration(10) * time.Microsecond)
		goto repeatCreate
	default:
		err = fmt.Errorf("Unknown server error on creating %s", structReqBody.Name)
		log.Error(err)
		return err
	}
}

func (vc *Client) deleteCreatedScalets() {
	defer vc.wg.Wait()

	for _, scalet := range vc.Scalets {
		vc.wg.Add(1)
		log.Info("Deleting created scalet", scalet.Name)
		go vc.deleteScalet(scalet.CtId)
	}
}

func (vc *Client) DeleteAllScalets() {
	var err error

	defer vc.wg.Wait()

	var request *http.Request
	if request, err = http.NewRequest("GET", API, bytes.NewBuffer(nil)); err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
	client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.RequestTimeout),
	}
	if resp, err = client.Do(request); err != nil {
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

func (vc *Client) deleteScalet(ctid int) {
	var err error

	defer vc.wg.Done()

	var request *http.Request
	if request, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%d", API, ctid), bytes.NewBuffer(nil)); err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("X-Token", cfg.XToken)

	var resp *http.Response
repeatDelete:
	client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.RequestTimeout),
	}
	if resp, err = client.Do(request); err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode > 199 && resp.StatusCode < 300:
		var respBody []byte
		if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Error(err)
			return
		}

		scalet := new(Scalet)
		if err = json.Unmarshal(respBody, scalet); err != nil {
			log.Error(err)
			return
		}
		log.Info(fmt.Sprintf("Deleted scalet [%d] %s", scalet.CtId, scalet.Name))
	case resp.StatusCode == 429:
		//time.Sleep(time.Duration(10) * time.Microsecond)
		//goto repeatDelete
	default:
		time.Sleep(time.Duration(10) * time.Microsecond)
		goto repeatDelete
	}
}
