package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashPwd(t *testing.T) {
	challenge := "1fa9f255"
	pwd := "kaese0815"
	expected := "1fa9f255-4ab017a7f5667a9fadf87380e6c484a7"

	res := hashPwd(challenge, pwd)

	assert.Equal(t, expected, res)
}

func TestSolveChallengeWithUsername(t *testing.T) {
	challenge := "1fa9f255"
	expected := "1fa9f255-4ab017a7f5667a9fadf87380e6c484a7"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Contains(t, req.RequestURI, expected)
		assert.Contains(t, req.RequestURI, "username=testuser")

		rw.Write([]byte(
			`<SessionInfo>
		<SID>1234567890123456</SID>
		<Challenge>951d45a8</Challenge>
		<BlockTime>0</BlockTime>
		<Rights/>
		</SessionInfo>`))
	}))
	defer server.Close()

	c := client{
		url:      server.URL,
		username: "testuser",
		password: "kaese0815",
		http:     server.Client(),
	}

	info, err := c.authenticate(challenge)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456", info.SID)
}

func TestSolveChallenge(t *testing.T) {
	challenge := "1fa9f255"
	expected := "1fa9f255-4ab017a7f5667a9fadf87380e6c484a7"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Contains(t, req.RequestURI, expected)
		assert.Contains(t, req.RequestURI, "login_sid.lua")
		assert.NotContains(t, req.RequestURI, "username=")

		rw.Write([]byte(
			`<SessionInfo>
		<SID>1234567890123456</SID>
		<Challenge>951d45a8</Challenge>
		<BlockTime>0</BlockTime>
		<Rights/>
		</SessionInfo>`))
	}))
	defer server.Close()

	c := client{
		url:      server.URL,
		password: "kaese0815",
		http:     server.Client(),
	}

	info, err := c.authenticate(challenge)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456", info.SID)
}

func TestWrongPassword(t *testing.T) {
	challenge := "1fa9f255"
	expected := "1fa9f255-4ab017a7f5667a9fadf87380e6c484a7"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Contains(t, req.RequestURI, expected)

		rw.Write([]byte(
			`<SessionInfo>
		<SID>0000000000000000</SID>
		<Challenge>951d45a8</Challenge>
		<BlockTime>0</BlockTime>
		<Rights/>
		</SessionInfo>`))
	}))
	defer server.Close()

	c := client{
		url:      server.URL,
		password: "kaese0815",
		http:     server.Client(),
	}

	sess, err := c.authenticate(challenge)
	assert.True(t, IsUnauthenticated(err))
	assert.Nil(t, sess)
}

func TestBlocked(t *testing.T) {
	challenge := "1fa9f255"
	expected := "1fa9f255-4ab017a7f5667a9fadf87380e6c484a7"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Contains(t, req.RequestURI, expected)

		rw.Write([]byte(
			`<SessionInfo>
		<SID>0000000000000000</SID>
		<Challenge>951d45a8</Challenge>
		<BlockTime>5</BlockTime>
		<Rights/>
		</SessionInfo>`))
	}))
	defer server.Close()

	c := client{
		url:      server.URL,
		password: "kaese0815",
		http:     server.Client(),
	}

	sess, err := c.authenticate(challenge)
	assert.True(t, IsUnauthenticated(err))
	assert.Equal(t, 5, sess.BlockTime)
}
