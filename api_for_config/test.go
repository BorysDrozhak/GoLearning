package main

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {

testPath := "/tmp/awscreds"
profile := "otherProfile"

updatedData := map[string]string{
"aws_access_key_id":     "updatedKeyID",
"aws_secret_access_key": "updatedKey",
"aws_session_token":     "updatedsessionToken",
"region":                "updatedregion",
}

oldcfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, testPath)
testcfg := ini.Empty()
err = testcfg.NewSections(profile)

for k,v := range updatedData {
_, err = testcfg.Section(profile).NewKey(k,v)
}

fmt.Println(oldcfg.Section(profile).KeysHash(),"\n",testcfg.Section(profile).KeysHash())

assert.NotEqual(t,
oldcfg.Section(profile).KeysHash(),
testcfg.Section(profile).KeysHash())

err = updateProfile(testPath, profile, updatedData)

updatedcfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, testPath)
assert.NoError(t, err)
assert.EqualValues(t,updatedcfg.Section(profile).KeysHash(),testcfg.Section(profile).KeysHash())

}