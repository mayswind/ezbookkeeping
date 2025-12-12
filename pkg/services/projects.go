package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// ProjectService represents project service
type ProjectService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a project service singleton instance
var (
	Projects = &ProjectService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalProjectCountByUid returns total project count of user
func (s *ProjectService) GetTotalProjectCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.Project{})

	return count, err
}

// GetAllProjectsByUid returns all project models of user
func (s *ProjectService) GetAllProjectsByUid(c core.Context, uid int64) ([]*models.Project, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var projects []*models.Project
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&projects)

	return projects, err
}

// GetProjectByProjectId returns a project model according to project id
func (s *ProjectService) GetProjectByProjectId(c core.Context, uid int64, projectId int64) (*models.Project, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if projectId <= 0 {
		return nil, errs.ErrProjectIdInvalid
	}

	project := &models.Project{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(projectId).Where("uid=? AND deleted=?", uid, false).Get(project)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrProjectNotFound
	}

	return project, nil
}

// GetProjectsByProjectIds returns project models according to project ids
func (s *ProjectService) GetProjectsByProjectIds(c core.Context, uid int64, projectIds []int64) (map[int64]*models.Project, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if projectIds == nil {
		return nil, errs.ErrProjectIdInvalid
	}

	var projects []*models.Project
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("project_id", projectIds).Find(&projects)

	if err != nil {
		return nil, err
	}

	projectMap := s.GetProjectMapByList(projects)
	return projectMap, err
}

// GetMaxDisplayOrder returns the max display order
func (s *ProjectService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	project := &models.Project{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(project)

	if err != nil {
		return 0, err
	}

	if has {
		return project.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateProject saves a new project model to database
func (s *ProjectService) CreateProject(c core.Context, project *models.Project) error {
	if project.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsProjectName(c, project.Uid, project.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrProjectNameAlreadyExists
	}

	project.ProjectId = s.GenerateUuid(uuid.UUID_TYPE_PROJECT)

	if project.ProjectId < 1 {
		return errs.ErrSystemIsBusy
	}

	project.Deleted = false
	project.CreatedUnixTime = time.Now().Unix()
	project.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(project.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(project)
		return err
	})
}

// ModifyProject saves an existed project model to database
func (s *ProjectService) ModifyProject(c core.Context, project *models.Project) error {
	if project.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsProjectNameIgnoreGivenId(c, project.Uid, project.Name, project.ProjectId)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrProjectNameAlreadyExists
	}

	project.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(project.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(project.ProjectId).Cols("name", "color", "comment", "hidden", "updated_unix_time").Where("uid=? AND deleted=?", project.Uid, false).Update(project)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrProjectNotFound
		}

		return err
	})
}

// HideProject updates hidden field of given projects
func (s *ProjectService) HideProject(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Project{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("project_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrProjectNotFound
		}

		return err
	})
}

// ModifyProjectDisplayOrders updates display order of given projects
func (s *ProjectService) ModifyProjectDisplayOrders(c core.Context, uid int64, projects []*models.Project) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(projects); i++ {
		projects[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(projects); i++ {
			project := projects[i]
			updatedRows, err := sess.ID(project.ProjectId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(project)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrProjectNotFound
			}
		}

		return nil
	})
}

// DeleteProject deletes an existed project from database
func (s *ProjectService) DeleteProject(c core.Context, uid int64, projectId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Project{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Check if project is used by any transaction
		exists, err := sess.Cols("uid", "project_id").Where("uid=? AND deleted=? AND project_id=?", uid, false, projectId).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrProjectInUseCannotBeDeleted
		}

		deletedRows, err := sess.ID(projectId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrProjectNotFound
		}

		return err
	})
}

// DeleteAllProjects deletes all existed projects from database
func (s *ProjectService) DeleteAllProjects(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Project{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

// ExistsProjectName returns whether the given project name exists
func (s *ProjectService) ExistsProjectName(c core.Context, uid int64, name string) (bool, error) {
	if name == "" {
		return false, errs.ErrProjectNameIsEmpty
	}

	return s.UserDataDB(uid).NewSession(c).Cols("name").Where("uid=? AND deleted=? AND name=?", uid, false, name).Exist(&models.Project{})
}

// ExistsProjectNameIgnoreGivenId returns whether the given project name exists (excluding specific project id)
func (s *ProjectService) ExistsProjectNameIgnoreGivenId(c core.Context, uid int64, name string, projectId int64) (bool, error) {
	if name == "" {
		return false, errs.ErrProjectNameIsEmpty
	}

	return s.UserDataDB(uid).NewSession(c).Cols("name").Where("uid=? AND deleted=? AND name=? AND project_id<>?", uid, false, name, projectId).Exist(&models.Project{})
}

// GetProjectMapByList returns a project map by a list
func (s *ProjectService) GetProjectMapByList(projects []*models.Project) map[int64]*models.Project {
	projectMap := make(map[int64]*models.Project)

	for i := 0; i < len(projects); i++ {
		project := projects[i]
		projectMap[project.ProjectId] = project
	}
	return projectMap
}

// GetVisibleProjectNameMapByList returns a visible project map by a list
func (s *ProjectService) GetVisibleProjectNameMapByList(projects []*models.Project) map[string]*models.Project {
	projectMap := make(map[string]*models.Project)

	for i := 0; i < len(projects); i++ {
		project := projects[i]

		if project.Hidden {
			continue
		}

		projectMap[project.Name] = project
	}

	return projectMap
}
