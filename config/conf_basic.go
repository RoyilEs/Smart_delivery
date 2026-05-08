package config

type Basic struct {
	SiteName             string `json:"siteName" yaml:"site_name"`                         //site_name: "恶臭站点"
	PickupTimeoutHours   int    `json:"pickupTimeoutHours" yaml:"pickup_time_out_hours"`   //pickup_time_out_hours: 48
	TemperatureThreshold int    `json:"temperatureThreshold" yaml:"temperature_threshold"` //temperature_threshold: 114514
	SupportPhone         string `json:"supportPhone" yaml:"support_phone"`                 //support_phone: 1145141919
}

//type SettingItem struct {
//	Label       string `json:"label" yaml:"label"`
//	Description string `json:"description" yaml:"description"`
//	Group       string `json:"group" yaml:"group"`
//	ValueType   string `json:"valueType" yaml:"value_type"`
//	Value       string `json:"value" yaml:"value"`
//}
