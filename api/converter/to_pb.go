package converter

import (
	pb "oss.navercorp.com/metis/metis-server/api"
	"oss.navercorp.com/metis/metis-server/server/database"
)

func ToProject(project *database.Project) *pb.Project {
	return &pb.Project{
		Name: project.Name,
	}
}

func ToProjects(projects []*database.Project) []*pb.Project {
	var pbProjects []*pb.Project
	for _, project := range projects {
		pbProjects = append(pbProjects, ToProject(project))
	}

	return pbProjects
}
