package model

var DefaultTenantId string
var NovaUrl, NovaVersion string
var KeystoneUrl, KeystoneVersion string
var NeutronUrl, NeutronVersion string
var CinderUrl, CinderVersion string
var GlanceUrl, GlanceVersion string
var RabbitMqIp, RabbitMqPort string
var MetricDBName string
var GMTTimeGap int64

type (

	/**
	Description : Openstack Admin Information to Get Openstack's Resources - ex: Hypervisor, Node, Tenant
	 */
	OpenstackProvider struct {
		Region 				string
		Domain 				string
		Username 			string
		UserId  			string
		Password 			string
		TenantName 			string
		AdminTenantId 		string
		KeystoneUrl 	 	string
		IdentityEndpoint	string
		RabbitmqUser 		string
		RabbitmqPass 		string
		RabbitmqTargetNode	string
	}


	/**
	Description :Openstack total hypervisor resources struct
	 */
	HypervisorResources struct {
		VmTotalLimit            int			`json:"vmTotalLimit"`
		VmTotal                 int			`json:"vmTotal"`
		VmRunning 				int			`json:"vmRunning"`
		VmState                 []VmState   `json:"vmState"`
		VcpuTotal 				float64 	`json:"vcpuTotal"`
		VcpuUsed 				float64 	`json:"vcpuUsed"`
		MemoryMbTotal 			float64		`json:"memoryMbTotal"`
		MemoryMbUsed 			float64 	`json:"memoryMbUsed"`
		MemoryMbFree 			float64 	`json:"memoryMbFree"`
		DiskGbTotal 			float64 	`json:"diskGbTotal"`
		DiskGbUsed 				float64 	`json:"diskGbUsed"`
		DiskGbFree 				float64 	`json:"diskGbFree"`
		DiskGbLeastAvailable	float64 	`json:"diskGbLeastAvailable"`
	}

	VmState struct{
		VmStateName		string	`json:"name"`
		VmCnt           int		`json:"vmCount"`
	}


	/**
	Description: Openstack Compute Node Information
	 */
	NodeResources struct {
		Id 					int 		`json:"nodeId"`
		Hostname 			string 		`json:"hostname"`
		HostIp	 			string 		`json:"hostIp"`
		Type 				string		`json:"type"`
		VcpusMax			int 		`json:"vcpusMax"`
		VcpusUsed 			int 		`json:"vcpusUsed"`
		MemoryMbMax 		int 		`json:"memoryMbMax"`
		MemoryMbUsed 		int 		`json:"memoryMbUsed"`
		MemoryMbFree  		int 		`json:"memoryMbFree"`
		DiskGbMax 			int 		`json:"diskGbMax"`
		DiskGbUsed 			int 		`json:"diskGbUsed"`
		DiskGbFree 			int 		`json:"diskGbFree"`
		DiskAvailableLeast	int 		`json:"diskAbilableLeast"`
		State 				string		`json:"state"`
		Status 				string		`json:"status"`
		TotalVms            int         `json:"totalVms"`
		RunningVms          int         `json:"runningVms"`
		CpuUsage            float64     `json:"cpuUsage"`
		MemUsage            float64     `json:"memUsage"`
		AgentStatus         string		`json:"agentStatus"`
	}

	ManageNodeResources struct {
		Hostname 		string 		`json:"hostname"`
		CpuUsage       	float64     `json:"cpuUsage"`
		MemUsage        float64     `json:"memoryUsage"`
		MemoryMbMax     float64     `json:"memoryMbMax"`
		MemoryMbUsed	float64     `json:"memoryUsedMb"`
		DiskUsage       float64		`json:"diskUsage"`
		DiskGbMax       float64     `json:"diskGbMax"`
		DiskGbUsed      float64     `json:"diskGbUsed"`
		AgentStatus     string      `json:"agentStatus"`
	}

	TenantSummaryInfo 	struct {
		//ParentId 				string 		`json:"parent_id"`
		//IsDomain 				bool 		`json:"is_domain"`
		Name 					string		`json:"name"`
		Id 						string 		`json:"id"`
		Description 			string 		`json:"description"`
		Enabled 				bool		`json:"enabled"`
		InstancesLimit			int 		`json:"instancesLimit"`
		InstancesUsed 			int 		`json:"instancesUsed"`
		VcpusLimit				int 		`json:"vcpusLimit"`
		VcpusUsed 				float64		`json:"vcpusUsed"`
		MemoryMbLimit 			int 		`json:"memoryMbLimit"`
		MemoryMbUsed 			float64		`json:"memoryMbUsed"`
		FloatingIpsLimit		int 		`json:"floatingIpsLimit"`
		FloatingIpsUsed			int 		`json:"floatingIpsUsed"`
		SecurityGroupsLimit		int 		`json:"securityGroupsLimit"`
		SecurityGroupsUsed		int 		`json:"securityGroupsUsed"`
		VolumeStorageLimit		int 		`json:"volumeStorageLimit"`
		VolumeStorageUsed		int 		`json:"volumeStorageUsed"`
		VolumeStorageLimitGb	int			`json:"volumeStorageLimitGb"`
		VolumeStorageUsedGb		int			`json:"volumeStorageUsedGb"`
	}


	/**
	Description : Openstack tenant struct
	 */
	TenantInfo struct {
		ParentId 		string 		`json:"tenantId"`
		DomainId 		string 		`json:"domainId"`
		Name 			string		`json:"name"`
		IsDomain 		bool 		`json:"isDomain"`
		Description 	string 		`json:"description"`
		Enabled 		bool		`json:"enabled"`
		Id 				string 		`json:"id"`
		Links 			map[string]interface{} `json:"links"`
	}


	/**
	Description: Openstack tenant resources usage including limit medatada
	 */
	TenantResourcesUsage struct {
		Instances 				int 		`json:"instances"`
		Vcpus 					float64 	`json:"vcpus"`
		MemoryMb 				float64 	`json:"memory_mb"`
		SecurityGroups			int 		`json:"security_groups"`
		FloatingIps 			int 		`json:"floating_ips"`
		TenantResourceLimit 	TenantResourcesLimit
		TenantNetworkLimit 		TenantNetworkLimit
		TenantStorageResource  	TenantStorageResources
	}


	/**
	Description: Openstack tenant resources limit metadata
	 */
	TenantResourcesLimit struct {
		/*
		[metadata_items:128 injected_files:5 injected_file_content_bytes:10240 server_groups:10 key_pairs:100
		injected_file_path_bytes:255 id:9c1a27e20412473b843dbf32bdec2390 instances:150 security_group_rules:20
		fixed_ips:-1 security_groups:10 server_group_members:10 ram:182400 floating_ips:10 cores:150]
		 */
		InstancesLimit		int `json:"instances_limit"`
		MemoryMbLimit		int `json:"memory_mb_limit"`
		CoresLimit			int `json:"cores_limit"`
		ServerGroupsLimit	int `json:"server_groups_limit"`
		KeyPairsLimit		int `json:"key_pairs_limit"`
	}


	/**
	Description: Openstack tenant storage limit metadata
	 */
	TenantStorageResources struct {
		/*
		map[per_volume_gigabytes:-1 gigabytes:1000 snapshots_SSD:-1 snapshots_SSD1:-1 gigabytes_SSD:-1 gigabytes_SSD1:-1
		backup_gigabytes:1000 volumes_SSD:-1 snapshots:10 volumes_SSD1:-1 volumes:30 id:9c1a27e20412473b843dbf32bdec2390 backups:10]
		 */
		VolumeLimitGb 	int 	`json:"volume_limit_gb"`
		VolumesLimit 	int		`json:"volumes_limit"`
		SnapshotsLimit	int 	`json:"snapshots_limit"`
		BackupsLimit 	int		`json:"backups_limit"`
		VolumeGb 		int 	`json:"volume_gb"`
		Volumes 		int 	`json:"volumes"`
		Snapshots 		int 	`json:"snapshots"`
		Backups 		int 	`json:"backups"`
	}


	/**
	Description: Openstack tenant network limit metadata
	 */
	TenantNetworkLimit struct{
		/*
		router:20 port:500 subnetpool:-1 security_group_rule:150 security_group:30 rbac_policy:10 subnet:100 network:100 floatingip:100]
		 */
		RouterLimit 			int 	`json:"router_limit"`
		SecurityGroupRuleLimit 	int		`json:"security_group_rule_limit"`
		SecurityGroupLimit		int 	`json:"security_group_limit"`
		FloatingIpsLimit		int 	`json:"floating_ips_limit"`
		SubnetLimit 			int 	`json:"subnet_limit"`
		NetworkLimit	 		int 	`json:"network_limit"`
		PortLimit 				int 	`json:"port_limit"`
	}


	/**
	Description: Openstack tenant's created instance info
	 */
	InstanceInfo struct {
		TenantId 	string 		`json:"tenant_id"`
		InstanceId 	string 		`json:"instance_id"`
		Zone        string      `json:"zone"`
		Name 		string 		`json:"name"`
		CpuUsage    float64 	`json:"cpuUsage"`
		MemoryUsage float64 	`json:"memoryUsage"`
		Address     []string    `json:"address"`
		Flavor 		string 		`json:"flavor"`
		Vcpus 		float64 	`json:"vcpus"`
		DiskGb 		float64 	`json:"disk_gb"`
		MemoryMb 	float64 	`json:"memory_mb"`
		State 		string 		`json:"state"`
		StartedAt 	string		`json:"started_at"`
		EndedAt		string 	 	`json:"ended_at"`
		Uptime 		float64 	`json:"uptime"`
	}


	/**
	Description: Tenant Floating IP Information
	 */
	FloatingIPInfo 	struct {
		Id 					string		`json:"id"`
		TenantId 			string 		`json:"tenant_id"`
		RouterId 			string 		`json:"router_id"`
		FloatingNetworkId 	string 		`json:"floating_network_id"`
		InnerIp 			string 		`json:"internal_ip"`
		FloatingIp 			string 		`json:"floating_ip"`
		PortId	 			string 		`json:"port_id"`
		Status 				string 		`json:"status"`
		Description  		string 		`json:"description"`
	}


	/**
	Description: Instance Detail Inforamation
	 */
	InstanceDetail struct {
		Id 					string 		`json:"id"`
		Name 				string 		`json:"name"`
		ProcessName 		string 		`json:"process_name"`
		AvailabilityZone 	string 		`json:"availability_zone"`
		SecurityGroups 		string 		`json:"security_groups"`
		CreatedDate 		string 		`json:"created_date"`
		Deployment 			DeploymentInfo
		Network 			[]NetworkInfo
		Flavor	 			FlavorInfo
		Image 				ImageInfo
	}


	/**
	Description: Instance Bosh Deployment Information
	 */
	DeploymentInfo struct {
		Name 		string		`json:"name"`
		Deployment 	string 		`json:"deployment"`
		Director 	string 		`json:"director"`
		Job	 		string 		`json:"job"`
	}


	/**
	Description: Instance Network Information
	 */
	NetworkInfo struct {
		Ip 			string 		`json:"ip"`
		Type 		string 		`json:"type"`
		Mac_addr 	string 		`json:"mac_addr"`
	}


	/**
	Description: Instance Flavor Information
	 */
	FlavorInfo struct {
		Name 		string 		`json:"name"`
		Vcpu 		int 		`json:"vcpu"`
		Memory 		int		 	`json:"memroy"`
		Disk		int 		`json:"disk"`
	}


	/**
	Description: Instance Image Information
	 */
	ImageInfo struct {
		Id 				string		`json:"id"`
		Name 			string 		`json:"name"`
		Version 		string 		`json:"version"`
		OsType 			string	 	`json:"os_type"`
		OsKind			string 		`json:"os_kind"`
		HypervisorType 	string 		`json:"hypervisor_type"`
	}


	/**
	Description: RabbitMQ GlobalCounts
 	*/
	RabbitMQGlobalResource struct {
		Connections 	int 	`json:"connections"`
		Channels 		int 	`json:"channels"`
		Queues 			int 	`json:"queues"`
		Consumers 		int 	`json:"consumers"`
		Exchanges 		int 	`json:"exchanges"`
		NodeResources 	RabbitMQNodeResources
	}

	/**
	Description : RabbitMQ Node's Resource Information
	 */
	RabbitMQNodeResources struct {
		DiskMbFree 			float64 	`json:"diskMbFree"`
		DiskMbLimit 		float64 	`json:"diskMbLimitFree"`
		MemoryMbUsed 		float64		`json:"memoryMbUsed"`
		MemoryMbLimit 		float64		`json:"memoryMbLimit"`
		FileDescriptorUsed	float64 	`json:"fileDescriptorUsed"`
		FileDescriptorLimit float64 	`json:"fileDescriptorTotal"`
		ErlangProcUsed  	int 		`json:"erlangProcessUsed"`
		ErlangProcLimit 	int 		`json:"erlangProcessLimit"`
		SocketUsed 			float64		`json:"socketsUsed"`
		SocketLimit 		float64 	`json:"socketsLimit"`
	}
)