package xml

import (
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"

	"fmt"
)

type Configs struct {
	Config []Config `xml:"config"`
}
type Config struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

func loadFromZkConfig() {
	encodePassword := base64.StdEncoding.EncodeToString([]byte("myPassword"))
	fmt.Printf("encode password is %s\n", string(encodePassword))

	content, err := ioutil.ReadFile("mongo.xml")
	if err != nil {
		fmt.Errorf("load zk config error: %s", err.Error())
	}
	var configs Configs
	err = xml.Unmarshal(content, &configs)
	if err != nil {
		fmt.Errorf("unmarshal zk config error: %s", err.Error())
	}
	//fmt.Println("conifgs = \n", configs)
	for _, v := range configs.Config {
		if v.Name == "mongodb.password" {
			decodedPassword, err := base64.StdEncoding.DecodeString(v.Value)
			if err != nil {
				fmt.Errorf("decode mongodb password error: %s", err.Error())
			}
			fmt.Printf("%s = %s\n", v.Name, string(decodedPassword))
		} else {
			fmt.Printf("%s = %s\n", v.Name, v.Value)
		}
	}
}
