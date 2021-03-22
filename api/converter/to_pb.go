package converter

import (
	pb "oss.navercorp.com/metis/metis-server/api"
	"oss.navercorp.com/metis/metis-server/server/database"
)

func ToProject(project *database.Project) *pb.Project {
	// pbCreatedAt, err := ptypes.TimestampProto(project.CreatedAt)
	// if err != nil {
	// 	return nil, err
	// }

	return &pb.Project{
		Id:   project.ID.String(),
		Name: project.Name,
		// CreatedAt: pbCreatedAt,
	}
}

func ToProjects(projects []*database.Project) []*pb.Project {
	var pbProjects []*pb.Project
	for _, project := range projects {
		pbProjects = append(pbProjects, ToProject(project))
	}

	return pbProjects
}
