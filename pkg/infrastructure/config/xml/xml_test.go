// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xml

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flycash/beego-config/pkg/infrastructure/config"
)

func TestXML(t *testing.T) {

	var (
		// xml parse should incluce in <config></config> tags
		xmlcontext = `<?xml version="1.0" encoding="UTF-8"?>
<config>
<appname>beeapi</appname>
<httpport>8080</httpport>
<mysqlport>3600</mysqlport>
<PI>3.1415976</PI>
<runmode>dev</runmode>
<autorender>false</autorender>
<copyrequestbody>true</copyrequestbody>
<path1>${GOPATH}</path1>
<path2>${GOPATH||/home/go}</path2>
<mysection>
<id>1</id>
<name>MySection</name>
</mysection>
</config>
`
		keyValue = map[string]interface{}{
			"appname":         "beeapi",
			"httpport":        8080,
			"mysqlport":       int64(3600),
			"PI":              3.1415976,
			"runmode":         "dev",
			"autorender":      false,
			"copyrequestbody": true,
			"path1":           os.Getenv("GOPATH"),
			"path2":           os.Getenv("GOPATH"),
			"error":           "",
			"emptystrings":    []string{},
		}
	)

	f, err := os.Create("testxml.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString(xmlcontext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("testxml.conf")

	xmlconf, err := config.NewConfig("xml", "testxml.conf")
	if err != nil {
		t.Fatal(err)
	}

	var xmlsection map[string]string
	xmlsection, err = xmlconf.GetSection(nil, "mysection")
	if err != nil {
		t.Fatal(err)
	}

	if len(xmlsection) == 0 {
		t.Error("section should not be empty")
	}

	for k, v := range keyValue {

		var (
			value interface{}
			err   error
		)

		switch v.(type) {
		case int:
			value, err = xmlconf.Int(nil, k)
		case int64:
			value, err = xmlconf.Int64(nil, k)
		case float64:
			value, err = xmlconf.Float(nil, k)
		case bool:
			value, err = xmlconf.Bool(nil, k)
		case []string:
			value, err = xmlconf.Strings(nil, k)
		case string:
			value, err = xmlconf.String(nil, k)
		default:
			value, err = xmlconf.DIY(nil, k)
		}
		if err != nil {
			t.Errorf("get key %q value fatal,%v err %s", k, v, err)
		} else if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", value) {
			t.Errorf("get key %q value, want %v got %v .", k, v, value)
		}

	}

	if err = xmlconf.Set(nil, "name", "astaxie"); err != nil {
		t.Fatal(err)
	}

	res, _ := xmlconf.String(context.Background(), "name")
	if res != "astaxie" {
		t.Fatal("get name error")
	}

	sub, err := xmlconf.Sub(context.Background(), "mysection")
	assert.Nil(t, err)
	assert.NotNil(t, sub)
	name, err := sub.String(context.Background(), "name")
	assert.Nil(t, err)
	assert.Equal(t, "MySection", name)

	id, err := sub.Int(context.Background(), "id")
	assert.Nil(t, err)
	assert.Equal(t, 1, id)

	sec := &Section{}

	err = sub.Unmarshaler(context.Background(), "", sec)
	assert.Nil(t, err)
	assert.Equal(t, "MySection", sec.Name)

	sec = &Section{}

	err = xmlconf.Unmarshaler(context.Background(), "mysection", sec)
	assert.Nil(t, err)
	assert.Equal(t, "MySection", sec.Name)

}

type Section struct {
	Name string `xml:"name"`
}
