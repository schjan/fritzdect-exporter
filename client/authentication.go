package client

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"unicode/utf16"
)

type SessionInfo struct {
	XMLName   xml.Name `xml:"SessionInfo"`
	Text      string   `xml:",chardata"`
	SID       string   `xml:"SID"`
	Challenge string   `xml:"Challenge"`
	BlockTime int      `xml:"BlockTime"`
	Rights    string   `xml:"Rights"`
}

func (c *client) login() error {
	info, err := c.getSessionInfo()
	if err != nil {
		return err
	}

	if info.SID != unauthenticatedSid {
		// logged in

		return nil
	}

	info, err = c.authenticate(info.Challenge)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getSessionInfo() (*SessionInfo, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/login_sid.lua", c.rootUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sess SessionInfo
	err = xml.Unmarshal(bts, &sess)
	if err != nil {
		return nil, err
	}

	if sess.BlockTime > 0 {
		return &SessionInfo{
			BlockTime: sess.BlockTime,
		}, unauthenticatedError
	}

	if sess.SID == unauthenticatedSid {
		return nil, unauthenticatedError
	}

	return &sess, nil
}

func (c *client) authenticate(challenge string) (*SessionInfo, error) {
	challengeStr := utf16.Encode([]rune(challenge + "-" + c.password))
	result := make([]byte, len(challengeStr)*2)
	for i, c := range challengeStr {
		result[i*2] = byte(c)
		result[i*2+1] = byte(c >> 8)
	}
	hash := md5.Sum(result)

	//todo add username if set
	req := fmt.Sprintf("%s/login_sid.lua?response=%s-%x", c.rootUrl, challenge, hash)

	resp, err := c.http.Get(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sess SessionInfo
	err = xml.Unmarshal(bts, &sess)
	if err != nil {
		return nil, err
	}

	if sess.BlockTime > 0 {
		return &SessionInfo{
			BlockTime: sess.BlockTime,
		}, unauthenticatedError
	}

	if sess.SID == unauthenticatedSid {
		return nil, unauthenticatedError
	}

	return &sess, nil
}
