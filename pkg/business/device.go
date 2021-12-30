package business

type DeviceSupport struct {
	Enabled bool
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
