package main

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	// // "reflect"
	// // "bufio"
	"encoding/json"
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

type SendMessage struct {
	header, body, footer string
}

type Command struct {
	message string
}

func (c *Command) encrypt(filename, pathToCert, serial, password string) {
	var mes SendMessage
	encData := c.getEncodedData(filename)

	mes.header = fmt.Sprintf("<xml><data>%s</data>%s, %s, %s</xml>", encData, pathToCert, serial, password)
	mes.footer = "</data></req>"

	c.message =  fmt.Sprintf("%v", mes)

}

func (c *Command) getEncodedData(filename string) string  {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString([]byte(data))

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
							panic("Wrong xml at parameter!")
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
							panic("Wrong xml at parameter!")
					}
				}
			default:
				panic("Wrong xml at parameter!")
		}
	}
	
}

func main() {

	var conf ConfigSt
	var comm Command
	conf.parseConfig(PATH_TO_JSONCONFIG)
	switch conf.globParams.command {
		case "encrypt":
			comm.encrypt(conf.globParams.datafile, conf.sndParams.pathToCert, conf.sndParams.serial, conf.sndParams.password)
			fmt.Println(comm.message)
		default:
			panic("Wrong XML structure")
	}
	

	return 
}
