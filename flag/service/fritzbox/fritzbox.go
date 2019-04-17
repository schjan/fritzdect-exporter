package fritzbox

import "github.com/schjan/fritzdect-exporter/flag/service/fritzbox/user"

type FritzBox struct {
	Url  string
	User user.User
}
