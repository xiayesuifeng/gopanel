package firewall

import "gitlab.com/xiayesuifeng/gopanel/api/server/router"

type Firewall struct {
}

func (f *Firewall) Name() string {
	return "firewall"
}

func (f *Firewall) Run(r router.Router) {

}
