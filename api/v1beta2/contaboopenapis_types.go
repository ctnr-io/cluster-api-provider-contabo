// Package models provides primitives to interact with the openapi HTTP API.
package v1beta2

const (
	BearerScopes = "bearer.Scopes"
)

// Defines values for ApiPermissionsResponseActions.
const (
	ApiPermissionsResponseActionsCREATE ApiPermissionsResponseActions = "CREATE"
	ApiPermissionsResponseActionsDELETE ApiPermissionsResponseActions = "DELETE"
	ApiPermissionsResponseActionsREAD   ApiPermissionsResponseActions = "READ"
	ApiPermissionsResponseActionsUPDATE ApiPermissionsResponseActions = "UPDATE"
)

// Defines values for AssignmentAuditResponseAction.
const (
	AssignmentAuditResponseActionCREATED AssignmentAuditResponseAction = "CREATED"
	AssignmentAuditResponseActionDELETED AssignmentAuditResponseAction = "DELETED"
)

// Defines values for AutoScalingTypeRequestState.
const (
	AutoScalingTypeRequestStateDisabled AutoScalingTypeRequestState = "disabled"
	AutoScalingTypeRequestStateEnabled  AutoScalingTypeRequestState = "enabled"
)

// Defines values for AutoScalingTypeResponseState.
const (
	AutoScalingTypeResponseStateDisabled AutoScalingTypeResponseState = "disabled"
	AutoScalingTypeResponseStateEnabled  AutoScalingTypeResponseState = "enabled"
	AutoScalingTypeResponseStateError    AutoScalingTypeResponseState = "error"
)

// Defines values for CreateCustomImageRequestOsType.
const (
	Linux   CreateCustomImageRequestOsType = "Linux"
	Windows CreateCustomImageRequestOsType = "Windows"
)

// Defines values for CreateInstanceRequestDefaultUser.
const (
	CreateInstanceRequestDefaultUserAdmin         CreateInstanceRequestDefaultUser = "admin"
	CreateInstanceRequestDefaultUserAdministrator CreateInstanceRequestDefaultUser = "administrator"
	CreateInstanceRequestDefaultUserRoot          CreateInstanceRequestDefaultUser = "root"
)

// Defines values for CreateInstanceRequestLicense.
const (
	CPanel100  CreateInstanceRequestLicense = "cPanel100"
	CPanel1000 CreateInstanceRequestLicense = "cPanel1000"
	CPanel150  CreateInstanceRequestLicense = "cPanel150"
	CPanel200  CreateInstanceRequestLicense = "cPanel200"
	CPanel250  CreateInstanceRequestLicense = "cPanel250"
	CPanel30   CreateInstanceRequestLicense = "cPanel30"
	CPanel300  CreateInstanceRequestLicense = "cPanel300"
	CPanel350  CreateInstanceRequestLicense = "cPanel350"
	CPanel400  CreateInstanceRequestLicense = "cPanel400"
	CPanel450  CreateInstanceRequestLicense = "cPanel450"
	CPanel5    CreateInstanceRequestLicense = "cPanel5"
	CPanel50   CreateInstanceRequestLicense = "cPanel50"
	CPanel500  CreateInstanceRequestLicense = "cPanel500"
	CPanel550  CreateInstanceRequestLicense = "cPanel550"
	CPanel600  CreateInstanceRequestLicense = "cPanel600"
	CPanel650  CreateInstanceRequestLicense = "cPanel650"
	CPanel700  CreateInstanceRequestLicense = "cPanel700"
	CPanel750  CreateInstanceRequestLicense = "cPanel750"
	CPanel800  CreateInstanceRequestLicense = "cPanel800"
	CPanel850  CreateInstanceRequestLicense = "cPanel850"
	CPanel900  CreateInstanceRequestLicense = "cPanel900"
	CPanel950  CreateInstanceRequestLicense = "cPanel950"
	PleskAdmin CreateInstanceRequestLicense = "PleskAdmin"
	PleskHost  CreateInstanceRequestLicense = "PleskHost"
	PleskPro   CreateInstanceRequestLicense = "PleskPro"
)

// Defines values for CreateInstanceRequestRegion.
const (
	AUS       CreateInstanceRequestRegion = "AUS"
	EU        CreateInstanceRequestRegion = "EU"
	IND       CreateInstanceRequestRegion = "IND"
	JPN       CreateInstanceRequestRegion = "JPN"
	SIN       CreateInstanceRequestRegion = "SIN"
	UK        CreateInstanceRequestRegion = "UK"
	USCentral CreateInstanceRequestRegion = "US-central"
	USEast    CreateInstanceRequestRegion = "US-east"
	USWest    CreateInstanceRequestRegion = "US-west"
)

// Defines values for CreateObjectStorageResponseDataStatus.
const (
	CreateObjectStorageResponseDataStatusCANCELLED            CreateObjectStorageResponseDataStatus = "CANCELLED"
	CreateObjectStorageResponseDataStatusDISABLED             CreateObjectStorageResponseDataStatus = "DISABLED"
	CreateObjectStorageResponseDataStatusERROR                CreateObjectStorageResponseDataStatus = "ERROR"
	CreateObjectStorageResponseDataStatusLIMITEXCEEDED        CreateObjectStorageResponseDataStatus = "LIMIT_EXCEEDED"
	CreateObjectStorageResponseDataStatusMANUALPROVISIONING   CreateObjectStorageResponseDataStatus = "MANUAL_PROVISIONING"
	CreateObjectStorageResponseDataStatusORDERPROCESSING      CreateObjectStorageResponseDataStatus = "ORDER_PROCESSING"
	CreateObjectStorageResponseDataStatusPENDINGPAYMENT       CreateObjectStorageResponseDataStatus = "PENDING_PAYMENT"
	CreateObjectStorageResponseDataStatusPRODUCTNOTAVAILABLE  CreateObjectStorageResponseDataStatus = "PRODUCT_NOT_AVAILABLE"
	CreateObjectStorageResponseDataStatusPROVISIONING         CreateObjectStorageResponseDataStatus = "PROVISIONING"
	CreateObjectStorageResponseDataStatusREADY                CreateObjectStorageResponseDataStatus = "READY"
	CreateObjectStorageResponseDataStatusUNKNOWN              CreateObjectStorageResponseDataStatus = "UNKNOWN"
	CreateObjectStorageResponseDataStatusUPGRADING            CreateObjectStorageResponseDataStatus = "UPGRADING"
	CreateObjectStorageResponseDataStatusVERIFICATIONREQUIRED CreateObjectStorageResponseDataStatus = "VERIFICATION_REQUIRED"
)

// Defines values for CreateSecretRequestType.
const (
	CreateSecretRequestTypePassword CreateSecretRequestType = "password"
	CreateSecretRequestTypeSsh      CreateSecretRequestType = "ssh"
)

// Defines values for CreateUserRequestLocale.
const (
	CreateUserRequestLocaleDe   CreateUserRequestLocale = "de"
	CreateUserRequestLocaleDeDE CreateUserRequestLocale = "de-DE"
	CreateUserRequestLocaleEn   CreateUserRequestLocale = "en"
	CreateUserRequestLocaleEnUS CreateUserRequestLocale = "en-US"
	CreateUserRequestLocaleEs   CreateUserRequestLocale = "es"
	CreateUserRequestLocaleEsES CreateUserRequestLocale = "es-ES"
	CreateUserRequestLocalePt   CreateUserRequestLocale = "pt"
	CreateUserRequestLocalePtBR CreateUserRequestLocale = "pt-BR"
)

// Defines values for DataCenterResponseCapabilities.
const (
	ObjectStorage     DataCenterResponseCapabilities = "Object-Storage"
	PrivateNetworking DataCenterResponseCapabilities = "Private-Networking"
	VDS               DataCenterResponseCapabilities = "VDS"
	VPS               DataCenterResponseCapabilities = "VPS"
)

// Defines values for ImageAuditResponseDataAction.
const (
	ImageAuditResponseDataActionCREATED ImageAuditResponseDataAction = "CREATED"
	ImageAuditResponseDataActionDELETED ImageAuditResponseDataAction = "DELETED"
	ImageAuditResponseDataActionUPDATED ImageAuditResponseDataAction = "UPDATED"
)

// Defines values for ImageResponseFormat.
const (
	ImageResponseFormatIso   ImageResponseFormat = "iso"
	ImageResponseFormatQcow2 ImageResponseFormat = "qcow2"
)

// Defines values for ImageResponseTenantId.
const (
	ImageResponseTenantIdDE  ImageResponseTenantId = "DE"
	ImageResponseTenantIdINT ImageResponseTenantId = "INT"
)

// Defines values for InstanceResponseDefaultUser.
const (
	InstanceResponseDefaultUserAdmin         InstanceResponseDefaultUser = "admin"
	InstanceResponseDefaultUserAdministrator InstanceResponseDefaultUser = "administrator"
	InstanceResponseDefaultUserRoot          InstanceResponseDefaultUser = "root"
)

// Defines values for InstanceResponseProductType.
const (
	InstanceResponseProductTypeHdd  InstanceResponseProductType = "hdd"
	InstanceResponseProductTypeNvme InstanceResponseProductType = "nvme"
	InstanceResponseProductTypeSsd  InstanceResponseProductType = "ssd"
	InstanceResponseProductTypeVds  InstanceResponseProductType = "vds"
)

// Defines values for InstanceResponseTenantId.
const (
	InstanceResponseTenantIdDE  InstanceResponseTenantId = "DE"
	InstanceResponseTenantIdINT InstanceResponseTenantId = "INT"
)

// Defines values for InstancesStatus.
const (
	InstancesStatusInstalling           InstancesStatus = "installing"
	InstancesStatusOk                   InstancesStatus = "ok"
	InstancesStatusReinstall            InstancesStatus = "reinstall"
	InstancesStatusReinstallationFailed InstancesStatus = "reinstallation failed"
	InstancesStatusRestart              InstancesStatus = "restart"
)

// Defines values for InstancesActionsAuditResponseAction.
const (
	InstancesActionsAuditResponseActionCREATED InstancesActionsAuditResponseAction = "CREATED"
	InstancesActionsAuditResponseActionDELETED InstancesActionsAuditResponseAction = "DELETED"
	InstancesActionsAuditResponseActionUPDATED InstancesActionsAuditResponseAction = "UPDATED"
)

// Defines values for InstancesAuditResponseAction.
const (
	InstancesAuditResponseActionCREATED InstancesAuditResponseAction = "CREATED"
	InstancesAuditResponseActionDELETED InstancesAuditResponseAction = "DELETED"
	InstancesAuditResponseActionUPDATED InstancesAuditResponseAction = "UPDATED"
)

// Defines values for ListImageResponseDataFormat.
const (
	ListImageResponseDataFormatIso   ListImageResponseDataFormat = "iso"
	ListImageResponseDataFormatQcow2 ListImageResponseDataFormat = "qcow2"
)

// Defines values for ListImageResponseDataTenantId.
const (
	ListImageResponseDataTenantIdDE  ListImageResponseDataTenantId = "DE"
	ListImageResponseDataTenantIdINT ListImageResponseDataTenantId = "INT"
)

// Defines values for ListInstancesResponseDataDefaultUser.
const (
	ListInstancesResponseDataDefaultUserAdmin         ListInstancesResponseDataDefaultUser = "admin"
	ListInstancesResponseDataDefaultUserAdministrator ListInstancesResponseDataDefaultUser = "administrator"
	ListInstancesResponseDataDefaultUserRoot          ListInstancesResponseDataDefaultUser = "root"
)

// Defines values for ListInstancesResponseDataProductType.
const (
	ListInstancesResponseDataProductTypeHdd  ListInstancesResponseDataProductType = "hdd"
	ListInstancesResponseDataProductTypeNvme ListInstancesResponseDataProductType = "nvme"
	ListInstancesResponseDataProductTypeSsd  ListInstancesResponseDataProductType = "ssd"
	ListInstancesResponseDataProductTypeVds  ListInstancesResponseDataProductType = "vds"
)

// Defines values for ListInstancesResponseDataTenantId.
const (
	ListInstancesResponseDataTenantIdDE  ListInstancesResponseDataTenantId = "DE"
	ListInstancesResponseDataTenantIdINT ListInstancesResponseDataTenantId = "INT"
)

// Defines values for ListVipResponseDataIpVersion.
const (
	ListVipResponseDataIpVersionV4 ListVipResponseDataIpVersion = "v4"
)

// Defines values for ListVipResponseDataResourceType.
const (
	ListVipResponseDataResourceTypeBareMetal ListVipResponseDataResourceType = "bare-metal"
	ListVipResponseDataResourceTypeInstances ListVipResponseDataResourceType = "instances"
	ListVipResponseDataResourceTypeNull      ListVipResponseDataResourceType = "null"
)

// Defines values for ListVipResponseDataType.
const (
	ListVipResponseDataTypeAdditional ListVipResponseDataType = "additional"
	ListVipResponseDataTypeFloating   ListVipResponseDataType = "floating"
)

// Defines values for ObjectStorageAuditResponseAction.
const (
	ObjectStorageAuditResponseActionCREATED ObjectStorageAuditResponseAction = "CREATED"
	ObjectStorageAuditResponseActionDELETED ObjectStorageAuditResponseAction = "DELETED"
	ObjectStorageAuditResponseActionUPDATED ObjectStorageAuditResponseAction = "UPDATED"
)

// Defines values for ObjectStorageResponseStatus.
const (
	ObjectStorageResponseStatusCANCELLED            ObjectStorageResponseStatus = "CANCELLED"
	ObjectStorageResponseStatusDISABLED             ObjectStorageResponseStatus = "DISABLED"
	ObjectStorageResponseStatusERROR                ObjectStorageResponseStatus = "ERROR"
	ObjectStorageResponseStatusLIMITEXCEEDED        ObjectStorageResponseStatus = "LIMIT_EXCEEDED"
	ObjectStorageResponseStatusMANUALPROVISIONING   ObjectStorageResponseStatus = "MANUAL_PROVISIONING"
	ObjectStorageResponseStatusORDERPROCESSING      ObjectStorageResponseStatus = "ORDER_PROCESSING"
	ObjectStorageResponseStatusPENDINGPAYMENT       ObjectStorageResponseStatus = "PENDING_PAYMENT"
	ObjectStorageResponseStatusPRODUCTNOTAVAILABLE  ObjectStorageResponseStatus = "PRODUCT_NOT_AVAILABLE"
	ObjectStorageResponseStatusPROVISIONING         ObjectStorageResponseStatus = "PROVISIONING"
	ObjectStorageResponseStatusREADY                ObjectStorageResponseStatus = "READY"
	ObjectStorageResponseStatusUNKNOWN              ObjectStorageResponseStatus = "UNKNOWN"
	ObjectStorageResponseStatusUPGRADING            ObjectStorageResponseStatus = "UPGRADING"
	ObjectStorageResponseStatusVERIFICATIONREQUIRED ObjectStorageResponseStatus = "VERIFICATION_REQUIRED"
)

// Defines values for PermissionRequestActions.
const (
	PermissionRequestActionsCREATE PermissionRequestActions = "CREATE"
	PermissionRequestActionsDELETE PermissionRequestActions = "DELETE"
	PermissionRequestActionsREAD   PermissionRequestActions = "READ"
	PermissionRequestActionsUPDATE PermissionRequestActions = "UPDATE"
)

// Defines values for PermissionResponseActions.
const (
	CREATE PermissionResponseActions = "CREATE"
	DELETE PermissionResponseActions = "DELETE"
	READ   PermissionResponseActions = "READ"
	UPDATE PermissionResponseActions = "UPDATE"
)

// Defines values for PrivateNetworkAuditResponseAction.
const (
	PrivateNetworkAuditResponseActionCREATED PrivateNetworkAuditResponseAction = "CREATED"
	PrivateNetworkAuditResponseActionDELETED PrivateNetworkAuditResponseAction = "DELETED"
	PrivateNetworkAuditResponseActionUPDATED PrivateNetworkAuditResponseAction = "UPDATED"
)

// Defines values for ReinstallInstanceRequestDefaultUser.
const (
	ReinstallInstanceRequestDefaultUserAdmin         ReinstallInstanceRequestDefaultUser = "admin"
	ReinstallInstanceRequestDefaultUserAdministrator ReinstallInstanceRequestDefaultUser = "administrator"
	ReinstallInstanceRequestDefaultUserRoot          ReinstallInstanceRequestDefaultUser = "root"
)

// Defines values for RoleAuditResponseAction.
const (
	RoleAuditResponseActionCREATED RoleAuditResponseAction = "CREATED"
	RoleAuditResponseActionDELETED RoleAuditResponseAction = "DELETED"
	RoleAuditResponseActionUPDATED RoleAuditResponseAction = "UPDATED"
)

// Defines values for SecretAuditResponseAction.
const (
	SecretAuditResponseActionCREATED SecretAuditResponseAction = "CREATED"
	SecretAuditResponseActionDELETED SecretAuditResponseAction = "DELETED"
	SecretAuditResponseActionUPDATED SecretAuditResponseAction = "UPDATED"
)

// Defines values for SecretResponseType.
const (
	SecretResponseTypePassword SecretResponseType = "password"
	SecretResponseTypeSsh      SecretResponseType = "ssh"
)

// Defines values for SnapshotsAuditResponseAction.
const (
	SnapshotsAuditResponseActionCREATED SnapshotsAuditResponseAction = "CREATED"
	SnapshotsAuditResponseActionDELETED SnapshotsAuditResponseAction = "DELETED"
	SnapshotsAuditResponseActionUPDATED SnapshotsAuditResponseAction = "UPDATED"
)

// Defines values for TagAuditResponseAction.
const (
	TagAuditResponseActionCREATED TagAuditResponseAction = "CREATED"
	TagAuditResponseActionDELETED TagAuditResponseAction = "DELETED"
	TagAuditResponseActionUPDATED TagAuditResponseAction = "UPDATED"
)

// Defines values for UpdateUserRequestLocale.
const (
	UpdateUserRequestLocaleDe   UpdateUserRequestLocale = "de"
	UpdateUserRequestLocaleDeDE UpdateUserRequestLocale = "de-DE"
	UpdateUserRequestLocaleEn   UpdateUserRequestLocale = "en"
	UpdateUserRequestLocaleEnUS UpdateUserRequestLocale = "en-US"
	UpdateUserRequestLocaleEs   UpdateUserRequestLocale = "es"
	UpdateUserRequestLocaleEsES UpdateUserRequestLocale = "es-ES"
	UpdateUserRequestLocalePt   UpdateUserRequestLocale = "pt"
	UpdateUserRequestLocalePtBR UpdateUserRequestLocale = "pt-BR"
)

// Defines values for UpgradeAutoScalingTypeState.
const (
	Disabled UpgradeAutoScalingTypeState = "disabled"
	Enabled  UpgradeAutoScalingTypeState = "enabled"
)

// Defines values for UpgradeObjectStorageResponseDataStatus.
const (
	CANCELLED            UpgradeObjectStorageResponseDataStatus = "CANCELLED"
	DISABLED             UpgradeObjectStorageResponseDataStatus = "DISABLED"
	ERROR                UpgradeObjectStorageResponseDataStatus = "ERROR"
	LIMITEXCEEDED        UpgradeObjectStorageResponseDataStatus = "LIMIT_EXCEEDED"
	MANUALPROVISIONING   UpgradeObjectStorageResponseDataStatus = "MANUAL_PROVISIONING"
	ORDERPROCESSING      UpgradeObjectStorageResponseDataStatus = "ORDER_PROCESSING"
	PENDINGPAYMENT       UpgradeObjectStorageResponseDataStatus = "PENDING_PAYMENT"
	PRODUCTNOTAVAILABLE  UpgradeObjectStorageResponseDataStatus = "PRODUCT_NOT_AVAILABLE"
	PROVISIONING         UpgradeObjectStorageResponseDataStatus = "PROVISIONING"
	READY                UpgradeObjectStorageResponseDataStatus = "READY"
	UNKNOWN              UpgradeObjectStorageResponseDataStatus = "UNKNOWN"
	UPGRADING            UpgradeObjectStorageResponseDataStatus = "UPGRADING"
	VERIFICATIONREQUIRED UpgradeObjectStorageResponseDataStatus = "VERIFICATION_REQUIRED"
)

// Defines values for UserAuditResponseAction.
const (
	UserAuditResponseActionCREATED UserAuditResponseAction = "CREATED"
	UserAuditResponseActionDELETED UserAuditResponseAction = "DELETED"
	UserAuditResponseActionUPDATED UserAuditResponseAction = "UPDATED"
)

// Defines values for UserResponseLocale.
const (
	De   UserResponseLocale = "de"
	DeDE UserResponseLocale = "de-DE"
	En   UserResponseLocale = "en"
	EnUS UserResponseLocale = "en-US"
	Es   UserResponseLocale = "es"
	EsES UserResponseLocale = "es-ES"
	Pt   UserResponseLocale = "pt"
	PtBR UserResponseLocale = "pt-BR"
)

// Defines values for VipAuditResponseAction.
const (
	VipAuditResponseActionCREATED VipAuditResponseAction = "CREATED"
	VipAuditResponseActionDELETED VipAuditResponseAction = "DELETED"
	VipAuditResponseActionUPDATED VipAuditResponseAction = "UPDATED"
)

// Defines values for VipResponseIpVersion.
const (
	VipResponseIpVersionV4 VipResponseIpVersion = "v4"
)

// Defines values for VipResponseResourceType.
const (
	VipResponseResourceTypeBareMetal VipResponseResourceType = "bare-metal"
	VipResponseResourceTypeInstances VipResponseResourceType = "instances"
	VipResponseResourceTypeNull      VipResponseResourceType = "null"
)

// Defines values for VipResponseType.
const (
	VipResponseTypeAdditional VipResponseType = "additional"
	VipResponseTypeFloating   VipResponseType = "floating"
)

// Defines values for InstanceStatus.
const (
	InstanceStatusError                InstanceStatus = "error"
	InstanceStatusInstalling           InstanceStatus = "installing"
	InstanceStatusManualProvisioning   InstanceStatus = "manual_provisioning"
	InstanceStatusOther                InstanceStatus = "other"
	InstanceStatusPendingPayment       InstanceStatus = "pending_payment"
	InstanceStatusProductNotAvailable  InstanceStatus = "product_not_available"
	InstanceStatusProvisioning         InstanceStatus = "provisioning"
	InstanceStatusRescue               InstanceStatus = "rescue"
	InstanceStatusResetPassword        InstanceStatus = "reset_password"
	InstanceStatusRunning              InstanceStatus = "running"
	InstanceStatusStopped              InstanceStatus = "stopped"
	InstanceStatusUninstalled          InstanceStatus = "uninstalled"
	InstanceStatusUnknown              InstanceStatus = "unknown"
	InstanceStatusVerificationRequired InstanceStatus = "verification_required"
)

// Defines values for RetrieveInstancesListParamsStatus.
const (
	Error                RetrieveInstancesListParamsStatus = "error"
	Installing           RetrieveInstancesListParamsStatus = "installing"
	ManualProvisioning   RetrieveInstancesListParamsStatus = "manual_provisioning"
	Other                RetrieveInstancesListParamsStatus = "other"
	PendingPayment       RetrieveInstancesListParamsStatus = "pending_payment"
	ProductNotAvailable  RetrieveInstancesListParamsStatus = "product_not_available"
	Provisioning         RetrieveInstancesListParamsStatus = "provisioning"
	Rescue               RetrieveInstancesListParamsStatus = "rescue"
	ResetPassword        RetrieveInstancesListParamsStatus = "reset_password"
	Running              RetrieveInstancesListParamsStatus = "running"
	Stopped              RetrieveInstancesListParamsStatus = "stopped"
	Uninstalled          RetrieveInstancesListParamsStatus = "uninstalled"
	Unknown              RetrieveInstancesListParamsStatus = "unknown"
	VerificationRequired RetrieveInstancesListParamsStatus = "verification_required"
)

// Defines values for RetrieveSecretListParamsType.
const (
	Password RetrieveSecretListParamsType = "password"
	Ssh      RetrieveSecretListParamsType = "ssh"
)

// Defines values for RetrieveVipListParamsResourceType.
const (
	RetrieveVipListParamsResourceTypeBareMetal RetrieveVipListParamsResourceType = "bare-metal"
	RetrieveVipListParamsResourceTypeInstances RetrieveVipListParamsResourceType = "instances"
	RetrieveVipListParamsResourceTypeNull      RetrieveVipListParamsResourceType = "null"
)

// Defines values for RetrieveVipListParamsIpVersion.
const (
	V4 RetrieveVipListParamsIpVersion = "v4"
)

// Defines values for RetrieveVipListParamsType.
const (
	Additional RetrieveVipListParamsType = "additional"
	Floating   RetrieveVipListParamsType = "floating"
)

// Defines values for UnassignIpParamsResourceType.
const (
	UnassignIpParamsResourceTypeBareMetal UnassignIpParamsResourceType = "bare-metal"
	UnassignIpParamsResourceTypeInstances UnassignIpParamsResourceType = "instances"
)

// Defines values for AssignIpParamsResourceType.
const (
	AssignIpParamsResourceTypeBareMetal AssignIpParamsResourceType = "bare-metal"
	AssignIpParamsResourceTypeInstances AssignIpParamsResourceType = "instances"
)

// AddOnRequest defines model for AddOnRequest.
type AddOnRequest struct {
	// Id Id of the Addon. Please refer to list [here](https://contabo.com/en/product-list/?show_ids=true).
	Id int64 `json:"id"`

	// Quantity The number of Addons you wish to aquire.
	Quantity int64 `json:"quantity"`
}

// AddOnResponse defines model for AddOnResponse.
type AddOnResponse struct {
	// Id Id of the Addon. Please refer to list [here](https://contabo.com/en/product-list/?show_ids=true).
	Id int64 `json:"id"`

	// Quantity The number of Addons you wish to aquire.
	Quantity int64 `json:"quantity"`
}

// AdditionalIp defines model for AdditionalIp.
type AdditionalIp struct {
	V4 IpV4 `json:"v4"`
}

// ApiPermissionsResponse defines model for ApiPermissionsResponse.
type ApiPermissionsResponse struct {
	// Actions Action allowed for the API endpoint. Basically `CREATE` corresponds to POST endpoints, `READ` to GET endpoints, `UPDATE` to PATCH / PUT endpoints and `DELETE` to DELETE endpoints.
	Actions []ApiPermissionsResponseActions `json:"actions"`

	// ApiName API endpoint. In order to get a list availbale api enpoints please refer to the GET api-permissions endpoint.
	ApiName string `json:"apiName"`
}

// ApiPermissionsResponseActions defines model for ApiPermissionsResponse.Actions.
type ApiPermissionsResponseActions string

// AssignInstancePrivateNetworkResponse defines model for AssignInstancePrivateNetworkResponse.
type AssignInstancePrivateNetworkResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks InstanceAssignmentSelfLinks `json:"_links"`
}

// AssignVipResponse defines model for AssignVipResponse.
type AssignVipResponse struct {
	UnderscoreLinks SelfLinks     `json:"_links"`
	Data            []VipResponse `json:"data"`
}

// AssignedTagResponse defines model for AssignedTagResponse.
type AssignedTagResponse struct {
	// TagId Tag's id
	TagId float32 `json:"tagId"`

	// TagName Tag's name
	TagName string `json:"tagName"`
}

// AssignmentAuditResponse defines model for AssignmentAuditResponse.
type AssignmentAuditResponse struct {
	// Action Audit Action
	Action AssignmentAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes Changes made for a specific Tag
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id float32 `json:"id"`

	// RequestId Request ID
	RequestId string `json:"requestId"`

	// ResourceId Resource's id
	ResourceId string `json:"resourceId"`

	// ResourceType Resource type. Resource type is one of `instance|image|object-storage`.
	ResourceType string `json:"resourceType"`

	// TagId Tag's id
	TagId int64 `json:"tagId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp Audit creation date
	Timestamp int64 `json:"timestamp"`

	// TraceId Trace ID
	TraceId string `json:"traceId"`

	// Username User Full Name
	Username string `json:"username"`
}

// AssignmentAuditResponseAction Audit Action
type AssignmentAuditResponseAction string

// AssignmentResponse defines model for AssignmentResponse.
type AssignmentResponse struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// ResourceId Resource id
	ResourceId string `json:"resourceId"`

	// ResourceName Resource name
	ResourceName string `json:"resourceName"`

	// ResourceType Resource type. Resource type is one of `instance|image|object-storage`.
	ResourceType string `json:"resourceType"`

	// TagId The identifier of the tag.
	TagId int64 `json:"tagId"`

	// TagName Tag's name
	TagName string `json:"tagName"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// AutoScalingTypeRequest defines model for AutoScalingTypeRequest.
type AutoScalingTypeRequest struct {
	// SizeLimitTB Autoscaling size limit for the current object storage.
	SizeLimitTB float64 `json:"sizeLimitTB"`

	// State State of the autoscaling for the current object storage.
	State AutoScalingTypeRequestState `json:"state"`
}

// AutoScalingTypeRequestState State of the autoscaling for the current object storage.
type AutoScalingTypeRequestState string

// AutoScalingTypeResponse defines model for AutoScalingTypeResponse.
type AutoScalingTypeResponse struct {
	// ErrorMessage Error message
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// SizeLimitTB Autoscaling size limit for the current object storage.
	SizeLimitTB float64 `json:"sizeLimitTB"`

	// State State of the autoscaling for the current object storage.
	State AutoScalingTypeResponseState `json:"state"`
}

// AutoScalingTypeResponseState State of the autoscaling for the current object storage.
type AutoScalingTypeResponseState string

// Backup defines model for Backup.
type Backup = map[string]string

// CancelInstanceRequest defines model for CancelInstanceRequest.
type CancelInstanceRequest struct {
	// CancelDate Date of cancellation
	CancelDate *int64 `json:"cancelDate,omitempty"`
}

// CancelInstanceResponse defines model for CancelInstanceResponse.
type CancelInstanceResponse struct {
	UnderscoreLinks SelfLinks                    `json:"_links"`
	Data            []CancelInstanceResponseData `json:"data"`
}

// CancelInstanceResponseData defines model for CancelInstanceResponseData.
type CancelInstanceResponseData struct {
	// CancelDate The date on which the instance will be cancelled
	CancelDate string `json:"cancelDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Instance's id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CancelObjectStorageRequest defines model for CancelObjectStorageRequest.
type CancelObjectStorageRequest struct {
	// CancelDate Date of cancellation
	CancelDate *int64 `json:"cancelDate,omitempty"`
}

// CancelObjectStorageResponse defines model for CancelObjectStorageResponse.
type CancelObjectStorageResponse struct {
	UnderscoreLinks SelfLinks                         `json:"_links"`
	Data            []CancelObjectStorageResponseData `json:"data"`
}

// CancelObjectStorageResponseData defines model for CancelObjectStorageResponseData.
type CancelObjectStorageResponseData struct {
	// CancelDate Cancellation date for object storage.
	CancelDate string `json:"cancelDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DisplayName Display name for object storage.
	DisplayName string `json:"displayName"`

	// ObjectStorageId Object Storage id
	ObjectStorageId string `json:"objectStorageId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// ClientResponse defines model for ClientResponse.
type ClientResponse struct {
	// ClientId IDM client id
	ClientId string `json:"clientId"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Id Client's id
	Id string `json:"id"`

	// Secret IDM client secret
	Secret string `json:"secret"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// ClientSecretResponse defines model for ClientSecretResponse.
type ClientSecretResponse struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Secret IDM client secret
	Secret string `json:"secret"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateAssignmentResponse defines model for CreateAssignmentResponse.
type CreateAssignmentResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks TagAssignmentSelfLinks `json:"_links"`
}

// CreateCustomImageFailResponse defines model for CreateCustomImageFailResponse.
type CreateCustomImageFailResponse struct {
	// Message Unsupported Media Type: Please provide a direct link to an .iso or .qcow2 image.
	Message string `json:"message"`

	// StatusCode statuscode:415
	StatusCode int `json:"statusCode"`
}

// CreateCustomImageRequest defines model for CreateCustomImageRequest.
type CreateCustomImageRequest struct {
	// Description Image Description
	Description *string `json:"description,omitempty"`

	// Name Image Name
	Name string `json:"name"`

	// OsType Provided type of operating system (OS). Please specify `Windows` for MS Windows and `Linux` for other OS. Specifying wrong OS type may lead to disfunctional cloud instance.
	OsType CreateCustomImageRequestOsType `json:"osType"`

	// Url URL from where the image has been downloaded / provided.
	Url string `json:"url"`

	// Version Version number to distinguish the contents of an image. Could be the version of the operating system for example.
	Version string `json:"version"`
}

// CreateCustomImageRequestOsType Provided type of operating system (OS). Please specify `Windows` for MS Windows and `Linux` for other OS. Specifying wrong OS type may lead to disfunctional cloud instance.
type CreateCustomImageRequestOsType string

// CreateCustomImageResponse defines model for CreateCustomImageResponse.
type CreateCustomImageResponse struct {
	UnderscoreLinks SelfLinks                       `json:"_links"`
	Data            []CreateCustomImageResponseData `json:"data"`
}

// CreateCustomImageResponseData defines model for CreateCustomImageResponseData.
type CreateCustomImageResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateInstanceAddons defines model for CreateInstanceAddons.
type CreateInstanceAddons struct {
	// AdditionalIps Set this attribute if you want to upgrade your instance with the Additional IPs addon. Please provide an empty object for the time being as value. There will be more configuration possible in the future.
	AdditionalIps *map[string]string `json:"additionalIps,omitempty"`
	AddonsIds     *[]AddOnRequest    `json:"addonsIds,omitempty"`

	// Backup Set this attribute if you want to upgrade your instance with the Automated backup addon.     Please provide an empty object for the time being as value. There will be more configuration possible     in the future.
	Backup *map[string]string `json:"backup,omitempty"`

	// CustomImage Set this attribute if you want to upgrade your instance with the Custom Images addon.   Please provide an empty object for the time being as value. There will be more configuration possible   in the future.
	CustomImage *map[string]string `json:"customImage,omitempty"`

	// ExtraStorage Set this attribute if you want to upgrade your instance with the Extra Storage addon.
	ExtraStorage *ExtraStorageRequest `json:"extraStorage,omitempty"`

	// PrivateNetworking Set this attribute if you want to upgrade your instance with the Private Networking addon.   Please provide an empty object for the time being as value. There will be more configuration possible   in the future.
	PrivateNetworking *map[string]string `json:"privateNetworking,omitempty"`
}

// CreateInstanceRequest defines model for CreateInstanceRequest.
type CreateInstanceRequest struct {
	// AddOns Set attributes in the addons object for the corresponding ones that need to be added to the instance
	AddOns *CreateInstanceAddons `json:"addOns,omitempty"`

	// ApplicationId Application ID
	ApplicationId *string `json:"applicationId,omitempty"`

	// DefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
	DefaultUser *CreateInstanceRequestDefaultUser `json:"defaultUser,omitempty"`

	// DisplayName The display name of the instance
	DisplayName *string `json:"displayName,omitempty"`

	// ImageId ImageId to be used to setup the compute instance. Default is Ubuntu 22.04
	ImageId *string `json:"imageId,omitempty"`

	// License Additional licence in order to enhance your chosen product, mainly needed for software licenses on your product (not needed for windows).
	License *CreateInstanceRequestLicense `json:"license,omitempty"`

	// Period Initial contract period in months. Available periods are: 1, 3, 6 and 12 months. Default to 1 month
	Period int64 `json:"period"`

	// ProductId Default is V92
	ProductId *string `json:"productId,omitempty"`

	// Region Instance Region where the compute instance should be located. Default is EU
	Region *CreateInstanceRequestRegion `json:"region,omitempty"`

	// RootPassword `secretId` of the password for the `defaultUser` with administrator/root privileges. For Linux/BSD please use SSH, for Windows RDP. Please refer to Secrets Management API.
	RootPassword *int64 `json:"rootPassword,omitempty"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys *[]int64 `json:"sshKeys,omitempty"`

	// UserData [Cloud-Init](https://cloud-init.io/) Config in order to customize during start of compute instance.
	UserData *string `json:"userData,omitempty"`
}

// CreateInstanceRequestDefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
type CreateInstanceRequestDefaultUser string

// CreateInstanceRequestLicense Additional licence in order to enhance your chosen product, mainly needed for software licenses on your product (not needed for windows).
type CreateInstanceRequestLicense string

// CreateInstanceRequestRegion Instance Region where the compute instance should be located. Default is EU
type CreateInstanceRequestRegion string

// CreateInstanceResponse defines model for CreateInstanceResponse.
type CreateInstanceResponse struct {
	UnderscoreLinks SelfLinks                    `json:"_links"`
	Data            []CreateInstanceResponseData `json:"data"`
}

// CreateInstanceResponseData defines model for CreateInstanceResponseData.
type CreateInstanceResponseData struct {
	AddOns []AddOnResponse `json:"addOns"`

	// CreatedDate Creation date for instance
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// InstanceId Instance's id
	InstanceId int64 `json:"instanceId"`

	// OsType Type of operating system (OS)
	OsType string `json:"osType"`

	// ProductId Product ID
	ProductId string `json:"productId"`

	// Region Instance Region where the compute instance should be located.
	Region string `json:"region"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys []int64 `json:"sshKeys"`

	// Status Instance's status
	Status InstanceStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateObjectStorageRequest defines model for CreateObjectStorageRequest.
type CreateObjectStorageRequest struct {
	// AutoScaling Autoscaling settings
	AutoScaling *AutoScalingTypeRequest `json:"autoScaling,omitempty"`

	// DisplayName Display name helps to differentiate between object storages, especially if they are in the same region. If display name is not provided, it will be generated. Display name can be changed any time.
	DisplayName *string `json:"displayName,omitempty"`

	// Region Region where the object storage should be located. Default is EU. Available regions: EU, US-central, SIN
	Region string `json:"region"`

	// TotalPurchasedSpaceTB Amount of purchased / requested object storage in TB.
	TotalPurchasedSpaceTB float64 `json:"totalPurchasedSpaceTB"`
}

// CreateObjectStorageResponse defines model for CreateObjectStorageResponse.
type CreateObjectStorageResponse struct {
	UnderscoreLinks SelfLinks                         `json:"_links"`
	Data            []CreateObjectStorageResponseData `json:"data"`
}

// CreateObjectStorageResponseData defines model for CreateObjectStorageResponseData.
type CreateObjectStorageResponseData struct {
	// AutoScaling Autoscaling settings
	AutoScaling AutoScalingTypeResponse `json:"autoScaling"`

	// CancelDate Cancellation date for object storage.
	CancelDate string `json:"cancelDate"`

	// CreatedDate Creation date for object storage.
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DataCenter The data center of the storage
	DataCenter string `json:"dataCenter"`

	// DisplayName Display name for object storage.
	DisplayName string `json:"displayName"`

	// ObjectStorageId Your object storage id
	ObjectStorageId string `json:"objectStorageId"`

	// Region The region where your object storage is located
	Region string `json:"region"`

	// S3TenantId Your S3 tenantId. Only required for public sharing.
	S3TenantId string `json:"s3TenantId"`

	// S3Url S3 URL to connect to your S3 compatible object storage
	S3Url string `json:"s3Url"`

	// Status The object storage status
	Status CreateObjectStorageResponseDataStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// TotalPurchasedSpaceTB Amount of purchased / requested object storage in TB.
	TotalPurchasedSpaceTB float64 `json:"totalPurchasedSpaceTB"`

	// UsedSpacePercentage Currently used space in percentage.
	UsedSpacePercentage float64 `json:"usedSpacePercentage"`

	// UsedSpaceTB Currently used space in TB.
	UsedSpaceTB float64 `json:"usedSpaceTB"`
}

// CreateObjectStorageResponseDataStatus The object storage status
type CreateObjectStorageResponseDataStatus string

// CreatePrivateNetworkRequest defines model for CreatePrivateNetworkRequest.
type CreatePrivateNetworkRequest struct {
	// Description The description of the Private Network. There is a limit of 255 characters per Private Network description.
	Description *string `json:"description,omitempty"`

	// Name The name of the Private Network. It may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per Private Network name.
	Name string `json:"name"`

	// Region Region where the Private Network should be located. Default is `EU`
	Region *string `json:"region,omitempty"`
}

// CreatePrivateNetworkResponse defines model for CreatePrivateNetworkResponse.
type CreatePrivateNetworkResponse struct {
	UnderscoreLinks SelfLinks                `json:"_links"`
	Data            []PrivateNetworkResponse `json:"data"`
}

// CreateRoleRequest defines model for CreateRoleRequest.
type CreateRoleRequest struct {
	// AccessAllResources Allow access to all resources. This will superseed all assigned resources in a role.
	AccessAllResources bool `json:"accessAllResources"`

	// Admin If user is admin he will have permissions to all API endpoints and resources. Enabling this will superseed all role definitions and `accessAllResources`.
	Admin bool `json:"admin"`

	// Name The name of the role. There is a limit of 255 characters per role.
	Name        string               `json:"name"`
	Permissions *[]PermissionRequest `json:"permissions,omitempty"`
}

// CreateRoleResponse defines model for CreateRoleResponse.
type CreateRoleResponse struct {
	UnderscoreLinks SelfLinks                `json:"_links"`
	Data            []CreateRoleResponseData `json:"data"`
}

// CreateRoleResponseData defines model for CreateRoleResponseData.
type CreateRoleResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// RoleId Role's id
	RoleId int64 `json:"roleId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateSecretRequest defines model for CreateSecretRequest.
type CreateSecretRequest struct {
	// Name The name of the secret that will keep the password
	Name string `json:"name"`

	// Type The type of the secret. Can be `password` or `ssh`
	Type CreateSecretRequestType `json:"type"`

	// Value The secret value that needs to be saved. In case of a password it must match a pattern with at least one upper and lower case character and either one number with two special characters `!@#$^&*?_~` or at least three numbers with one special character `!@#$^&*?_~`. This is expressed in the following regular expression: `^((?=.*?[A-Z]{1,})(?=.*?[a-z]{1,}))(((?=(?:[^d]*d){1})(?=([^^&*?_~]*[!@#$^&*?_~]){2,}))|((?=(?:[^d]*d){3})(?=.*?[!@#$^&*?_~]+))).{8,}$`
	Value string `json:"value"`
}

// CreateSecretRequestType The type of the secret. Can be `password` or `ssh`
type CreateSecretRequestType string

// CreateSecretResponse defines model for CreateSecretResponse.
type CreateSecretResponse struct {
	UnderscoreLinks SelfLinks        `json:"_links"`
	Data            []SecretResponse `json:"data"`
}

// CreateSnapshotRequest defines model for CreateSnapshotRequest.
type CreateSnapshotRequest struct {
	// Description The description of the snapshot. There is a limit of 255 characters per snapshot.
	Description *string `json:"description,omitempty"`

	// Name The name of the snapshot. It may contain letters, numbers, spaces, dashes. There is a limit of 30 characters per snapshot.
	Name string `json:"name"`
}

// CreateSnapshotResponse defines model for CreateSnapshotResponse.
type CreateSnapshotResponse struct {
	UnderscoreLinks SelfLinks          `json:"_links"`
	Data            []SnapshotResponse `json:"data"`
}

// CreateTagRequest defines model for CreateTagRequest.
type CreateTagRequest struct {
	// Color The color of the tag. Color can be specified using hexadecimal value. Default color is #0A78C3
	Color string `json:"color"`

	// Description The description of the Tag name.
	Description *string `json:"description,omitempty"`

	// Name The name of the tag. Tags may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per tag.
	Name string `json:"name"`
}

// CreateTagResponse defines model for CreateTagResponse.
type CreateTagResponse struct {
	UnderscoreLinks SelfLinks               `json:"_links"`
	Data            []CreateTagResponseData `json:"data"`
}

// CreateTagResponseData defines model for CreateTagResponseData.
type CreateTagResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// TagId Tag's id
	TagId int64 `json:"tagId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateTicketRequest defines model for CreateTicketRequest.
type CreateTicketRequest struct {
	// Note The ticket note
	Note string `json:"note"`

	// Sender Customer email
	Sender string `json:"sender"`

	// Subject The ticket subject
	Subject string `json:"subject"`
}

// CreateTicketResponse defines model for CreateTicketResponse.
type CreateTicketResponse struct {
	UnderscoreLinks SelfLinks                  `json:"_links"`
	Data            []CreateTicketResponseData `json:"data"`
}

// CreateTicketResponseData defines model for CreateTicketResponseData.
type CreateTicketResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	// Email The email of the user to which activation and forgot password links are being sent to. There is a limit of 255 characters per email.
	Email string `json:"email"`

	// Enabled If user is not enabled, he can't login and thus use services any longer.
	Enabled bool `json:"enabled"`

	// FirstName The name of the user. Names may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	FirstName *string `json:"firstName,omitempty"`

	// LastName The last name of the user. Users may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	LastName *string `json:"lastName,omitempty"`

	// Locale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
	Locale CreateUserRequestLocale `json:"locale"`

	// Roles The roles as list of `roleId`s of the user.
	Roles *[]int64 `json:"roles,omitempty"`

	// Totp Enable or disable two-factor authentication (2FA) via time based OTP.
	Totp bool `json:"totp"`
}

// CreateUserRequestLocale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
type CreateUserRequestLocale string

// CreateUserResponse defines model for CreateUserResponse.
type CreateUserResponse struct {
	UnderscoreLinks SelfLinks                `json:"_links"`
	Data            []CreateUserResponseData `json:"data"`
}

// CreateUserResponseData defines model for CreateUserResponseData.
type CreateUserResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// UserId User's id
	UserId string `json:"userId"`
}

// CredentialData defines model for CredentialData.
type CredentialData struct {
	// AccessKey Access key ID.
	AccessKey string `json:"accessKey"`

	// CredentialId Object Storage Credential ID
	CredentialId float32 `json:"credentialId"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DisplayName Object Storage Name.
	DisplayName string `json:"displayName"`

	// ObjectStorageId Object Storage ID.
	ObjectStorageId string `json:"objectStorageId"`

	// Region Object Storage Region.
	Region string `json:"region"`

	// SecretKey Secret key ID.
	SecretKey string `json:"secretKey"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// CustomImagesStatsResponse defines model for CustomImagesStatsResponse.
type CustomImagesStatsResponse struct {
	UnderscoreLinks SelfLinks                       `json:"_links"`
	Data            []CustomImagesStatsResponseData `json:"data"`
}

// CustomImagesStatsResponseData defines model for CustomImagesStatsResponseData.
type CustomImagesStatsResponseData struct {
	// CurrentImagesCount The number of existing custom images
	CurrentImagesCount float32 `json:"currentImagesCount"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// FreeSizeMb Free disk space in MB
	FreeSizeMb float32 `json:"freeSizeMb"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// TotalSizeMb Total available disk space in MB
	TotalSizeMb float32 `json:"totalSizeMb"`

	// UsedSizeMb Used disk space in MB
	UsedSizeMb float32 `json:"usedSizeMb"`
}

// DataCenterResponse defines model for DataCenterResponse.
type DataCenterResponse struct {
	Capabilities []DataCenterResponseCapabilities `json:"capabilities"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Name Name of the data center
	Name string `json:"name"`

	// RegionName Name of the region
	RegionName string `json:"regionName"`

	// RegionSlug Slug of the region
	RegionSlug string `json:"regionSlug"`

	// S3Url S3 URL of the data center
	S3Url string `json:"s3Url"`

	// Slug Slug of the data center
	Slug string `json:"slug"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// DataCenterResponseCapabilities defines model for DataCenterResponse.Capabilities.
type DataCenterResponseCapabilities string

// ExtraStorageRequest defines model for ExtraStorageRequest.
type ExtraStorageRequest struct {
	// Nvme Specify the size in TB and the quantity
	Nvme *[]string `json:"nvme,omitempty"`

	// Ssd Specify the size in TB and the quantity
	Ssd *[]string `json:"ssd,omitempty"`
}

// FindAssignmentResponse defines model for FindAssignmentResponse.
type FindAssignmentResponse struct {
	UnderscoreLinks TagAssignmentSelfLinks `json:"_links"`
	Data            []AssignmentResponse   `json:"data"`
}

// FindClientResponse defines model for FindClientResponse.
type FindClientResponse struct {
	UnderscoreLinks SelfLinks        `json:"_links"`
	Data            []ClientResponse `json:"data"`
}

// FindCredentialResponse defines model for FindCredentialResponse.
type FindCredentialResponse struct {
	UnderscoreLinks SelfLinks        `json:"_links"`
	Data            []CredentialData `json:"data"`
}

// FindImageResponse defines model for FindImageResponse.
type FindImageResponse struct {
	UnderscoreLinks SelfLinks       `json:"_links"`
	Data            []ImageResponse `json:"data"`
}

// FindInstanceResponse defines model for FindInstanceResponse.
type FindInstanceResponse struct {
	UnderscoreLinks SelfLinks          `json:"_links"`
	Data            []InstanceResponse `json:"data"`
}

// FindObjectStorageResponse defines model for FindObjectStorageResponse.
type FindObjectStorageResponse struct {
	UnderscoreLinks SelfLinks               `json:"_links"`
	Data            []ObjectStorageResponse `json:"data"`
}

// FindPrivateNetworkResponse defines model for FindPrivateNetworkResponse.
type FindPrivateNetworkResponse struct {
	UnderscoreLinks SelfLinks                `json:"_links"`
	Data            []PrivateNetworkResponse `json:"data"`
}

// FindRoleResponse defines model for FindRoleResponse.
type FindRoleResponse struct {
	UnderscoreLinks SelfLinks      `json:"_links"`
	Data            []RoleResponse `json:"data"`
}

// FindSecretResponse defines model for FindSecretResponse.
type FindSecretResponse struct {
	UnderscoreLinks SelfLinks        `json:"_links"`
	Data            []SecretResponse `json:"data"`
}

// FindSnapshotResponse defines model for FindSnapshotResponse.
type FindSnapshotResponse struct {
	UnderscoreLinks SelfLinks          `json:"_links"`
	Data            []SnapshotResponse `json:"data"`
}

// FindTagResponse defines model for FindTagResponse.
type FindTagResponse struct {
	UnderscoreLinks SelfLinks     `json:"_links"`
	Data            []TagResponse `json:"data"`
}

// FindUserIsPasswordSetResponse defines model for FindUserIsPasswordSetResponse.
type FindUserIsPasswordSetResponse struct {
	UnderscoreLinks SelfLinks                   `json:"_links"`
	Data            []UserIsPasswordSetResponse `json:"data"`
}

// FindUserResponse defines model for FindUserResponse.
type FindUserResponse struct {
	UnderscoreLinks SelfLinks      `json:"_links"`
	Data            []UserResponse `json:"data"`
}

// FindVipResponse defines model for FindVipResponse.
type FindVipResponse struct {
	UnderscoreLinks SelfLinks     `json:"_links"`
	Data            []VipResponse `json:"data"`
}

// GenerateClientSecretResponse defines model for GenerateClientSecretResponse.
type GenerateClientSecretResponse struct {
	UnderscoreLinks SelfLinks              `json:"_links"`
	Data            []ClientSecretResponse `json:"data"`
}

// ImageAuditResponse defines model for ImageAuditResponse.
type ImageAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta           `json:"_pagination"`
	Data                 []ImageAuditResponseData `json:"data"`
}

// ImageAuditResponseData defines model for ImageAuditResponseData.
type ImageAuditResponseData struct {
	// Action Type of the action.
	Action ImageAuditResponseDataAction `json:"action"`

	// ChangedBy Id of user who performed the change
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Id The ID of the audit entry.
	Id int64 `json:"id"`

	// ImageId The identifier of the image
	ImageId string `json:"imageId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// ImageAuditResponseDataAction Type of the action.
type ImageAuditResponseDataAction string

// ImageResponse defines model for ImageResponse.
type ImageResponse struct {
	// CreationDate The creation date time for the image
	CreationDate int64 `json:"creationDate"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Description Image Description
	Description string `json:"description"`

	// ErrorMessage Image download error message
	ErrorMessage string `json:"errorMessage"`

	// Format Image format
	Format ImageResponseFormat `json:"format"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// LastModifiedDate The last modified date time for the image
	LastModifiedDate int64 `json:"lastModifiedDate"`

	// Name Image Name
	Name string `json:"name"`

	// OsType Type of operating system (OS)
	OsType string `json:"osType"`

	// SizeMb Image Size in MB
	SizeMb float32 `json:"sizeMb"`

	// StandardImage Flag indicating that image is either a standard (true) or a custom image (false)
	StandardImage bool `json:"standardImage"`

	// Status Image status (e.g. if image is still downloading)
	Status string `json:"status"`

	// TenantId Your customer tenant id
	TenantId ImageResponseTenantId `json:"tenantId"`

	// UploadedSizeMb Image Uploaded Size in MB
	UploadedSizeMb float32 `json:"uploadedSizeMb"`

	// Url URL from where the image has been downloaded / provided.
	Url string `json:"url"`

	// Version Version number to distinguish the contents of an image. Could be the version of the operating system for example.
	Version string `json:"version"`
}

// ImageResponseFormat Image format
type ImageResponseFormat string

// ImageResponseTenantId Your customer tenant id
type ImageResponseTenantId string

// InstanceAssignmentSelfLinks defines model for InstanceAssignmentSelfLinks.
type InstanceAssignmentSelfLinks struct {
	// Instance Link to assigned instance.
	Instance string `json:"instance"`

	// Self Link to current resource.
	Self string `json:"self"`

	// VirtualPrivateCloud Link to related Private Network.
	VirtualPrivateCloud string `json:"virtualPrivateCloud"`
}

// InstanceRescueActionResponse defines model for InstanceRescueActionResponse.
type InstanceRescueActionResponse struct {
	UnderscoreLinks SelfLinks                          `json:"_links"`
	Data            []InstanceRescueActionResponseData `json:"data"`
}

// InstanceRescueActionResponseData defines model for InstanceRescueActionResponseData.
type InstanceRescueActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// InstanceResetPasswordActionResponse defines model for InstanceResetPasswordActionResponse.
type InstanceResetPasswordActionResponse struct {
	UnderscoreLinks SelfLinks                                 `json:"_links"`
	Data            []InstanceResetPasswordActionResponseData `json:"data"`
}

// InstanceResetPasswordActionResponseData defines model for InstanceResetPasswordActionResponseData.
type InstanceResetPasswordActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// InstanceResponse defines model for InstanceResponse.
type InstanceResponse struct {
	AddOns        []AddOnResponse `json:"addOns"`
	AdditionalIps []AdditionalIp  `json:"additionalIps"`

	// CancelDate The date on which the instance will be cancelled
	CancelDate string `json:"cancelDate"`

	// CpuCores CPU core count
	CpuCores int64 `json:"cpuCores"`

	// CreatedDate The creation date for the instance
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// DataCenter The data center where your Private Network is located
	DataCenter string `json:"dataCenter"`

	// DefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
	DefaultUser *InstanceResponseDefaultUser `json:"defaultUser,omitempty"`

	// DiskMb Image Disk size in MB
	DiskMb float32 `json:"diskMb"`

	// DisplayName Instance display name
	DisplayName string `json:"displayName"`

	// ErrorMessage Message in case of an error.
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// InstanceId Instance ID
	InstanceId int64    `json:"instanceId"`
	IpConfig   IpConfig `json:"ipConfig"`

	// MacAddress MAC Address
	MacAddress string `json:"macAddress"`

	// Name Instance Name
	Name string `json:"name"`

	// OsType Type of operating system (OS)
	OsType string `json:"osType"`

	// ProductId Product ID
	ProductId string `json:"productId"`

	// ProductName Instance's Product Name
	ProductName string `json:"productName"`

	// ProductType Instance's category depending on Product Id
	ProductType InstanceResponseProductType `json:"productType"`

	// RamMb Image RAM size in MB
	RamMb float32 `json:"ramMb"`

	// Region Instance region where the compute instance should be located.
	Region string `json:"region"`

	// RegionName The name of the region where the instance is located.
	RegionName string `json:"regionName"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys []int64 `json:"sshKeys"`

	// Status Instance's status
	Status InstanceStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId InstanceResponseTenantId `json:"tenantId"`

	// VHostId ID of host system
	VHostId int64 `json:"vHostId"`

	// VHostName Name of host system
	VHostName string `json:"vHostName"`

	// VHostNumber Number of host system
	VHostNumber int64 `json:"vHostNumber"`
}

// InstanceResponseDefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
type InstanceResponseDefaultUser string

// InstanceResponseProductType Instance's category depending on Product Id
type InstanceResponseProductType string

// InstanceResponseTenantId Your customer tenant id
type InstanceResponseTenantId string

// InstanceRestartActionResponse defines model for InstanceRestartActionResponse.
type InstanceRestartActionResponse struct {
	UnderscoreLinks SelfLinks                           `json:"_links"`
	Data            []InstanceRestartActionResponseData `json:"data"`
}

// InstanceRestartActionResponseData defines model for InstanceRestartActionResponseData.
type InstanceRestartActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// InstanceShutdownActionResponse defines model for InstanceShutdownActionResponse.
type InstanceShutdownActionResponse struct {
	UnderscoreLinks SelfLinks                            `json:"_links"`
	Data            []InstanceShutdownActionResponseData `json:"data"`
}

// InstanceShutdownActionResponseData defines model for InstanceShutdownActionResponseData.
type InstanceShutdownActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// InstanceStartActionResponse defines model for InstanceStartActionResponse.
type InstanceStartActionResponse struct {
	UnderscoreLinks SelfLinks                         `json:"_links"`
	Data            []InstanceStartActionResponseData `json:"data"`
}

// InstanceStartActionResponseData defines model for InstanceStartActionResponseData.
type InstanceStartActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// InstanceStopActionResponse defines model for InstanceStopActionResponse.
type InstanceStopActionResponse struct {
	UnderscoreLinks SelfLinks                        `json:"_links"`
	Data            []InstanceStopActionResponseData `json:"data"`
}

// InstanceStopActionResponseData defines model for InstanceStopActionResponseData.
type InstanceStopActionResponseData struct {
	// Action Action that was triggered
	Action string `json:"action"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Compute instance / resource id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// Instances defines model for Instances.
type Instances struct {
	// DisplayName Instance display name
	DisplayName string `json:"displayName"`

	// ErrorMessage Message in case of an error.
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// InstanceId Instance id
	InstanceId int64    `json:"instanceId"`
	IpConfig   IpConfig `json:"ipConfig"`

	// Name Instance name
	Name            string          `json:"name"`
	PrivateIpConfig PrivateIpConfig `json:"privateIpConfig"`

	// ProductId Product id
	ProductId string `json:"productId"`

	// Status State of the instance in the Private Network
	Status InstancesStatus `json:"status"`
}

// InstancesStatus State of the instance in the Private Network
type InstancesStatus string

// InstancesActionsAuditResponse defines model for InstancesActionsAuditResponse.
type InstancesActionsAuditResponse struct {
	// Action Type of the action.
	Action InstancesActionsAuditResponseAction `json:"action"`

	// ChangedBy Id of user who performed the change
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Id The ID of the audit entry.
	Id int64 `json:"id"`

	// InstanceId The identifier of the instancesActions
	InstanceId int64 `json:"instanceId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// InstancesActionsAuditResponseAction Type of the action.
type InstancesActionsAuditResponseAction string

// InstancesActionsRescueRequest defines model for InstancesActionsRescueRequest.
type InstancesActionsRescueRequest struct {
	// RootPassword `secretId` of the password to login into rescue system for the `root` user.
	RootPassword *int64 `json:"rootPassword,omitempty"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into rescue system as `root` user.
	SshKeys *[]int64 `json:"sshKeys,omitempty"`

	// UserData [Cloud-Init](https://cloud-init.io/) Config in order to customize during start of compute instance.
	UserData *string `json:"userData,omitempty"`
}

// InstancesAuditResponse defines model for InstancesAuditResponse.
type InstancesAuditResponse struct {
	// Action Type of the action.
	Action InstancesAuditResponseAction `json:"action"`

	// ChangedBy Id of user who performed the change
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Id The ID of the audit entry.
	Id int64 `json:"id"`

	// InstanceId The identifier of the instance
	InstanceId int64 `json:"instanceId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// InstancesAuditResponseAction Type of the action.
type InstancesAuditResponseAction string

// InstancesResetPasswordActionsRequest defines model for InstancesResetPasswordActionsRequest.
type InstancesResetPasswordActionsRequest struct {
	// RootPassword `secretId` of the password for the `defaultUser` with administrator/root privileges. For Linux/BSD please use SSH, for Windows RDP. Please refer to Secrets Management API.
	RootPassword *int64 `json:"rootPassword,omitempty"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys *[]int64 `json:"sshKeys,omitempty"`

	// UserData [Cloud-Init](https://cloud-init.io/) Config in order to customize during start of compute instance.
	UserData *string `json:"userData,omitempty"`
}

// IpConfig defines model for IpConfig.
type IpConfig struct {
	V4 IpV4 `json:"v4"`
	V6 IpV6 `json:"v6"`
}

// IpV4 defines model for IpV4.
type IpV4 struct {
	// Gateway Gateway
	Gateway string `json:"gateway"`

	// Ip IP Address
	Ip string `json:"ip"`

	// NetmaskCidr Netmask CIDR
	NetmaskCidr int32 `json:"netmaskCidr"`
}

// IpV41 defines model for IpV41.
type IpV41 struct {
	// Broadcast Broadcast address
	Broadcast string `json:"broadcast"`

	// Gateway Gateway
	Gateway string `json:"gateway"`

	// Ip IP address
	Ip string `json:"ip"`

	// Net Net address
	Net string `json:"net"`

	// NetmaskCidr Netmask CIDR
	NetmaskCidr int64 `json:"netmaskCidr"`
}

// IpV6 defines model for IpV6.
type IpV6 struct {
	// Gateway Gateway
	Gateway string `json:"gateway"`

	// Ip IP Address
	Ip string `json:"ip"`

	// NetmaskCidr Netmask CIDR
	NetmaskCidr int32 `json:"netmaskCidr"`
}

// Links defines model for Links.
type Links struct {
	// First Link to first page, if applicable.
	First string `json:"first"`

	// Last Link to last page, if applicable.
	Last string `json:"last"`

	// Next Link to next page, if applicable.
	Next *string `json:"next,omitempty"`

	// Previous Link to previous page, if applicable.
	Previous *string `json:"previous,omitempty"`

	// Self Link to current resource.
	Self string `json:"self"`
}

// ListApiPermissionResponse defines model for ListApiPermissionResponse.
type ListApiPermissionResponse struct {
	UnderscoreLinks Links                    `json:"_links"`
	Data            []ApiPermissionsResponse `json:"data"`
}

// ListAssignmentAuditsResponse defines model for ListAssignmentAuditsResponse.
type ListAssignmentAuditsResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta            `json:"_pagination"`
	Data                 []AssignmentAuditResponse `json:"data"`
}

// ListAssignmentResponse defines model for ListAssignmentResponse.
type ListAssignmentResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta       `json:"_pagination"`
	Data                 []AssignmentResponse `json:"data"`
}

// ListCredentialResponse defines model for ListCredentialResponse.
type ListCredentialResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta   `json:"_pagination"`
	Data                 []CredentialData `json:"data"`
}

// ListDataCenterResponse defines model for ListDataCenterResponse.
type ListDataCenterResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta       `json:"_pagination"`
	Data                 []DataCenterResponse `json:"data"`
}

// ListImageResponse defines model for ListImageResponse.
type ListImageResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta          `json:"_pagination"`
	Data                 []ListImageResponseData `json:"data"`
}

// ListImageResponseData defines model for ListImageResponseData.
type ListImageResponseData struct {
	// CreationDate The creation date time for the image
	CreationDate int64 `json:"creationDate"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Description Image Description
	Description string `json:"description"`

	// ErrorMessage Image download error message
	ErrorMessage string `json:"errorMessage"`

	// Format Image format
	Format ListImageResponseDataFormat `json:"format"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// LastModifiedDate The last modified date time for the image
	LastModifiedDate int64 `json:"lastModifiedDate"`

	// Name Image Name
	Name string `json:"name"`

	// OsType Type of operating system (OS)
	OsType string `json:"osType"`

	// SizeMb Image Size in MB
	SizeMb float32 `json:"sizeMb"`

	// StandardImage Flag indicating that image is either a standard (true) or a custom image (false)
	StandardImage bool `json:"standardImage"`

	// Status Image status (e.g. if image is still downloading)
	Status string `json:"status"`

	// Tags The tags assigned to the image
	Tags []AssignedTagResponse `json:"tags"`

	// TenantId Your customer tenant id
	TenantId ListImageResponseDataTenantId `json:"tenantId"`

	// UploadedSizeMb Image Uploaded Size in MB
	UploadedSizeMb float32 `json:"uploadedSizeMb"`

	// Url URL from where the image has been downloaded / provided.
	Url string `json:"url"`

	// Version Version number to distinguish the contents of an image. Could be the version of the operating system for example.
	Version string `json:"version"`
}

// ListImageResponseDataFormat Image format
type ListImageResponseDataFormat string

// ListImageResponseDataTenantId Your customer tenant id
type ListImageResponseDataTenantId string

// ListInstancesActionsAuditResponse defines model for ListInstancesActionsAuditResponse.
type ListInstancesActionsAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta                  `json:"_pagination"`
	Data                 []InstancesActionsAuditResponse `json:"data"`
}

// ListInstancesAuditResponse defines model for ListInstancesAuditResponse.
type ListInstancesAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta           `json:"_pagination"`
	Data                 []InstancesAuditResponse `json:"data"`
}

// ListInstancesResponse defines model for ListInstancesResponse.
type ListInstancesResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta              `json:"_pagination"`
	Data                 []ListInstancesResponseData `json:"data"`
}

// ListInstancesResponseData defines model for ListInstancesResponseData.
type ListInstancesResponseData struct {
	AddOns        []AddOnResponse `json:"addOns"`
	AdditionalIps []AdditionalIp  `json:"additionalIps"`

	// CancelDate The date on which the instance will be cancelled
	CancelDate string `json:"cancelDate"`

	// CpuCores CPU core count
	CpuCores int64 `json:"cpuCores"`

	// CreatedDate The creation date for the instance
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// DataCenter The data center where your Private Network is located
	DataCenter string `json:"dataCenter"`

	// DefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
	DefaultUser *ListInstancesResponseDataDefaultUser `json:"defaultUser,omitempty"`

	// DiskMb Image Disk size in MB
	DiskMb float32 `json:"diskMb"`

	// DisplayName Instance display name
	DisplayName string `json:"displayName"`

	// ErrorMessage Message in case of an error.
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// InstanceId Instance ID
	InstanceId int64    `json:"instanceId"`
	IpConfig   IpConfig `json:"ipConfig"`

	// MacAddress MAC Address
	MacAddress string `json:"macAddress"`

	// Name Instance Name
	Name string `json:"name"`

	// OsType Type of operating system (OS)
	OsType string `json:"osType"`

	// ProductId Product ID
	ProductId string `json:"productId"`

	// ProductName Instance's Product Name
	ProductName string `json:"productName"`

	// ProductType Instance's category depending on Product Id
	ProductType ListInstancesResponseDataProductType `json:"productType"`

	// RamMb Image RAM size in MB
	RamMb float32 `json:"ramMb"`

	// Region Instance region where the compute instance should be located.
	Region string `json:"region"`

	// RegionName The name of the region where the instance is located.
	RegionName string `json:"regionName"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys []int64 `json:"sshKeys"`

	// Status Instance's status
	Status InstanceStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId ListInstancesResponseDataTenantId `json:"tenantId"`

	// VHostId ID of host system
	VHostId int64 `json:"vHostId"`

	// VHostName Name of host system
	VHostName string `json:"vHostName"`

	// VHostNumber Number of host system
	VHostNumber int64 `json:"vHostNumber"`
}

// ListInstancesResponseDataDefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
type ListInstancesResponseDataDefaultUser string

// ListInstancesResponseDataProductType Instance's category depending on Product Id
type ListInstancesResponseDataProductType string

// ListInstancesResponseDataTenantId Your customer tenant id
type ListInstancesResponseDataTenantId string

// ListObjectStorageAuditResponse defines model for ListObjectStorageAuditResponse.
type ListObjectStorageAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta               `json:"_pagination"`
	Data                 []ObjectStorageAuditResponse `json:"data"`
}

// ListObjectStorageResponse defines model for ListObjectStorageResponse.
type ListObjectStorageResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta          `json:"_pagination"`
	Data                 []ObjectStorageResponse `json:"data"`
}

// ListPrivateNetworkAuditResponse defines model for ListPrivateNetworkAuditResponse.
type ListPrivateNetworkAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta                `json:"_pagination"`
	Data                 []PrivateNetworkAuditResponse `json:"data"`
}

// ListPrivateNetworkResponse defines model for ListPrivateNetworkResponse.
type ListPrivateNetworkResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta                   `json:"_pagination"`
	Data                 []ListPrivateNetworkResponseData `json:"data"`
}

// ListPrivateNetworkResponseData defines model for ListPrivateNetworkResponseData.
type ListPrivateNetworkResponseData struct {
	// AvailableIps The total available IPs of the Private Network
	AvailableIps int64 `json:"availableIps"`

	// Cidr The cidr range of the Private Network
	Cidr string `json:"cidr"`

	// CreatedDate The creation date of the Private Network
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DataCenter The data center where your Private Network is located
	DataCenter string `json:"dataCenter"`

	// Description The description of the Private Network
	Description string      `json:"description"`
	Instances   []Instances `json:"instances"`

	// Name The name of the Private Network
	Name string `json:"name"`

	// PrivateNetworkId Private Network's id
	PrivateNetworkId int64 `json:"privateNetworkId"`

	// Region The slug of the region where your Private Network is located
	Region string `json:"region"`

	// RegionName The region where your Private Network is located
	RegionName string `json:"regionName"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// ListRoleAuditResponse defines model for ListRoleAuditResponse.
type ListRoleAuditResponse struct {
	UnderscoreLinks Links               `json:"_links"`
	Data            []RoleAuditResponse `json:"data"`
}

// ListRoleResponse defines model for ListRoleResponse.
type ListRoleResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta `json:"_pagination"`
	Data                 []RoleResponse `json:"data"`
}

// ListSecretAuditResponse defines model for ListSecretAuditResponse.
type ListSecretAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta        `json:"_pagination"`
	Data                 []SecretAuditResponse `json:"data"`
}

// ListSecretResponse defines model for ListSecretResponse.
type ListSecretResponse struct {
	UnderscoreLinks SelfLinks `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta   `json:"_pagination"`
	Data                 []SecretResponse `json:"data"`
}

// ListSnapshotResponse defines model for ListSnapshotResponse.
type ListSnapshotResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta     `json:"_pagination"`
	Data                 []SnapshotResponse `json:"data"`
}

// ListSnapshotsAuditResponse defines model for ListSnapshotsAuditResponse.
type ListSnapshotsAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta           `json:"_pagination"`
	Data                 []SnapshotsAuditResponse `json:"data"`
}

// ListTagAuditsResponse defines model for ListTagAuditsResponse.
type ListTagAuditsResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta     `json:"_pagination"`
	Data                 []TagAuditResponse `json:"data"`
}

// ListTagResponse defines model for ListTagResponse.
type ListTagResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta `json:"_pagination"`
	Data                 []TagResponse  `json:"data"`
}

// ListUserAuditResponse defines model for ListUserAuditResponse.
type ListUserAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta      `json:"_pagination"`
	Data                 []UserAuditResponse `json:"data"`
}

// ListUserResponse defines model for ListUserResponse.
type ListUserResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta `json:"_pagination"`
	Data                 []UserResponse `json:"data"`
}

// ListVipAuditResponse defines model for ListVipAuditResponse.
type ListVipAuditResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta     `json:"_pagination"`
	Data                 []VipAuditResponse `json:"data"`
}

// ListVipResponse defines model for ListVipResponse.
type ListVipResponse struct {
	UnderscoreLinks Links `json:"_links"`

	// UnderscorePagination Data about pagination like how many results, pages, page size.
	UnderscorePagination PaginationMeta        `json:"_pagination"`
	Data                 []ListVipResponseData `json:"data"`
}

// ListVipResponseData defines model for ListVipResponseData.
type ListVipResponseData struct {
	// CustomerId Customer's Id.
	CustomerId string `json:"customerId"`

	// DataCenter data center.
	DataCenter string `json:"dataCenter"`

	// IpVersion Version of Ip.
	IpVersion ListVipResponseDataIpVersion `json:"ipVersion"`

	// Region Region
	Region string `json:"region"`

	// ResourceDisplayName Resource display name.
	ResourceDisplayName string `json:"resourceDisplayName"`

	// ResourceId Resource Id.
	ResourceId string `json:"resourceId"`

	// ResourceName Resource name.
	ResourceName string `json:"resourceName"`

	// ResourceType The resourceType using the VIP.
	ResourceType *ListVipResponseDataResourceType `json:"resourceType,omitempty"`

	// TenantId Tenant Id.
	TenantId string `json:"tenantId"`

	// Type The VIP type.
	Type *ListVipResponseDataType `json:"type,omitempty"`
	V4   *IpV41                   `json:"v4,omitempty"`

	// VipId Vip uuid.
	VipId string `json:"vipId"`
}

// ListVipResponseDataIpVersion Version of Ip.
type ListVipResponseDataIpVersion string

// ListVipResponseDataResourceType The resourceType using the VIP.
type ListVipResponseDataResourceType string

// ListVipResponseDataType The VIP type.
type ListVipResponseDataType string

// ObjectStorageAuditResponse defines model for ObjectStorageAuditResponse.
type ObjectStorageAuditResponse struct {
	// Action Type of the action.
	Action ObjectStorageAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id int64 `json:"id"`

	// ObjectStorageId Object Storage Id
	ObjectStorageId string `json:"objectStorageId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// ObjectStorageAuditResponseAction Type of the action.
type ObjectStorageAuditResponseAction string

// ObjectStorageResponse defines model for ObjectStorageResponse.
type ObjectStorageResponse struct {
	// AutoScaling Autoscaling settings
	AutoScaling AutoScalingTypeResponse `json:"autoScaling"`

	// CancelDate Cancellation date for object storage.
	CancelDate string `json:"cancelDate"`

	// CreatedDate Creation date for object storage.
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DataCenter Data center your object storage is located
	DataCenter string `json:"dataCenter"`

	// DisplayName Display name for object storage.
	DisplayName string `json:"displayName"`

	// ObjectStorageId Your object storage id
	ObjectStorageId string `json:"objectStorageId"`

	// Region The region where your object storage is located
	Region string `json:"region"`

	// S3TenantId Your S3 tenantId. Only required for public sharing.
	S3TenantId string `json:"s3TenantId"`

	// S3Url S3 URL to connect to your S3 compatible object storage
	S3Url string `json:"s3Url"`

	// Status The object storage status
	Status ObjectStorageResponseStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// TotalPurchasedSpaceTB Amount of purchased / requested object storage in TB.
	TotalPurchasedSpaceTB float64 `json:"totalPurchasedSpaceTB"`
}

// ObjectStorageResponseStatus The object storage status
type ObjectStorageResponseStatus string

// ObjectStoragesStatsResponse defines model for ObjectStoragesStatsResponse.
type ObjectStoragesStatsResponse struct {
	UnderscoreLinks SelfLinks                         `json:"_links"`
	Data            []ObjectStoragesStatsResponseData `json:"data"`
}

// ObjectStoragesStatsResponseData defines model for ObjectStoragesStatsResponseData.
type ObjectStoragesStatsResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// NumberOfObjects Number of all objects (i.e. files and folders) in object storage.
	NumberOfObjects int64 `json:"numberOfObjects"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// UsedSpacePercentage Currently used space in percentage.
	UsedSpacePercentage float64 `json:"usedSpacePercentage"`

	// UsedSpaceTB Currently used space in TB.
	UsedSpaceTB float64 `json:"usedSpaceTB"`
}

// PaginationMeta defines model for PaginationMeta.
type PaginationMeta struct {
	// Page Current number of page.
	Page float32 `json:"page"`

	// Size Number of elements per page.
	Size float32 `json:"size"`

	// TotalElements Number of overall matched elements.
	TotalElements float32 `json:"totalElements"`

	// TotalPages Overall number of pages.
	TotalPages float32 `json:"totalPages"`
}

// PatchInstanceRequest defines model for PatchInstanceRequest.
type PatchInstanceRequest struct {
	// DisplayName The display name of the instance
	DisplayName *string `json:"displayName,omitempty"`
}

// PatchInstanceResponse defines model for PatchInstanceResponse.
type PatchInstanceResponse struct {
	UnderscoreLinks SelfLinks                   `json:"_links"`
	Data            []PatchInstanceResponseData `json:"data"`
}

// PatchInstanceResponseData defines model for PatchInstanceResponseData.
type PatchInstanceResponseData struct {
	// CreatedDate Creation date of the instance
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Instance's id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// PatchObjectStorageRequest defines model for PatchObjectStorageRequest.
type PatchObjectStorageRequest struct {
	// DisplayName Display name helps to differentiate between object storages, especially if they are in the same region.
	DisplayName string `json:"displayName"`
}

// PatchPrivateNetworkRequest defines model for PatchPrivateNetworkRequest.
type PatchPrivateNetworkRequest struct {
	// Description The description of the Private Network. There is a limit of 255 characters per Private Network.
	Description *string `json:"description,omitempty"`

	// Name The name of the Private Network. It may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per Private Network.
	Name *string `json:"name,omitempty"`
}

// PatchPrivateNetworkResponse defines model for PatchPrivateNetworkResponse.
type PatchPrivateNetworkResponse struct {
	UnderscoreLinks SelfLinks                `json:"_links"`
	Data            []PrivateNetworkResponse `json:"data"`
}

// PermissionRequest defines model for PermissionRequest.
type PermissionRequest struct {
	// Actions Action allowed for the API endpoint. Basically `CREATE` corresponds to POST endpoints, `READ` to GET endpoints, `UPDATE` to PATCH / PUT endpoints and `DELETE` to DELETE endpoints.
	Actions []PermissionRequestActions `json:"actions"`

	// ApiName The name of the role. There is a limit of 255 characters per role.
	ApiName string `json:"apiName"`

	// Resources The IDs of tags. Only if those tags are assgined to a resource the user with that role will be able to access the resource.
	Resources *[]int64 `json:"resources,omitempty"`
}

// PermissionRequestActions defines model for PermissionRequest.Actions.
type PermissionRequestActions string

// PermissionResponse defines model for PermissionResponse.
type PermissionResponse struct {
	// Actions Action allowed for the API endpoint. Basically `CREATE` corresponds to POST endpoints, `READ` to GET endpoints, `UPDATE` to PATCH / PUT endpoints and `DELETE` to DELETE endpoints.
	Actions []PermissionResponseActions `json:"actions"`

	// ApiName API endpoint. In order to get a list availbale api enpoints please refer to the GET api-permissions endpoint.
	ApiName   string                         `json:"apiName"`
	Resources *[]ResourcePermissionsResponse `json:"resources,omitempty"`
}

// PermissionResponseActions defines model for PermissionResponse.Actions.
type PermissionResponseActions string

// PrivateIpConfig defines model for PrivateIpConfig.
type PrivateIpConfig struct {
	V4 []IpV4 `json:"v4"`
}

// PrivateNetworkAuditResponse defines model for PrivateNetworkAuditResponse.
type PrivateNetworkAuditResponse struct {
	// Action Type of the action.
	Action PrivateNetworkAuditResponseAction `json:"action"`

	// ChangedBy User id
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id int64 `json:"id"`

	// PrivateNetworkId The identifier of the Private Network
	PrivateNetworkId float32 `json:"privateNetworkId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username User name which did the change.
	Username string `json:"username"`
}

// PrivateNetworkAuditResponseAction Type of the action.
type PrivateNetworkAuditResponseAction string

// PrivateNetworkResponse defines model for PrivateNetworkResponse.
type PrivateNetworkResponse struct {
	// AvailableIps The total available IPs of the Private Network
	AvailableIps int64 `json:"availableIps"`

	// Cidr The cidr range of the Private Network
	Cidr string `json:"cidr"`

	// CreatedDate The creation date of the Private Network
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DataCenter The data center where your Private Network is located
	DataCenter string `json:"dataCenter"`

	// Description The description of the Private Network
	Description string      `json:"description"`
	Instances   []Instances `json:"instances"`

	// Name The name of the Private Network
	Name string `json:"name"`

	// PrivateNetworkId Private Network's id
	PrivateNetworkId int64 `json:"privateNetworkId"`

	// Region The slug of the region where your Private Network is located
	Region string `json:"region"`

	// RegionName The region where your Private Network is located
	RegionName string `json:"regionName"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// PrivateNetworkingUpgradeRequest defines model for PrivateNetworkingUpgradeRequest.
type PrivateNetworkingUpgradeRequest = map[string]string

// ReinstallInstanceRequest defines model for ReinstallInstanceRequest.
type ReinstallInstanceRequest struct {
	// ApplicationId Application ID
	ApplicationId *string `json:"applicationId,omitempty"`

	// DefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
	DefaultUser *ReinstallInstanceRequestDefaultUser `json:"defaultUser,omitempty"`

	// ImageId ImageId to be used to setup the compute instance.
	ImageId string `json:"imageId"`

	// RootPassword `secretId` of the password for the `defaultUser` with administrator/root privileges. For Linux/BSD please use SSH, for Windows RDP. Please refer to Secrets Management API.
	RootPassword *int64 `json:"rootPassword,omitempty"`

	// SshKeys Array of `secretId`s of public SSH keys for logging into as `defaultUser` with administrator/root privileges. Applies to Linux/BSD systems. Please refer to Secrets Management API.
	SshKeys *[]int64 `json:"sshKeys,omitempty"`

	// UserData [Cloud-Init](https://cloud-init.io/) Config in order to customize during start of compute instance.
	UserData *string `json:"userData,omitempty"`
}

// ReinstallInstanceRequestDefaultUser Default user name created for login during (re-)installation with administrative privileges. Allowed values for Linux/BSD are `admin` (use sudo to apply administrative privileges like root) or `root`. Allowed values for Windows are `admin` (has administrative privileges like administrator) or `administrator`.
type ReinstallInstanceRequestDefaultUser string

// ReinstallInstanceResponse defines model for ReinstallInstanceResponse.
type ReinstallInstanceResponse struct {
	UnderscoreLinks SelfLinks                       `json:"_links"`
	Data            []ReinstallInstanceResponseData `json:"data"`
}

// ReinstallInstanceResponseData defines model for ReinstallInstanceResponseData.
type ReinstallInstanceResponseData struct {
	// CreatedDate Creation date for instance
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// InstanceId Instance's id
	InstanceId int64 `json:"instanceId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// ResourcePermissionsResponse defines model for ResourcePermissionsResponse.
type ResourcePermissionsResponse struct {
	// TagId Tag's id
	TagId int64 `json:"tagId"`

	// TagName Tag name. The resriction is based on the resources which have been assigned to that tag. If no resource has been assigned to the given tag, no access will be possible.
	TagName string `json:"tagName"`
}

// RoleAuditResponse defines model for RoleAuditResponse.
type RoleAuditResponse struct {
	// Action Type of the action.
	Action RoleAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id int64 `json:"id"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// RoleId The identifier of the role
	RoleId float32 `json:"roleId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// RoleAuditResponseAction Type of the action.
type RoleAuditResponseAction string

// RoleResponse defines model for RoleResponse.
type RoleResponse struct {
	// AccessAllResources Access All Resources
	AccessAllResources bool `json:"accessAllResources"`

	// Admin Admin
	Admin bool `json:"admin"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Name Role's name
	Name        string                `json:"name"`
	Permissions *[]PermissionResponse `json:"permissions,omitempty"`

	// RoleId Role's id
	RoleId int64 `json:"roleId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// Type Role type can be either `default` or `custom`. The `default` roles cannot be modified or deleted.
	Type string `json:"type"`
}

// RollbackSnapshotRequest defines model for RollbackSnapshotRequest.
type RollbackSnapshotRequest = map[string]string

// RollbackSnapshotResponse defines model for RollbackSnapshotResponse.
type RollbackSnapshotResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks SelfLinks `json:"_links"`
}

// SecretAuditResponse defines model for SecretAuditResponse.
type SecretAuditResponse struct {
	// Action Type of the action.
	Action SecretAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id float32 `json:"id"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// SecretId Secret's id
	SecretId float32 `json:"secretId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// SecretAuditResponseAction Type of the action.
type SecretAuditResponseAction string

// SecretResponse defines model for SecretResponse.
type SecretResponse struct {
	// CreatedAt The creation date for the secret
	CreatedAt int64 `json:"createdAt"`

	// CustomerId Your Customer number
	CustomerId string `json:"customerId"`

	// Name The name assigned to the password/ssh
	Name string `json:"name"`

	// SecretId Secret's id
	SecretId float32 `json:"secretId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// Type The type of the secret. This will be available only when retrieving secrets
	Type SecretResponseType `json:"type"`

	// UpdatedAt The last update date for the secret
	UpdatedAt int64 `json:"updatedAt"`

	// Value The value of the secret. This will be available only when retrieving a single secret
	Value string `json:"value"`
}

// SecretResponseType The type of the secret. This will be available only when retrieving secrets
type SecretResponseType string

// SelfLinks defines model for SelfLinks.
type SelfLinks struct {
	// Self Link to current resource.
	Self string `json:"self"`
}

// SnapshotResponse defines model for SnapshotResponse.
type SnapshotResponse struct {
	// AutoDeleteDate The date when the snapshot will be auto-deleted
	AutoDeleteDate int64 `json:"autoDeleteDate"`

	// CreatedDate The date when the snapshot was created
	CreatedDate int64 `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Description The description of the snapshot.
	Description string `json:"description"`

	// ImageId Image Id the snapshot was taken on
	ImageId string `json:"imageId"`

	// ImageName Image name the snapshot was taken on
	ImageName string `json:"imageName"`

	// InstanceId The instance identifier associated with the snapshot
	InstanceId int64 `json:"instanceId"`

	// Name The name of the snapshot.
	Name string `json:"name"`

	// SnapshotId Snapshot's id
	SnapshotId string `json:"snapshotId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// SnapshotsAuditResponse defines model for SnapshotsAuditResponse.
type SnapshotsAuditResponse struct {
	// Action Type of the action.
	Action SnapshotsAuditResponseAction `json:"action"`

	// ChangedBy Id of user who performed the change
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer ID
	CustomerId string `json:"customerId"`

	// Id The ID of the audit entry.
	Id int64 `json:"id"`

	// InstanceId The identifier of the instance
	InstanceId int64 `json:"instanceId"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// SnapshotId The identifier of the snapshot
	SnapshotId string `json:"snapshotId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// SnapshotsAuditResponseAction Type of the action.
type SnapshotsAuditResponseAction string

// TagAssignmentSelfLinks defines model for TagAssignmentSelfLinks.
type TagAssignmentSelfLinks struct {
	// UnderscoreResource Link to assigned resource
	UnderscoreResource string `json:"_resource"`

	// Self Link to current resource.
	Self string `json:"self"`

	// Tag Link to related tag.
	Tag string `json:"tag"`
}

// TagAuditResponse defines model for TagAuditResponse.
type TagAuditResponse struct {
	// Action Type of the action.
	Action TagAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id float32 `json:"id"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TagId The identifier of the audit entry.
	TagId int64 `json:"tagId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// TagAuditResponseAction Type of the action.
type TagAuditResponseAction string

// TagResponse defines model for TagResponse.
type TagResponse struct {
	// Color Tag's color
	Color string `json:"color"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Description The description of the Tag
	Description string `json:"description"`

	// Name Tag's name
	Name string `json:"name"`

	// TagId Tag's id
	TagId int64 `json:"tagId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// UnassignInstancePrivateNetworkResponse defines model for UnassignInstancePrivateNetworkResponse.
type UnassignInstancePrivateNetworkResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks InstanceAssignmentSelfLinks `json:"_links"`
}

// UpdateCustomImageRequest defines model for UpdateCustomImageRequest.
type UpdateCustomImageRequest struct {
	// Description Image Description
	Description *string `json:"description,omitempty"`

	// Name Image Name
	Name *string `json:"name,omitempty"`
}

// UpdateCustomImageResponse defines model for UpdateCustomImageResponse.
type UpdateCustomImageResponse struct {
	UnderscoreLinks SelfLinks                       `json:"_links"`
	Data            []UpdateCustomImageResponseData `json:"data"`
}

// UpdateCustomImageResponseData defines model for UpdateCustomImageResponseData.
type UpdateCustomImageResponseData struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// ImageId Image's id
	ImageId string `json:"imageId"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// UpdateRoleRequest defines model for UpdateRoleRequest.
type UpdateRoleRequest struct {
	// AccessAllResources Allow access to all resources. This will superseed all assigned resources in a role.
	AccessAllResources bool `json:"accessAllResources"`

	// Admin If user is admin he will have permissions to all API endpoints and resources. Enabling this will superseed all role definitions and `accessAllResources`.
	Admin bool `json:"admin"`

	// Name The name of the role. There is a limit of 255 characters per role.
	Name        string               `json:"name"`
	Permissions *[]PermissionRequest `json:"permissions,omitempty"`
}

// UpdateRoleResponse defines model for UpdateRoleResponse.
type UpdateRoleResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks SelfLinks `json:"_links"`
}

// UpdateSecretRequest defines model for UpdateSecretRequest.
type UpdateSecretRequest struct {
	// Name The name of the secret to be saved
	Name *string `json:"name,omitempty"`

	// Value The value of the secret to be saved
	Value *string `json:"value,omitempty"`
}

// UpdateSecretResponse defines model for UpdateSecretResponse.
type UpdateSecretResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks SelfLinks `json:"_links"`
}

// UpdateSnapshotRequest defines model for UpdateSnapshotRequest.
type UpdateSnapshotRequest struct {
	// Description The description of the snapshot. There is a limit of 255 characters per snapshot.
	Description *string `json:"description,omitempty"`

	// Name The name of the snapshot. Tags may contain only letters, numbers, spaces, dashes. There is a limit of 30 characters per snapshot.
	Name *string `json:"name,omitempty"`
}

// UpdateSnapshotResponse defines model for UpdateSnapshotResponse.
type UpdateSnapshotResponse struct {
	UnderscoreLinks SelfLinks          `json:"_links"`
	Data            []SnapshotResponse `json:"data"`
}

// UpdateTagRequest defines model for UpdateTagRequest.
type UpdateTagRequest struct {
	// Color The color of the tag. Color can be specified using hexadecimal value. Default color is #0A78C3
	Color *string `json:"color,omitempty"`

	// Description The description of the Tag name.
	Description *string `json:"description,omitempty"`

	// Name The name of the tag. Tags may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per tag.
	Name *string `json:"name,omitempty"`
}

// UpdateTagResponse defines model for UpdateTagResponse.
type UpdateTagResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks SelfLinks `json:"_links"`
}

// UpdateUserRequest defines model for UpdateUserRequest.
type UpdateUserRequest struct {
	// Email The email of the user to which activation and forgot password links are being sent to. There is a limit of 255 characters per email.
	Email *string `json:"email,omitempty"`

	// Enabled If user is not enabled, he can't login and thus use services any longer.
	Enabled *bool `json:"enabled,omitempty"`

	// FirstName The name of the user. Names may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	FirstName *string `json:"firstName,omitempty"`

	// LastName The last name of the user. Users may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	LastName *string `json:"lastName,omitempty"`

	// Locale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
	Locale *UpdateUserRequestLocale `json:"locale,omitempty"`

	// Roles The roles as list of `roleId`s of the user.
	Roles *[]int64 `json:"roles,omitempty"`

	// Totp Enable or disable two-factor authentication (2FA) via time based OTP.
	Totp *bool `json:"totp,omitempty"`
}

// UpdateUserRequestLocale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
type UpdateUserRequestLocale string

// UpdateUserResponse defines model for UpdateUserResponse.
type UpdateUserResponse struct {
	// UnderscoreLinks Links for easy navigation.
	UnderscoreLinks SelfLinks `json:"_links"`
}

// UpgradeAutoScalingType defines model for UpgradeAutoScalingType.
type UpgradeAutoScalingType struct {
	// SizeLimitTB Autoscaling size limit for the current object storage.
	SizeLimitTB *float64 `json:"sizeLimitTB,omitempty"`

	// State State of the autoscaling for the current object storage.
	State *UpgradeAutoScalingTypeState `json:"state,omitempty"`
}

// UpgradeAutoScalingTypeState State of the autoscaling for the current object storage.
type UpgradeAutoScalingTypeState string

// UpgradeInstanceRequest defines model for UpgradeInstanceRequest.
type UpgradeInstanceRequest struct {
	// Backup Set this attribute if you want to upgrade your instance with the Automated Backup addon.   Please provide an empty object for the time being as value. There will be more configuration possible   in the future.
	Backup *Backup `json:"backup,omitempty"`

	// PrivateNetworking Set this attribute if you want to upgrade your instance with the Private Networking addon. Please provide an empty object for the time being as value. There will be more configuration possible in the future.
	PrivateNetworking *PrivateNetworkingUpgradeRequest `json:"privateNetworking,omitempty"`
}

// UpgradeObjectStorageRequest defines model for UpgradeObjectStorageRequest.
type UpgradeObjectStorageRequest struct {
	// AutoScaling New monthly object storage size limit for autoscaling if enabled.
	AutoScaling *UpgradeAutoScalingType `json:"autoScaling,omitempty"`

	// TotalPurchasedSpaceTB New total object storage limit. If this number is larger than before you will also be billed for the added storage space. No downgrade possible.
	TotalPurchasedSpaceTB *float64 `json:"totalPurchasedSpaceTB,omitempty"`
}

// UpgradeObjectStorageResponse defines model for UpgradeObjectStorageResponse.
type UpgradeObjectStorageResponse struct {
	UnderscoreLinks SelfLinks                          `json:"_links"`
	Data            []UpgradeObjectStorageResponseData `json:"data"`
}

// UpgradeObjectStorageResponseData defines model for UpgradeObjectStorageResponseData.
type UpgradeObjectStorageResponseData struct {
	// AutoScaling The autoscaling limit of the object storage.
	AutoScaling AutoScalingTypeResponse `json:"autoScaling"`

	// CreatedDate Creation date for object storage.
	CreatedDate string `json:"createdDate"`

	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// DataCenter Data center of the object storage.
	DataCenter string `json:"dataCenter"`

	// DisplayName Display name for object storage.
	DisplayName string `json:"displayName"`

	// ObjectStorageId Object storage id
	ObjectStorageId string `json:"objectStorageId"`

	// Region The region where your object storage is located
	Region string `json:"region"`

	// S3Url S3 URL to connect to your S3 compatible object storage
	S3Url string `json:"s3Url"`

	// Status The object storage status
	Status UpgradeObjectStorageResponseDataStatus `json:"status"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// TotalPurchasedSpaceTB Total purchased object storage space in TB.
	TotalPurchasedSpaceTB float64 `json:"totalPurchasedSpaceTB"`
}

// UpgradeObjectStorageResponseDataStatus The object storage status
type UpgradeObjectStorageResponseDataStatus string

// UserAuditResponse defines model for UserAuditResponse.
type UserAuditResponse struct {
	// Action Type of the action.
	Action UserAuditResponseAction `json:"action"`

	// ChangedBy User ID
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id int64 `json:"id"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// UserId The identifier of the user
	UserId string `json:"userId"`

	// Username Name of the user which led to the change.
	Username string `json:"username"`
}

// UserAuditResponseAction Type of the action.
type UserAuditResponseAction string

// UserIsPasswordSetResponse defines model for UserIsPasswordSetResponse.
type UserIsPasswordSetResponse struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// IsPasswordSet Indicates if the user has set a password for his account
	IsPasswordSet bool `json:"isPasswordSet"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	// CustomerId Your customer number
	CustomerId string `json:"customerId"`

	// Email The email of the user to which activation and forgot password links are being sent to. There is a limit of 255 characters per email.
	Email string `json:"email"`

	// EmailVerified User email verification status.
	EmailVerified bool `json:"emailVerified"`

	// Enabled If uses is not enabled, he can't login and thus use services any longer.
	Enabled bool `json:"enabled"`

	// FirstName The first name of the user. Users may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	FirstName string `json:"firstName"`

	// LastName The last name of the user. Users may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per user.
	LastName string `json:"lastName"`

	// Locale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
	Locale UserResponseLocale `json:"locale"`

	// Owner If user is owner he will have permissions to all API endpoints and resources. Enabling this will superseed all role definitions and `accessAllResources`.
	Owner bool `json:"owner"`

	// Roles The roles as list of `roleId`s of the user.
	Roles []RoleResponse `json:"roles"`

	// TenantId Your customer tenant id
	TenantId string `json:"tenantId"`

	// Totp Enable or disable two-factor authentication (2FA) via time based OTP.
	Totp bool `json:"totp"`

	// UserId The identifier of the user.
	UserId string `json:"userId"`
}

// UserResponseLocale The locale of the user. This can be `de-DE`, `de`, `en-US`, `en`, `es-ES`, `es`, `pt-BR`, `pt`.
type UserResponseLocale string

// VipAuditResponse defines model for VipAuditResponse.
type VipAuditResponse struct {
	// Action Type of the action.
	Action VipAuditResponseAction `json:"action"`

	// ChangedBy User id
	ChangedBy string `json:"changedBy"`

	// Changes List of actual changes.
	Changes *map[string]string `json:"changes,omitempty"`

	// CustomerId Customer number
	CustomerId string `json:"customerId"`

	// Id The identifier of the audit entry.
	Id int64 `json:"id"`

	// RequestId The requestId of the API call which led to the change.
	RequestId string `json:"requestId"`

	// TenantId Customer tenant id
	TenantId string `json:"tenantId"`

	// Timestamp When the change took place.
	Timestamp int64 `json:"timestamp"`

	// TraceId The traceId of the API call which led to the change.
	TraceId string `json:"traceId"`

	// Username User name which did the change.
	Username string `json:"username"`

	// VipId The identifier of the VIP
	VipId string `json:"vipId"`
}

// VipAuditResponseAction Type of the action.
type VipAuditResponseAction string

// VipResponse defines model for VipResponse.
type VipResponse struct {
	// CustomerId Customer's Id.
	CustomerId string `json:"customerId"`

	// DataCenter data center.
	DataCenter string `json:"dataCenter"`

	// IpVersion Version of Ip.
	IpVersion VipResponseIpVersion `json:"ipVersion"`

	// Region Region
	Region string `json:"region"`

	// ResourceDisplayName Resource display name.
	ResourceDisplayName string `json:"resourceDisplayName"`

	// ResourceId Resource Id.
	ResourceId string `json:"resourceId"`

	// ResourceName Resource name.
	ResourceName string `json:"resourceName"`

	// ResourceType The resourceType using the VIP.
	ResourceType *VipResponseResourceType `json:"resourceType,omitempty"`

	// TenantId Tenant Id.
	TenantId string `json:"tenantId"`

	// Type The VIP type.
	Type *VipResponseType `json:"type,omitempty"`
	V4   *IpV41           `json:"v4,omitempty"`

	// VipId Vip uuid.
	VipId string `json:"vipId"`
}

// VipResponseIpVersion Version of Ip.
type VipResponseIpVersion string

// VipResponseResourceType The resourceType using the VIP.
type VipResponseResourceType string

// VipResponseType The VIP type.
type VipResponseType string

// InstanceStatus Instance's status
type InstanceStatus string

// RetrieveImageListParams defines parameters for RetrieveImageList.
type RetrieveImageListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name The name of the image
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// StandardImage Flag indicating that image is either a standard (true) or a custom image (false)
	StandardImage *bool `form:"standardImage,omitempty" json:"standardImage,omitempty"`

	// Search full text search on image name or image os type
	Search *string `form:"search,omitempty" json:"search,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateCustomImageParams defines parameters for CreateCustomImage.
type CreateCustomImageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveImageAuditsListParams defines parameters for RetrieveImageAuditsList.
type RetrieveImageAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ImageId The identifier of the image.
	ImageId *string `form:"imageId,omitempty" json:"imageId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy UserId of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveCustomImagesStatsParams defines parameters for RetrieveCustomImagesStats.
type RetrieveCustomImagesStatsParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteImageParams defines parameters for DeleteImage.
type DeleteImageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveImageParams defines parameters for RetrieveImage.
type RetrieveImageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateImageParams defines parameters for UpdateImage.
type UpdateImageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveInstancesListParams defines parameters for RetrieveInstancesList.
type RetrieveInstancesListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name The name of the instance
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// DisplayName The display name of the instance
	DisplayName *string `form:"displayName,omitempty" json:"displayName,omitempty"`

	// DataCenter The data center of the instance
	DataCenter *string `form:"dataCenter,omitempty" json:"dataCenter,omitempty"`

	// Region The Region of the instance
	Region *string `form:"region,omitempty" json:"region,omitempty"`

	// InstanceId The identifier of the instance (deprecated)
	InstanceId *int64 `form:"instanceId,omitempty" json:"instanceId,omitempty"`

	// InstanceIds Comma separated instances identifiers
	InstanceIds *string `form:"instanceIds,omitempty" json:"instanceIds,omitempty"`

	// Status The status of the instance
	Status *RetrieveInstancesListParamsStatus `form:"status,omitempty" json:"status,omitempty"`

	// ProductIds Identifiers of the instance products
	ProductIds *string `form:"productIds,omitempty" json:"productIds,omitempty"`

	// AddOnIds Identifiers of Addons the instances have
	AddOnIds *string `form:"addOnIds,omitempty" json:"addOnIds,omitempty"`

	// ProductTypes Comma separated instance's category depending on Product Id
	ProductTypes *string `form:"productTypes,omitempty" json:"productTypes,omitempty"`

	// IpConfig Filter instances that have an ip config
	IpConfig *bool `form:"ipConfig,omitempty" json:"ipConfig,omitempty"`

	// Search Full text search when listing the instances. Can be searched by `name`, `displayName`, `ipAddress`
	Search *string `form:"search,omitempty" json:"search,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveInstancesListParamsStatus defines parameters for RetrieveInstancesList.
type RetrieveInstancesListParamsStatus string

// CreateInstanceParams defines parameters for CreateInstance.
type CreateInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveInstancesActionsAuditsListParams defines parameters for RetrieveInstancesActionsAuditsList.
type RetrieveInstancesActionsAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// InstanceId The identifier of the instancesActions.
	InstanceId *int64 `form:"instanceId,omitempty" json:"instanceId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveInstancesAuditsListParams defines parameters for RetrieveInstancesAuditsList.
type RetrieveInstancesAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// InstanceId The identifier of the instances.
	InstanceId *int64 `form:"instanceId,omitempty" json:"instanceId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveInstanceParams defines parameters for RetrieveInstance.
type RetrieveInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// PatchInstanceParams defines parameters for PatchInstance.
type PatchInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ReinstallInstanceParams defines parameters for ReinstallInstance.
type ReinstallInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RescueParams defines parameters for Rescue.
type RescueParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ResetPasswordActionParams defines parameters for ResetPasswordAction.
type ResetPasswordActionParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RestartParams defines parameters for Restart.
type RestartParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ShutdownParams defines parameters for Shutdown.
type ShutdownParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// StartParams defines parameters for Start.
type StartParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// StopParams defines parameters for Stop.
type StopParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CancelInstanceParams defines parameters for CancelInstance.
type CancelInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSnapshotListParams defines parameters for RetrieveSnapshotList.
type RetrieveSnapshotListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name Filter as substring match for snapshots names.
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateSnapshotParams defines parameters for CreateSnapshot.
type CreateSnapshotParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteSnapshotParams defines parameters for DeleteSnapshot.
type DeleteSnapshotParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSnapshotParams defines parameters for RetrieveSnapshot.
type RetrieveSnapshotParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateSnapshotParams defines parameters for UpdateSnapshot.
type UpdateSnapshotParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RollbackSnapshotParams defines parameters for RollbackSnapshot.
type RollbackSnapshotParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpgradeInstanceParams defines parameters for UpgradeInstance.
type UpgradeInstanceParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSnapshotsAuditsListParams defines parameters for RetrieveSnapshotsAuditsList.
type RetrieveSnapshotsAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// InstanceId The identifier of the instance
	InstanceId *int64 `form:"instanceId,omitempty" json:"instanceId,omitempty"`

	// SnapshotId The identifier of the snapshot
	SnapshotId *string `form:"snapshotId,omitempty" json:"snapshotId,omitempty"`

	// RequestId The requestId of the API call which led to the change
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateTicketParams defines parameters for CreateTicket.
type CreateTicketParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveDataCenterListParams defines parameters for RetrieveDataCenterList.
type RetrieveDataCenterListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Slug Filter as match for data centers.
	Slug *string `form:"slug,omitempty" json:"slug,omitempty"`

	// Name Filter for Object Storages regions.
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// RegionName Filter for Object Storage region names.
	RegionName *string `form:"regionName,omitempty" json:"regionName,omitempty"`

	// RegionSlug Filter for Object Storage region slugs.
	RegionSlug *string `form:"regionSlug,omitempty" json:"regionSlug,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveObjectStorageListParams defines parameters for RetrieveObjectStorageList.
type RetrieveObjectStorageListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// DataCenterName Filter for Object Storage locations.
	DataCenterName *string `form:"dataCenterName,omitempty" json:"dataCenterName,omitempty"`

	// S3TenantId Filter for Object Storage S3 tenantId.
	S3TenantId *string `form:"s3TenantId,omitempty" json:"s3TenantId,omitempty"`

	// Region Filter for Object Storage by regions. Available regions: EU, US-central, SIN
	Region *string `form:"region,omitempty" json:"region,omitempty"`

	// DisplayName Filter for Object Storage by display name.
	DisplayName *string `form:"displayName,omitempty" json:"displayName,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateObjectStorageParams defines parameters for CreateObjectStorage.
type CreateObjectStorageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveObjectStorageAuditsListParams defines parameters for RetrieveObjectStorageAuditsList.
type RetrieveObjectStorageAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ObjectStorageId The identifier of the object storage.
	ObjectStorageId *string `form:"objectStorageId,omitempty" json:"objectStorageId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveObjectStorageParams defines parameters for RetrieveObjectStorage.
type RetrieveObjectStorageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateObjectStorageParams defines parameters for UpdateObjectStorage.
type UpdateObjectStorageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CancelObjectStorageParams defines parameters for CancelObjectStorage.
type CancelObjectStorageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpgradeObjectStorageParams defines parameters for UpgradeObjectStorage.
type UpgradeObjectStorageParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveObjectStoragesStatsParams defines parameters for RetrieveObjectStoragesStats.
type RetrieveObjectStoragesStatsParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrievePrivateNetworkListParams defines parameters for RetrievePrivateNetworkList.
type RetrievePrivateNetworkListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name The name of the Private Network
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// InstanceIds Comma separated instances identifiers
	InstanceIds *string `form:"instanceIds,omitempty" json:"instanceIds,omitempty"`

	// Region The slug of the region where your Private Network is located
	Region *string `form:"region,omitempty" json:"region,omitempty"`

	// DataCenter The data center where your Private Network is located
	DataCenter *string `form:"dataCenter,omitempty" json:"dataCenter,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreatePrivateNetworkParams defines parameters for CreatePrivateNetwork.
type CreatePrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrievePrivateNetworkAuditsListParams defines parameters for RetrievePrivateNetworkAuditsList.
type RetrievePrivateNetworkAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// PrivateNetworkId The identifier of the Private Network
	PrivateNetworkId *int64 `form:"privateNetworkId,omitempty" json:"privateNetworkId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy User name which did the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeletePrivateNetworkParams defines parameters for DeletePrivateNetwork.
type DeletePrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrievePrivateNetworkParams defines parameters for RetrievePrivateNetwork.
type RetrievePrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// PatchPrivateNetworkParams defines parameters for PatchPrivateNetwork.
type PatchPrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UnassignInstancePrivateNetworkParams defines parameters for UnassignInstancePrivateNetwork.
type UnassignInstancePrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// AssignInstancePrivateNetworkParams defines parameters for AssignInstancePrivateNetwork.
type AssignInstancePrivateNetworkParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveRoleListParams defines parameters for RetrieveRoleList.
type RetrieveRoleListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name The name of the role
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// ApiName The name of api
	ApiName *string `form:"apiName,omitempty" json:"apiName,omitempty"`

	// TagName The name of the tag
	TagName *string `form:"tagName,omitempty" json:"tagName,omitempty"`

	// Type The type of the tag. Can be either `default` or `custom`
	Type *string `form:"type,omitempty" json:"type,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateRoleParams defines parameters for CreateRole.
type CreateRoleParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveApiPermissionsListParams defines parameters for RetrieveApiPermissionsList.
type RetrieveApiPermissionsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ApiName The name of api
	ApiName *string `form:"apiName,omitempty" json:"apiName,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveRoleAuditsListParams defines parameters for RetrieveRoleAuditsList.
type RetrieveRoleAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// RoleId The identifier of the role.
	RoleId *int64 `form:"roleId,omitempty" json:"roleId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteRoleParams defines parameters for DeleteRole.
type DeleteRoleParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveRoleParams defines parameters for RetrieveRole.
type RetrieveRoleParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateRoleParams defines parameters for UpdateRole.
type UpdateRoleParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSecretListParams defines parameters for RetrieveSecretList.
type RetrieveSecretListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name Filter secrets by name
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Type Filter secrets by type
	Type *RetrieveSecretListParamsType `form:"type,omitempty" json:"type,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSecretListParamsType defines parameters for RetrieveSecretList.
type RetrieveSecretListParamsType string

// CreateSecretParams defines parameters for CreateSecret.
type CreateSecretParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSecretAuditsListParams defines parameters for RetrieveSecretAuditsList.
type RetrieveSecretAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// SecretId The id of the secret.
	SecretId *int64 `form:"secretId,omitempty" json:"secretId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteSecretParams defines parameters for DeleteSecret.
type DeleteSecretParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveSecretParams defines parameters for RetrieveSecret.
type RetrieveSecretParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateSecretParams defines parameters for UpdateSecret.
type UpdateSecretParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveTagListParams defines parameters for RetrieveTagList.
type RetrieveTagListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Name Filter as substring match for tag names. Tags may contain letters, numbers, colons, dashes, and underscores. There is a limit of 255 characters per tag.
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateTagParams defines parameters for CreateTag.
type CreateTagParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveAssignmentsAuditsListParams defines parameters for RetrieveAssignmentsAuditsList.
type RetrieveAssignmentsAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// TagId The identifier of the tag.
	TagId *int64 `form:"tagId,omitempty" json:"tagId,omitempty"`

	// ResourceId The identifier of the resource.
	ResourceId *string `form:"resourceId,omitempty" json:"resourceId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy UserId of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveTagAuditsListParams defines parameters for RetrieveTagAuditsList.
type RetrieveTagAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// TagId The identifier of the tag.
	TagId *int64 `form:"tagId,omitempty" json:"tagId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy UserId of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteTagParams defines parameters for DeleteTag.
type DeleteTagParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveTagParams defines parameters for RetrieveTag.
type RetrieveTagParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateTagParams defines parameters for UpdateTag.
type UpdateTagParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveAssignmentListParams defines parameters for RetrieveAssignmentList.
type RetrieveAssignmentListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ResourceType Filter as substring match for assignment resource type. Resource type is one of `instance|image|object-storage`.
	ResourceType *string `form:"resourceType,omitempty" json:"resourceType,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteAssignmentParams defines parameters for DeleteAssignment.
type DeleteAssignmentParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveAssignmentParams defines parameters for RetrieveAssignment.
type RetrieveAssignmentParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateAssignmentParams defines parameters for CreateAssignment.
type CreateAssignmentParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveUserListParams defines parameters for RetrieveUserList.
type RetrieveUserListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// Email Filter as substring match for user emails.
	Email *string `form:"email,omitempty" json:"email,omitempty"`

	// Enabled Filter if user is enabled or not.
	Enabled *bool `form:"enabled,omitempty" json:"enabled,omitempty"`

	// Owner Filter if user is owner or not.
	Owner *bool `form:"owner,omitempty" json:"owner,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// CreateUserParams defines parameters for CreateUser.
type CreateUserParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveUserAuditsListParams defines parameters for RetrieveUserAuditsList.
type RetrieveUserAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// UserId The identifier of the user.
	UserId *string `form:"userId,omitempty" json:"userId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy changedBy of the user which led to the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveUserClientParams defines parameters for RetrieveUserClient.
type RetrieveUserClientParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// GenerateClientSecretParams defines parameters for GenerateClientSecret.
type GenerateClientSecretParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveUserIsPasswordSetParams defines parameters for RetrieveUserIsPasswordSet.
type RetrieveUserIsPasswordSetParams struct {
	// UserId The user ID for checking if password is set for him
	UserId *string `form:"userId,omitempty" json:"userId,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// DeleteUserParams defines parameters for DeleteUser.
type DeleteUserParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveUserParams defines parameters for RetrieveUser.
type RetrieveUserParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UpdateUserParams defines parameters for UpdateUser.
type UpdateUserParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ListObjectStorageCredentialsParams defines parameters for ListObjectStorageCredentials.
type ListObjectStorageCredentialsParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ObjectStorageId The identifier of the S3 object storage
	ObjectStorageId *string `form:"objectStorageId,omitempty" json:"objectStorageId,omitempty"`

	// RegionName Filter for Object Storage by regions. Available regions: Asia (Singapore), European Union, United States (Central)
	RegionName *string `form:"regionName,omitempty" json:"regionName,omitempty"`

	// DisplayName Filter for Object Storage by his displayName.
	DisplayName *string `form:"displayName,omitempty" json:"displayName,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// GetObjectStorageCredentialsParams defines parameters for GetObjectStorageCredentials.
type GetObjectStorageCredentialsParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RegenerateObjectStorageCredentialsParams defines parameters for RegenerateObjectStorageCredentials.
type RegenerateObjectStorageCredentialsParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ResendEmailVerificationParams defines parameters for ResendEmailVerification.
type ResendEmailVerificationParams struct {
	// RedirectUrl The redirect url used for email verification
	RedirectUrl *string `form:"redirectUrl,omitempty" json:"redirectUrl,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// ResetPasswordParams defines parameters for ResetPassword.
type ResetPasswordParams struct {
	// RedirectUrl The redirect url used for resetting password
	RedirectUrl *string `form:"redirectUrl,omitempty" json:"redirectUrl,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveVipListParams defines parameters for RetrieveVipList.
type RetrieveVipListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// ResourceId The resourceId using the VIP.
	ResourceId *string `form:"resourceId,omitempty" json:"resourceId,omitempty"`

	// ResourceType The resourceType using the VIP.
	ResourceType *RetrieveVipListParamsResourceType `form:"resourceType,omitempty" json:"resourceType,omitempty"`

	// ResourceName The name of the resource.
	ResourceName *string `form:"resourceName,omitempty" json:"resourceName,omitempty"`

	// ResourceDisplayName The display name of the resource.
	ResourceDisplayName *string `form:"resourceDisplayName,omitempty" json:"resourceDisplayName,omitempty"`

	// IpVersion The VIP version.
	IpVersion *RetrieveVipListParamsIpVersion `form:"ipVersion,omitempty" json:"ipVersion,omitempty"`

	// Ips Comma separated IPs
	Ips *string `form:"ips,omitempty" json:"ips,omitempty"`

	// Ip The ip of the VIP
	Ip *string `form:"ip,omitempty" json:"ip,omitempty"`

	// Type The VIP type.
	Type *RetrieveVipListParamsType `form:"type,omitempty" json:"type,omitempty"`

	// DataCenter The dataCenter of the VIP.
	DataCenter *string `form:"dataCenter,omitempty" json:"dataCenter,omitempty"`

	// Region The region of the VIP.
	Region *string `form:"region,omitempty" json:"region,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveVipListParamsResourceType defines parameters for RetrieveVipList.
type RetrieveVipListParamsResourceType string

// RetrieveVipListParamsIpVersion defines parameters for RetrieveVipList.
type RetrieveVipListParamsIpVersion string

// RetrieveVipListParamsType defines parameters for RetrieveVipList.
type RetrieveVipListParamsType string

// RetrieveVipAuditsListParams defines parameters for RetrieveVipAuditsList.
type RetrieveVipAuditsListParams struct {
	// Page Number of page to be fetched.
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of elements per page.
	Size *int64 `form:"size,omitempty" json:"size,omitempty"`

	// OrderBy Specify fields and ordering (ASC for ascending, DESC for descending) in following format `field:ASC|DESC`.
	OrderBy *[]string `form:"orderBy,omitempty" json:"orderBy,omitempty"`

	// VipId The identifier of the VIP.
	VipId *string `form:"vipId,omitempty" json:"vipId,omitempty"`

	// RequestId The requestId of the API call which led to the change.
	RequestId *string `form:"requestId,omitempty" json:"requestId,omitempty"`

	// ChangedBy User name which did the change.
	ChangedBy *string `form:"changedBy,omitempty" json:"changedBy,omitempty"`

	// StartDate Start of search time range.
	StartDate *string `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate End of search time range.
	EndDate *string `form:"endDate,omitempty" json:"endDate,omitempty"`

	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// RetrieveVipParams defines parameters for RetrieveVip.
type RetrieveVipParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UnassignIpParams defines parameters for UnassignIp.
type UnassignIpParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// UnassignIpParamsResourceType defines parameters for UnassignIp.
type UnassignIpParamsResourceType string

// AssignIpParams defines parameters for AssignIp.
type AssignIpParams struct {
	// XRequestId [Uuid4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) to identify individual requests for support cases. You can use [uuidgenerator](https://www.uuidgenerator.net/version4) to generate them manually.
	XRequestId string `json:"x-request-id"`

	// XTraceId Identifier to trace group of requests.
	XTraceId *string `json:"x-trace-id,omitempty"`
}

// AssignIpParamsResourceType defines parameters for AssignIp.
type AssignIpParamsResourceType string

// CreateCustomImageJSONRequestBody defines body for CreateCustomImage for application/json ContentType.
type CreateCustomImageJSONRequestBody = CreateCustomImageRequest

// UpdateImageJSONRequestBody defines body for UpdateImage for application/json ContentType.
type UpdateImageJSONRequestBody = UpdateCustomImageRequest

// CreateInstanceJSONRequestBody defines body for CreateInstance for application/json ContentType.
type CreateInstanceJSONRequestBody = CreateInstanceRequest

// PatchInstanceJSONRequestBody defines body for PatchInstance for application/json ContentType.
type PatchInstanceJSONRequestBody = PatchInstanceRequest

// ReinstallInstanceJSONRequestBody defines body for ReinstallInstance for application/json ContentType.
type ReinstallInstanceJSONRequestBody = ReinstallInstanceRequest

// RescueJSONRequestBody defines body for Rescue for application/json ContentType.
type RescueJSONRequestBody = InstancesActionsRescueRequest

// ResetPasswordActionJSONRequestBody defines body for ResetPasswordAction for application/json ContentType.
type ResetPasswordActionJSONRequestBody = InstancesResetPasswordActionsRequest

// CancelInstanceJSONRequestBody defines body for CancelInstance for application/json ContentType.
type CancelInstanceJSONRequestBody = CancelInstanceRequest

// CreateSnapshotJSONRequestBody defines body for CreateSnapshot for application/json ContentType.
type CreateSnapshotJSONRequestBody = CreateSnapshotRequest

// UpdateSnapshotJSONRequestBody defines body for UpdateSnapshot for application/json ContentType.
type UpdateSnapshotJSONRequestBody = UpdateSnapshotRequest

// RollbackSnapshotJSONRequestBody defines body for RollbackSnapshot for application/json ContentType.
type RollbackSnapshotJSONRequestBody = RollbackSnapshotRequest

// UpgradeInstanceJSONRequestBody defines body for UpgradeInstance for application/json ContentType.
type UpgradeInstanceJSONRequestBody = UpgradeInstanceRequest

// CreateTicketJSONRequestBody defines body for CreateTicket for application/json ContentType.
type CreateTicketJSONRequestBody = CreateTicketRequest

// CreateObjectStorageJSONRequestBody defines body for CreateObjectStorage for application/json ContentType.
type CreateObjectStorageJSONRequestBody = CreateObjectStorageRequest

// UpdateObjectStorageJSONRequestBody defines body for UpdateObjectStorage for application/json ContentType.
type UpdateObjectStorageJSONRequestBody = PatchObjectStorageRequest

// CancelObjectStorageJSONRequestBody defines body for CancelObjectStorage for application/json ContentType.
type CancelObjectStorageJSONRequestBody = CancelObjectStorageRequest

// UpgradeObjectStorageJSONRequestBody defines body for UpgradeObjectStorage for application/json ContentType.
type UpgradeObjectStorageJSONRequestBody = UpgradeObjectStorageRequest

// CreatePrivateNetworkJSONRequestBody defines body for CreatePrivateNetwork for application/json ContentType.
type CreatePrivateNetworkJSONRequestBody = CreatePrivateNetworkRequest

// PatchPrivateNetworkJSONRequestBody defines body for PatchPrivateNetwork for application/json ContentType.
type PatchPrivateNetworkJSONRequestBody = PatchPrivateNetworkRequest

// CreateRoleJSONRequestBody defines body for CreateRole for application/json ContentType.
type CreateRoleJSONRequestBody = CreateRoleRequest

// UpdateRoleJSONRequestBody defines body for UpdateRole for application/json ContentType.
type UpdateRoleJSONRequestBody = UpdateRoleRequest

// CreateSecretJSONRequestBody defines body for CreateSecret for application/json ContentType.
type CreateSecretJSONRequestBody = CreateSecretRequest

// UpdateSecretJSONRequestBody defines body for UpdateSecret for application/json ContentType.
type UpdateSecretJSONRequestBody = UpdateSecretRequest

// CreateTagJSONRequestBody defines body for CreateTag for application/json ContentType.
type CreateTagJSONRequestBody = CreateTagRequest

// UpdateTagJSONRequestBody defines body for UpdateTag for application/json ContentType.
type UpdateTagJSONRequestBody = UpdateTagRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserRequest

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = UpdateUserRequest
