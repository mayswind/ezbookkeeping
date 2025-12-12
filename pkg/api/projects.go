package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// ProjectsApi represents project api
type ProjectsApi struct {
	projects *services.ProjectService
}

// Initialize a project api singleton instance
var (
	Projects = &ProjectsApi{
		projects: services.Projects,
	}
)

// ProjectListHandler returns project list of current user
func (a *ProjectsApi) ProjectListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	projects, err := a.projects.GetAllProjectsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[projects.ProjectListHandler] failed to get projects for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	projectResps := make(models.ProjectInfoResponseSlice, len(projects))

	for i := 0; i < len(projects); i++ {
		projectResps[i] = projects[i].ToProjectInfoResponse()
	}

	sort.Sort(projectResps)

	return projectResps, nil
}

// ProjectGetHandler returns one specific project of current user
func (a *ProjectsApi) ProjectGetHandler(c *core.WebContext) (any, *errs.Error) {
	var projectGetReq models.ProjectGetRequest
	err := c.ShouldBindQuery(&projectGetReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	project, err := a.projects.GetProjectByProjectId(c, uid, projectGetReq.Id)

	if err != nil {
		log.Errorf(c, "[projects.ProjectGetHandler] failed to get project \"id:%d\" for user \"uid:%d\", because %s", projectGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	projectResp := project.ToProjectInfoResponse()

	return projectResp, nil
}

// ProjectCreateHandler saves a new project by request parameters for current user
func (a *ProjectsApi) ProjectCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var projectCreateReq models.ProjectCreateRequest
	err := c.ShouldBindJSON(&projectCreateReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.projects.GetMaxDisplayOrder(c, uid)

	if err != nil {
		log.Errorf(c, "[projects.ProjectCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	project := a.createNewProjectModel(uid, &projectCreateReq, maxOrderId+1)

	err = a.projects.CreateProject(c, project)

	if err != nil {
		log.Errorf(c, "[projects.ProjectCreateHandler] failed to create project \"id:%d\" for user \"uid:%d\", because %s", project.ProjectId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[projects.ProjectCreateHandler] user \"uid:%d\" has created a new project \"id:%d\" successfully", uid, project.ProjectId)

	projectResp := project.ToProjectInfoResponse()

	return projectResp, nil
}

// ProjectModifyHandler saves an existed project by request parameters for current user
func (a *ProjectsApi) ProjectModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var projectModifyReq models.ProjectModifyRequest
	err := c.ShouldBindJSON(&projectModifyReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	project, err := a.projects.GetProjectByProjectId(c, uid, projectModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[projects.ProjectModifyHandler] failed to get project \"id:%d\" for user \"uid:%d\", because %s", projectModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	newProject := &models.Project{
		ProjectId: project.ProjectId,
		Uid:       uid,
		Name:      projectModifyReq.Name,
		Color:     projectModifyReq.Color,
		Comment:   projectModifyReq.Comment,
		Hidden:    projectModifyReq.Hidden,
	}

	if newProject.Name == project.Name &&
		newProject.Color == project.Color &&
		newProject.Comment == project.Comment &&
		newProject.Hidden == project.Hidden {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.projects.ModifyProject(c, newProject)

	if err != nil {
		log.Errorf(c, "[projects.ProjectModifyHandler] failed to update project \"id:%d\" for user \"uid:%d\", because %s", projectModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[projects.ProjectModifyHandler] user \"uid:%d\" has updated project \"id:%d\" successfully", uid, projectModifyReq.Id)

	project.Name = newProject.Name
	project.Color = newProject.Color
	project.Comment = newProject.Comment
	project.Hidden = newProject.Hidden
	projectResp := project.ToProjectInfoResponse()

	return projectResp, nil
}

// ProjectHideHandler hides a project by request parameters for current user
func (a *ProjectsApi) ProjectHideHandler(c *core.WebContext) (any, *errs.Error) {
	var projectHideReq models.ProjectHideRequest
	err := c.ShouldBindJSON(&projectHideReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.projects.HideProject(c, uid, []int64{projectHideReq.Id}, projectHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[projects.ProjectHideHandler] failed to hide project \"id:%d\" for user \"uid:%d\", because %s", projectHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[projects.ProjectHideHandler] user \"uid:%d\" has hidden project \"id:%d\"", uid, projectHideReq.Id)
	return true, nil
}

// ProjectMoveHandler moves display order of existed projects by request parameters for current user
func (a *ProjectsApi) ProjectMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var projectMoveReq models.ProjectMoveRequest
	err := c.ShouldBindJSON(&projectMoveReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	projects := make([]*models.Project, len(projectMoveReq.NewDisplayOrders))

	for i := 0; i < len(projectMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := projectMoveReq.NewDisplayOrders[i]
		project := &models.Project{
			Uid:          uid,
			ProjectId:    newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		projects[i] = project
	}

	err = a.projects.ModifyProjectDisplayOrders(c, uid, projects)

	if err != nil {
		log.Errorf(c, "[projects.ProjectMoveHandler] failed to move projects for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[projects.ProjectMoveHandler] user \"uid:%d\" has moved projects", uid)
	return true, nil
}

// ProjectDeleteHandler deletes an existed project by request parameters for current user
func (a *ProjectsApi) ProjectDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var projectDeleteReq models.ProjectDeleteRequest
	err := c.ShouldBindJSON(&projectDeleteReq)

	if err != nil {
		log.Warnf(c, "[projects.ProjectDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.projects.DeleteProject(c, uid, projectDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[projects.ProjectDeleteHandler] failed to delete project \"id:%d\" for user \"uid:%d\", because %s", projectDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[projects.ProjectDeleteHandler] user \"uid:%d\" has deleted project \"id:%d\"", uid, projectDeleteReq.Id)
	return true, nil
}

func (a *ProjectsApi) createNewProjectModel(uid int64, projectCreateReq *models.ProjectCreateRequest, order int32) *models.Project {
	return &models.Project{
		Uid:          uid,
		Name:         projectCreateReq.Name,
		Color:        projectCreateReq.Color,
		Comment:      projectCreateReq.Comment,
		DisplayOrder: order,
	}
}
