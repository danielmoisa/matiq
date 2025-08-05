package common

type DataConnector interface {
	ValidateResourceOptions(resourceOptions map[string]interface{}) (ValidateResult, error)
	ValidateActionTemplate(actionOptions map[string]interface{}) (ValidateResult, error)
	TestConnection(resourceOptions map[string]interface{}) (ConnectionResult, error)
	GetMetaInfo(resourceOptions map[string]interface{}) (MetaInfoResult, error)
	Run(resourceOptions map[string]interface{}, actionOptions map[string]interface{}, rawActionOptions map[string]interface{}) (RuntimeResult, error)
}
