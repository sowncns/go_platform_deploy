package project

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/sowncns/k3s-deploy-platform/internal/project/dto"
)

var (
	ErrNotFounded       = errors.New("project not found")
	ErrForbidden      = errors.New("permission denied")
	ErrInvalidEnvKey  = errors.New("invalid environment variable key")
)

type ProjectService interface {
	CreateProject(ctx context.Context, userID uint, req dto.CreateProjectRequest) (*Project, error)
	UpdateProject(ctx context.Context, projectID uint, userID uint, req *dto.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, projectID uint, userID uint) error
	
	GetProject(ctx context.Context, projectID uint, userID uint) (*Project, error)
	ListProjects(ctx context.Context, userID uint) ([]*Project, error)
	
	PatchEnvVars(ctx context.Context ,projectID uint , userID uint, newVars map[string]string) error
}


type Service struct {
	repo ProjectRepository
	
}


func NewService(repo ProjectRepository) *Service {
	return &Service{
		repo: repo, 
		
	}
}

func (s *Service) GetProject(ctx context.Context, projectID uint, userID uint) (*Project, error) {
	project, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrNotFound
	}
	return project, nil
}

func (s *Service)PatchEnvVars(ctx context.Context ,projectID uint, userID uint, rawVars map[string]string) error{
		proj, err := s.repo.FindProjectByID(ctx,projectID)
		if err !=nil {
			return err
		}
		if proj.UserID != userID{
			return  ErrNotFounded
		}
		cleanedVars := make(map[string]string, len(rawVars))

		for k, v := range rawVars {
			key := strings.TrimSpace(k)
			if key == "" {
				return ErrInvalidEnvKey
			}
			// Convert mọi kiểu dữ liệu (8080, true) thành "8080", "true"
			cleanedVars[key] = fmt.Sprintf("%v", v)
		}
			if err := s.repo.PatchEnvVars(ctx, projectID, userID, cleanedVars); err != nil {
			return fmt.Errorf("update env vars in db: %w", err)
		}

		return nil

}

func (s *Service) ListProjects(ctx context.Context, userID uint) ([]*Project, error) {
	
	if userID == 0 {
		log.Print("[ERROR] user_id không hợp lệ")
		return nil, fmt.Errorf("userID cannot be zero")
	}

	return s.repo.ListProjectsByUserID(ctx, userID)
}

func (s *Service) DeleteProject(ctx context.Context, projectID uint, userID uint) error {
	project, err := s.repo.FindProjectByID(ctx, projectID)
	if err !=nil {
		return err
	}
	if project.UserID!= userID{
		return nil
	}
	return s.repo.DeleteProject(ctx, projectID, userID)
}


func (s *Service) CreateProject(ctx context.Context, userID uint, req dto.CreateProjectRequest) (*Project, error) {
	project := &Project{
		UserID: userID,
		Name:   req.Name,
		Branch: req.Branch,
		GitURL: req.GitURL,
		Env:    req.Env,
	}
	err := s.repo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, err
}

func (s *Service) UpdateProject(ctx context.Context, projectID uint, userID uint, req *dto.UpdateProjectRequest) error {
	project, err := s.GetProject(ctx, projectID, userID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.GitURL != nil {
		project.GitURL = *req.GitURL
	}
	if req.Branch != nil {
		project.Branch = *req.Branch
	}
	if req.Env != nil {
		project.Env = req.Env
	}

	return s.repo.UpdateProject(ctx, project)
}
