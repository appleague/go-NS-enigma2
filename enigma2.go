package enigma2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	logger "log"
	"net/http"
	"strings"
	"time"
)

// EnableLogging allows you to turn on debug logging
var EnableLogging = true
var match = ""
var port = 80
var err = ""
var resp = ""

// STB is a settop box, identified by it's hostname/IP address.
type STB struct {
	Host            string // The hostname or IP address of the STB
	ApplicationID   string // ApplicationName is displayed on the screen the first time the ApplicationID is used.
	ApplicationName string // ApplicationName is displayed on the screen the first time the ApplicationID is used.
}

// OnlineState allows you to monitor the on/off state. The returned channel will send a boolean indicating when the STB goes online/offline.
func (stb *STB) OnlineState(interval time.Duration) chan bool {
	fmt.Println("Monitoring power state of STB %s", stb.Host)
	var lastState *bool

	stateChannel := make(chan bool, 1)

	go func() {
		for {
			online := stb.Online(interval)

			if lastState == nil || *lastState != online {
				lastState = &online

				select {
				case stateChannel <- online:
					lastState = &online
				default:
					// Nothing is listening
				}

			}

			time.Sleep(interval)

		}
	}()

	return stateChannel
}

func (stb *STB) Online(timeout time.Duration) bool {
	//get response from web powerstate call
	res, err := http.Get(fmt.Sprintf("http://%s:%d/web/powerstate", stb.Host, port))
	if err != nil {
		fmt.Println(fmt.Sprintf("FATAL: %s", err))
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(fmt.Sprintf("FATAL: %s", err))
	}

	var ret bool
	match := "<e2instandby>false</e2instandby>"

	if strings.Contains(fmt.Sprintf("%s", resp), match) {
		ret = true
	} else {
		ret = false
	}

	return ret
}

func (stb *STB) SendMessage(msg string) error {
	res, err := http.Get(fmt.Sprintf("http://%s:%d/web/message?text=%s&type=1&timeout=5", stb.Host, port, msg))
	if err != nil {
		fmt.Println(fmt.Sprintf("FATAL: %s", err))
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(fmt.Sprintf("FATAL: %s", err))
	}

	match = "<e2statetext>Message sent successfully!</e2statetext>"
	if strings.Contains(fmt.Sprintf("%s", resp), match) {
		fmt.Println(fmt.Sprintf("INFO: Sending message %s to STB with a IP %s", msg, stb.Host))
		return nil
	} else {
		fmt.Println(fmt.Sprintf("ERROR: Sending message failed: %s", err))
		return err
	}
}

func (stb *STB) SendCommand(cmd string) error {
	fmt.Println(fmt.Sprintf("INFO: Sending command %s to STB with a IP %s", match, stb.Host))
	//VOLUP STB
	match = "VOLUP"
	if strings.Contains(fmt.Sprintf("%s", cmd), match) {
		res, err := http.Get(fmt.Sprintf("http://10.0.0.7/send.htm?remote=denon&command=volup"))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match = "IR Code sent"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println("INFO: Command %s sent successfully to %s", match, stb.Host)
			return nil
		} else {
			fmt.Println("ERROR: Command sent error: %s", err)
			return err
		}
	}

	//VOLDOWN STB
	match = "VOLDOWN"
	if strings.Contains(fmt.Sprintf("%s", cmd), match) {
		res, err := http.Get(fmt.Sprintf("http://10.0.0.7/send.htm?remote=denon&command=voldown"))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match = "IR Code sent"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println(fmt.Sprintf("INFO: Command %s sent successfully to %s", match, stb.Host))
			return nil
		} else {
			fmt.Println(fmt.Sprintf("ERROR: Command sent error %s", err))
			return err
		}
	}

	//MUTE STB
	match = "MUTE"
	if strings.Contains(fmt.Sprintf("%s", cmd), match) {
		res, err := http.Get(fmt.Sprintf("http://10.0.0.7/send.htm?remote=denon&command=mute"))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match = "IR Code sent"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println(fmt.Sprintf("INFO: Command %s sent successfully to %s", match, stb.Host))
			return nil
		} else {
			fmt.Println(fmt.Sprintf("ERROR: Command sent error %s", err))
			return err
		}
	}

	//TOGGLEONOFF STB
	match = "TOGGLEONOFF"
	if strings.Contains(fmt.Sprintf("%s", cmd), match) {
		/*
			res, err := http.Get(fmt.Sprintf("http://%s:%d/web/powerstate?newstate=0", stb.Host, port))
			if err != nil {
				fmt.Println("FATAL: %s", err)
			}
			resp, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				fmt.Println("FATAL: %s", err)
			}

			match = "<e2instandby>"
			if strings.Contains(fmt.Sprintf("%s", resp), match) {
				fmt.Println("ERROR: Command sent error")
				return err
			} else {
				fmt.Println("INFO: Command %s sent successfully to %s", match, stb.Host)
				return nil
			}
		*/
		res, err := http.Get(fmt.Sprintf("http://%s:%d/web/powerstate", stb.Host, port))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match = "<e2instandby>false</e2instandby>"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println(fmt.Sprintf("INFO: Command 'TV schauen aus' started"))
			//TV schauen aus
			var postdata = []byte("var value=dom.GetObject('BidCos-RF.JEQ0038959:1.STATE').State(0);")
			preq, err := http.NewRequest("POST", "http://10.0.0.6:8181/tclrega.exe", bytes.NewBuffer(postdata))
			preq.Header.Set("Content-Type", "text/xml")

			client := &http.Client{}
			resp, err := client.Do(preq)
			if err != nil {
				fmt.Println(fmt.Sprintf("ERROR: %s", err))
			}
			defer resp.Body.Close()

			res, err = http.Get("http://10.0.0.7/send.htm?remote=lgtv&command=onoff")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			res, err = http.Get("http://10.0.0.7/send.htm?remote=denon&command=off")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			res, err = http.Get("http://10.0.0.20/web/powerstate?newstate=5")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			return nil
		} else {
			fmt.Println("INFO: Command 'TV schauen an' started")
			//TV schauen an
			var postdata = []byte("var value=dom.GetObject('BidCos-RF.JEQ0038959:1.STATE').State(1);")
			preq, err := http.NewRequest("POST", "http://10.0.0.6:8181/tclrega.exe", bytes.NewBuffer(postdata))
			preq.Header.Set("Content-Type", "text/xml")

			client := &http.Client{}
			resp, err := client.Do(preq)
			if err != nil {
				fmt.Println(fmt.Sprintf("ERROR: %s", err))
			}
			defer resp.Body.Close()

			res, err = http.Get("http://10.0.0.7/send.htm?remote=lgtv&command=onoff")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			res, err = http.Get("http://10.0.0.7/send.htm?remote=denon&command=on")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			res, err = http.Get("http://10.0.0.20/web/powerstate?newstate=4")
			if err != nil {
				fmt.Println(fmt.Sprintf("FATAL: %s", err))
			}
			res.Body.Close()

			return nil
		}
	}

	//PowerOff STB
	match = "POWEROFF"
	if strings.Contains(fmt.Sprintf("%s", cmd), match) {
		res, err := http.Get(fmt.Sprintf("http://%s:%d/web/powerstate?newstate=5", stb.Host, port))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match = "<e2instandby>false</e2instandby>"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println(fmt.Sprintf("ERROR: Command sent error"))
			return err
		} else {
			fmt.Println(fmt.Sprintf("INFO: Command %s sent successfully to %s", match, stb.Host))
			return nil
		}

	}

	//PowerOn STB
	match = "POWERON"
	if strings.Contains(cmd, match) {
		res, err := http.Get(fmt.Sprintf("http://%s:%d/web/powerstate?newstate=4", stb.Host, port))
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}
		resp, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("FATAL: %s", err))
		}

		match := "<e2instandby>true</e2instandby>"
		if strings.Contains(fmt.Sprintf("%s", resp), match) {
			fmt.Println(fmt.Sprintf("ERROR: Command sent error %s", err))
			return err
		} else {
			fmt.Println(fmt.Sprintf("INFO: Command %s sent successfully to %s", match, stb.Host))
			return nil
		}

	}
	return nil
}

func log(msg string, args ...interface{}) {
	if EnableLogging {
		logger.Printf("enigma2: "+msg, args...)
		fmt.Println(msg)
	}
}
