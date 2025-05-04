package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VanLavr/Diploma-fin/internal/controllers/dto"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/utils/auth"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type GroupHandler struct {
	groupUsecase logic.GroupUsecase
}

func NewGroupHandler(groupUsecase logic.GroupUsecase) *GroupHandler {
	return &GroupHandler{
		groupUsecase: groupUsecase,
	}
}

func (this GroupHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/group", this.CreateGroup)                 // + admin
	group.PUT("/group", this.UpdateGroup)                  // + admin
	group.DELETE("/group/:id", this.DeleteGroup)           // + admin
	group.GET("/group/all/:limit/:offset", this.GetGroups) // + admin
	group.GET("/group/:id", this.GetGroup)                 // + admin
}

func (g GroupHandler) CreateGroup(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	var r dto.CreateGroupDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	id, err := g.groupUsecase.CreateGroup(c.Request.Context(), dto.TypesGroupFromCreateGroupDTO(r))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreateExamResponseDTO{
		Err:  nil,
		Data: int(id),
	})
}
func (g GroupHandler) UpdateGroup(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	var r dto.UpdateGroupDTO
	if err := c.Bind(&r); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := g.groupUsecase.UpdateGroup(c.Request.Context(), dto.TypesGroupFromUpdateGroupDTO(r)); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: nil,
	})
}
func (g GroupHandler) DeleteGroup(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := g.groupUsecase.DeleteGroup(c.Request.Context(), int64(id)); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetAllDebtsDTO{
		Err:  nil,
		Data: nil,
	})
}
func (g GroupHandler) GetGroups(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	lim := c.Param("limit")
	limit, err := strconv.Atoi(lim)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	offsetStr := c.Param("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	groups, err := g.groupUsecase.GetGroups(c.Request.Context(), int64(limit), int64(offset))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := make([]dto.Group, 0, len(groups))
	for _, group := range groups {
		result = append(result, dto.GroupDTOFromTypes(group))
	}

	c.JSON(http.StatusOK, dto.GetAllGroupsDTO{
		Err:  nil,
		Data: result,
	})
}
func (g GroupHandler) GetGroup(c *gin.Context) {
	if c.Value(auth.RoleKey) != auth.AdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrUserDoesNotHaveRights})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	group, err := g.groupUsecase.GetGroupByID(c.Request.Context(), int64(id))
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	result := dto.GroupDTOFromTypes(*group)

	c.JSON(http.StatusOK, dto.GetGroupDTO{
		Err:  nil,
		Data: result,
	})
}
