package model

type GlobalSetting struct {
	SectionName  string `json:"sectionName"`
	SettingName  string `json:"settingName"`
	SettingValue string `json:"settingValue"`
	SettingType  int    `json:"settingType"`
}
