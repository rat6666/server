package client

//NewClientFactoryFuncMock ...
func NewClientFactoryFuncMock(clientMock Client, err error) func() (Client, error) {
	return func() (Client, error) {
		return clientMock, err
	}
}
