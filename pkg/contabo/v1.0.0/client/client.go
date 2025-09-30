// Package client provides primitives to interact with the openapi HTTP API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	models "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// RetrieveImageList request
	RetrieveImageList(ctx context.Context, params *models.RetrieveImageListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateCustomImageWithBody request with any body
	CreateCustomImageWithBody(ctx context.Context, params *models.CreateCustomImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateCustomImage(ctx context.Context, params *models.CreateCustomImageParams, body models.CreateCustomImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveImageAuditsList request
	RetrieveImageAuditsList(ctx context.Context, params *models.RetrieveImageAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveCustomImagesStats request
	RetrieveCustomImagesStats(ctx context.Context, params *models.RetrieveCustomImagesStatsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteImage request
	DeleteImage(ctx context.Context, imageId string, params *models.DeleteImageParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveImage request
	RetrieveImage(ctx context.Context, imageId string, params *models.RetrieveImageParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateImageWithBody request with any body
	UpdateImageWithBody(ctx context.Context, imageId string, params *models.UpdateImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateImage(ctx context.Context, imageId string, params *models.UpdateImageParams, body models.UpdateImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveInstancesList request
	RetrieveInstancesList(ctx context.Context, params *models.RetrieveInstancesListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateInstanceWithBody request with any body
	CreateInstanceWithBody(ctx context.Context, params *models.CreateInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateInstance(ctx context.Context, params *models.CreateInstanceParams, body models.CreateInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveInstancesActionsAuditsList request
	RetrieveInstancesActionsAuditsList(ctx context.Context, params *models.RetrieveInstancesActionsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveInstancesAuditsList request
	RetrieveInstancesAuditsList(ctx context.Context, params *models.RetrieveInstancesAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveInstance request
	RetrieveInstance(ctx context.Context, instanceId int64, params *models.RetrieveInstanceParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PatchInstanceWithBody request with any body
	PatchInstanceWithBody(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PatchInstance(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, body models.PatchInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ReinstallInstanceWithBody request with any body
	ReinstallInstanceWithBody(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ReinstallInstance(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, body models.ReinstallInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RescueWithBody request with any body
	RescueWithBody(ctx context.Context, instanceId int64, params *models.RescueParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Rescue(ctx context.Context, instanceId int64, params *models.RescueParams, body models.RescueJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResetPasswordActionWithBody request with any body
	ResetPasswordActionWithBody(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ResetPasswordAction(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, body models.ResetPasswordActionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Restart request
	Restart(ctx context.Context, instanceId int64, params *models.RestartParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Shutdown request
	Shutdown(ctx context.Context, instanceId int64, params *models.ShutdownParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Start request
	Start(ctx context.Context, instanceId int64, params *models.StartParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Stop request
	Stop(ctx context.Context, instanceId int64, params *models.StopParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CancelInstanceWithBody request with any body
	CancelInstanceWithBody(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CancelInstance(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, body models.CancelInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSnapshotList request
	RetrieveSnapshotList(ctx context.Context, instanceId int64, params *models.RetrieveSnapshotListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateSnapshotWithBody request with any body
	CreateSnapshotWithBody(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSnapshot(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, body models.CreateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteSnapshot request
	DeleteSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.DeleteSnapshotParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSnapshot request
	RetrieveSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.RetrieveSnapshotParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSnapshotWithBody request with any body
	UpdateSnapshotWithBody(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, body models.UpdateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RollbackSnapshotWithBody request with any body
	RollbackSnapshotWithBody(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RollbackSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, body models.RollbackSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpgradeInstanceWithBody request with any body
	UpgradeInstanceWithBody(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpgradeInstance(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, body models.UpgradeInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSnapshotsAuditsList request
	RetrieveSnapshotsAuditsList(ctx context.Context, params *models.RetrieveSnapshotsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTicketWithBody request with any body
	CreateTicketWithBody(ctx context.Context, params *models.CreateTicketParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTicket(ctx context.Context, params *models.CreateTicketParams, body models.CreateTicketJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveDataCenterList request
	RetrieveDataCenterList(ctx context.Context, params *models.RetrieveDataCenterListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveObjectStorageList request
	RetrieveObjectStorageList(ctx context.Context, params *models.RetrieveObjectStorageListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateObjectStorageWithBody request with any body
	CreateObjectStorageWithBody(ctx context.Context, params *models.CreateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateObjectStorage(ctx context.Context, params *models.CreateObjectStorageParams, body models.CreateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveObjectStorageAuditsList request
	RetrieveObjectStorageAuditsList(ctx context.Context, params *models.RetrieveObjectStorageAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveObjectStorage request
	RetrieveObjectStorage(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStorageParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateObjectStorageWithBody request with any body
	UpdateObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateObjectStorage(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, body models.UpdateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CancelObjectStorageWithBody request with any body
	CancelObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CancelObjectStorage(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, body models.CancelObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpgradeObjectStorageWithBody request with any body
	UpgradeObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpgradeObjectStorage(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, body models.UpgradeObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveObjectStoragesStats request
	RetrieveObjectStoragesStats(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStoragesStatsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrievePrivateNetworkList request
	RetrievePrivateNetworkList(ctx context.Context, params *models.RetrievePrivateNetworkListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreatePrivateNetworkWithBody request with any body
	CreatePrivateNetworkWithBody(ctx context.Context, params *models.CreatePrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreatePrivateNetwork(ctx context.Context, params *models.CreatePrivateNetworkParams, body models.CreatePrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrievePrivateNetworkAuditsList request
	RetrievePrivateNetworkAuditsList(ctx context.Context, params *models.RetrievePrivateNetworkAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeletePrivateNetwork request
	DeletePrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.DeletePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrievePrivateNetwork request
	RetrievePrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.RetrievePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PatchPrivateNetworkWithBody request with any body
	PatchPrivateNetworkWithBody(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PatchPrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, body models.PatchPrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UnassignInstancePrivateNetwork request
	UnassignInstancePrivateNetwork(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.UnassignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AssignInstancePrivateNetwork request
	AssignInstancePrivateNetwork(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.AssignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveRoleList request
	RetrieveRoleList(ctx context.Context, params *models.RetrieveRoleListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateRoleWithBody request with any body
	CreateRoleWithBody(ctx context.Context, params *models.CreateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateRole(ctx context.Context, params *models.CreateRoleParams, body models.CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveApiPermissionsList request
	RetrieveApiPermissionsList(ctx context.Context, params *models.RetrieveApiPermissionsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveRoleAuditsList request
	RetrieveRoleAuditsList(ctx context.Context, params *models.RetrieveRoleAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteRole request
	DeleteRole(ctx context.Context, roleId int64, params *models.DeleteRoleParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveRole request
	RetrieveRole(ctx context.Context, roleId int64, params *models.RetrieveRoleParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateRoleWithBody request with any body
	UpdateRoleWithBody(ctx context.Context, roleId int64, params *models.UpdateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateRole(ctx context.Context, roleId int64, params *models.UpdateRoleParams, body models.UpdateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSecretList request
	RetrieveSecretList(ctx context.Context, params *models.RetrieveSecretListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateSecretWithBody request with any body
	CreateSecretWithBody(ctx context.Context, params *models.CreateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSecret(ctx context.Context, params *models.CreateSecretParams, body models.CreateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSecretAuditsList request
	RetrieveSecretAuditsList(ctx context.Context, params *models.RetrieveSecretAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteSecret request
	DeleteSecret(ctx context.Context, secretId int64, params *models.DeleteSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveSecret request
	RetrieveSecret(ctx context.Context, secretId int64, params *models.RetrieveSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSecretWithBody request with any body
	UpdateSecretWithBody(ctx context.Context, secretId int64, params *models.UpdateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSecret(ctx context.Context, secretId int64, params *models.UpdateSecretParams, body models.UpdateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveTagList request
	RetrieveTagList(ctx context.Context, params *models.RetrieveTagListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTagWithBody request with any body
	CreateTagWithBody(ctx context.Context, params *models.CreateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTag(ctx context.Context, params *models.CreateTagParams, body models.CreateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveAssignmentsAuditsList request
	RetrieveAssignmentsAuditsList(ctx context.Context, params *models.RetrieveAssignmentsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveTagAuditsList request
	RetrieveTagAuditsList(ctx context.Context, params *models.RetrieveTagAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTag request
	DeleteTag(ctx context.Context, tagId int64, params *models.DeleteTagParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveTag request
	RetrieveTag(ctx context.Context, tagId int64, params *models.RetrieveTagParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTagWithBody request with any body
	UpdateTagWithBody(ctx context.Context, tagId int64, params *models.UpdateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTag(ctx context.Context, tagId int64, params *models.UpdateTagParams, body models.UpdateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveAssignmentList request
	RetrieveAssignmentList(ctx context.Context, tagId int64, params *models.RetrieveAssignmentListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteAssignment request
	DeleteAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.DeleteAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveAssignment request
	RetrieveAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.RetrieveAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateAssignment request
	CreateAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.CreateAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveUserList request
	RetrieveUserList(ctx context.Context, params *models.RetrieveUserListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateUserWithBody request with any body
	CreateUserWithBody(ctx context.Context, params *models.CreateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateUser(ctx context.Context, params *models.CreateUserParams, body models.CreateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveUserAuditsList request
	RetrieveUserAuditsList(ctx context.Context, params *models.RetrieveUserAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveUserClient request
	RetrieveUserClient(ctx context.Context, params *models.RetrieveUserClientParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GenerateClientSecret request
	GenerateClientSecret(ctx context.Context, params *models.GenerateClientSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveUserIsPasswordSet request
	RetrieveUserIsPasswordSet(ctx context.Context, params *models.RetrieveUserIsPasswordSetParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteUser request
	DeleteUser(ctx context.Context, userId string, params *models.DeleteUserParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveUser request
	RetrieveUser(ctx context.Context, userId string, params *models.RetrieveUserParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateUserWithBody request with any body
	UpdateUserWithBody(ctx context.Context, userId string, params *models.UpdateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateUser(ctx context.Context, userId string, params *models.UpdateUserParams, body models.UpdateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListObjectStorageCredentials request
	ListObjectStorageCredentials(ctx context.Context, userId string, params *models.ListObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetObjectStorageCredentials request
	GetObjectStorageCredentials(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.GetObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RegenerateObjectStorageCredentials request
	RegenerateObjectStorageCredentials(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.RegenerateObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResendEmailVerification request
	ResendEmailVerification(ctx context.Context, userId string, params *models.ResendEmailVerificationParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResetPassword request
	ResetPassword(ctx context.Context, userId string, params *models.ResetPasswordParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveVipList request
	RetrieveVipList(ctx context.Context, params *models.RetrieveVipListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveVipAuditsList request
	RetrieveVipAuditsList(ctx context.Context, params *models.RetrieveVipAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RetrieveVip request
	RetrieveVip(ctx context.Context, ip string, params *models.RetrieveVipParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UnassignIp request
	UnassignIp(ctx context.Context, ip string, resourceType models.UnassignIpParamsResourceType, resourceId int64, params *models.UnassignIpParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AssignIp request
	AssignIp(ctx context.Context, ip string, resourceType models.AssignIpParamsResourceType, resourceId int64, params *models.AssignIpParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) RetrieveImageList(ctx context.Context, params *models.RetrieveImageListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveImageListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateCustomImageWithBody(ctx context.Context, params *models.CreateCustomImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateCustomImageRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateCustomImage(ctx context.Context, params *models.CreateCustomImageParams, body models.CreateCustomImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateCustomImageRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveImageAuditsList(ctx context.Context, params *models.RetrieveImageAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveImageAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveCustomImagesStats(ctx context.Context, params *models.RetrieveCustomImagesStatsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveCustomImagesStatsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteImage(ctx context.Context, imageId string, params *models.DeleteImageParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteImageRequest(c.Server, imageId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveImage(ctx context.Context, imageId string, params *models.RetrieveImageParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveImageRequest(c.Server, imageId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateImageWithBody(ctx context.Context, imageId string, params *models.UpdateImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateImageRequestWithBody(c.Server, imageId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateImage(ctx context.Context, imageId string, params *models.UpdateImageParams, body models.UpdateImageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateImageRequest(c.Server, imageId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveInstancesList(ctx context.Context, params *models.RetrieveInstancesListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveInstancesListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateInstanceWithBody(ctx context.Context, params *models.CreateInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateInstanceRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateInstance(ctx context.Context, params *models.CreateInstanceParams, body models.CreateInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateInstanceRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveInstancesActionsAuditsList(ctx context.Context, params *models.RetrieveInstancesActionsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveInstancesActionsAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveInstancesAuditsList(ctx context.Context, params *models.RetrieveInstancesAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveInstancesAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveInstance(ctx context.Context, instanceId int64, params *models.RetrieveInstanceParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveInstanceRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchInstanceWithBody(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchInstanceRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchInstance(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, body models.PatchInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchInstanceRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReinstallInstanceWithBody(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReinstallInstanceRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReinstallInstance(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, body models.ReinstallInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReinstallInstanceRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RescueWithBody(ctx context.Context, instanceId int64, params *models.RescueParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRescueRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Rescue(ctx context.Context, instanceId int64, params *models.RescueParams, body models.RescueJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRescueRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ResetPasswordActionWithBody(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewResetPasswordActionRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ResetPasswordAction(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, body models.ResetPasswordActionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewResetPasswordActionRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Restart(ctx context.Context, instanceId int64, params *models.RestartParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRestartRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Shutdown(ctx context.Context, instanceId int64, params *models.ShutdownParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewShutdownRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Start(ctx context.Context, instanceId int64, params *models.StartParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStartRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Stop(ctx context.Context, instanceId int64, params *models.StopParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStopRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelInstanceWithBody(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelInstanceRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelInstance(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, body models.CancelInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelInstanceRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSnapshotList(ctx context.Context, instanceId int64, params *models.RetrieveSnapshotListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSnapshotListRequest(c.Server, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSnapshotWithBody(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSnapshotRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSnapshot(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, body models.CreateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSnapshotRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.DeleteSnapshotParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteSnapshotRequest(c.Server, instanceId, snapshotId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.RetrieveSnapshotParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSnapshotRequest(c.Server, instanceId, snapshotId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateSnapshotWithBody(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateSnapshotRequestWithBody(c.Server, instanceId, snapshotId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, body models.UpdateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateSnapshotRequest(c.Server, instanceId, snapshotId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RollbackSnapshotWithBody(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRollbackSnapshotRequestWithBody(c.Server, instanceId, snapshotId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RollbackSnapshot(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, body models.RollbackSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRollbackSnapshotRequest(c.Server, instanceId, snapshotId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpgradeInstanceWithBody(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpgradeInstanceRequestWithBody(c.Server, instanceId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpgradeInstance(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, body models.UpgradeInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpgradeInstanceRequest(c.Server, instanceId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSnapshotsAuditsList(ctx context.Context, params *models.RetrieveSnapshotsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSnapshotsAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTicketWithBody(ctx context.Context, params *models.CreateTicketParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTicketRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTicket(ctx context.Context, params *models.CreateTicketParams, body models.CreateTicketJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTicketRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveDataCenterList(ctx context.Context, params *models.RetrieveDataCenterListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveDataCenterListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveObjectStorageList(ctx context.Context, params *models.RetrieveObjectStorageListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveObjectStorageListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateObjectStorageWithBody(ctx context.Context, params *models.CreateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateObjectStorageRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateObjectStorage(ctx context.Context, params *models.CreateObjectStorageParams, body models.CreateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateObjectStorageRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveObjectStorageAuditsList(ctx context.Context, params *models.RetrieveObjectStorageAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveObjectStorageAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveObjectStorage(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStorageParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveObjectStorageRequest(c.Server, objectStorageId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateObjectStorageRequestWithBody(c.Server, objectStorageId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateObjectStorage(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, body models.UpdateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateObjectStorageRequest(c.Server, objectStorageId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelObjectStorageRequestWithBody(c.Server, objectStorageId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelObjectStorage(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, body models.CancelObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelObjectStorageRequest(c.Server, objectStorageId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpgradeObjectStorageWithBody(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpgradeObjectStorageRequestWithBody(c.Server, objectStorageId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpgradeObjectStorage(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, body models.UpgradeObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpgradeObjectStorageRequest(c.Server, objectStorageId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveObjectStoragesStats(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStoragesStatsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveObjectStoragesStatsRequest(c.Server, objectStorageId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrievePrivateNetworkList(ctx context.Context, params *models.RetrievePrivateNetworkListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrievePrivateNetworkListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreatePrivateNetworkWithBody(ctx context.Context, params *models.CreatePrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreatePrivateNetworkRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreatePrivateNetwork(ctx context.Context, params *models.CreatePrivateNetworkParams, body models.CreatePrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreatePrivateNetworkRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrievePrivateNetworkAuditsList(ctx context.Context, params *models.RetrievePrivateNetworkAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrievePrivateNetworkAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeletePrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.DeletePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeletePrivateNetworkRequest(c.Server, privateNetworkId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrievePrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.RetrievePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrievePrivateNetworkRequest(c.Server, privateNetworkId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchPrivateNetworkWithBody(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchPrivateNetworkRequestWithBody(c.Server, privateNetworkId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchPrivateNetwork(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, body models.PatchPrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchPrivateNetworkRequest(c.Server, privateNetworkId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UnassignInstancePrivateNetwork(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.UnassignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUnassignInstancePrivateNetworkRequest(c.Server, privateNetworkId, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AssignInstancePrivateNetwork(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.AssignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAssignInstancePrivateNetworkRequest(c.Server, privateNetworkId, instanceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveRoleList(ctx context.Context, params *models.RetrieveRoleListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveRoleListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRoleWithBody(ctx context.Context, params *models.CreateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRoleRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateRole(ctx context.Context, params *models.CreateRoleParams, body models.CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRoleRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveApiPermissionsList(ctx context.Context, params *models.RetrieveApiPermissionsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveApiPermissionsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveRoleAuditsList(ctx context.Context, params *models.RetrieveRoleAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveRoleAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteRole(ctx context.Context, roleId int64, params *models.DeleteRoleParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteRoleRequest(c.Server, roleId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveRole(ctx context.Context, roleId int64, params *models.RetrieveRoleParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveRoleRequest(c.Server, roleId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateRoleWithBody(ctx context.Context, roleId int64, params *models.UpdateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateRoleRequestWithBody(c.Server, roleId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateRole(ctx context.Context, roleId int64, params *models.UpdateRoleParams, body models.UpdateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateRoleRequest(c.Server, roleId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSecretList(ctx context.Context, params *models.RetrieveSecretListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSecretListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSecretWithBody(ctx context.Context, params *models.CreateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSecretRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSecret(ctx context.Context, params *models.CreateSecretParams, body models.CreateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSecretRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSecretAuditsList(ctx context.Context, params *models.RetrieveSecretAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSecretAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteSecret(ctx context.Context, secretId int64, params *models.DeleteSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteSecretRequest(c.Server, secretId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveSecret(ctx context.Context, secretId int64, params *models.RetrieveSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveSecretRequest(c.Server, secretId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateSecretWithBody(ctx context.Context, secretId int64, params *models.UpdateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateSecretRequestWithBody(c.Server, secretId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateSecret(ctx context.Context, secretId int64, params *models.UpdateSecretParams, body models.UpdateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateSecretRequest(c.Server, secretId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveTagList(ctx context.Context, params *models.RetrieveTagListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveTagListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTagWithBody(ctx context.Context, params *models.CreateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTagRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTag(ctx context.Context, params *models.CreateTagParams, body models.CreateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTagRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveAssignmentsAuditsList(ctx context.Context, params *models.RetrieveAssignmentsAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveAssignmentsAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveTagAuditsList(ctx context.Context, params *models.RetrieveTagAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveTagAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteTag(ctx context.Context, tagId int64, params *models.DeleteTagParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteTagRequest(c.Server, tagId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveTag(ctx context.Context, tagId int64, params *models.RetrieveTagParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveTagRequest(c.Server, tagId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTagWithBody(ctx context.Context, tagId int64, params *models.UpdateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTagRequestWithBody(c.Server, tagId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTag(ctx context.Context, tagId int64, params *models.UpdateTagParams, body models.UpdateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTagRequest(c.Server, tagId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveAssignmentList(ctx context.Context, tagId int64, params *models.RetrieveAssignmentListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveAssignmentListRequest(c.Server, tagId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.DeleteAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteAssignmentRequest(c.Server, tagId, resourceType, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.RetrieveAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveAssignmentRequest(c.Server, tagId, resourceType, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateAssignment(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.CreateAssignmentParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateAssignmentRequest(c.Server, tagId, resourceType, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveUserList(ctx context.Context, params *models.RetrieveUserListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveUserListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateUserWithBody(ctx context.Context, params *models.CreateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateUserRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateUser(ctx context.Context, params *models.CreateUserParams, body models.CreateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateUserRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveUserAuditsList(ctx context.Context, params *models.RetrieveUserAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveUserAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveUserClient(ctx context.Context, params *models.RetrieveUserClientParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveUserClientRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GenerateClientSecret(ctx context.Context, params *models.GenerateClientSecretParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGenerateClientSecretRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveUserIsPasswordSet(ctx context.Context, params *models.RetrieveUserIsPasswordSetParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveUserIsPasswordSetRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteUser(ctx context.Context, userId string, params *models.DeleteUserParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteUserRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveUser(ctx context.Context, userId string, params *models.RetrieveUserParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveUserRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateUserWithBody(ctx context.Context, userId string, params *models.UpdateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateUserRequestWithBody(c.Server, userId, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateUser(ctx context.Context, userId string, params *models.UpdateUserParams, body models.UpdateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateUserRequest(c.Server, userId, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListObjectStorageCredentials(ctx context.Context, userId string, params *models.ListObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListObjectStorageCredentialsRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetObjectStorageCredentials(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.GetObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetObjectStorageCredentialsRequest(c.Server, userId, objectStorageId, credentialId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RegenerateObjectStorageCredentials(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.RegenerateObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegenerateObjectStorageCredentialsRequest(c.Server, userId, objectStorageId, credentialId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ResendEmailVerification(ctx context.Context, userId string, params *models.ResendEmailVerificationParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewResendEmailVerificationRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ResetPassword(ctx context.Context, userId string, params *models.ResetPasswordParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewResetPasswordRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveVipList(ctx context.Context, params *models.RetrieveVipListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveVipListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveVipAuditsList(ctx context.Context, params *models.RetrieveVipAuditsListParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveVipAuditsListRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RetrieveVip(ctx context.Context, ip string, params *models.RetrieveVipParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRetrieveVipRequest(c.Server, ip, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UnassignIp(ctx context.Context, ip string, resourceType models.UnassignIpParamsResourceType, resourceId int64, params *models.UnassignIpParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUnassignIpRequest(c.Server, ip, resourceType, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AssignIp(ctx context.Context, ip string, resourceType models.AssignIpParamsResourceType, resourceId int64, params *models.AssignIpParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAssignIpRequest(c.Server, ip, resourceType, resourceId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewRetrieveImageListRequest generates requests for RetrieveImageList
func NewRetrieveImageListRequest(server string, params *models.RetrieveImageListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/images"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StandardImage != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "standardImage", runtime.ParamLocationQuery, *params.StandardImage); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Search != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "search", runtime.ParamLocationQuery, *params.Search); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateCustomImageRequest calls the generic CreateCustomImage builder with application/json body
func NewCreateCustomImageRequest(server string, params *models.CreateCustomImageParams, body models.CreateCustomImageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateCustomImageRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateCustomImageRequestWithBody generates requests for CreateCustomImage with any type of body
func NewCreateCustomImageRequestWithBody(server string, params *models.CreateCustomImageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/images"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveImageAuditsListRequest generates requests for RetrieveImageAuditsList
func NewRetrieveImageAuditsListRequest(server string, params *models.RetrieveImageAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/images/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ImageId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "imageId", runtime.ParamLocationQuery, *params.ImageId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveCustomImagesStatsRequest generates requests for RetrieveCustomImagesStats
func NewRetrieveCustomImagesStatsRequest(server string, params *models.RetrieveCustomImagesStatsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/images/stats"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteImageRequest generates requests for DeleteImage
func NewDeleteImageRequest(server string, imageId string, params *models.DeleteImageParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "imageId", runtime.ParamLocationPath, imageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/images/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveImageRequest generates requests for RetrieveImage
func NewRetrieveImageRequest(server string, imageId string, params *models.RetrieveImageParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "imageId", runtime.ParamLocationPath, imageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/images/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateImageRequest calls the generic UpdateImage builder with application/json body
func NewUpdateImageRequest(server string, imageId string, params *models.UpdateImageParams, body models.UpdateImageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateImageRequestWithBody(server, imageId, params, "application/json", bodyReader)
}

// NewUpdateImageRequestWithBody generates requests for UpdateImage with any type of body
func NewUpdateImageRequestWithBody(server string, imageId string, params *models.UpdateImageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "imageId", runtime.ParamLocationPath, imageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/images/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveInstancesListRequest generates requests for RetrieveInstancesList
func NewRetrieveInstancesListRequest(server string, params *models.RetrieveInstancesListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/instances"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DisplayName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "displayName", runtime.ParamLocationQuery, *params.DisplayName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DataCenter != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "dataCenter", runtime.ParamLocationQuery, *params.DataCenter); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Region != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "region", runtime.ParamLocationQuery, *params.Region); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceId", runtime.ParamLocationQuery, *params.InstanceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceIds != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceIds", runtime.ParamLocationQuery, *params.InstanceIds); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Status != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "status", runtime.ParamLocationQuery, *params.Status); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ProductIds != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "productIds", runtime.ParamLocationQuery, *params.ProductIds); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.AddOnIds != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "addOnIds", runtime.ParamLocationQuery, *params.AddOnIds); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ProductTypes != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "productTypes", runtime.ParamLocationQuery, *params.ProductTypes); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.IpConfig != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "ipConfig", runtime.ParamLocationQuery, *params.IpConfig); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Search != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "search", runtime.ParamLocationQuery, *params.Search); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateInstanceRequest calls the generic CreateInstance builder with application/json body
func NewCreateInstanceRequest(server string, params *models.CreateInstanceParams, body models.CreateInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateInstanceRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateInstanceRequestWithBody generates requests for CreateInstance with any type of body
func NewCreateInstanceRequestWithBody(server string, params *models.CreateInstanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/instances"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveInstancesActionsAuditsListRequest generates requests for RetrieveInstancesActionsAuditsList
func NewRetrieveInstancesActionsAuditsListRequest(server string, params *models.RetrieveInstancesActionsAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/instances/actions/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceId", runtime.ParamLocationQuery, *params.InstanceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveInstancesAuditsListRequest generates requests for RetrieveInstancesAuditsList
func NewRetrieveInstancesAuditsListRequest(server string, params *models.RetrieveInstancesAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/instances/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceId", runtime.ParamLocationQuery, *params.InstanceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveInstanceRequest generates requests for RetrieveInstance
func NewRetrieveInstanceRequest(server string, instanceId int64, params *models.RetrieveInstanceParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewPatchInstanceRequest calls the generic PatchInstance builder with application/json body
func NewPatchInstanceRequest(server string, instanceId int64, params *models.PatchInstanceParams, body models.PatchInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPatchInstanceRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewPatchInstanceRequestWithBody generates requests for PatchInstance with any type of body
func NewPatchInstanceRequestWithBody(server string, instanceId int64, params *models.PatchInstanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewReinstallInstanceRequest calls the generic ReinstallInstance builder with application/json body
func NewReinstallInstanceRequest(server string, instanceId int64, params *models.ReinstallInstanceParams, body models.ReinstallInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewReinstallInstanceRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewReinstallInstanceRequestWithBody generates requests for ReinstallInstance with any type of body
func NewReinstallInstanceRequestWithBody(server string, instanceId int64, params *models.ReinstallInstanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRescueRequest calls the generic Rescue builder with application/json body
func NewRescueRequest(server string, instanceId int64, params *models.RescueParams, body models.RescueJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRescueRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewRescueRequestWithBody generates requests for Rescue with any type of body
func NewRescueRequestWithBody(server string, instanceId int64, params *models.RescueParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/rescue", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewResetPasswordActionRequest calls the generic ResetPasswordAction builder with application/json body
func NewResetPasswordActionRequest(server string, instanceId int64, params *models.ResetPasswordActionParams, body models.ResetPasswordActionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewResetPasswordActionRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewResetPasswordActionRequestWithBody generates requests for ResetPasswordAction with any type of body
func NewResetPasswordActionRequestWithBody(server string, instanceId int64, params *models.ResetPasswordActionParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/resetPassword", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRestartRequest generates requests for Restart
func NewRestartRequest(server string, instanceId int64, params *models.RestartParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/restart", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewShutdownRequest generates requests for Shutdown
func NewShutdownRequest(server string, instanceId int64, params *models.ShutdownParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/shutdown", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewStartRequest generates requests for Start
func NewStartRequest(server string, instanceId int64, params *models.StartParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/start", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewStopRequest generates requests for Stop
func NewStopRequest(server string, instanceId int64, params *models.StopParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/actions/stop", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCancelInstanceRequest calls the generic CancelInstance builder with application/json body
func NewCancelInstanceRequest(server string, instanceId int64, params *models.CancelInstanceParams, body models.CancelInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCancelInstanceRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewCancelInstanceRequestWithBody generates requests for CancelInstance with any type of body
func NewCancelInstanceRequestWithBody(server string, instanceId int64, params *models.CancelInstanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/cancel", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSnapshotListRequest generates requests for RetrieveSnapshotList
func NewRetrieveSnapshotListRequest(server string, instanceId int64, params *models.RetrieveSnapshotListParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateSnapshotRequest calls the generic CreateSnapshot builder with application/json body
func NewCreateSnapshotRequest(server string, instanceId int64, params *models.CreateSnapshotParams, body models.CreateSnapshotJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateSnapshotRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewCreateSnapshotRequestWithBody generates requests for CreateSnapshot with any type of body
func NewCreateSnapshotRequestWithBody(server string, instanceId int64, params *models.CreateSnapshotParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteSnapshotRequest generates requests for DeleteSnapshot
func NewDeleteSnapshotRequest(server string, instanceId int64, snapshotId string, params *models.DeleteSnapshotParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "snapshotId", runtime.ParamLocationPath, snapshotId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSnapshotRequest generates requests for RetrieveSnapshot
func NewRetrieveSnapshotRequest(server string, instanceId int64, snapshotId string, params *models.RetrieveSnapshotParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "snapshotId", runtime.ParamLocationPath, snapshotId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateSnapshotRequest calls the generic UpdateSnapshot builder with application/json body
func NewUpdateSnapshotRequest(server string, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, body models.UpdateSnapshotJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateSnapshotRequestWithBody(server, instanceId, snapshotId, params, "application/json", bodyReader)
}

// NewUpdateSnapshotRequestWithBody generates requests for UpdateSnapshot with any type of body
func NewUpdateSnapshotRequestWithBody(server string, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "snapshotId", runtime.ParamLocationPath, snapshotId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRollbackSnapshotRequest calls the generic RollbackSnapshot builder with application/json body
func NewRollbackSnapshotRequest(server string, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, body models.RollbackSnapshotJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRollbackSnapshotRequestWithBody(server, instanceId, snapshotId, params, "application/json", bodyReader)
}

// NewRollbackSnapshotRequestWithBody generates requests for RollbackSnapshot with any type of body
func NewRollbackSnapshotRequestWithBody(server string, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "snapshotId", runtime.ParamLocationPath, snapshotId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/snapshots/%s/rollback", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpgradeInstanceRequest calls the generic UpgradeInstance builder with application/json body
func NewUpgradeInstanceRequest(server string, instanceId int64, params *models.UpgradeInstanceParams, body models.UpgradeInstanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpgradeInstanceRequestWithBody(server, instanceId, params, "application/json", bodyReader)
}

// NewUpgradeInstanceRequestWithBody generates requests for UpgradeInstance with any type of body
func NewUpgradeInstanceRequestWithBody(server string, instanceId int64, params *models.UpgradeInstanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/compute/instances/%s/upgrade", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSnapshotsAuditsListRequest generates requests for RetrieveSnapshotsAuditsList
func NewRetrieveSnapshotsAuditsListRequest(server string, params *models.RetrieveSnapshotsAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/compute/snapshots/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceId", runtime.ParamLocationQuery, *params.InstanceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.SnapshotId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "snapshotId", runtime.ParamLocationQuery, *params.SnapshotId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateTicketRequest calls the generic CreateTicket builder with application/json body
func NewCreateTicketRequest(server string, params *models.CreateTicketParams, body models.CreateTicketJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateTicketRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateTicketRequestWithBody generates requests for CreateTicket with any type of body
func NewCreateTicketRequestWithBody(server string, params *models.CreateTicketParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/create-ticket"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveDataCenterListRequest generates requests for RetrieveDataCenterList
func NewRetrieveDataCenterListRequest(server string, params *models.RetrieveDataCenterListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/data-centers"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Slug != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "slug", runtime.ParamLocationQuery, *params.Slug); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RegionName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "regionName", runtime.ParamLocationQuery, *params.RegionName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RegionSlug != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "regionSlug", runtime.ParamLocationQuery, *params.RegionSlug); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveObjectStorageListRequest generates requests for RetrieveObjectStorageList
func NewRetrieveObjectStorageListRequest(server string, params *models.RetrieveObjectStorageListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/object-storages"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DataCenterName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "dataCenterName", runtime.ParamLocationQuery, *params.DataCenterName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.S3TenantId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "s3TenantId", runtime.ParamLocationQuery, *params.S3TenantId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Region != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "region", runtime.ParamLocationQuery, *params.Region); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DisplayName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "displayName", runtime.ParamLocationQuery, *params.DisplayName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateObjectStorageRequest calls the generic CreateObjectStorage builder with application/json body
func NewCreateObjectStorageRequest(server string, params *models.CreateObjectStorageParams, body models.CreateObjectStorageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateObjectStorageRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateObjectStorageRequestWithBody generates requests for CreateObjectStorage with any type of body
func NewCreateObjectStorageRequestWithBody(server string, params *models.CreateObjectStorageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/object-storages"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveObjectStorageAuditsListRequest generates requests for RetrieveObjectStorageAuditsList
func NewRetrieveObjectStorageAuditsListRequest(server string, params *models.RetrieveObjectStorageAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/object-storages/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ObjectStorageId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "objectStorageId", runtime.ParamLocationQuery, *params.ObjectStorageId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveObjectStorageRequest generates requests for RetrieveObjectStorage
func NewRetrieveObjectStorageRequest(server string, objectStorageId string, params *models.RetrieveObjectStorageParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/object-storages/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateObjectStorageRequest calls the generic UpdateObjectStorage builder with application/json body
func NewUpdateObjectStorageRequest(server string, objectStorageId string, params *models.UpdateObjectStorageParams, body models.UpdateObjectStorageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateObjectStorageRequestWithBody(server, objectStorageId, params, "application/json", bodyReader)
}

// NewUpdateObjectStorageRequestWithBody generates requests for UpdateObjectStorage with any type of body
func NewUpdateObjectStorageRequestWithBody(server string, objectStorageId string, params *models.UpdateObjectStorageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/object-storages/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCancelObjectStorageRequest calls the generic CancelObjectStorage builder with application/json body
func NewCancelObjectStorageRequest(server string, objectStorageId string, params *models.CancelObjectStorageParams, body models.CancelObjectStorageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCancelObjectStorageRequestWithBody(server, objectStorageId, params, "application/json", bodyReader)
}

// NewCancelObjectStorageRequestWithBody generates requests for CancelObjectStorage with any type of body
func NewCancelObjectStorageRequestWithBody(server string, objectStorageId string, params *models.CancelObjectStorageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/object-storages/%s/cancel", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpgradeObjectStorageRequest calls the generic UpgradeObjectStorage builder with application/json body
func NewUpgradeObjectStorageRequest(server string, objectStorageId string, params *models.UpgradeObjectStorageParams, body models.UpgradeObjectStorageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpgradeObjectStorageRequestWithBody(server, objectStorageId, params, "application/json", bodyReader)
}

// NewUpgradeObjectStorageRequestWithBody generates requests for UpgradeObjectStorage with any type of body
func NewUpgradeObjectStorageRequestWithBody(server string, objectStorageId string, params *models.UpgradeObjectStorageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/object-storages/%s/resize", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveObjectStoragesStatsRequest generates requests for RetrieveObjectStoragesStats
func NewRetrieveObjectStoragesStatsRequest(server string, objectStorageId string, params *models.RetrieveObjectStoragesStatsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/object-storages/%s/stats", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrievePrivateNetworkListRequest generates requests for RetrievePrivateNetworkList
func NewRetrievePrivateNetworkListRequest(server string, params *models.RetrievePrivateNetworkListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/private-networks"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.InstanceIds != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "instanceIds", runtime.ParamLocationQuery, *params.InstanceIds); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Region != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "region", runtime.ParamLocationQuery, *params.Region); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DataCenter != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "dataCenter", runtime.ParamLocationQuery, *params.DataCenter); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreatePrivateNetworkRequest calls the generic CreatePrivateNetwork builder with application/json body
func NewCreatePrivateNetworkRequest(server string, params *models.CreatePrivateNetworkParams, body models.CreatePrivateNetworkJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreatePrivateNetworkRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreatePrivateNetworkRequestWithBody generates requests for CreatePrivateNetwork with any type of body
func NewCreatePrivateNetworkRequestWithBody(server string, params *models.CreatePrivateNetworkParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/private-networks"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrievePrivateNetworkAuditsListRequest generates requests for RetrievePrivateNetworkAuditsList
func NewRetrievePrivateNetworkAuditsListRequest(server string, params *models.RetrievePrivateNetworkAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/private-networks/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PrivateNetworkId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "privateNetworkId", runtime.ParamLocationQuery, *params.PrivateNetworkId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeletePrivateNetworkRequest generates requests for DeletePrivateNetwork
func NewDeletePrivateNetworkRequest(server string, privateNetworkId int64, params *models.DeletePrivateNetworkParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "privateNetworkId", runtime.ParamLocationPath, privateNetworkId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/private-networks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrievePrivateNetworkRequest generates requests for RetrievePrivateNetwork
func NewRetrievePrivateNetworkRequest(server string, privateNetworkId int64, params *models.RetrievePrivateNetworkParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "privateNetworkId", runtime.ParamLocationPath, privateNetworkId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/private-networks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewPatchPrivateNetworkRequest calls the generic PatchPrivateNetwork builder with application/json body
func NewPatchPrivateNetworkRequest(server string, privateNetworkId int64, params *models.PatchPrivateNetworkParams, body models.PatchPrivateNetworkJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPatchPrivateNetworkRequestWithBody(server, privateNetworkId, params, "application/json", bodyReader)
}

// NewPatchPrivateNetworkRequestWithBody generates requests for PatchPrivateNetwork with any type of body
func NewPatchPrivateNetworkRequestWithBody(server string, privateNetworkId int64, params *models.PatchPrivateNetworkParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "privateNetworkId", runtime.ParamLocationPath, privateNetworkId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/private-networks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUnassignInstancePrivateNetworkRequest generates requests for UnassignInstancePrivateNetwork
func NewUnassignInstancePrivateNetworkRequest(server string, privateNetworkId int64, instanceId int64, params *models.UnassignInstancePrivateNetworkParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "privateNetworkId", runtime.ParamLocationPath, privateNetworkId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/private-networks/%s/instances/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewAssignInstancePrivateNetworkRequest generates requests for AssignInstancePrivateNetwork
func NewAssignInstancePrivateNetworkRequest(server string, privateNetworkId int64, instanceId int64, params *models.AssignInstancePrivateNetworkParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "privateNetworkId", runtime.ParamLocationPath, privateNetworkId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "instanceId", runtime.ParamLocationPath, instanceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/private-networks/%s/instances/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveRoleListRequest generates requests for RetrieveRoleList
func NewRetrieveRoleListRequest(server string, params *models.RetrieveRoleListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/roles"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ApiName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "apiName", runtime.ParamLocationQuery, *params.ApiName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.TagName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "tagName", runtime.ParamLocationQuery, *params.TagName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Type != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "type", runtime.ParamLocationQuery, *params.Type); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateRoleRequest calls the generic CreateRole builder with application/json body
func NewCreateRoleRequest(server string, params *models.CreateRoleParams, body models.CreateRoleJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateRoleRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateRoleRequestWithBody generates requests for CreateRole with any type of body
func NewCreateRoleRequestWithBody(server string, params *models.CreateRoleParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/roles"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveApiPermissionsListRequest generates requests for RetrieveApiPermissionsList
func NewRetrieveApiPermissionsListRequest(server string, params *models.RetrieveApiPermissionsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/roles/api-permissions"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ApiName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "apiName", runtime.ParamLocationQuery, *params.ApiName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveRoleAuditsListRequest generates requests for RetrieveRoleAuditsList
func NewRetrieveRoleAuditsListRequest(server string, params *models.RetrieveRoleAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/roles/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RoleId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "roleId", runtime.ParamLocationQuery, *params.RoleId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteRoleRequest generates requests for DeleteRole
func NewDeleteRoleRequest(server string, roleId int64, params *models.DeleteRoleParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roleId", runtime.ParamLocationPath, roleId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/roles/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveRoleRequest generates requests for RetrieveRole
func NewRetrieveRoleRequest(server string, roleId int64, params *models.RetrieveRoleParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roleId", runtime.ParamLocationPath, roleId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/roles/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateRoleRequest calls the generic UpdateRole builder with application/json body
func NewUpdateRoleRequest(server string, roleId int64, params *models.UpdateRoleParams, body models.UpdateRoleJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateRoleRequestWithBody(server, roleId, params, "application/json", bodyReader)
}

// NewUpdateRoleRequestWithBody generates requests for UpdateRole with any type of body
func NewUpdateRoleRequestWithBody(server string, roleId int64, params *models.UpdateRoleParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "roleId", runtime.ParamLocationPath, roleId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/roles/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSecretListRequest generates requests for RetrieveSecretList
func NewRetrieveSecretListRequest(server string, params *models.RetrieveSecretListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/secrets"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Type != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "type", runtime.ParamLocationQuery, *params.Type); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateSecretRequest calls the generic CreateSecret builder with application/json body
func NewCreateSecretRequest(server string, params *models.CreateSecretParams, body models.CreateSecretJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateSecretRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateSecretRequestWithBody generates requests for CreateSecret with any type of body
func NewCreateSecretRequestWithBody(server string, params *models.CreateSecretParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/secrets"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSecretAuditsListRequest generates requests for RetrieveSecretAuditsList
func NewRetrieveSecretAuditsListRequest(server string, params *models.RetrieveSecretAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/secrets/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.SecretId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "secretId", runtime.ParamLocationQuery, *params.SecretId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteSecretRequest generates requests for DeleteSecret
func NewDeleteSecretRequest(server string, secretId int64, params *models.DeleteSecretParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "secretId", runtime.ParamLocationPath, secretId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/secrets/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveSecretRequest generates requests for RetrieveSecret
func NewRetrieveSecretRequest(server string, secretId int64, params *models.RetrieveSecretParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "secretId", runtime.ParamLocationPath, secretId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/secrets/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateSecretRequest calls the generic UpdateSecret builder with application/json body
func NewUpdateSecretRequest(server string, secretId int64, params *models.UpdateSecretParams, body models.UpdateSecretJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateSecretRequestWithBody(server, secretId, params, "application/json", bodyReader)
}

// NewUpdateSecretRequestWithBody generates requests for UpdateSecret with any type of body
func NewUpdateSecretRequestWithBody(server string, secretId int64, params *models.UpdateSecretParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "secretId", runtime.ParamLocationPath, secretId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/secrets/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveTagListRequest generates requests for RetrieveTagList
func NewRetrieveTagListRequest(server string, params *models.RetrieveTagListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/tags"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Name != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "name", runtime.ParamLocationQuery, *params.Name); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateTagRequest calls the generic CreateTag builder with application/json body
func NewCreateTagRequest(server string, params *models.CreateTagParams, body models.CreateTagJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateTagRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateTagRequestWithBody generates requests for CreateTag with any type of body
func NewCreateTagRequestWithBody(server string, params *models.CreateTagParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/tags"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveAssignmentsAuditsListRequest generates requests for RetrieveAssignmentsAuditsList
func NewRetrieveAssignmentsAuditsListRequest(server string, params *models.RetrieveAssignmentsAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/tags/assignments/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.TagId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "tagId", runtime.ParamLocationQuery, *params.TagId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceId", runtime.ParamLocationQuery, *params.ResourceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveTagAuditsListRequest generates requests for RetrieveTagAuditsList
func NewRetrieveTagAuditsListRequest(server string, params *models.RetrieveTagAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/tags/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.TagId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "tagId", runtime.ParamLocationQuery, *params.TagId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteTagRequest generates requests for DeleteTag
func NewDeleteTagRequest(server string, tagId int64, params *models.DeleteTagParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveTagRequest generates requests for RetrieveTag
func NewRetrieveTagRequest(server string, tagId int64, params *models.RetrieveTagParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateTagRequest calls the generic UpdateTag builder with application/json body
func NewUpdateTagRequest(server string, tagId int64, params *models.UpdateTagParams, body models.UpdateTagJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateTagRequestWithBody(server, tagId, params, "application/json", bodyReader)
}

// NewUpdateTagRequestWithBody generates requests for UpdateTag with any type of body
func NewUpdateTagRequestWithBody(server string, tagId int64, params *models.UpdateTagParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveAssignmentListRequest generates requests for RetrieveAssignmentList
func NewRetrieveAssignmentListRequest(server string, tagId int64, params *models.RetrieveAssignmentListParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s/assignments", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceType != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceType", runtime.ParamLocationQuery, *params.ResourceType); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteAssignmentRequest generates requests for DeleteAssignment
func NewDeleteAssignmentRequest(server string, tagId int64, resourceType string, resourceId string, params *models.DeleteAssignmentParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "resourceType", runtime.ParamLocationPath, resourceType)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s/assignments/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveAssignmentRequest generates requests for RetrieveAssignment
func NewRetrieveAssignmentRequest(server string, tagId int64, resourceType string, resourceId string, params *models.RetrieveAssignmentParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "resourceType", runtime.ParamLocationPath, resourceType)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s/assignments/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateAssignmentRequest generates requests for CreateAssignment
func NewCreateAssignmentRequest(server string, tagId int64, resourceType string, resourceId string, params *models.CreateAssignmentParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "tagId", runtime.ParamLocationPath, tagId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "resourceType", runtime.ParamLocationPath, resourceType)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/tags/%s/assignments/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveUserListRequest generates requests for RetrieveUserList
func NewRetrieveUserListRequest(server string, params *models.RetrieveUserListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Email != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "email", runtime.ParamLocationQuery, *params.Email); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Enabled != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "enabled", runtime.ParamLocationQuery, *params.Enabled); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Owner != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "owner", runtime.ParamLocationQuery, *params.Owner); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewCreateUserRequest calls the generic CreateUser builder with application/json body
func NewCreateUserRequest(server string, params *models.CreateUserParams, body models.CreateUserJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateUserRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateUserRequestWithBody generates requests for CreateUser with any type of body
func NewCreateUserRequestWithBody(server string, params *models.CreateUserParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveUserAuditsListRequest generates requests for RetrieveUserAuditsList
func NewRetrieveUserAuditsListRequest(server string, params *models.RetrieveUserAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.UserId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "userId", runtime.ParamLocationQuery, *params.UserId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveUserClientRequest generates requests for RetrieveUserClient
func NewRetrieveUserClientRequest(server string, params *models.RetrieveUserClientParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users/client"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewGenerateClientSecretRequest generates requests for GenerateClientSecret
func NewGenerateClientSecretRequest(server string, params *models.GenerateClientSecretParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users/client/secret"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveUserIsPasswordSetRequest generates requests for RetrieveUserIsPasswordSet
func NewRetrieveUserIsPasswordSetRequest(server string, params *models.RetrieveUserIsPasswordSetParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/users/is-password-set"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.UserId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "userId", runtime.ParamLocationQuery, *params.UserId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewDeleteUserRequest generates requests for DeleteUser
func NewDeleteUserRequest(server string, userId string, params *models.DeleteUserParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveUserRequest generates requests for RetrieveUser
func NewRetrieveUserRequest(server string, userId string, params *models.RetrieveUserParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUpdateUserRequest calls the generic UpdateUser builder with application/json body
func NewUpdateUserRequest(server string, userId string, params *models.UpdateUserParams, body models.UpdateUserJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateUserRequestWithBody(server, userId, params, "application/json", bodyReader)
}

// NewUpdateUserRequestWithBody generates requests for UpdateUser with any type of body
func NewUpdateUserRequestWithBody(server string, userId string, params *models.UpdateUserParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewListObjectStorageCredentialsRequest generates requests for ListObjectStorageCredentials
func NewListObjectStorageCredentialsRequest(server string, userId string, params *models.ListObjectStorageCredentialsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s/object-storages/credentials", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ObjectStorageId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "objectStorageId", runtime.ParamLocationQuery, *params.ObjectStorageId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RegionName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "regionName", runtime.ParamLocationQuery, *params.RegionName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DisplayName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "displayName", runtime.ParamLocationQuery, *params.DisplayName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewGetObjectStorageCredentialsRequest generates requests for GetObjectStorageCredentials
func NewGetObjectStorageCredentialsRequest(server string, userId string, objectStorageId string, credentialId int64, params *models.GetObjectStorageCredentialsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "credentialId", runtime.ParamLocationPath, credentialId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s/object-storages/%s/credentials/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRegenerateObjectStorageCredentialsRequest generates requests for RegenerateObjectStorageCredentials
func NewRegenerateObjectStorageCredentialsRequest(server string, userId string, objectStorageId string, credentialId int64, params *models.RegenerateObjectStorageCredentialsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "objectStorageId", runtime.ParamLocationPath, objectStorageId)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "credentialId", runtime.ParamLocationPath, credentialId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s/object-storages/%s/credentials/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewResendEmailVerificationRequest generates requests for ResendEmailVerification
func NewResendEmailVerificationRequest(server string, userId string, params *models.ResendEmailVerificationParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s/resend-email-verification", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.RedirectUrl != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "redirectUrl", runtime.ParamLocationQuery, *params.RedirectUrl); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewResetPasswordRequest generates requests for ResetPassword
func NewResetPasswordRequest(server string, userId string, params *models.ResetPasswordParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/users/%s/reset-password", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.RedirectUrl != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "redirectUrl", runtime.ParamLocationQuery, *params.RedirectUrl); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveVipListRequest generates requests for RetrieveVipList
func NewRetrieveVipListRequest(server string, params *models.RetrieveVipListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/vips"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceId", runtime.ParamLocationQuery, *params.ResourceId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceType != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceType", runtime.ParamLocationQuery, *params.ResourceType); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceName", runtime.ParamLocationQuery, *params.ResourceName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ResourceDisplayName != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "resourceDisplayName", runtime.ParamLocationQuery, *params.ResourceDisplayName); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.IpVersion != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "ipVersion", runtime.ParamLocationQuery, *params.IpVersion); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Ips != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "ips", runtime.ParamLocationQuery, *params.Ips); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Ip != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "ip", runtime.ParamLocationQuery, *params.Ip); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Type != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "type", runtime.ParamLocationQuery, *params.Type); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.DataCenter != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "dataCenter", runtime.ParamLocationQuery, *params.DataCenter); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Region != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "region", runtime.ParamLocationQuery, *params.Region); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveVipAuditsListRequest generates requests for RetrieveVipAuditsList
func NewRetrieveVipAuditsListRequest(server string, params *models.RetrieveVipAuditsListParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/v1/vips/audits"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Size != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "size", runtime.ParamLocationQuery, *params.Size); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.OrderBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orderBy", runtime.ParamLocationQuery, *params.OrderBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.VipId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "vipId", runtime.ParamLocationQuery, *params.VipId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.RequestId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "requestId", runtime.ParamLocationQuery, *params.RequestId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.ChangedBy != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "changedBy", runtime.ParamLocationQuery, *params.ChangedBy); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.StartDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "startDate", runtime.ParamLocationQuery, *params.StartDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.EndDate != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "endDate", runtime.ParamLocationQuery, *params.EndDate); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewRetrieveVipRequest generates requests for RetrieveVip
func NewRetrieveVipRequest(server string, ip string, params *models.RetrieveVipParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "ip", runtime.ParamLocationPath, ip)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/vips/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewUnassignIpRequest generates requests for UnassignIp
func NewUnassignIpRequest(server string, ip string, resourceType models.UnassignIpParamsResourceType, resourceId int64, params *models.UnassignIpParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "ip", runtime.ParamLocationPath, ip)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "resourceType", runtime.ParamLocationPath, resourceType)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/vips/%s/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

// NewAssignIpRequest generates requests for AssignIp
func NewAssignIpRequest(server string, ip string, resourceType models.AssignIpParamsResourceType, resourceId int64, params *models.AssignIpParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "ip", runtime.ParamLocationPath, ip)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "resourceType", runtime.ParamLocationPath, resourceType)
	if err != nil {
		return nil, err
	}

	var pathParam2 string

	pathParam2, err = runtime.StyleParamWithLocation("simple", false, "resourceId", runtime.ParamLocationPath, resourceId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/vips/%s/%s/%s", pathParam0, pathParam1, pathParam2)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "x-request-id", runtime.ParamLocationHeader, params.XRequestId)
		if err != nil {
			return nil, err
		}

		req.Header.Set("x-request-id", headerParam0)

		if params.XTraceId != nil {
			var headerParam1 string

			headerParam1, err = runtime.StyleParamWithLocation("simple", false, "x-trace-id", runtime.ParamLocationHeader, *params.XTraceId)
			if err != nil {
				return nil, err
			}

			req.Header.Set("x-trace-id", headerParam1)
		}

	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// RetrieveImageListWithResponse request
	RetrieveImageListWithResponse(ctx context.Context, params *models.RetrieveImageListParams, reqEditors ...RequestEditorFn) (*RetrieveImageListResponse, error)

	// CreateCustomImageWithBodyWithResponse request with any body
	CreateCustomImageWithBodyWithResponse(ctx context.Context, params *models.CreateCustomImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateCustomImageResponse, error)

	CreateCustomImageWithResponse(ctx context.Context, params *models.CreateCustomImageParams, body models.CreateCustomImageJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateCustomImageResponse, error)

	// RetrieveImageAuditsListWithResponse request
	RetrieveImageAuditsListWithResponse(ctx context.Context, params *models.RetrieveImageAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveImageAuditsListResponse, error)

	// RetrieveCustomImagesStatsWithResponse request
	RetrieveCustomImagesStatsWithResponse(ctx context.Context, params *models.RetrieveCustomImagesStatsParams, reqEditors ...RequestEditorFn) (*RetrieveCustomImagesStatsResponse, error)

	// DeleteImageWithResponse request
	DeleteImageWithResponse(ctx context.Context, imageId string, params *models.DeleteImageParams, reqEditors ...RequestEditorFn) (*DeleteImageResponse, error)

	// RetrieveImageWithResponse request
	RetrieveImageWithResponse(ctx context.Context, imageId string, params *models.RetrieveImageParams, reqEditors ...RequestEditorFn) (*RetrieveImageResponse, error)

	// UpdateImageWithBodyWithResponse request with any body
	UpdateImageWithBodyWithResponse(ctx context.Context, imageId string, params *models.UpdateImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateImageResponse, error)

	UpdateImageWithResponse(ctx context.Context, imageId string, params *models.UpdateImageParams, body models.UpdateImageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateImageResponse, error)

	// RetrieveInstancesListWithResponse request
	RetrieveInstancesListWithResponse(ctx context.Context, params *models.RetrieveInstancesListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesListResponse, error)

	// CreateInstanceWithBodyWithResponse request with any body
	CreateInstanceWithBodyWithResponse(ctx context.Context, params *models.CreateInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateInstanceResponse, error)

	CreateInstanceWithResponse(ctx context.Context, params *models.CreateInstanceParams, body models.CreateInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateInstanceResponse, error)

	// RetrieveInstancesActionsAuditsListWithResponse request
	RetrieveInstancesActionsAuditsListWithResponse(ctx context.Context, params *models.RetrieveInstancesActionsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesActionsAuditsListResponse, error)

	// RetrieveInstancesAuditsListWithResponse request
	RetrieveInstancesAuditsListWithResponse(ctx context.Context, params *models.RetrieveInstancesAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesAuditsListResponse, error)

	// RetrieveInstanceWithResponse request
	RetrieveInstanceWithResponse(ctx context.Context, instanceId int64, params *models.RetrieveInstanceParams, reqEditors ...RequestEditorFn) (*RetrieveInstanceResponse, error)

	// PatchInstanceWithBodyWithResponse request with any body
	PatchInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchInstanceResponse, error)

	PatchInstanceWithResponse(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, body models.PatchInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchInstanceResponse, error)

	// ReinstallInstanceWithBodyWithResponse request with any body
	ReinstallInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ReinstallInstanceResponse, error)

	ReinstallInstanceWithResponse(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, body models.ReinstallInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*ReinstallInstanceResponse, error)

	// RescueWithBodyWithResponse request with any body
	RescueWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.RescueParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RescueResponse, error)

	RescueWithResponse(ctx context.Context, instanceId int64, params *models.RescueParams, body models.RescueJSONRequestBody, reqEditors ...RequestEditorFn) (*RescueResponse, error)

	// ResetPasswordActionWithBodyWithResponse request with any body
	ResetPasswordActionWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ResetPasswordActionResponse, error)

	ResetPasswordActionWithResponse(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, body models.ResetPasswordActionJSONRequestBody, reqEditors ...RequestEditorFn) (*ResetPasswordActionResponse, error)

	// RestartWithResponse request
	RestartWithResponse(ctx context.Context, instanceId int64, params *models.RestartParams, reqEditors ...RequestEditorFn) (*RestartResponse, error)

	// ShutdownWithResponse request
	ShutdownWithResponse(ctx context.Context, instanceId int64, params *models.ShutdownParams, reqEditors ...RequestEditorFn) (*ShutdownResponse, error)

	// StartWithResponse request
	StartWithResponse(ctx context.Context, instanceId int64, params *models.StartParams, reqEditors ...RequestEditorFn) (*StartResponse, error)

	// StopWithResponse request
	StopWithResponse(ctx context.Context, instanceId int64, params *models.StopParams, reqEditors ...RequestEditorFn) (*StopResponse, error)

	// CancelInstanceWithBodyWithResponse request with any body
	CancelInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelInstanceResponse, error)

	CancelInstanceWithResponse(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, body models.CancelInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelInstanceResponse, error)

	// RetrieveSnapshotListWithResponse request
	RetrieveSnapshotListWithResponse(ctx context.Context, instanceId int64, params *models.RetrieveSnapshotListParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotListResponse, error)

	// CreateSnapshotWithBodyWithResponse request with any body
	CreateSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSnapshotResponse, error)

	CreateSnapshotWithResponse(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, body models.CreateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSnapshotResponse, error)

	// DeleteSnapshotWithResponse request
	DeleteSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.DeleteSnapshotParams, reqEditors ...RequestEditorFn) (*DeleteSnapshotResponse, error)

	// RetrieveSnapshotWithResponse request
	RetrieveSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RetrieveSnapshotParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotResponse, error)

	// UpdateSnapshotWithBodyWithResponse request with any body
	UpdateSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSnapshotResponse, error)

	UpdateSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, body models.UpdateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSnapshotResponse, error)

	// RollbackSnapshotWithBodyWithResponse request with any body
	RollbackSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RollbackSnapshotResponse, error)

	RollbackSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, body models.RollbackSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*RollbackSnapshotResponse, error)

	// UpgradeInstanceWithBodyWithResponse request with any body
	UpgradeInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpgradeInstanceResponse, error)

	UpgradeInstanceWithResponse(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, body models.UpgradeInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*UpgradeInstanceResponse, error)

	// RetrieveSnapshotsAuditsListWithResponse request
	RetrieveSnapshotsAuditsListWithResponse(ctx context.Context, params *models.RetrieveSnapshotsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotsAuditsListResponse, error)

	// CreateTicketWithBodyWithResponse request with any body
	CreateTicketWithBodyWithResponse(ctx context.Context, params *models.CreateTicketParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTicketResponse, error)

	CreateTicketWithResponse(ctx context.Context, params *models.CreateTicketParams, body models.CreateTicketJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTicketResponse, error)

	// RetrieveDataCenterListWithResponse request
	RetrieveDataCenterListWithResponse(ctx context.Context, params *models.RetrieveDataCenterListParams, reqEditors ...RequestEditorFn) (*RetrieveDataCenterListResponse, error)

	// RetrieveObjectStorageListWithResponse request
	RetrieveObjectStorageListWithResponse(ctx context.Context, params *models.RetrieveObjectStorageListParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageListResponse, error)

	// CreateObjectStorageWithBodyWithResponse request with any body
	CreateObjectStorageWithBodyWithResponse(ctx context.Context, params *models.CreateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateObjectStorageResponse, error)

	CreateObjectStorageWithResponse(ctx context.Context, params *models.CreateObjectStorageParams, body models.CreateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateObjectStorageResponse, error)

	// RetrieveObjectStorageAuditsListWithResponse request
	RetrieveObjectStorageAuditsListWithResponse(ctx context.Context, params *models.RetrieveObjectStorageAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageAuditsListResponse, error)

	// RetrieveObjectStorageWithResponse request
	RetrieveObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStorageParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageResponse, error)

	// UpdateObjectStorageWithBodyWithResponse request with any body
	UpdateObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateObjectStorageResponse, error)

	UpdateObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, body models.UpdateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateObjectStorageResponse, error)

	// CancelObjectStorageWithBodyWithResponse request with any body
	CancelObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelObjectStorageResponse, error)

	CancelObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, body models.CancelObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelObjectStorageResponse, error)

	// UpgradeObjectStorageWithBodyWithResponse request with any body
	UpgradeObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpgradeObjectStorageResponse, error)

	UpgradeObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, body models.UpgradeObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpgradeObjectStorageResponse, error)

	// RetrieveObjectStoragesStatsWithResponse request
	RetrieveObjectStoragesStatsWithResponse(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStoragesStatsParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStoragesStatsResponse, error)

	// RetrievePrivateNetworkListWithResponse request
	RetrievePrivateNetworkListWithResponse(ctx context.Context, params *models.RetrievePrivateNetworkListParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkListResponse, error)

	// CreatePrivateNetworkWithBodyWithResponse request with any body
	CreatePrivateNetworkWithBodyWithResponse(ctx context.Context, params *models.CreatePrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreatePrivateNetworkResponse, error)

	CreatePrivateNetworkWithResponse(ctx context.Context, params *models.CreatePrivateNetworkParams, body models.CreatePrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*CreatePrivateNetworkResponse, error)

	// RetrievePrivateNetworkAuditsListWithResponse request
	RetrievePrivateNetworkAuditsListWithResponse(ctx context.Context, params *models.RetrievePrivateNetworkAuditsListParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkAuditsListResponse, error)

	// DeletePrivateNetworkWithResponse request
	DeletePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.DeletePrivateNetworkParams, reqEditors ...RequestEditorFn) (*DeletePrivateNetworkResponse, error)

	// RetrievePrivateNetworkWithResponse request
	RetrievePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.RetrievePrivateNetworkParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkResponse, error)

	// PatchPrivateNetworkWithBodyWithResponse request with any body
	PatchPrivateNetworkWithBodyWithResponse(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchPrivateNetworkResponse, error)

	PatchPrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, body models.PatchPrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchPrivateNetworkResponse, error)

	// UnassignInstancePrivateNetworkWithResponse request
	UnassignInstancePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.UnassignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*UnassignInstancePrivateNetworkResponse, error)

	// AssignInstancePrivateNetworkWithResponse request
	AssignInstancePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.AssignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*AssignInstancePrivateNetworkResponse, error)

	// RetrieveRoleListWithResponse request
	RetrieveRoleListWithResponse(ctx context.Context, params *models.RetrieveRoleListParams, reqEditors ...RequestEditorFn) (*RetrieveRoleListResponse, error)

	// CreateRoleWithBodyWithResponse request with any body
	CreateRoleWithBodyWithResponse(ctx context.Context, params *models.CreateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error)

	CreateRoleWithResponse(ctx context.Context, params *models.CreateRoleParams, body models.CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error)

	// RetrieveApiPermissionsListWithResponse request
	RetrieveApiPermissionsListWithResponse(ctx context.Context, params *models.RetrieveApiPermissionsListParams, reqEditors ...RequestEditorFn) (*RetrieveApiPermissionsListResponse, error)

	// RetrieveRoleAuditsListWithResponse request
	RetrieveRoleAuditsListWithResponse(ctx context.Context, params *models.RetrieveRoleAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveRoleAuditsListResponse, error)

	// DeleteRoleWithResponse request
	DeleteRoleWithResponse(ctx context.Context, roleId int64, params *models.DeleteRoleParams, reqEditors ...RequestEditorFn) (*DeleteRoleResponse, error)

	// RetrieveRoleWithResponse request
	RetrieveRoleWithResponse(ctx context.Context, roleId int64, params *models.RetrieveRoleParams, reqEditors ...RequestEditorFn) (*RetrieveRoleResponse, error)

	// UpdateRoleWithBodyWithResponse request with any body
	UpdateRoleWithBodyWithResponse(ctx context.Context, roleId int64, params *models.UpdateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateRoleResponse, error)

	UpdateRoleWithResponse(ctx context.Context, roleId int64, params *models.UpdateRoleParams, body models.UpdateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateRoleResponse, error)

	// RetrieveSecretListWithResponse request
	RetrieveSecretListWithResponse(ctx context.Context, params *models.RetrieveSecretListParams, reqEditors ...RequestEditorFn) (*RetrieveSecretListResponse, error)

	// CreateSecretWithBodyWithResponse request with any body
	CreateSecretWithBodyWithResponse(ctx context.Context, params *models.CreateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSecretResponse, error)

	CreateSecretWithResponse(ctx context.Context, params *models.CreateSecretParams, body models.CreateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSecretResponse, error)

	// RetrieveSecretAuditsListWithResponse request
	RetrieveSecretAuditsListWithResponse(ctx context.Context, params *models.RetrieveSecretAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveSecretAuditsListResponse, error)

	// DeleteSecretWithResponse request
	DeleteSecretWithResponse(ctx context.Context, secretId int64, params *models.DeleteSecretParams, reqEditors ...RequestEditorFn) (*DeleteSecretResponse, error)

	// RetrieveSecretWithResponse request
	RetrieveSecretWithResponse(ctx context.Context, secretId int64, params *models.RetrieveSecretParams, reqEditors ...RequestEditorFn) (*RetrieveSecretResponse, error)

	// UpdateSecretWithBodyWithResponse request with any body
	UpdateSecretWithBodyWithResponse(ctx context.Context, secretId int64, params *models.UpdateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSecretResponse, error)

	UpdateSecretWithResponse(ctx context.Context, secretId int64, params *models.UpdateSecretParams, body models.UpdateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSecretResponse, error)

	// RetrieveTagListWithResponse request
	RetrieveTagListWithResponse(ctx context.Context, params *models.RetrieveTagListParams, reqEditors ...RequestEditorFn) (*RetrieveTagListResponse, error)

	// CreateTagWithBodyWithResponse request with any body
	CreateTagWithBodyWithResponse(ctx context.Context, params *models.CreateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTagResponse, error)

	CreateTagWithResponse(ctx context.Context, params *models.CreateTagParams, body models.CreateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTagResponse, error)

	// RetrieveAssignmentsAuditsListWithResponse request
	RetrieveAssignmentsAuditsListWithResponse(ctx context.Context, params *models.RetrieveAssignmentsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentsAuditsListResponse, error)

	// RetrieveTagAuditsListWithResponse request
	RetrieveTagAuditsListWithResponse(ctx context.Context, params *models.RetrieveTagAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveTagAuditsListResponse, error)

	// DeleteTagWithResponse request
	DeleteTagWithResponse(ctx context.Context, tagId int64, params *models.DeleteTagParams, reqEditors ...RequestEditorFn) (*DeleteTagResponse, error)

	// RetrieveTagWithResponse request
	RetrieveTagWithResponse(ctx context.Context, tagId int64, params *models.RetrieveTagParams, reqEditors ...RequestEditorFn) (*RetrieveTagResponse, error)

	// UpdateTagWithBodyWithResponse request with any body
	UpdateTagWithBodyWithResponse(ctx context.Context, tagId int64, params *models.UpdateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTagResponse, error)

	UpdateTagWithResponse(ctx context.Context, tagId int64, params *models.UpdateTagParams, body models.UpdateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTagResponse, error)

	// RetrieveAssignmentListWithResponse request
	RetrieveAssignmentListWithResponse(ctx context.Context, tagId int64, params *models.RetrieveAssignmentListParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentListResponse, error)

	// DeleteAssignmentWithResponse request
	DeleteAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.DeleteAssignmentParams, reqEditors ...RequestEditorFn) (*DeleteAssignmentResponse, error)

	// RetrieveAssignmentWithResponse request
	RetrieveAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.RetrieveAssignmentParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentResponse, error)

	// CreateAssignmentWithResponse request
	CreateAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.CreateAssignmentParams, reqEditors ...RequestEditorFn) (*CreateAssignmentResponse, error)

	// RetrieveUserListWithResponse request
	RetrieveUserListWithResponse(ctx context.Context, params *models.RetrieveUserListParams, reqEditors ...RequestEditorFn) (*RetrieveUserListResponse, error)

	// CreateUserWithBodyWithResponse request with any body
	CreateUserWithBodyWithResponse(ctx context.Context, params *models.CreateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateUserResponse, error)

	CreateUserWithResponse(ctx context.Context, params *models.CreateUserParams, body models.CreateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateUserResponse, error)

	// RetrieveUserAuditsListWithResponse request
	RetrieveUserAuditsListWithResponse(ctx context.Context, params *models.RetrieveUserAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveUserAuditsListResponse, error)

	// RetrieveUserClientWithResponse request
	RetrieveUserClientWithResponse(ctx context.Context, params *models.RetrieveUserClientParams, reqEditors ...RequestEditorFn) (*RetrieveUserClientResponse, error)

	// GenerateClientSecretWithResponse request
	GenerateClientSecretWithResponse(ctx context.Context, params *models.GenerateClientSecretParams, reqEditors ...RequestEditorFn) (*GenerateClientSecretResponse, error)

	// RetrieveUserIsPasswordSetWithResponse request
	RetrieveUserIsPasswordSetWithResponse(ctx context.Context, params *models.RetrieveUserIsPasswordSetParams, reqEditors ...RequestEditorFn) (*RetrieveUserIsPasswordSetResponse, error)

	// DeleteUserWithResponse request
	DeleteUserWithResponse(ctx context.Context, userId string, params *models.DeleteUserParams, reqEditors ...RequestEditorFn) (*DeleteUserResponse, error)

	// RetrieveUserWithResponse request
	RetrieveUserWithResponse(ctx context.Context, userId string, params *models.RetrieveUserParams, reqEditors ...RequestEditorFn) (*RetrieveUserResponse, error)

	// UpdateUserWithBodyWithResponse request with any body
	UpdateUserWithBodyWithResponse(ctx context.Context, userId string, params *models.UpdateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserResponse, error)

	UpdateUserWithResponse(ctx context.Context, userId string, params *models.UpdateUserParams, body models.UpdateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserResponse, error)

	// ListObjectStorageCredentialsWithResponse request
	ListObjectStorageCredentialsWithResponse(ctx context.Context, userId string, params *models.ListObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*ListObjectStorageCredentialsResponse, error)

	// GetObjectStorageCredentialsWithResponse request
	GetObjectStorageCredentialsWithResponse(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.GetObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*GetObjectStorageCredentialsResponse, error)

	// RegenerateObjectStorageCredentialsWithResponse request
	RegenerateObjectStorageCredentialsWithResponse(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.RegenerateObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*RegenerateObjectStorageCredentialsResponse, error)

	// ResendEmailVerificationWithResponse request
	ResendEmailVerificationWithResponse(ctx context.Context, userId string, params *models.ResendEmailVerificationParams, reqEditors ...RequestEditorFn) (*ResendEmailVerificationResponse, error)

	// ResetPasswordWithResponse request
	ResetPasswordWithResponse(ctx context.Context, userId string, params *models.ResetPasswordParams, reqEditors ...RequestEditorFn) (*ResetPasswordResponse, error)

	// RetrieveVipListWithResponse request
	RetrieveVipListWithResponse(ctx context.Context, params *models.RetrieveVipListParams, reqEditors ...RequestEditorFn) (*RetrieveVipListResponse, error)

	// RetrieveVipAuditsListWithResponse request
	RetrieveVipAuditsListWithResponse(ctx context.Context, params *models.RetrieveVipAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveVipAuditsListResponse, error)

	// RetrieveVipWithResponse request
	RetrieveVipWithResponse(ctx context.Context, ip string, params *models.RetrieveVipParams, reqEditors ...RequestEditorFn) (*RetrieveVipResponse, error)

	// UnassignIpWithResponse request
	UnassignIpWithResponse(ctx context.Context, ip string, resourceType models.UnassignIpParamsResourceType, resourceId int64, params *models.UnassignIpParams, reqEditors ...RequestEditorFn) (*UnassignIpResponse, error)

	// AssignIpWithResponse request
	AssignIpWithResponse(ctx context.Context, ip string, resourceType models.AssignIpParamsResourceType, resourceId int64, params *models.AssignIpParams, reqEditors ...RequestEditorFn) (*AssignIpResponse, error)
}

type RetrieveImageListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListImageResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveImageListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveImageListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateCustomImageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateCustomImageResponse
	JSON415      *models.CreateCustomImageFailResponse
}

// Status returns HTTPResponse.Status
func (r CreateCustomImageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateCustomImageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveImageAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ImageAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveImageAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveImageAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveCustomImagesStatsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.CustomImagesStatsResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveCustomImagesStatsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveCustomImagesStatsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteImageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteImageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteImageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveImageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindImageResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveImageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveImageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateImageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateCustomImageResponse
}

// Status returns HTTPResponse.Status
func (r UpdateImageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateImageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveInstancesListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListInstancesResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveInstancesListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveInstancesListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateInstanceResponse
}

// Status returns HTTPResponse.Status
func (r CreateInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveInstancesActionsAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListInstancesActionsAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveInstancesActionsAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveInstancesActionsAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveInstancesAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListInstancesAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveInstancesAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveInstancesAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindInstanceResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PatchInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.PatchInstanceResponse
}

// Status returns HTTPResponse.Status
func (r PatchInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PatchInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ReinstallInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ReinstallInstanceResponse
}

// Status returns HTTPResponse.Status
func (r ReinstallInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReinstallInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RescueResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceRescueActionResponse
}

// Status returns HTTPResponse.Status
func (r RescueResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RescueResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResetPasswordActionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceResetPasswordActionResponse
}

// Status returns HTTPResponse.Status
func (r ResetPasswordActionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResetPasswordActionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RestartResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceRestartActionResponse
}

// Status returns HTTPResponse.Status
func (r RestartResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RestartResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ShutdownResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceShutdownActionResponse
}

// Status returns HTTPResponse.Status
func (r ShutdownResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ShutdownResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type StartResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceStartActionResponse
}

// Status returns HTTPResponse.Status
func (r StartResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r StartResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type StopResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.InstanceStopActionResponse
}

// Status returns HTTPResponse.Status
func (r StopResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r StopResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CancelInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CancelInstanceResponse
}

// Status returns HTTPResponse.Status
func (r CancelInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CancelInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSnapshotListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListSnapshotResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSnapshotListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSnapshotListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateSnapshotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateSnapshotResponse
}

// Status returns HTTPResponse.Status
func (r CreateSnapshotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateSnapshotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteSnapshotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteSnapshotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteSnapshotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSnapshotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindSnapshotResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSnapshotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSnapshotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSnapshotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateSnapshotResponse
}

// Status returns HTTPResponse.Status
func (r UpdateSnapshotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSnapshotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RollbackSnapshotResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.RollbackSnapshotResponse
}

// Status returns HTTPResponse.Status
func (r RollbackSnapshotResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RollbackSnapshotResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpgradeInstanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.PatchInstanceResponse
}

// Status returns HTTPResponse.Status
func (r UpgradeInstanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpgradeInstanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSnapshotsAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListSnapshotsAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSnapshotsAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSnapshotsAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTicketResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateTicketResponse
}

// Status returns HTTPResponse.Status
func (r CreateTicketResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTicketResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveDataCenterListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListDataCenterResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveDataCenterListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveDataCenterListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveObjectStorageListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveObjectStorageListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveObjectStorageListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateObjectStorageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r CreateObjectStorageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateObjectStorageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveObjectStorageAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListObjectStorageAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveObjectStorageAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveObjectStorageAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveObjectStorageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveObjectStorageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveObjectStorageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateObjectStorageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.CancelObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r UpdateObjectStorageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateObjectStorageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CancelObjectStorageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.CancelObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r CancelObjectStorageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CancelObjectStorageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpgradeObjectStorageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpgradeObjectStorageResponse
}

// Status returns HTTPResponse.Status
func (r UpgradeObjectStorageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpgradeObjectStorageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveObjectStoragesStatsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ObjectStoragesStatsResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveObjectStoragesStatsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveObjectStoragesStatsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrievePrivateNetworkListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListPrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r RetrievePrivateNetworkListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrievePrivateNetworkListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreatePrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreatePrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r CreatePrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreatePrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrievePrivateNetworkAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListPrivateNetworkAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrievePrivateNetworkAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrievePrivateNetworkAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeletePrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeletePrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeletePrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrievePrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindPrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r RetrievePrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrievePrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PatchPrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.PatchPrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r PatchPrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PatchPrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UnassignInstancePrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.UnassignInstancePrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r UnassignInstancePrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UnassignInstancePrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AssignInstancePrivateNetworkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.AssignInstancePrivateNetworkResponse
}

// Status returns HTTPResponse.Status
func (r AssignInstancePrivateNetworkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AssignInstancePrivateNetworkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveRoleListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListRoleResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveRoleListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveRoleListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateRoleResponse
}

// Status returns HTTPResponse.Status
func (r CreateRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveApiPermissionsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListApiPermissionResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveApiPermissionsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveApiPermissionsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveRoleAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListRoleAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveRoleAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveRoleAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindRoleResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateRoleResponse
}

// Status returns HTTPResponse.Status
func (r UpdateRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSecretListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListSecretResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSecretListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSecretListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateSecretResponse
}

// Status returns HTTPResponse.Status
func (r CreateSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSecretAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListSecretAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSecretAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSecretAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindSecretResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateSecretResponse
}

// Status returns HTTPResponse.Status
func (r UpdateSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveTagListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListTagResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveTagListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveTagListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTagResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateTagResponse
}

// Status returns HTTPResponse.Status
func (r CreateTagResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTagResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveAssignmentsAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListAssignmentAuditsResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveAssignmentsAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveAssignmentsAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveTagAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListTagAuditsResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveTagAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveTagAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTagResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteTagResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTagResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveTagResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindTagResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveTagResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveTagResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTagResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateTagResponse
}

// Status returns HTTPResponse.Status
func (r UpdateTagResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTagResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveAssignmentListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListAssignmentResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveAssignmentListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveAssignmentListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteAssignmentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteAssignmentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteAssignmentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveAssignmentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindAssignmentResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveAssignmentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveAssignmentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateAssignmentResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateAssignmentResponse
}

// Status returns HTTPResponse.Status
func (r CreateAssignmentResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateAssignmentResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveUserListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListUserResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveUserListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveUserListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *models.CreateUserResponse
}

// Status returns HTTPResponse.Status
func (r CreateUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveUserAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListUserAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveUserAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveUserAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveUserClientResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindClientResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveUserClientResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveUserClientResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GenerateClientSecretResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.GenerateClientSecretResponse
}

// Status returns HTTPResponse.Status
func (r GenerateClientSecretResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GenerateClientSecretResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveUserIsPasswordSetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindUserIsPasswordSetResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveUserIsPasswordSetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveUserIsPasswordSetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindUserResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.UpdateUserResponse
}

// Status returns HTTPResponse.Status
func (r UpdateUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListObjectStorageCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListCredentialResponse
}

// Status returns HTTPResponse.Status
func (r ListObjectStorageCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListObjectStorageCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetObjectStorageCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindCredentialResponse
}

// Status returns HTTPResponse.Status
func (r GetObjectStorageCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetObjectStorageCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RegenerateObjectStorageCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindCredentialResponse
}

// Status returns HTTPResponse.Status
func (r RegenerateObjectStorageCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RegenerateObjectStorageCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResendEmailVerificationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ResendEmailVerificationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResendEmailVerificationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResetPasswordResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r ResetPasswordResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResetPasswordResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveVipListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListVipResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveVipListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveVipListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveVipAuditsListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.ListVipAuditResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveVipAuditsListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveVipAuditsListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RetrieveVipResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.FindVipResponse
}

// Status returns HTTPResponse.Status
func (r RetrieveVipResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RetrieveVipResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UnassignIpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r UnassignIpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UnassignIpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AssignIpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *models.AssignVipResponse
}

// Status returns HTTPResponse.Status
func (r AssignIpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AssignIpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// RetrieveImageListWithResponse request returning *RetrieveImageListResponse
func (c *ClientWithResponses) RetrieveImageListWithResponse(ctx context.Context, params *models.RetrieveImageListParams, reqEditors ...RequestEditorFn) (*RetrieveImageListResponse, error) {
	rsp, err := c.RetrieveImageList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveImageListResponse(rsp)
}

// CreateCustomImageWithBodyWithResponse request with arbitrary body returning *CreateCustomImageResponse
func (c *ClientWithResponses) CreateCustomImageWithBodyWithResponse(ctx context.Context, params *models.CreateCustomImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateCustomImageResponse, error) {
	rsp, err := c.CreateCustomImageWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateCustomImageResponse(rsp)
}

func (c *ClientWithResponses) CreateCustomImageWithResponse(ctx context.Context, params *models.CreateCustomImageParams, body models.CreateCustomImageJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateCustomImageResponse, error) {
	rsp, err := c.CreateCustomImage(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateCustomImageResponse(rsp)
}

// RetrieveImageAuditsListWithResponse request returning *RetrieveImageAuditsListResponse
func (c *ClientWithResponses) RetrieveImageAuditsListWithResponse(ctx context.Context, params *models.RetrieveImageAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveImageAuditsListResponse, error) {
	rsp, err := c.RetrieveImageAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveImageAuditsListResponse(rsp)
}

// RetrieveCustomImagesStatsWithResponse request returning *RetrieveCustomImagesStatsResponse
func (c *ClientWithResponses) RetrieveCustomImagesStatsWithResponse(ctx context.Context, params *models.RetrieveCustomImagesStatsParams, reqEditors ...RequestEditorFn) (*RetrieveCustomImagesStatsResponse, error) {
	rsp, err := c.RetrieveCustomImagesStats(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveCustomImagesStatsResponse(rsp)
}

// DeleteImageWithResponse request returning *DeleteImageResponse
func (c *ClientWithResponses) DeleteImageWithResponse(ctx context.Context, imageId string, params *models.DeleteImageParams, reqEditors ...RequestEditorFn) (*DeleteImageResponse, error) {
	rsp, err := c.DeleteImage(ctx, imageId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteImageResponse(rsp)
}

// RetrieveImageWithResponse request returning *RetrieveImageResponse
func (c *ClientWithResponses) RetrieveImageWithResponse(ctx context.Context, imageId string, params *models.RetrieveImageParams, reqEditors ...RequestEditorFn) (*RetrieveImageResponse, error) {
	rsp, err := c.RetrieveImage(ctx, imageId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveImageResponse(rsp)
}

// UpdateImageWithBodyWithResponse request with arbitrary body returning *UpdateImageResponse
func (c *ClientWithResponses) UpdateImageWithBodyWithResponse(ctx context.Context, imageId string, params *models.UpdateImageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateImageResponse, error) {
	rsp, err := c.UpdateImageWithBody(ctx, imageId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateImageResponse(rsp)
}

func (c *ClientWithResponses) UpdateImageWithResponse(ctx context.Context, imageId string, params *models.UpdateImageParams, body models.UpdateImageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateImageResponse, error) {
	rsp, err := c.UpdateImage(ctx, imageId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateImageResponse(rsp)
}

// RetrieveInstancesListWithResponse request returning *RetrieveInstancesListResponse
func (c *ClientWithResponses) RetrieveInstancesListWithResponse(ctx context.Context, params *models.RetrieveInstancesListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesListResponse, error) {
	rsp, err := c.RetrieveInstancesList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveInstancesListResponse(rsp)
}

// CreateInstanceWithBodyWithResponse request with arbitrary body returning *CreateInstanceResponse
func (c *ClientWithResponses) CreateInstanceWithBodyWithResponse(ctx context.Context, params *models.CreateInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateInstanceResponse, error) {
	rsp, err := c.CreateInstanceWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateInstanceResponse(rsp)
}

func (c *ClientWithResponses) CreateInstanceWithResponse(ctx context.Context, params *models.CreateInstanceParams, body models.CreateInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateInstanceResponse, error) {
	rsp, err := c.CreateInstance(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateInstanceResponse(rsp)
}

// RetrieveInstancesActionsAuditsListWithResponse request returning *RetrieveInstancesActionsAuditsListResponse
func (c *ClientWithResponses) RetrieveInstancesActionsAuditsListWithResponse(ctx context.Context, params *models.RetrieveInstancesActionsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesActionsAuditsListResponse, error) {
	rsp, err := c.RetrieveInstancesActionsAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveInstancesActionsAuditsListResponse(rsp)
}

// RetrieveInstancesAuditsListWithResponse request returning *RetrieveInstancesAuditsListResponse
func (c *ClientWithResponses) RetrieveInstancesAuditsListWithResponse(ctx context.Context, params *models.RetrieveInstancesAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveInstancesAuditsListResponse, error) {
	rsp, err := c.RetrieveInstancesAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveInstancesAuditsListResponse(rsp)
}

// RetrieveInstanceWithResponse request returning *RetrieveInstanceResponse
func (c *ClientWithResponses) RetrieveInstanceWithResponse(ctx context.Context, instanceId int64, params *models.RetrieveInstanceParams, reqEditors ...RequestEditorFn) (*RetrieveInstanceResponse, error) {
	rsp, err := c.RetrieveInstance(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveInstanceResponse(rsp)
}

// PatchInstanceWithBodyWithResponse request with arbitrary body returning *PatchInstanceResponse
func (c *ClientWithResponses) PatchInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchInstanceResponse, error) {
	rsp, err := c.PatchInstanceWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchInstanceResponse(rsp)
}

func (c *ClientWithResponses) PatchInstanceWithResponse(ctx context.Context, instanceId int64, params *models.PatchInstanceParams, body models.PatchInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchInstanceResponse, error) {
	rsp, err := c.PatchInstance(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchInstanceResponse(rsp)
}

// ReinstallInstanceWithBodyWithResponse request with arbitrary body returning *ReinstallInstanceResponse
func (c *ClientWithResponses) ReinstallInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ReinstallInstanceResponse, error) {
	rsp, err := c.ReinstallInstanceWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReinstallInstanceResponse(rsp)
}

func (c *ClientWithResponses) ReinstallInstanceWithResponse(ctx context.Context, instanceId int64, params *models.ReinstallInstanceParams, body models.ReinstallInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*ReinstallInstanceResponse, error) {
	rsp, err := c.ReinstallInstance(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReinstallInstanceResponse(rsp)
}

// RescueWithBodyWithResponse request with arbitrary body returning *RescueResponse
func (c *ClientWithResponses) RescueWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.RescueParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RescueResponse, error) {
	rsp, err := c.RescueWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRescueResponse(rsp)
}

func (c *ClientWithResponses) RescueWithResponse(ctx context.Context, instanceId int64, params *models.RescueParams, body models.RescueJSONRequestBody, reqEditors ...RequestEditorFn) (*RescueResponse, error) {
	rsp, err := c.Rescue(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRescueResponse(rsp)
}

// ResetPasswordActionWithBodyWithResponse request with arbitrary body returning *ResetPasswordActionResponse
func (c *ClientWithResponses) ResetPasswordActionWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ResetPasswordActionResponse, error) {
	rsp, err := c.ResetPasswordActionWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResetPasswordActionResponse(rsp)
}

func (c *ClientWithResponses) ResetPasswordActionWithResponse(ctx context.Context, instanceId int64, params *models.ResetPasswordActionParams, body models.ResetPasswordActionJSONRequestBody, reqEditors ...RequestEditorFn) (*ResetPasswordActionResponse, error) {
	rsp, err := c.ResetPasswordAction(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResetPasswordActionResponse(rsp)
}

// RestartWithResponse request returning *RestartResponse
func (c *ClientWithResponses) RestartWithResponse(ctx context.Context, instanceId int64, params *models.RestartParams, reqEditors ...RequestEditorFn) (*RestartResponse, error) {
	rsp, err := c.Restart(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRestartResponse(rsp)
}

// ShutdownWithResponse request returning *ShutdownResponse
func (c *ClientWithResponses) ShutdownWithResponse(ctx context.Context, instanceId int64, params *models.ShutdownParams, reqEditors ...RequestEditorFn) (*ShutdownResponse, error) {
	rsp, err := c.Shutdown(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseShutdownResponse(rsp)
}

// StartWithResponse request returning *StartResponse
func (c *ClientWithResponses) StartWithResponse(ctx context.Context, instanceId int64, params *models.StartParams, reqEditors ...RequestEditorFn) (*StartResponse, error) {
	rsp, err := c.Start(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStartResponse(rsp)
}

// StopWithResponse request returning *StopResponse
func (c *ClientWithResponses) StopWithResponse(ctx context.Context, instanceId int64, params *models.StopParams, reqEditors ...RequestEditorFn) (*StopResponse, error) {
	rsp, err := c.Stop(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStopResponse(rsp)
}

// CancelInstanceWithBodyWithResponse request with arbitrary body returning *CancelInstanceResponse
func (c *ClientWithResponses) CancelInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelInstanceResponse, error) {
	rsp, err := c.CancelInstanceWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelInstanceResponse(rsp)
}

func (c *ClientWithResponses) CancelInstanceWithResponse(ctx context.Context, instanceId int64, params *models.CancelInstanceParams, body models.CancelInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelInstanceResponse, error) {
	rsp, err := c.CancelInstance(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelInstanceResponse(rsp)
}

// RetrieveSnapshotListWithResponse request returning *RetrieveSnapshotListResponse
func (c *ClientWithResponses) RetrieveSnapshotListWithResponse(ctx context.Context, instanceId int64, params *models.RetrieveSnapshotListParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotListResponse, error) {
	rsp, err := c.RetrieveSnapshotList(ctx, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSnapshotListResponse(rsp)
}

// CreateSnapshotWithBodyWithResponse request with arbitrary body returning *CreateSnapshotResponse
func (c *ClientWithResponses) CreateSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSnapshotResponse, error) {
	rsp, err := c.CreateSnapshotWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSnapshotResponse(rsp)
}

func (c *ClientWithResponses) CreateSnapshotWithResponse(ctx context.Context, instanceId int64, params *models.CreateSnapshotParams, body models.CreateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSnapshotResponse, error) {
	rsp, err := c.CreateSnapshot(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSnapshotResponse(rsp)
}

// DeleteSnapshotWithResponse request returning *DeleteSnapshotResponse
func (c *ClientWithResponses) DeleteSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.DeleteSnapshotParams, reqEditors ...RequestEditorFn) (*DeleteSnapshotResponse, error) {
	rsp, err := c.DeleteSnapshot(ctx, instanceId, snapshotId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteSnapshotResponse(rsp)
}

// RetrieveSnapshotWithResponse request returning *RetrieveSnapshotResponse
func (c *ClientWithResponses) RetrieveSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RetrieveSnapshotParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotResponse, error) {
	rsp, err := c.RetrieveSnapshot(ctx, instanceId, snapshotId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSnapshotResponse(rsp)
}

// UpdateSnapshotWithBodyWithResponse request with arbitrary body returning *UpdateSnapshotResponse
func (c *ClientWithResponses) UpdateSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSnapshotResponse, error) {
	rsp, err := c.UpdateSnapshotWithBody(ctx, instanceId, snapshotId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSnapshotResponse(rsp)
}

func (c *ClientWithResponses) UpdateSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.UpdateSnapshotParams, body models.UpdateSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSnapshotResponse, error) {
	rsp, err := c.UpdateSnapshot(ctx, instanceId, snapshotId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSnapshotResponse(rsp)
}

// RollbackSnapshotWithBodyWithResponse request with arbitrary body returning *RollbackSnapshotResponse
func (c *ClientWithResponses) RollbackSnapshotWithBodyWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RollbackSnapshotResponse, error) {
	rsp, err := c.RollbackSnapshotWithBody(ctx, instanceId, snapshotId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRollbackSnapshotResponse(rsp)
}

func (c *ClientWithResponses) RollbackSnapshotWithResponse(ctx context.Context, instanceId int64, snapshotId string, params *models.RollbackSnapshotParams, body models.RollbackSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*RollbackSnapshotResponse, error) {
	rsp, err := c.RollbackSnapshot(ctx, instanceId, snapshotId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRollbackSnapshotResponse(rsp)
}

// UpgradeInstanceWithBodyWithResponse request with arbitrary body returning *UpgradeInstanceResponse
func (c *ClientWithResponses) UpgradeInstanceWithBodyWithResponse(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpgradeInstanceResponse, error) {
	rsp, err := c.UpgradeInstanceWithBody(ctx, instanceId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpgradeInstanceResponse(rsp)
}

func (c *ClientWithResponses) UpgradeInstanceWithResponse(ctx context.Context, instanceId int64, params *models.UpgradeInstanceParams, body models.UpgradeInstanceJSONRequestBody, reqEditors ...RequestEditorFn) (*UpgradeInstanceResponse, error) {
	rsp, err := c.UpgradeInstance(ctx, instanceId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpgradeInstanceResponse(rsp)
}

// RetrieveSnapshotsAuditsListWithResponse request returning *RetrieveSnapshotsAuditsListResponse
func (c *ClientWithResponses) RetrieveSnapshotsAuditsListWithResponse(ctx context.Context, params *models.RetrieveSnapshotsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveSnapshotsAuditsListResponse, error) {
	rsp, err := c.RetrieveSnapshotsAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSnapshotsAuditsListResponse(rsp)
}

// CreateTicketWithBodyWithResponse request with arbitrary body returning *CreateTicketResponse
func (c *ClientWithResponses) CreateTicketWithBodyWithResponse(ctx context.Context, params *models.CreateTicketParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTicketResponse, error) {
	rsp, err := c.CreateTicketWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTicketResponse(rsp)
}

func (c *ClientWithResponses) CreateTicketWithResponse(ctx context.Context, params *models.CreateTicketParams, body models.CreateTicketJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTicketResponse, error) {
	rsp, err := c.CreateTicket(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTicketResponse(rsp)
}

// RetrieveDataCenterListWithResponse request returning *RetrieveDataCenterListResponse
func (c *ClientWithResponses) RetrieveDataCenterListWithResponse(ctx context.Context, params *models.RetrieveDataCenterListParams, reqEditors ...RequestEditorFn) (*RetrieveDataCenterListResponse, error) {
	rsp, err := c.RetrieveDataCenterList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveDataCenterListResponse(rsp)
}

// RetrieveObjectStorageListWithResponse request returning *RetrieveObjectStorageListResponse
func (c *ClientWithResponses) RetrieveObjectStorageListWithResponse(ctx context.Context, params *models.RetrieveObjectStorageListParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageListResponse, error) {
	rsp, err := c.RetrieveObjectStorageList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveObjectStorageListResponse(rsp)
}

// CreateObjectStorageWithBodyWithResponse request with arbitrary body returning *CreateObjectStorageResponse
func (c *ClientWithResponses) CreateObjectStorageWithBodyWithResponse(ctx context.Context, params *models.CreateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateObjectStorageResponse, error) {
	rsp, err := c.CreateObjectStorageWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateObjectStorageResponse(rsp)
}

func (c *ClientWithResponses) CreateObjectStorageWithResponse(ctx context.Context, params *models.CreateObjectStorageParams, body models.CreateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateObjectStorageResponse, error) {
	rsp, err := c.CreateObjectStorage(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateObjectStorageResponse(rsp)
}

// RetrieveObjectStorageAuditsListWithResponse request returning *RetrieveObjectStorageAuditsListResponse
func (c *ClientWithResponses) RetrieveObjectStorageAuditsListWithResponse(ctx context.Context, params *models.RetrieveObjectStorageAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageAuditsListResponse, error) {
	rsp, err := c.RetrieveObjectStorageAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveObjectStorageAuditsListResponse(rsp)
}

// RetrieveObjectStorageWithResponse request returning *RetrieveObjectStorageResponse
func (c *ClientWithResponses) RetrieveObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStorageParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStorageResponse, error) {
	rsp, err := c.RetrieveObjectStorage(ctx, objectStorageId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveObjectStorageResponse(rsp)
}

// UpdateObjectStorageWithBodyWithResponse request with arbitrary body returning *UpdateObjectStorageResponse
func (c *ClientWithResponses) UpdateObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateObjectStorageResponse, error) {
	rsp, err := c.UpdateObjectStorageWithBody(ctx, objectStorageId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateObjectStorageResponse(rsp)
}

func (c *ClientWithResponses) UpdateObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.UpdateObjectStorageParams, body models.UpdateObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateObjectStorageResponse, error) {
	rsp, err := c.UpdateObjectStorage(ctx, objectStorageId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateObjectStorageResponse(rsp)
}

// CancelObjectStorageWithBodyWithResponse request with arbitrary body returning *CancelObjectStorageResponse
func (c *ClientWithResponses) CancelObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CancelObjectStorageResponse, error) {
	rsp, err := c.CancelObjectStorageWithBody(ctx, objectStorageId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelObjectStorageResponse(rsp)
}

func (c *ClientWithResponses) CancelObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.CancelObjectStorageParams, body models.CancelObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*CancelObjectStorageResponse, error) {
	rsp, err := c.CancelObjectStorage(ctx, objectStorageId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelObjectStorageResponse(rsp)
}

// UpgradeObjectStorageWithBodyWithResponse request with arbitrary body returning *UpgradeObjectStorageResponse
func (c *ClientWithResponses) UpgradeObjectStorageWithBodyWithResponse(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpgradeObjectStorageResponse, error) {
	rsp, err := c.UpgradeObjectStorageWithBody(ctx, objectStorageId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpgradeObjectStorageResponse(rsp)
}

func (c *ClientWithResponses) UpgradeObjectStorageWithResponse(ctx context.Context, objectStorageId string, params *models.UpgradeObjectStorageParams, body models.UpgradeObjectStorageJSONRequestBody, reqEditors ...RequestEditorFn) (*UpgradeObjectStorageResponse, error) {
	rsp, err := c.UpgradeObjectStorage(ctx, objectStorageId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpgradeObjectStorageResponse(rsp)
}

// RetrieveObjectStoragesStatsWithResponse request returning *RetrieveObjectStoragesStatsResponse
func (c *ClientWithResponses) RetrieveObjectStoragesStatsWithResponse(ctx context.Context, objectStorageId string, params *models.RetrieveObjectStoragesStatsParams, reqEditors ...RequestEditorFn) (*RetrieveObjectStoragesStatsResponse, error) {
	rsp, err := c.RetrieveObjectStoragesStats(ctx, objectStorageId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveObjectStoragesStatsResponse(rsp)
}

// RetrievePrivateNetworkListWithResponse request returning *RetrievePrivateNetworkListResponse
func (c *ClientWithResponses) RetrievePrivateNetworkListWithResponse(ctx context.Context, params *models.RetrievePrivateNetworkListParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkListResponse, error) {
	rsp, err := c.RetrievePrivateNetworkList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrievePrivateNetworkListResponse(rsp)
}

// CreatePrivateNetworkWithBodyWithResponse request with arbitrary body returning *CreatePrivateNetworkResponse
func (c *ClientWithResponses) CreatePrivateNetworkWithBodyWithResponse(ctx context.Context, params *models.CreatePrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreatePrivateNetworkResponse, error) {
	rsp, err := c.CreatePrivateNetworkWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreatePrivateNetworkResponse(rsp)
}

func (c *ClientWithResponses) CreatePrivateNetworkWithResponse(ctx context.Context, params *models.CreatePrivateNetworkParams, body models.CreatePrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*CreatePrivateNetworkResponse, error) {
	rsp, err := c.CreatePrivateNetwork(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreatePrivateNetworkResponse(rsp)
}

// RetrievePrivateNetworkAuditsListWithResponse request returning *RetrievePrivateNetworkAuditsListResponse
func (c *ClientWithResponses) RetrievePrivateNetworkAuditsListWithResponse(ctx context.Context, params *models.RetrievePrivateNetworkAuditsListParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkAuditsListResponse, error) {
	rsp, err := c.RetrievePrivateNetworkAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrievePrivateNetworkAuditsListResponse(rsp)
}

// DeletePrivateNetworkWithResponse request returning *DeletePrivateNetworkResponse
func (c *ClientWithResponses) DeletePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.DeletePrivateNetworkParams, reqEditors ...RequestEditorFn) (*DeletePrivateNetworkResponse, error) {
	rsp, err := c.DeletePrivateNetwork(ctx, privateNetworkId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeletePrivateNetworkResponse(rsp)
}

// RetrievePrivateNetworkWithResponse request returning *RetrievePrivateNetworkResponse
func (c *ClientWithResponses) RetrievePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.RetrievePrivateNetworkParams, reqEditors ...RequestEditorFn) (*RetrievePrivateNetworkResponse, error) {
	rsp, err := c.RetrievePrivateNetwork(ctx, privateNetworkId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrievePrivateNetworkResponse(rsp)
}

// PatchPrivateNetworkWithBodyWithResponse request with arbitrary body returning *PatchPrivateNetworkResponse
func (c *ClientWithResponses) PatchPrivateNetworkWithBodyWithResponse(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchPrivateNetworkResponse, error) {
	rsp, err := c.PatchPrivateNetworkWithBody(ctx, privateNetworkId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchPrivateNetworkResponse(rsp)
}

func (c *ClientWithResponses) PatchPrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, params *models.PatchPrivateNetworkParams, body models.PatchPrivateNetworkJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchPrivateNetworkResponse, error) {
	rsp, err := c.PatchPrivateNetwork(ctx, privateNetworkId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchPrivateNetworkResponse(rsp)
}

// UnassignInstancePrivateNetworkWithResponse request returning *UnassignInstancePrivateNetworkResponse
func (c *ClientWithResponses) UnassignInstancePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.UnassignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*UnassignInstancePrivateNetworkResponse, error) {
	rsp, err := c.UnassignInstancePrivateNetwork(ctx, privateNetworkId, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUnassignInstancePrivateNetworkResponse(rsp)
}

// AssignInstancePrivateNetworkWithResponse request returning *AssignInstancePrivateNetworkResponse
func (c *ClientWithResponses) AssignInstancePrivateNetworkWithResponse(ctx context.Context, privateNetworkId int64, instanceId int64, params *models.AssignInstancePrivateNetworkParams, reqEditors ...RequestEditorFn) (*AssignInstancePrivateNetworkResponse, error) {
	rsp, err := c.AssignInstancePrivateNetwork(ctx, privateNetworkId, instanceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAssignInstancePrivateNetworkResponse(rsp)
}

// RetrieveRoleListWithResponse request returning *RetrieveRoleListResponse
func (c *ClientWithResponses) RetrieveRoleListWithResponse(ctx context.Context, params *models.RetrieveRoleListParams, reqEditors ...RequestEditorFn) (*RetrieveRoleListResponse, error) {
	rsp, err := c.RetrieveRoleList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveRoleListResponse(rsp)
}

// CreateRoleWithBodyWithResponse request with arbitrary body returning *CreateRoleResponse
func (c *ClientWithResponses) CreateRoleWithBodyWithResponse(ctx context.Context, params *models.CreateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error) {
	rsp, err := c.CreateRoleWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoleResponse(rsp)
}

func (c *ClientWithResponses) CreateRoleWithResponse(ctx context.Context, params *models.CreateRoleParams, body models.CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error) {
	rsp, err := c.CreateRole(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoleResponse(rsp)
}

// RetrieveApiPermissionsListWithResponse request returning *RetrieveApiPermissionsListResponse
func (c *ClientWithResponses) RetrieveApiPermissionsListWithResponse(ctx context.Context, params *models.RetrieveApiPermissionsListParams, reqEditors ...RequestEditorFn) (*RetrieveApiPermissionsListResponse, error) {
	rsp, err := c.RetrieveApiPermissionsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveApiPermissionsListResponse(rsp)
}

// RetrieveRoleAuditsListWithResponse request returning *RetrieveRoleAuditsListResponse
func (c *ClientWithResponses) RetrieveRoleAuditsListWithResponse(ctx context.Context, params *models.RetrieveRoleAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveRoleAuditsListResponse, error) {
	rsp, err := c.RetrieveRoleAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveRoleAuditsListResponse(rsp)
}

// DeleteRoleWithResponse request returning *DeleteRoleResponse
func (c *ClientWithResponses) DeleteRoleWithResponse(ctx context.Context, roleId int64, params *models.DeleteRoleParams, reqEditors ...RequestEditorFn) (*DeleteRoleResponse, error) {
	rsp, err := c.DeleteRole(ctx, roleId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteRoleResponse(rsp)
}

// RetrieveRoleWithResponse request returning *RetrieveRoleResponse
func (c *ClientWithResponses) RetrieveRoleWithResponse(ctx context.Context, roleId int64, params *models.RetrieveRoleParams, reqEditors ...RequestEditorFn) (*RetrieveRoleResponse, error) {
	rsp, err := c.RetrieveRole(ctx, roleId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveRoleResponse(rsp)
}

// UpdateRoleWithBodyWithResponse request with arbitrary body returning *UpdateRoleResponse
func (c *ClientWithResponses) UpdateRoleWithBodyWithResponse(ctx context.Context, roleId int64, params *models.UpdateRoleParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateRoleResponse, error) {
	rsp, err := c.UpdateRoleWithBody(ctx, roleId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateRoleResponse(rsp)
}

func (c *ClientWithResponses) UpdateRoleWithResponse(ctx context.Context, roleId int64, params *models.UpdateRoleParams, body models.UpdateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateRoleResponse, error) {
	rsp, err := c.UpdateRole(ctx, roleId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateRoleResponse(rsp)
}

// RetrieveSecretListWithResponse request returning *RetrieveSecretListResponse
func (c *ClientWithResponses) RetrieveSecretListWithResponse(ctx context.Context, params *models.RetrieveSecretListParams, reqEditors ...RequestEditorFn) (*RetrieveSecretListResponse, error) {
	rsp, err := c.RetrieveSecretList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSecretListResponse(rsp)
}

// CreateSecretWithBodyWithResponse request with arbitrary body returning *CreateSecretResponse
func (c *ClientWithResponses) CreateSecretWithBodyWithResponse(ctx context.Context, params *models.CreateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSecretResponse, error) {
	rsp, err := c.CreateSecretWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSecretResponse(rsp)
}

func (c *ClientWithResponses) CreateSecretWithResponse(ctx context.Context, params *models.CreateSecretParams, body models.CreateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSecretResponse, error) {
	rsp, err := c.CreateSecret(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSecretResponse(rsp)
}

// RetrieveSecretAuditsListWithResponse request returning *RetrieveSecretAuditsListResponse
func (c *ClientWithResponses) RetrieveSecretAuditsListWithResponse(ctx context.Context, params *models.RetrieveSecretAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveSecretAuditsListResponse, error) {
	rsp, err := c.RetrieveSecretAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSecretAuditsListResponse(rsp)
}

// DeleteSecretWithResponse request returning *DeleteSecretResponse
func (c *ClientWithResponses) DeleteSecretWithResponse(ctx context.Context, secretId int64, params *models.DeleteSecretParams, reqEditors ...RequestEditorFn) (*DeleteSecretResponse, error) {
	rsp, err := c.DeleteSecret(ctx, secretId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteSecretResponse(rsp)
}

// RetrieveSecretWithResponse request returning *RetrieveSecretResponse
func (c *ClientWithResponses) RetrieveSecretWithResponse(ctx context.Context, secretId int64, params *models.RetrieveSecretParams, reqEditors ...RequestEditorFn) (*RetrieveSecretResponse, error) {
	rsp, err := c.RetrieveSecret(ctx, secretId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveSecretResponse(rsp)
}

// UpdateSecretWithBodyWithResponse request with arbitrary body returning *UpdateSecretResponse
func (c *ClientWithResponses) UpdateSecretWithBodyWithResponse(ctx context.Context, secretId int64, params *models.UpdateSecretParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSecretResponse, error) {
	rsp, err := c.UpdateSecretWithBody(ctx, secretId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSecretResponse(rsp)
}

func (c *ClientWithResponses) UpdateSecretWithResponse(ctx context.Context, secretId int64, params *models.UpdateSecretParams, body models.UpdateSecretJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSecretResponse, error) {
	rsp, err := c.UpdateSecret(ctx, secretId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSecretResponse(rsp)
}

// RetrieveTagListWithResponse request returning *RetrieveTagListResponse
func (c *ClientWithResponses) RetrieveTagListWithResponse(ctx context.Context, params *models.RetrieveTagListParams, reqEditors ...RequestEditorFn) (*RetrieveTagListResponse, error) {
	rsp, err := c.RetrieveTagList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveTagListResponse(rsp)
}

// CreateTagWithBodyWithResponse request with arbitrary body returning *CreateTagResponse
func (c *ClientWithResponses) CreateTagWithBodyWithResponse(ctx context.Context, params *models.CreateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTagResponse, error) {
	rsp, err := c.CreateTagWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTagResponse(rsp)
}

func (c *ClientWithResponses) CreateTagWithResponse(ctx context.Context, params *models.CreateTagParams, body models.CreateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTagResponse, error) {
	rsp, err := c.CreateTag(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTagResponse(rsp)
}

// RetrieveAssignmentsAuditsListWithResponse request returning *RetrieveAssignmentsAuditsListResponse
func (c *ClientWithResponses) RetrieveAssignmentsAuditsListWithResponse(ctx context.Context, params *models.RetrieveAssignmentsAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentsAuditsListResponse, error) {
	rsp, err := c.RetrieveAssignmentsAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveAssignmentsAuditsListResponse(rsp)
}

// RetrieveTagAuditsListWithResponse request returning *RetrieveTagAuditsListResponse
func (c *ClientWithResponses) RetrieveTagAuditsListWithResponse(ctx context.Context, params *models.RetrieveTagAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveTagAuditsListResponse, error) {
	rsp, err := c.RetrieveTagAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveTagAuditsListResponse(rsp)
}

// DeleteTagWithResponse request returning *DeleteTagResponse
func (c *ClientWithResponses) DeleteTagWithResponse(ctx context.Context, tagId int64, params *models.DeleteTagParams, reqEditors ...RequestEditorFn) (*DeleteTagResponse, error) {
	rsp, err := c.DeleteTag(ctx, tagId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTagResponse(rsp)
}

// RetrieveTagWithResponse request returning *RetrieveTagResponse
func (c *ClientWithResponses) RetrieveTagWithResponse(ctx context.Context, tagId int64, params *models.RetrieveTagParams, reqEditors ...RequestEditorFn) (*RetrieveTagResponse, error) {
	rsp, err := c.RetrieveTag(ctx, tagId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveTagResponse(rsp)
}

// UpdateTagWithBodyWithResponse request with arbitrary body returning *UpdateTagResponse
func (c *ClientWithResponses) UpdateTagWithBodyWithResponse(ctx context.Context, tagId int64, params *models.UpdateTagParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTagResponse, error) {
	rsp, err := c.UpdateTagWithBody(ctx, tagId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTagResponse(rsp)
}

func (c *ClientWithResponses) UpdateTagWithResponse(ctx context.Context, tagId int64, params *models.UpdateTagParams, body models.UpdateTagJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTagResponse, error) {
	rsp, err := c.UpdateTag(ctx, tagId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTagResponse(rsp)
}

// RetrieveAssignmentListWithResponse request returning *RetrieveAssignmentListResponse
func (c *ClientWithResponses) RetrieveAssignmentListWithResponse(ctx context.Context, tagId int64, params *models.RetrieveAssignmentListParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentListResponse, error) {
	rsp, err := c.RetrieveAssignmentList(ctx, tagId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveAssignmentListResponse(rsp)
}

// DeleteAssignmentWithResponse request returning *DeleteAssignmentResponse
func (c *ClientWithResponses) DeleteAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.DeleteAssignmentParams, reqEditors ...RequestEditorFn) (*DeleteAssignmentResponse, error) {
	rsp, err := c.DeleteAssignment(ctx, tagId, resourceType, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteAssignmentResponse(rsp)
}

// RetrieveAssignmentWithResponse request returning *RetrieveAssignmentResponse
func (c *ClientWithResponses) RetrieveAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.RetrieveAssignmentParams, reqEditors ...RequestEditorFn) (*RetrieveAssignmentResponse, error) {
	rsp, err := c.RetrieveAssignment(ctx, tagId, resourceType, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveAssignmentResponse(rsp)
}

// CreateAssignmentWithResponse request returning *CreateAssignmentResponse
func (c *ClientWithResponses) CreateAssignmentWithResponse(ctx context.Context, tagId int64, resourceType string, resourceId string, params *models.CreateAssignmentParams, reqEditors ...RequestEditorFn) (*CreateAssignmentResponse, error) {
	rsp, err := c.CreateAssignment(ctx, tagId, resourceType, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateAssignmentResponse(rsp)
}

// RetrieveUserListWithResponse request returning *RetrieveUserListResponse
func (c *ClientWithResponses) RetrieveUserListWithResponse(ctx context.Context, params *models.RetrieveUserListParams, reqEditors ...RequestEditorFn) (*RetrieveUserListResponse, error) {
	rsp, err := c.RetrieveUserList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveUserListResponse(rsp)
}

// CreateUserWithBodyWithResponse request with arbitrary body returning *CreateUserResponse
func (c *ClientWithResponses) CreateUserWithBodyWithResponse(ctx context.Context, params *models.CreateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateUserResponse, error) {
	rsp, err := c.CreateUserWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateUserResponse(rsp)
}

func (c *ClientWithResponses) CreateUserWithResponse(ctx context.Context, params *models.CreateUserParams, body models.CreateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateUserResponse, error) {
	rsp, err := c.CreateUser(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateUserResponse(rsp)
}

// RetrieveUserAuditsListWithResponse request returning *RetrieveUserAuditsListResponse
func (c *ClientWithResponses) RetrieveUserAuditsListWithResponse(ctx context.Context, params *models.RetrieveUserAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveUserAuditsListResponse, error) {
	rsp, err := c.RetrieveUserAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveUserAuditsListResponse(rsp)
}

// RetrieveUserClientWithResponse request returning *RetrieveUserClientResponse
func (c *ClientWithResponses) RetrieveUserClientWithResponse(ctx context.Context, params *models.RetrieveUserClientParams, reqEditors ...RequestEditorFn) (*RetrieveUserClientResponse, error) {
	rsp, err := c.RetrieveUserClient(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveUserClientResponse(rsp)
}

// GenerateClientSecretWithResponse request returning *GenerateClientSecretResponse
func (c *ClientWithResponses) GenerateClientSecretWithResponse(ctx context.Context, params *models.GenerateClientSecretParams, reqEditors ...RequestEditorFn) (*GenerateClientSecretResponse, error) {
	rsp, err := c.GenerateClientSecret(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGenerateClientSecretResponse(rsp)
}

// RetrieveUserIsPasswordSetWithResponse request returning *RetrieveUserIsPasswordSetResponse
func (c *ClientWithResponses) RetrieveUserIsPasswordSetWithResponse(ctx context.Context, params *models.RetrieveUserIsPasswordSetParams, reqEditors ...RequestEditorFn) (*RetrieveUserIsPasswordSetResponse, error) {
	rsp, err := c.RetrieveUserIsPasswordSet(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveUserIsPasswordSetResponse(rsp)
}

// DeleteUserWithResponse request returning *DeleteUserResponse
func (c *ClientWithResponses) DeleteUserWithResponse(ctx context.Context, userId string, params *models.DeleteUserParams, reqEditors ...RequestEditorFn) (*DeleteUserResponse, error) {
	rsp, err := c.DeleteUser(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteUserResponse(rsp)
}

// RetrieveUserWithResponse request returning *RetrieveUserResponse
func (c *ClientWithResponses) RetrieveUserWithResponse(ctx context.Context, userId string, params *models.RetrieveUserParams, reqEditors ...RequestEditorFn) (*RetrieveUserResponse, error) {
	rsp, err := c.RetrieveUser(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveUserResponse(rsp)
}

// UpdateUserWithBodyWithResponse request with arbitrary body returning *UpdateUserResponse
func (c *ClientWithResponses) UpdateUserWithBodyWithResponse(ctx context.Context, userId string, params *models.UpdateUserParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserResponse, error) {
	rsp, err := c.UpdateUserWithBody(ctx, userId, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserResponse(rsp)
}

func (c *ClientWithResponses) UpdateUserWithResponse(ctx context.Context, userId string, params *models.UpdateUserParams, body models.UpdateUserJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserResponse, error) {
	rsp, err := c.UpdateUser(ctx, userId, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserResponse(rsp)
}

// ListObjectStorageCredentialsWithResponse request returning *ListObjectStorageCredentialsResponse
func (c *ClientWithResponses) ListObjectStorageCredentialsWithResponse(ctx context.Context, userId string, params *models.ListObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*ListObjectStorageCredentialsResponse, error) {
	rsp, err := c.ListObjectStorageCredentials(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListObjectStorageCredentialsResponse(rsp)
}

// GetObjectStorageCredentialsWithResponse request returning *GetObjectStorageCredentialsResponse
func (c *ClientWithResponses) GetObjectStorageCredentialsWithResponse(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.GetObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*GetObjectStorageCredentialsResponse, error) {
	rsp, err := c.GetObjectStorageCredentials(ctx, userId, objectStorageId, credentialId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetObjectStorageCredentialsResponse(rsp)
}

// RegenerateObjectStorageCredentialsWithResponse request returning *RegenerateObjectStorageCredentialsResponse
func (c *ClientWithResponses) RegenerateObjectStorageCredentialsWithResponse(ctx context.Context, userId string, objectStorageId string, credentialId int64, params *models.RegenerateObjectStorageCredentialsParams, reqEditors ...RequestEditorFn) (*RegenerateObjectStorageCredentialsResponse, error) {
	rsp, err := c.RegenerateObjectStorageCredentials(ctx, userId, objectStorageId, credentialId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegenerateObjectStorageCredentialsResponse(rsp)
}

// ResendEmailVerificationWithResponse request returning *ResendEmailVerificationResponse
func (c *ClientWithResponses) ResendEmailVerificationWithResponse(ctx context.Context, userId string, params *models.ResendEmailVerificationParams, reqEditors ...RequestEditorFn) (*ResendEmailVerificationResponse, error) {
	rsp, err := c.ResendEmailVerification(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResendEmailVerificationResponse(rsp)
}

// ResetPasswordWithResponse request returning *ResetPasswordResponse
func (c *ClientWithResponses) ResetPasswordWithResponse(ctx context.Context, userId string, params *models.ResetPasswordParams, reqEditors ...RequestEditorFn) (*ResetPasswordResponse, error) {
	rsp, err := c.ResetPassword(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResetPasswordResponse(rsp)
}

// RetrieveVipListWithResponse request returning *RetrieveVipListResponse
func (c *ClientWithResponses) RetrieveVipListWithResponse(ctx context.Context, params *models.RetrieveVipListParams, reqEditors ...RequestEditorFn) (*RetrieveVipListResponse, error) {
	rsp, err := c.RetrieveVipList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveVipListResponse(rsp)
}

// RetrieveVipAuditsListWithResponse request returning *RetrieveVipAuditsListResponse
func (c *ClientWithResponses) RetrieveVipAuditsListWithResponse(ctx context.Context, params *models.RetrieveVipAuditsListParams, reqEditors ...RequestEditorFn) (*RetrieveVipAuditsListResponse, error) {
	rsp, err := c.RetrieveVipAuditsList(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveVipAuditsListResponse(rsp)
}

// RetrieveVipWithResponse request returning *RetrieveVipResponse
func (c *ClientWithResponses) RetrieveVipWithResponse(ctx context.Context, ip string, params *models.RetrieveVipParams, reqEditors ...RequestEditorFn) (*RetrieveVipResponse, error) {
	rsp, err := c.RetrieveVip(ctx, ip, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRetrieveVipResponse(rsp)
}

// UnassignIpWithResponse request returning *UnassignIpResponse
func (c *ClientWithResponses) UnassignIpWithResponse(ctx context.Context, ip string, resourceType models.UnassignIpParamsResourceType, resourceId int64, params *models.UnassignIpParams, reqEditors ...RequestEditorFn) (*UnassignIpResponse, error) {
	rsp, err := c.UnassignIp(ctx, ip, resourceType, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUnassignIpResponse(rsp)
}

// AssignIpWithResponse request returning *AssignIpResponse
func (c *ClientWithResponses) AssignIpWithResponse(ctx context.Context, ip string, resourceType models.AssignIpParamsResourceType, resourceId int64, params *models.AssignIpParams, reqEditors ...RequestEditorFn) (*AssignIpResponse, error) {
	rsp, err := c.AssignIp(ctx, ip, resourceType, resourceId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAssignIpResponse(rsp)
}

// ParseRetrieveImageListResponse parses an HTTP response from a RetrieveImageListWithResponse call
func ParseRetrieveImageListResponse(rsp *http.Response) (*RetrieveImageListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveImageListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListImageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateCustomImageResponse parses an HTTP response from a CreateCustomImageWithResponse call
func ParseCreateCustomImageResponse(rsp *http.Response) (*CreateCustomImageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateCustomImageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateCustomImageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 415:
		var dest models.CreateCustomImageFailResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON415 = &dest

	}

	return response, nil
}

// ParseRetrieveImageAuditsListResponse parses an HTTP response from a RetrieveImageAuditsListWithResponse call
func ParseRetrieveImageAuditsListResponse(rsp *http.Response) (*RetrieveImageAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveImageAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ImageAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveCustomImagesStatsResponse parses an HTTP response from a RetrieveCustomImagesStatsWithResponse call
func ParseRetrieveCustomImagesStatsResponse(rsp *http.Response) (*RetrieveCustomImagesStatsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveCustomImagesStatsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.CustomImagesStatsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteImageResponse parses an HTTP response from a DeleteImageWithResponse call
func ParseDeleteImageResponse(rsp *http.Response) (*DeleteImageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteImageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveImageResponse parses an HTTP response from a RetrieveImageWithResponse call
func ParseRetrieveImageResponse(rsp *http.Response) (*RetrieveImageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveImageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindImageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateImageResponse parses an HTTP response from a UpdateImageWithResponse call
func ParseUpdateImageResponse(rsp *http.Response) (*UpdateImageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateImageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateCustomImageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveInstancesListResponse parses an HTTP response from a RetrieveInstancesListWithResponse call
func ParseRetrieveInstancesListResponse(rsp *http.Response) (*RetrieveInstancesListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}
	response := &RetrieveInstancesListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListInstancesResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateInstanceResponse parses an HTTP response from a CreateInstanceWithResponse call
func ParseCreateInstanceResponse(rsp *http.Response) (*CreateInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveInstancesActionsAuditsListResponse parses an HTTP response from a RetrieveInstancesActionsAuditsListWithResponse call
func ParseRetrieveInstancesActionsAuditsListResponse(rsp *http.Response) (*RetrieveInstancesActionsAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveInstancesActionsAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListInstancesActionsAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveInstancesAuditsListResponse parses an HTTP response from a RetrieveInstancesAuditsListWithResponse call
func ParseRetrieveInstancesAuditsListResponse(rsp *http.Response) (*RetrieveInstancesAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveInstancesAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListInstancesAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveInstanceResponse parses an HTTP response from a RetrieveInstanceWithResponse call
func ParseRetrieveInstanceResponse(rsp *http.Response) (*RetrieveInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePatchInstanceResponse parses an HTTP response from a PatchInstanceWithResponse call
func ParsePatchInstanceResponse(rsp *http.Response) (*PatchInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PatchInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.PatchInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseReinstallInstanceResponse parses an HTTP response from a ReinstallInstanceWithResponse call
func ParseReinstallInstanceResponse(rsp *http.Response) (*ReinstallInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ReinstallInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ReinstallInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRescueResponse parses an HTTP response from a RescueWithResponse call
func ParseRescueResponse(rsp *http.Response) (*RescueResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RescueResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceRescueActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseResetPasswordActionResponse parses an HTTP response from a ResetPasswordActionWithResponse call
func ParseResetPasswordActionResponse(rsp *http.Response) (*ResetPasswordActionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ResetPasswordActionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceResetPasswordActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRestartResponse parses an HTTP response from a RestartWithResponse call
func ParseRestartResponse(rsp *http.Response) (*RestartResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RestartResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceRestartActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseShutdownResponse parses an HTTP response from a ShutdownWithResponse call
func ParseShutdownResponse(rsp *http.Response) (*ShutdownResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ShutdownResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceShutdownActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseStartResponse parses an HTTP response from a StartWithResponse call
func ParseStartResponse(rsp *http.Response) (*StartResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StartResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceStartActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseStopResponse parses an HTTP response from a StopWithResponse call
func ParseStopResponse(rsp *http.Response) (*StopResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StopResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.InstanceStopActionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseCancelInstanceResponse parses an HTTP response from a CancelInstanceWithResponse call
func ParseCancelInstanceResponse(rsp *http.Response) (*CancelInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CancelInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CancelInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveSnapshotListResponse parses an HTTP response from a RetrieveSnapshotListWithResponse call
func ParseRetrieveSnapshotListResponse(rsp *http.Response) (*RetrieveSnapshotListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSnapshotListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListSnapshotResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateSnapshotResponse parses an HTTP response from a CreateSnapshotWithResponse call
func ParseCreateSnapshotResponse(rsp *http.Response) (*CreateSnapshotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateSnapshotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateSnapshotResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseDeleteSnapshotResponse parses an HTTP response from a DeleteSnapshotWithResponse call
func ParseDeleteSnapshotResponse(rsp *http.Response) (*DeleteSnapshotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteSnapshotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveSnapshotResponse parses an HTTP response from a RetrieveSnapshotWithResponse call
func ParseRetrieveSnapshotResponse(rsp *http.Response) (*RetrieveSnapshotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSnapshotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindSnapshotResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateSnapshotResponse parses an HTTP response from a UpdateSnapshotWithResponse call
func ParseUpdateSnapshotResponse(rsp *http.Response) (*UpdateSnapshotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateSnapshotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateSnapshotResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRollbackSnapshotResponse parses an HTTP response from a RollbackSnapshotWithResponse call
func ParseRollbackSnapshotResponse(rsp *http.Response) (*RollbackSnapshotResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RollbackSnapshotResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.RollbackSnapshotResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpgradeInstanceResponse parses an HTTP response from a UpgradeInstanceWithResponse call
func ParseUpgradeInstanceResponse(rsp *http.Response) (*UpgradeInstanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpgradeInstanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.PatchInstanceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveSnapshotsAuditsListResponse parses an HTTP response from a RetrieveSnapshotsAuditsListWithResponse call
func ParseRetrieveSnapshotsAuditsListResponse(rsp *http.Response) (*RetrieveSnapshotsAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSnapshotsAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListSnapshotsAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateTicketResponse parses an HTTP response from a CreateTicketWithResponse call
func ParseCreateTicketResponse(rsp *http.Response) (*CreateTicketResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateTicketResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateTicketResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveDataCenterListResponse parses an HTTP response from a RetrieveDataCenterListWithResponse call
func ParseRetrieveDataCenterListResponse(rsp *http.Response) (*RetrieveDataCenterListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveDataCenterListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListDataCenterResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveObjectStorageListResponse parses an HTTP response from a RetrieveObjectStorageListWithResponse call
func ParseRetrieveObjectStorageListResponse(rsp *http.Response) (*RetrieveObjectStorageListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveObjectStorageListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateObjectStorageResponse parses an HTTP response from a CreateObjectStorageWithResponse call
func ParseCreateObjectStorageResponse(rsp *http.Response) (*CreateObjectStorageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateObjectStorageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveObjectStorageAuditsListResponse parses an HTTP response from a RetrieveObjectStorageAuditsListWithResponse call
func ParseRetrieveObjectStorageAuditsListResponse(rsp *http.Response) (*RetrieveObjectStorageAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveObjectStorageAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListObjectStorageAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveObjectStorageResponse parses an HTTP response from a RetrieveObjectStorageWithResponse call
func ParseRetrieveObjectStorageResponse(rsp *http.Response) (*RetrieveObjectStorageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveObjectStorageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateObjectStorageResponse parses an HTTP response from a UpdateObjectStorageWithResponse call
func ParseUpdateObjectStorageResponse(rsp *http.Response) (*UpdateObjectStorageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateObjectStorageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.CancelObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCancelObjectStorageResponse parses an HTTP response from a CancelObjectStorageWithResponse call
func ParseCancelObjectStorageResponse(rsp *http.Response) (*CancelObjectStorageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CancelObjectStorageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.CancelObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpgradeObjectStorageResponse parses an HTTP response from a UpgradeObjectStorageWithResponse call
func ParseUpgradeObjectStorageResponse(rsp *http.Response) (*UpgradeObjectStorageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpgradeObjectStorageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpgradeObjectStorageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveObjectStoragesStatsResponse parses an HTTP response from a RetrieveObjectStoragesStatsWithResponse call
func ParseRetrieveObjectStoragesStatsResponse(rsp *http.Response) (*RetrieveObjectStoragesStatsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveObjectStoragesStatsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ObjectStoragesStatsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrievePrivateNetworkListResponse parses an HTTP response from a RetrievePrivateNetworkListWithResponse call
func ParseRetrievePrivateNetworkListResponse(rsp *http.Response) (*RetrievePrivateNetworkListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrievePrivateNetworkListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListPrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreatePrivateNetworkResponse parses an HTTP response from a CreatePrivateNetworkWithResponse call
func ParseCreatePrivateNetworkResponse(rsp *http.Response) (*CreatePrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreatePrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreatePrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrievePrivateNetworkAuditsListResponse parses an HTTP response from a RetrievePrivateNetworkAuditsListWithResponse call
func ParseRetrievePrivateNetworkAuditsListResponse(rsp *http.Response) (*RetrievePrivateNetworkAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrievePrivateNetworkAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListPrivateNetworkAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeletePrivateNetworkResponse parses an HTTP response from a DeletePrivateNetworkWithResponse call
func ParseDeletePrivateNetworkResponse(rsp *http.Response) (*DeletePrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeletePrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrievePrivateNetworkResponse parses an HTTP response from a RetrievePrivateNetworkWithResponse call
func ParseRetrievePrivateNetworkResponse(rsp *http.Response) (*RetrievePrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrievePrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindPrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePatchPrivateNetworkResponse parses an HTTP response from a PatchPrivateNetworkWithResponse call
func ParsePatchPrivateNetworkResponse(rsp *http.Response) (*PatchPrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PatchPrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.PatchPrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUnassignInstancePrivateNetworkResponse parses an HTTP response from a UnassignInstancePrivateNetworkWithResponse call
func ParseUnassignInstancePrivateNetworkResponse(rsp *http.Response) (*UnassignInstancePrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UnassignInstancePrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.UnassignInstancePrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseAssignInstancePrivateNetworkResponse parses an HTTP response from a AssignInstancePrivateNetworkWithResponse call
func ParseAssignInstancePrivateNetworkResponse(rsp *http.Response) (*AssignInstancePrivateNetworkResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AssignInstancePrivateNetworkResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.AssignInstancePrivateNetworkResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveRoleListResponse parses an HTTP response from a RetrieveRoleListWithResponse call
func ParseRetrieveRoleListResponse(rsp *http.Response) (*RetrieveRoleListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveRoleListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListRoleResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateRoleResponse parses an HTTP response from a CreateRoleWithResponse call
func ParseCreateRoleResponse(rsp *http.Response) (*CreateRoleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateRoleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateRoleResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveApiPermissionsListResponse parses an HTTP response from a RetrieveApiPermissionsListWithResponse call
func ParseRetrieveApiPermissionsListResponse(rsp *http.Response) (*RetrieveApiPermissionsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveApiPermissionsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListApiPermissionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveRoleAuditsListResponse parses an HTTP response from a RetrieveRoleAuditsListWithResponse call
func ParseRetrieveRoleAuditsListResponse(rsp *http.Response) (*RetrieveRoleAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveRoleAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListRoleAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteRoleResponse parses an HTTP response from a DeleteRoleWithResponse call
func ParseDeleteRoleResponse(rsp *http.Response) (*DeleteRoleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteRoleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveRoleResponse parses an HTTP response from a RetrieveRoleWithResponse call
func ParseRetrieveRoleResponse(rsp *http.Response) (*RetrieveRoleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveRoleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindRoleResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateRoleResponse parses an HTTP response from a UpdateRoleWithResponse call
func ParseUpdateRoleResponse(rsp *http.Response) (*UpdateRoleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateRoleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateRoleResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveSecretListResponse parses an HTTP response from a RetrieveSecretListWithResponse call
func ParseRetrieveSecretListResponse(rsp *http.Response) (*RetrieveSecretListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSecretListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListSecretResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateSecretResponse parses an HTTP response from a CreateSecretWithResponse call
func ParseCreateSecretResponse(rsp *http.Response) (*CreateSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateSecretResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveSecretAuditsListResponse parses an HTTP response from a RetrieveSecretAuditsListWithResponse call
func ParseRetrieveSecretAuditsListResponse(rsp *http.Response) (*RetrieveSecretAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSecretAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListSecretAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteSecretResponse parses an HTTP response from a DeleteSecretWithResponse call
func ParseDeleteSecretResponse(rsp *http.Response) (*DeleteSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveSecretResponse parses an HTTP response from a RetrieveSecretWithResponse call
func ParseRetrieveSecretResponse(rsp *http.Response) (*RetrieveSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindSecretResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateSecretResponse parses an HTTP response from a UpdateSecretWithResponse call
func ParseUpdateSecretResponse(rsp *http.Response) (*UpdateSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateSecretResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveTagListResponse parses an HTTP response from a RetrieveTagListWithResponse call
func ParseRetrieveTagListResponse(rsp *http.Response) (*RetrieveTagListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveTagListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListTagResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateTagResponse parses an HTTP response from a CreateTagWithResponse call
func ParseCreateTagResponse(rsp *http.Response) (*CreateTagResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateTagResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateTagResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveAssignmentsAuditsListResponse parses an HTTP response from a RetrieveAssignmentsAuditsListWithResponse call
func ParseRetrieveAssignmentsAuditsListResponse(rsp *http.Response) (*RetrieveAssignmentsAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveAssignmentsAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListAssignmentAuditsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveTagAuditsListResponse parses an HTTP response from a RetrieveTagAuditsListWithResponse call
func ParseRetrieveTagAuditsListResponse(rsp *http.Response) (*RetrieveTagAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveTagAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListTagAuditsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteTagResponse parses an HTTP response from a DeleteTagWithResponse call
func ParseDeleteTagResponse(rsp *http.Response) (*DeleteTagResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteTagResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveTagResponse parses an HTTP response from a RetrieveTagWithResponse call
func ParseRetrieveTagResponse(rsp *http.Response) (*RetrieveTagResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveTagResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindTagResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateTagResponse parses an HTTP response from a UpdateTagWithResponse call
func ParseUpdateTagResponse(rsp *http.Response) (*UpdateTagResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateTagResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateTagResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveAssignmentListResponse parses an HTTP response from a RetrieveAssignmentListWithResponse call
func ParseRetrieveAssignmentListResponse(rsp *http.Response) (*RetrieveAssignmentListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveAssignmentListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListAssignmentResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteAssignmentResponse parses an HTTP response from a DeleteAssignmentWithResponse call
func ParseDeleteAssignmentResponse(rsp *http.Response) (*DeleteAssignmentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteAssignmentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveAssignmentResponse parses an HTTP response from a RetrieveAssignmentWithResponse call
func ParseRetrieveAssignmentResponse(rsp *http.Response) (*RetrieveAssignmentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveAssignmentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindAssignmentResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateAssignmentResponse parses an HTTP response from a CreateAssignmentWithResponse call
func ParseCreateAssignmentResponse(rsp *http.Response) (*CreateAssignmentResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateAssignmentResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateAssignmentResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveUserListResponse parses an HTTP response from a RetrieveUserListWithResponse call
func ParseRetrieveUserListResponse(rsp *http.Response) (*RetrieveUserListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveUserListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListUserResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateUserResponse parses an HTTP response from a CreateUserWithResponse call
func ParseCreateUserResponse(rsp *http.Response) (*CreateUserResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest models.CreateUserResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseRetrieveUserAuditsListResponse parses an HTTP response from a RetrieveUserAuditsListWithResponse call
func ParseRetrieveUserAuditsListResponse(rsp *http.Response) (*RetrieveUserAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveUserAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListUserAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveUserClientResponse parses an HTTP response from a RetrieveUserClientWithResponse call
func ParseRetrieveUserClientResponse(rsp *http.Response) (*RetrieveUserClientResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveUserClientResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindClientResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGenerateClientSecretResponse parses an HTTP response from a GenerateClientSecretWithResponse call
func ParseGenerateClientSecretResponse(rsp *http.Response) (*GenerateClientSecretResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GenerateClientSecretResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.GenerateClientSecretResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveUserIsPasswordSetResponse parses an HTTP response from a RetrieveUserIsPasswordSetWithResponse call
func ParseRetrieveUserIsPasswordSetResponse(rsp *http.Response) (*RetrieveUserIsPasswordSetResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveUserIsPasswordSetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindUserIsPasswordSetResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteUserResponse parses an HTTP response from a DeleteUserWithResponse call
func ParseDeleteUserResponse(rsp *http.Response) (*DeleteUserResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveUserResponse parses an HTTP response from a RetrieveUserWithResponse call
func ParseRetrieveUserResponse(rsp *http.Response) (*RetrieveUserResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindUserResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateUserResponse parses an HTTP response from a UpdateUserWithResponse call
func ParseUpdateUserResponse(rsp *http.Response) (*UpdateUserResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.UpdateUserResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseListObjectStorageCredentialsResponse parses an HTTP response from a ListObjectStorageCredentialsWithResponse call
func ParseListObjectStorageCredentialsResponse(rsp *http.Response) (*ListObjectStorageCredentialsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListObjectStorageCredentialsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListCredentialResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetObjectStorageCredentialsResponse parses an HTTP response from a GetObjectStorageCredentialsWithResponse call
func ParseGetObjectStorageCredentialsResponse(rsp *http.Response) (*GetObjectStorageCredentialsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetObjectStorageCredentialsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindCredentialResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRegenerateObjectStorageCredentialsResponse parses an HTTP response from a RegenerateObjectStorageCredentialsWithResponse call
func ParseRegenerateObjectStorageCredentialsResponse(rsp *http.Response) (*RegenerateObjectStorageCredentialsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RegenerateObjectStorageCredentialsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindCredentialResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseResendEmailVerificationResponse parses an HTTP response from a ResendEmailVerificationWithResponse call
func ParseResendEmailVerificationResponse(rsp *http.Response) (*ResendEmailVerificationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ResendEmailVerificationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseResetPasswordResponse parses an HTTP response from a ResetPasswordWithResponse call
func ParseResetPasswordResponse(rsp *http.Response) (*ResetPasswordResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ResetPasswordResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseRetrieveVipListResponse parses an HTTP response from a RetrieveVipListWithResponse call
func ParseRetrieveVipListResponse(rsp *http.Response) (*RetrieveVipListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveVipListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListVipResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveVipAuditsListResponse parses an HTTP response from a RetrieveVipAuditsListWithResponse call
func ParseRetrieveVipAuditsListResponse(rsp *http.Response) (*RetrieveVipAuditsListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveVipAuditsListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.ListVipAuditResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseRetrieveVipResponse parses an HTTP response from a RetrieveVipWithResponse call
func ParseRetrieveVipResponse(rsp *http.Response) (*RetrieveVipResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RetrieveVipResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.FindVipResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUnassignIpResponse parses an HTTP response from a UnassignIpWithResponse call
func ParseUnassignIpResponse(rsp *http.Response) (*UnassignIpResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UnassignIpResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseAssignIpResponse parses an HTTP response from a AssignIpWithResponse call
func ParseAssignIpResponse(rsp *http.Response) (*AssignIpResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AssignIpResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest models.AssignVipResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
