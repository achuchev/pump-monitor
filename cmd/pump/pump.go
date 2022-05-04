package pump

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	c "github.com/achuchev/pump-monitor/cmd/common"
	"github.com/achuchev/pump-monitor/cmd/utils"
	log "github.com/sirupsen/logrus"
)

var (
	jar = NewJar()

	httpClient = http.Client{Jar: jar}
)

type LiveData struct {
	Data      []float32 `json:"data"`
	Positions []string  `json:"positions"`
}

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func Authenticate() error {
	jar.cookies = map[string][]*http.Cookie{}
	log.Tracef(utils.GetURL("controller_login/setLogin"))
	log.Tracef("username: %s", c.FlagsConfig.ServiceUsername)
	log.Tracef("password: %s", c.FlagsConfig.ServicePassword)

	resp, _ := httpClient.PostForm(utils.GetURL("controller_login/setLogin"), url.Values{
		"login":    {c.FlagsConfig.ServiceUsername},
		"password": {c.FlagsConfig.ServicePassword},
	})

	if resp.StatusCode != 200 {
		return fmt.Errorf("authentication failed: %v", resp.Status)
	}

	//FIXME: check if the authentication was successful by checking the cookies

	b, _ := ioutil.ReadAll(resp.Body)
	log.Tracef("Auth response: %s", string(b))

	resp.Body.Close()
	return nil
}

func GetLiveData(DeviceID string, selection []string, positions []string) (LiveData, error) {
	var liveData LiveData

	resp, _ := httpClient.PostForm(utils.GetURL("controller_mydevice/getLiveData"), url.Values{
		"DeviceID":    {DeviceID},
		"selection[]": selection,
		"positions[]": positions,
	})
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return liveData, fmt.Errorf("failed to get live data: %v", string(b))
	}

	resp.Body.Close()

	if b == nil {
		return liveData, fmt.Errorf("failed to get live data")
	}

	err := json.Unmarshal(b, &liveData)
	if err != nil {
		log.Errorf("failed to unmarshal live data: %s", string(b))
		return liveData, fmt.Errorf("failed to get live data: %v", err)
	}
	return liveData, nil
}
