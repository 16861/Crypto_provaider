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

// type Send_message interface {

// }

type GlobalParams struct {
	datafile, res_file, res_data string
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

type Message struct {
    Name string
    Food string
}

func encrypt(data, pathToCert, serial, password string) SendMessage {
	var mes SendMessage

	mes.header = fmt.Sprintf("<xml><data>%s</data>%s, %s, %s</xml>", data, pathToCert, serial, password)
	mes.footer = "</data></req>"

	return mes

}

func parseConfig(filename string, conf *ConfigSt) {
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
	parseConfig(PATH_TO_JSONCONFIG, &conf)
	fmt.Println("Filepath: ",conf.sndParams.serial)

	return 
	b := [] byte(`{
    "global": {
        "datafile": "data",
        "res_file": "res.txt",
        "res_data": "res_data"
    },
    "params": {
        "serial": "dsadsa",
        "path_to_cert": "/home/espadon",
        "password": "123"
    }
}`)

	// var conf ConfigSt
	var  f interface {}
	err := json.Unmarshal(b, &f)
	mss := f.(map[string]interface{})
	// err := json.Unmarshal(b, &conf)
	for k, v := range mss {
		fmt.Println(k, v)
		
		switch vv := v.(type) {
			case string:
				fmt.Println(k, " is string ", v, vv)
			case []interface{}:
				for k_i, v_i := range vv {
					fmt.Println(k_i, v_i)
				}
			default:
				var vss = v.(map[string]interface{})
				for k_i, v_i := range vss {
					fmt.Println(k_i, v_i)
				}
		}
	}

	b1 := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	var m Message
	err1 := json.Unmarshal(b1, &m)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(m)


	
	// text := "4qd223"
	data, err := ioutil.ReadFile(PATH_TO_DATAFILE)
	if err != nil {
		panic(err)
	}
	

	encData := base64.StdEncoding.EncodeToString([]byte(string(data)))
	fmt.Println(encrypt(encData, "/home/espadon", "sf34223f", "123"))


	// data, err := base64.StdEncoding.DecodeString(str)

}
