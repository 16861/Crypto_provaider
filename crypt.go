package main

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"encoding/json"
    "net/http"
    "bytes"
	"encoding/xml"
)

const (
	PATH_TO_DATAFILE    string = "data"
	PATH_TO_JSONCONFIG    string = "config.json"
	PATH_TO_RESFILE     string = "res.txt"
	PATH_TO_RESDATAFILE string = "resdata"
)


type GlobalParams struct {
	datafile, res_file, res_data, command string
}

type SendParams struct {
	serial, pathToCert, password string
}

type ConfigSt struct {
	globParams GlobalParams
	sndParams SendParams
}

type Command struct {
	message string
    datafile string
    response string
    conf ConfigSt
}

func (c *Command) encrypt() {
	encData := c.getEncodedData(c.conf.globParams.datafile)
    crtData, err := ioutil.ReadFile(c.conf.sndParams.pathToCert)
    if err != nil {
		fmt.Println("Can't open cert file: ", c.conf.sndParams.pathToCert)
        panic(err)
	} 
    encCert := base64.StdEncoding.EncodeToString([]byte(crtData))
    encPassword := base64.StdEncoding.EncodeToString([]byte(c.conf.sndParams.password))

	c.message = fmt.Sprintf(`<xml>%s %s %s %s</xml>`, encCert, c.conf.sndParams.serial, encPassword, encData)
}

func (c *Command) getSkFiles(typeOfCerts string) {
    var keyUsage string
    switch typeOfCerts {
        case "forcrypt":
            keyUsage = "56"
        case "forsign":
            keyUsage = "192"
        default:
            panic("Wrong params typeOfCerts!")
    }

    encPassword := base64.StdEncoding.EncodeToString([]byte(c.conf.sndParams.password))
    c.message = fmt.Sprintf(`<xml>%s %s</xml>`, encPassword, keyUsage)

}

func (c *Command) setSign() {
    encData := c.getEncodedData(c.conf.globParams.datafile)
    encPassword := base64.StdEncoding.EncodeToString([]byte(c.conf.sndParams.password))
    c.message = fmt.Sprintf(`<xml>%s %s %s</xml>`, c.conf.sndParams.serial, encPassword, encData)

}

func (c *Command) getSigned() {
    encData := c.getEncodedData(c.conf.globParams.datafile)
    encPassword := base64.StdEncoding.EncodeToString([]byte(c.conf.sndParams.password))
    c.message = fmt.Sprintf(`<xml>%s %s</xml>`, c.conf.sndParams.serial, encPassword, encData)
}

func (c *Command) getEncodedData(filename string) string  {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (c *Command) sendRequest() {
	fmt.Println("Sending request to Cryptoprovider...")
    url := "http://127.0.0.1:19744"
    r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(c.message)))
    r.Header.Set("Content-Type", "text/xml")

    client := &http.Client{}
    resp, err := client.Do(r)
    if err != nil {
        panic(err)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    err = ioutil.WriteFile(c.conf.globParams.res_file, body, 0644)
}

func (c *Command) parseResponse() {
	data, err := ioutil.ReadFile(c.conf.globParams.res_file)
	if err != nil {
		panic("Can't open response file!")
	}
	var res ResposeFromCrypt
	xml.Unmarshal(data, &res)
	switch res.Result.Code {
		case 0:
			fmt.Println("Code 0. Cryptoprovider successive process request!")
			strData := base64.StdEncoding.Decode(base64Text, []byte(res.Data))
			err := ioutil.WriteFile(c.conf.globParams.res_data, strData, 0644)
			fmt.Println("Data is stored to ", c.conf.globParams.res_data)
		default:
			fmt.Println("Code ", res.Result.Code, ". Cryptoprovider failed to process request!")
	}
}

func (conf *ConfigSt) parseConfig(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Can't open from file: ", filename)
	}

	var outerLevel interface{}
	err = json.Unmarshal(data, &outerLevel)
	mOuterLevel := outerLevel.(map[string]interface{})
	for key, val := range mOuterLevel {
		switch key {
			case "global":
				mInnerLevel := val.(map[string]interface{})
				for innKey, innVal := range mInnerLevel {
					switch innKey {
						case "datafile":
							// reflect.ValueOf(&conf.globParams).Elem().FieldByName("datafile").SetString(innVal.(string))
							conf.globParams.datafile = innVal.(string)
						case "res_file":
							conf.globParams.res_file = innVal.(string)
						case "res_data":
							conf.globParams.res_data = innVal.(string)
						case "command":
							conf.globParams.command = innVal.(string)
						default:
							ms := fmt.Sprintf("Wrong parameter in config file! Tag: %s", innKey)
							panic(ms)
					}
				}
			case "params":
				mInnerLevel := val.(map[string]interface{})
				for innKey, innVal := range mInnerLevel {
					switch innKey {
						case "serial":
							conf.sndParams.serial = innVal.(string)
						case "path_to_cert":
							conf.sndParams.pathToCert = innVal.(string)
						case "password":
							conf.sndParams.password = innVal.(string)
						default:
							ms := fmt.Sprintf("Wrong parameter in config file! Tag: %s", innKey)
							panic(ms)
					}
				}
			default:
				ms := fmt.Sprintf("Wrong parameter in config file! Tag: %s", key)
				panic(ms)
		}
	}
	
}

type ReturnData struct {
	XMLName xml.Name `xml:"Return"`
	Code int `xml:"code"`
	Message string `xml:"Message"`
}

type ResposeFromCrypt struct {
	XMLName xml.Name `xml:"Result"`
	Data string `xml:"Data"`
	Ret ReturnData
}

func main() {
	var comm Command
	comm.conf.parseConfig(PATH_TO_JSONCONFIG)
	switch comm.conf.globParams.command {
		case "encrypt":
			comm.encrypt()
        case "getSkFiles":
            comm.getSkFiles("forsign")
        case "setSign":
            comm.setSign()
        case "getSigned":
            comm.getSigned()
		default:
			ms := fmt.Sprintf("Wrong command! Command in config: %s. Available commands: encrypt, getSkfilem setSign, getSigned", comm.conf.globParams.command)
			panic(ms)
	}

	comm.sendRequest()
	comm.parseResponse()

	return 
}