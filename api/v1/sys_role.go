package v1

import (
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/utils"
)

// FindRoleByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Role
// @Description FindRoleByIds
// @Param ids path string true "ids"
// @Router /role/list/{ids} [GET]
func FindRoleByIds(c *gin.Context) {
	ids := req.UintIds(c)
	s := cache_service.New(c)
	list := s.FindRoleByIds(ids)
	resp.SuccessWithData(list, &[]response.Role{})
}

// FindRole
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Role
// @Description FindRole
// @Param params query request.Role true "params"
// @Router /role/list [GET]
func FindRole(c *gin.Context) {
	var r request.Role
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	// bind current user role sort(low level cannot view high level)
	r.CurrentRoleSort = *user.Role.Sort

	s := cache_service.New(c)
	list := s.FindRole(&r)
	resp.SuccessWithPageData(list, &[]response.Role{}, r.Page)
}

// CreateRole
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Role
// @Description CreateRole
// @Param params body request.CreateRole true "params"
// @Router /role/create [POST]
func CreateRole(c *gin.Context) {
	user := GetCurrentUser(c)
	var r request.CreateRole
	req.ShouldBind(c, &r)
	req.Validate(c, r, r.FieldTrans())

	if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
		resp.CheckErr("sort must >= %d", *user.Role.Sort)
	}

	s := service.New(c)
	err := s.Q.Create(r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

// UpdateRoleById
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Role
// @Description UpdateRoleById
// @Param id path uint true "id"
// @Param params body request.UpdateRole true "params"
// @Router /role/update/{id} [PATCH]
func UpdateRoleById(c *gin.Context) {
	var r request.UpdateRole
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	if r.Sort != nil {
		user := GetCurrentUser(c)
		if r.Sort != nil && *user.Role.Sort > uint(*r.Sort) {
			resp.CheckErr("sort must >= %d", *user.Role.Sort)
		}
	}

	user := GetCurrentUser(c)
	if r.Status != nil && uint(*r.Status) == models.SysRoleStatusDisabled && id == user.RoleId {
		resp.CheckErr("cannot disable your role")
	}

	s := service.New(c)
	err := s.Q.UpdateById(id, r, new(models.SysRole))
	resp.CheckErr(err)
	resp.Success()
}

func RouterFindRoleKeywordByRoleIds(c *gin.Context, roleIds []uint) []string {
	s := cache_service.New(c)
	roles := s.FindRoleByIds(roleIds)
	keywords := make([]string, 0)
	for _, role := range roles {
		keywords = append(keywords, role.Keyword)
	}
	return keywords
}

// BatchDeleteRoleByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags Role
// @Description BatchDeleteRoleByIds
// @Param ids body req.Ids true "ids"
// @Router /role/delete/batch [DELETE]
func BatchDeleteRoleByIds(c *gin.Context) {
	var r req.Ids
	req.ShouldBind(c, &r)
	user := GetCurrentUser(c)
	if utils.ContainsUint(r.Uints(), user.RoleId) {
		resp.CheckErr("cannot delete your role")
	}

	s := service.New(c)
	err := s.DeleteRoleByIds(r.Uints())
	resp.CheckErr(err)
	resp.Success()
}
