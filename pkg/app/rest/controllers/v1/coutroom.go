package v1

import (
	"strconv"

	"eventSourcedBooks/pkg/app/rest/dto/req"
	"eventSourcedBooks/pkg/app/rest/dto/res"
	"eventSourcedBooks/pkg/domain/courtroom"
	"eventSourcedBooks/pkg/infra/logger"
	"eventSourcedBooks/pkg/infra/standard"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CourtroomController struct {
	svc  *courtroom.CourtroomService
	lgrf *logger.LoggerFactory
}

// CreateCourtroom godoc
// @Summary Create courtroom
// @Schemes
// @Description Create courtroom
// @Tags role
// @Accept json
// @Param role body req.CreateUpdateCourtroom true "create courtroom body"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 201 {object} []res.Courtroom{}
// @Router /api/v1/courtrooms/ [post]
func (ctrl *CourtroomController) createCourtroom(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing create courtroom")

	rbody := req.CreateUpdateCourtroom{}
	err := ctx.Bind(&rbody)
	if err != nil {
		lgr.Error("failed to unmarshal payload")
		ctx.Error(err)
		ctx.Error(standard.NewUnbindableError())
		return
	}
	if rbody.Title == nil ||
		rbody.Description == nil ||
		rbody.DateTimeStart == nil ||
		rbody.DateTimeEnd == nil ||
		rbody.Location == nil {
		lgr.Error("body fields are missing", zap.Any("body", rbody))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}

	// Smol hack to fix nil on metadata categories being nil
	if rbody.Metadata == nil {
		rbody.Metadata = map[string]interface{}{}
	}
	if rbody.Categories == nil {
		rbody.Categories = []string{}
	}

	resp, err := ctrl.svc.CreateCourtroom(
		ctx,
		*rbody.Title,
		*rbody.Description,
		*rbody.DateTimeStart,
		*rbody.DateTimeEnd,
		*rbody.Location,
		rbody.Categories,
		rbody.Metadata,
	)
	if err != nil {
		lgr.Error("failed to create courtroom")
		ctx.Error(err)
		return
	}
	ctx.JSON(201, res.MapCourtroomToDto(*resp))
}

// UpdateCourtroom godoc
// @Summary Update courtroom
// @Schemes
// @Description Update courtroom
// @Tags role
// @Accept json
// @Param role body req.CreateUpdateCourtroom true "create courtroom body"
// @Param courtroomId path string true "Courtroom Id"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} []res.Courtroom{}
// @Router /api/v1/courtrooms/{courtroomId} [patch]
func (ctrl *CourtroomController) updateCourtroom(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing update courtroom")

	rbody := req.CreateUpdateCourtroom{}
	err := ctx.Bind(&rbody)
	if err != nil {
		lgr.Error("failed to unmarshal payload")
		ctx.Error(standard.NewUnbindableError())
		return
	}

	if rbody.Title == nil &&
		rbody.Description == nil &&
		rbody.DateTimeStart == nil &&
		rbody.DateTimeEnd == nil &&
		rbody.Location == nil &&
		rbody.Categories == nil &&
		rbody.Metadata == nil {
		lgr.Error("body fields are missing", zap.Any("body", rbody))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}

	id := ctx.Param("courtroomId")
	if id == "" {
		lgr.Error("params are missing", zap.String("courtroomId", id))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}
	resp, err := ctrl.svc.UpdateCourtroom(
		ctx,
		id,
		rbody.Title,
		rbody.Description,
		rbody.DateTimeStart,
		rbody.DateTimeEnd,
		rbody.Location,
		rbody.Categories,
		rbody.Metadata,
	)
	if err != nil {
		lgr.Error("failed to update courtroom")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDto(*resp))
}

// GetCourtroom godoc
// @Summary Get courtroom
// @Schemes
// @Description Get courtroom
// @Tags role
// @Param courtroomId path string true "Courtroom Id"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} []res.Courtroom{}
// @Router /api/v1/courtrooms/{courtroomId} [get]
func (ctrl *CourtroomController) getCourtroom(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing get courtroom by id")

	id := ctx.Param("courtroomId")
	if id == "" {
		lgr.Error("params are missing", zap.String("courtroomId", id))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}
	resp, err := ctrl.svc.GetCourtroom(ctx, id)
	if err != nil {
		lgr.Error("failed to get courtroom")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDto(*resp))
}

// DeleteCourtroom godoc
// @Summary Delete courtroom
// @Schemes
// @Description Delete courtroom
// @Tags role
// @Param courtroomId path string true "Courtroom Id"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} []res.Courtroom{}
// @Router /api/v1/courtrooms/{courtroomId} [delete]
func (ctrl *CourtroomController) deleteCourtroom(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing delete courtroom")

	id := ctx.Param("courtroomId")
	if id == "" {
		lgr.Error("params are missing", zap.String("courtroomId", id))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}
	resp, err := ctrl.svc.DeleteCourtroom(ctx, id)
	if err != nil {
		lgr.Error("failed to delete courtroom")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDto(*resp))
}

// QueryCourtroom godoc
// @Summary Query courtrooms
// @Schemes
// @Description Query courtrooms
// @Tags role
// @Param role body []req.SearchQuery true "query courtroom body"
// @Param page query string false "Page Number"
// @Param items query string false "Items per Page"
// @Param orderby query string false "Property to order by"
// @Param desc query boolean false "Order by descending"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} []res.Courtroom{}
// @Router /api/v1/courtrooms/query [post]
func (ctrl *CourtroomController) queryCourtroom(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing query courtroom")

	pageQry := ctx.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageQry)
	if err != nil {
		lgr.Error("invalid query property", zap.String("pageQry", pageQry))
		ctx.Error(standard.NewInvalidQueryPropertyError())
		return
	}
	itemsQry := ctx.DefaultQuery("items", "100")
	items, err := strconv.Atoi(itemsQry)
	if err != nil {
		lgr.Error("invalid query property", zap.String("itemsQry", itemsQry))
		ctx.Error(standard.NewInvalidQueryPropertyError())
		return
	}
	orderby := ctx.Query("orderby")

	descQry := ctx.DefaultQuery("desc", "false")
	desc, err := strconv.ParseBool(descQry)
	if err != nil {
		lgr.Error("invalid query property", zap.String("descQry", descQry))
		ctx.Error(standard.NewInvalidQueryPropertyError())
		return
	}

	rbody := []req.SearchQuery{}
	err = ctx.Bind(&rbody)
	if err != nil {
		lgr.Error("failed to unmarshal payload")
		ctx.Error(standard.NewUnbindableError())
		return
	}

	resp, err := ctrl.svc.QueryCourtroom(
		ctx,
		page,
		items,
		req.MapDtoToSearchQueryList(rbody),
		orderby,
		desc,
	)
	if err != nil {
		lgr.Error("failed to query courtroom")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDtoList(resp))
}

// UpdateCourtroomMeta godoc
// @Summary Update courtroom metadata
// @Schemes
// @Description Update courtroom metadata
// @Tags role
// @Accept json
// @Param courtroomId path string true "Courtroom Id"
// @Param role body map[string]interface{} true "update key values"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} res.Courtroom{}
// @Router /api/v1/courtrooms/{courtroomId}/meta [patch]
func (ctrl *CourtroomController) updateCourtroomMeta(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing update courtroom meta")

	rbody := map[string]any{}
	err := ctx.Bind(&rbody)
	if err != nil {
		lgr.Error("failed to unmarshal payload")
		ctx.Error(err)
		ctx.Error(standard.NewUnbindableError())
		return
	}
	id := ctx.Param("courtroomId")
	if id == "" {
		lgr.Error("params are missing", zap.String("courtroomId", id))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}

	resp, err := ctrl.svc.UpdateCourtroomMeta(
		ctx,
		id,
		rbody,
	)
	if err != nil {
		lgr.Error("failed to update courtroom meta")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDto(*resp))
}

// DeleteCourtroomMeta godoc
// @Summary Delete courtroom metadata
// @Schemes
// @Description Delete courtroom metadata
// @Tags role
// @Accept json
// @Param courtroomId path string true "Courtroom Id"
// @Param role body []string true "Delete Keys"
// @Param Authorization header string true "Bearer"
// @Produce json
// @Success 200 {object} res.Courtroom{}
// @Router /api/v1/courtrooms/{courtroomId}/meta [delete]
func (ctrl *CourtroomController) deleteCourtroomMeta(
	ctx *gin.Context,
) {
	lgr := ctrl.lgrf.NewLogger(ctx)
	lgr.Info("processing delete courtroom meta")

	rbody := []string{}
	err := ctx.Bind(&rbody)
	if err != nil {
		lgr.Error("failed to unmarshal payload")
		ctx.Error(err)
		ctx.Error(standard.NewUnbindableError())
		return
	}
	id := ctx.Param("courtroomId")
	if id == "" {
		lgr.Error("params are missing", zap.String("courtroomId", id))
		ctx.Error(standard.NewPropertyMissingError())
		return
	}

	resp, err := ctrl.svc.DeleteCourtroomMeta(
		ctx,
		id,
		rbody,
	)
	if err != nil {
		lgr.Error("failed to delete courtroom meta")
		ctx.Error(err)
		return
	}
	ctx.JSON(200, res.MapCourtroomToDto(*resp))
}

func (ctrl *CourtroomController) RegisterRoutes(grp *gin.RouterGroup) {
	grp.POST("", ctrl.createCourtroom)
	grp.PATCH("/:courtroomId", ctrl.updateCourtroom)
	grp.GET("/:courtroomId", ctrl.getCourtroom)
	grp.DELETE("/:courtroomId", ctrl.deleteCourtroom)
	grp.POST("/query", ctrl.queryCourtroom)
	grp.PATCH("/:courtroomId/meta", ctrl.updateCourtroomMeta)
	grp.DELETE("/:courtroomId/meta", ctrl.deleteCourtroomMeta)
}

func NewCourtroomController(
	svc *courtroom.CourtroomService,
	lgrf *logger.LoggerFactory,
) *CourtroomController {
	return &CourtroomController{
		svc:  svc,
		lgrf: lgrf,
	}
}
