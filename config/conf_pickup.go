package config

type Pickup struct {
	CodeLength           int    `json:"codeLength" yaml:"code_length"`                      // code_length: 6
	SmsReminder          bool   `json:"smsReminder" yaml:"sms_reminder"`                    // sms_reminder: true
	AllowAnonymousLookup bool   `json:"allowAnonymousLookup" yaml:"allow_anonymous_lookup"` // allow_anonymous_lookup: true
	AnimationSpeed       string `json:"animationSpeed" yaml:"animation_speed"`              // animation_speed: normal # slow, normal, fast
}
