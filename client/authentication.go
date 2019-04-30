package client

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"gopkg.in/resty.v1"
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

const (
	loginSidUrl = "login_sid.lua"
)

func (c *client) Login() error {
	info, err := c.getSessionInfo()
	if err != nil {
		if !IsUnauthenticated(err) {
			return err
		}
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
	resp, err := c.r().Get(loginSidUrl)
	if err != nil {
		return nil, err
	}

	var sess SessionInfo
	err = xml.Unmarshal(resp.Body(), &sess)
	if err != nil {
		return nil, err
	}

	if sess.BlockTime > 0 {
		return &SessionInfo{
			BlockTime: sess.BlockTime,
		}, unauthenticatedError
	}

	if sess.SID == unauthenticatedSid {
		return &sess, unauthenticatedError
	}

	return &sess, nil
}

func (c *client) authenticate(challenge string) (*SessionInfo, error) {
	r := c.r().SetQueryParam("response", hashPwd(challenge, c.password))
	if c.username != "" {
		r = r.SetQueryParam("username", c.username)
	}

	resp, err := r.Get(loginSidUrl)
	if err != nil {
		return nil, err
	}

	var sess SessionInfo
	err = xml.Unmarshal(resp.Body(), &sess)
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

	c.sid = sess.SID

	return &sess, nil
}

func hashPwd(challenge, password string) string {
	challengeStr := utf16.Encode([]rune(challenge + "-" + password))
	result := make([]byte, len(challengeStr)*2)
	for i, c := range challengeStr {
		result[i*2] = byte(c)
		result[i*2+1] = byte(c >> 8)
	}
	hash := md5.Sum(result)

	return fmt.Sprintf("%s-%x", challenge, hash)
}

func (c *client) r() *resty.Request {
	r := resty.NewWithClient(c.http).
		SetHostURL(c.url)

	if c.sid != "" {
		r.SetQueryParam("sid", c.sid)
	}

	return r.R()
}
