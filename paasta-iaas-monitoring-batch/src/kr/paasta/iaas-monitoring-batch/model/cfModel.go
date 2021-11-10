package model

import "time"

type (
	CFConfig struct {
		ClientId string
		ClientPw string
		UserId   string
		UserPw   string
		Host     string
		ApiHost  string
		Port     string
	}
	UaaToken struct {
		Token            string
		Scope            string
		Expire           int64
		Refresh          string
		ExpireTime       time.Time
		Error            string
		ErrorDescription string
	}
	//GET process
	//ProcessResource struct {
	//	Resources Process  `json:"resources"`
	//}
	//Process struct {
	//	Guid                     string                 `json:"guid"` //?
	//	CreatedAt                string                 `json:"created_at"` //?
	//	UpdatedAt                string                 `json:"updated_at"`  //?
	//	Memory                   int                    `json:"memory_in_mb"`
	//	Instances                int                    `json:"instances"`
	//	DiskQuota                int                    `json:"disk_in_mb"`
	//
	//
	//}

	ProcessResource struct {
		Pagination string `json:"pagination"`
		Resources  []Resource
	}

	Resource struct {
		Guid      string `json:"guid"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Memory    int    `json:"memory_in_mb"`
		Instances int    `json:"instances"`
		DiskQuota int    `json:"disk_in_mb"`
	}
	//SpaceResource struct {
	//	Meta   Meta  `json:"metadata"`
	//	Entity Space `json:"entity"`
	//}
	//
	//Space struct {
	//	Guid                 string      `json:"guid"`
	//	CreatedAt            string      `json:"created_at"`
	//	UpdatedAt            string      `json:"updated_at"`
	//	Name                 string      `json:"name"`
	//	OrganizationGuid     string      `json:"organization_guid"`
	//	OrgURL               string      `json:"organization_url"`
	//	OrgData              OrgResource `json:"organization"`
	//	QuotaDefinitionGuid  string      `json:"space_quota_definition_guid"`
	//	IsolationSegmentGuid string      `json:"isolation_segment_guid"`
	//	AllowSSH             bool        `json:"allow_ssh"`
	//
	//}
	//
	//OrgResource struct {
	//	Meta   Meta `json:"metadata"`
	//	Entity Org  `json:"entity"`
	//}
	//
	//Org struct {
	//	Guid                        string `json:"guid"`
	//	CreatedAt                   string `json:"created_at"`
	//	UpdatedAt                   string `json:"updated_at"`
	//	Name                        string `json:"name"`
	//	Status                      string `json:"status"`
	//	QuotaDefinitionGuid         string `json:"quota_definition_guid"`
	//	DefaultIsolationSegmentGuid string `json:"default_isolation_segment_guid"`
	//
	//}
	//
	//
	//Meta struct {
	//	Guid      string `json:"guid"`
	//	Url       string `json:"url"`
	//	CreatedAt string `json:"created_at"`
	//	UpdatedAt string `json:"updated_at"`
	//}

	ScaleProcess struct {
		Memory    int `json:"memory_in_mb,omitempty"`
		Instances int `json:"instances,omitempty"`
		DiskQuota int `json:"disk_in_mb,omitempty"`
	}
	//UpdateResponse struct {
	//	Metadata Meta                 `json:"metadata"`
	//	Entity   UpdateResponseEntity `json:"entity"`
	//}

	//AppState string

	//AppUpdateResource struct {
	//	Name                     string                 `json:"name,omitempty"`
	//	Memory                   int                    `json:"memory,omitempty"`
	//	Instances                int                    `json:"instances,omitempty"`
	//	DiskQuota                int                    `json:"disk_quota,omitempty"`
	//	SpaceGuid                string                 `json:"space_guid,omitempty"`
	//	StackGuid                string                 `json:"stack_guid,omitempty"`
	//	State                    AppState               `json:"state,omitempty"`
	//	Command                  string                 `json:"command,omitempty"`
	//	Buildpack                string                 `json:"buildpack,omitempty"`
	//	HealthCheckHttpEndpoint  string                 `json:"health_check_http_endpoint,omitempty"`
	//	HealthCheckType          string                 `json:"health_check_type,omitempty"`
	//	HealthCheckTimeout       int                    `json:"health_check_timeout,omitempty"`
	//	Diego                    bool                   `json:"diego,omitempty"`
	//	EnableSSH                bool                   `json:"enable_ssh,omitempty"`
	//	DockerImage              string                 `json:"docker_image,omitempty"`
	//	DockerCredentials        map[string]interface{} `json:"docker_credentials_json,omitempty"`
	//	Environment              map[string]interface{} `json:"environment_json,omitempty"`
	//	StagingFailedReason      string                 `json:"staging_failed_reason,omitempty"`
	//	StagingFailedDescription string                 `json:"staging_failed_description,omitempty"`
	//	Ports                    []int                  `json:"ports,omitempty"`
	//}
	//
	//UpdateResponseEntity struct {
	//	Name                     string                 `json:"name"`
	//	Production               bool                   `json:"production"`
	//	SpaceGuid                string                 `json:"space_guid"`
	//	StackGuid                string                 `json:"stack_guid"`
	//	Buildpack                string                 `json:"buildpack"`
	//	DetectedBuildpack        string                 `json:"detected_buildpack"`
	//	DetectedBuildpackGuid    string                 `json:"detected_buildpack_guid"`
	//	Environment              map[string]interface{} `json:"environment_json"`
	//	Memory                   int                    `json:"memory"`
	//	Instances                int                    `json:"instances"`
	//	DiskQuota                int                    `json:"disk_quota"`
	//	State                    string                 `json:"state"`
	//	Version                  string                 `json:"version"`
	//	Command                  string                 `json:"command"`
	//	Console                  bool                   `json:"console"`
	//	Debug                    string                 `json:"debug"`
	//	StagingTaskId            string                 `json:"staging_task_id"`
	//	PackageState             string                 `json:"package_state"`
	//	HealthCheckHttpEndpoint  string                 `json:"health_check_http_endpoint"`
	//	HealthCheckType          string                 `json:"health_check_type"`
	//	HealthCheckTimeout       int                    `json:"health_check_timeout"`
	//	StagingFailedReason      string                 `json:"staging_failed_reason"`
	//	StagingFailedDescription string                 `json:"staging_failed_description"`
	//	Diego                    bool                   `json:"diego,omitempty"`
	//	DockerImage              string                 `json:"docker_image"`
	//	DockerCredentials        struct {
	//	Username string `json:"username"`
	//	Password string `json:"password"`
	//	} `json:"docker_credentials"`
	//	PackageUpdatedAt     string `json:"package_updated_at"`
	//	DetectedStartCommand string `json:"detected_start_command"`
	//	EnableSSH            bool   `json:"enable_ssh"`
	//	Ports                []int  `json:"ports"`
	//	SpaceURL             string `json:"space_url"`
	//	StackURL             string `json:"stack_url"`
	//	RoutesURL            string `json:"routes_url"`
	//	EventsURL            string `json:"events_url"`
	//	ServiceBindingsUrl   string `json:"service_bindings_url"`
	//	RouteMappingsUrl     string `json:"route_mappings_url"`
	//}
)
