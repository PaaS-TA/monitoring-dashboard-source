package service

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"monitoring-portal/paas/model"
)

type CloudFoundryService struct {
	cfClient *cfclient.Client
}



func GetCloudFoundryService(cfClient *cfclient.Client) *CloudFoundryService {
	return &CloudFoundryService{
		cfClient: cfClient,
	}
}

func (s *CloudFoundryService) GetPaasDiagram() (model.Diagram){
	var orgList []model.Diagram
	orgs, _ := s.cfClient.ListOrgs()
	for _, org := range orgs {
		var spaceList []model.Diagram
		spaces, _ := s.cfClient.ListSpacesByOrgGuid(org.Guid)
		for _, space := range spaces {
			var appList []model.Diagram
			apps, _ := s.cfClient.ListAppsBySpaceGuid(space.Guid)
			for _, app := range apps {
				appData := model.Diagram {
					Id: app.Guid,
					Name: app.Name,
					Title: app.Command,
				}
				appList = append(appList, appData)
			}
			spaceData := model.Diagram{
				Id: space.Guid,
				Name: space.Name,
				Title: space.Name,
				Children: appList,
			}
			spaceList = append(spaceList, spaceData)
		}
		orgData := model.Diagram{
			Id : org.Guid,
			Name: org.Name,
			Title: org.Name,
			Children: spaceList,
		}
		orgList = append(orgList, orgData)
	}

	rootNode := model.Diagram{
		Id: "0",
		Name: "PaaS-TA",
		Title: "PaaS-TA",
		Children: orgList,
	}
	return rootNode
}