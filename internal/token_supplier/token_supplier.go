package tokensupplier

import (
	"sync"
)

type Fetcher func() (token string, err error)
type Supplier struct {
	mx      sync.Mutex
	token   string
	fetcher Fetcher
}

func MustNewSupplier(fetcher Fetcher) *Supplier {
	sp := &Supplier{fetcher: fetcher}
	if err := sp.RefreshToken(); err != nil {
		panic(err)
	}
	return sp
}
func NewSupplier(fetcher Fetcher) (*Supplier, error) {
	sp := &Supplier{fetcher: fetcher}
	return sp, sp.RefreshToken()
}

func (ts *Supplier) Get() string {
	ts.mx.Lock()
	defer ts.mx.Unlock()
	return ts.token
}
func (ts *Supplier) set(value string) {
	ts.mx.Lock()
	defer ts.mx.Unlock()
	ts.token = value

}

func (ts *Supplier) RefreshToken() error {
	token, err := ts.fetcher()
	if err != nil {
		return err
	}
	ts.set(token)
	return nil
}
