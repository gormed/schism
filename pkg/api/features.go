package api

import (
	"gitlab.void-ptr.org/go/schism/pkg/business"
)

// FeatureSet of the api
type FeatureSet struct {
	Devices business.DeviceSupport `json:"devices"`
	Data    business.DataSupport   `json:"data"`
}

// Features supported
var Features = FeatureSet{
	Devices: business.DeviceSupport{
		Enabled: true,
	},
	Data: business.DataSupport{
		Enabled: true,
	},
}
