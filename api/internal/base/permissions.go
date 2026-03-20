package base

type PermissionType int

const (
	Log PermissionType = iota + 1
	ReadOnly
	Admin
	Test
)

func (ak *ApiKey) ConfigurePermissions() {
	switch ak.PermissionString {
	case "log":
		ak.Permission = Log
	case "readonly":
		ak.Permission = ReadOnly
	case "admin":
		ak.Permission = Admin
	case "test":
		ak.Permission = Test
	}

}

type ApiKey struct {
	AppId            string `json:"app_id" mapstructure:"app_id"`
	KeyId            string `json:"key_id" mapstructure:"key_id"`
	Key              string `json:"key" mapstructure:"key"`
	PermissionString string `json:"permission" mapstructure:"permission"`
	Permission       PermissionType
}

type ApiKeys []ApiKey
