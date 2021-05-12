package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "oss.navercorp.com/metis/metis-server/api"
	"oss.navercorp.com/metis/metis-server/server/types"
)

// ToProject converts the given model to Protobuf message.
func ToProject(project *types.Project) *pb.Project {
	return &pb.Project{
		Id:        project.ID.String(),
		Name:      project.Name,
		CreatedAt: timestamppb.New(project.CreatedAt),
	}
}

// ToProjects converts the given model to Protobuf message.
func ToProjects(projects []*types.Project) []*pb.Project {
	var pbProjects []*pb.Project
	for _, project := range projects {
		pbProjects = append(pbProjects, ToProject(project))
	}

	return pbProjects
}
