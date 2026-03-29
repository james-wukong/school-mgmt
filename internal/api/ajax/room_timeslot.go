package ajax

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func AjaxRoomSemesterTSHanlder(dbConn *gorm.DB) context.Handler {

	// 1. Authenticate the request
	// 3. Return JSON response
	// This matches the 'res.code' check in your JavaScript
	return func(ctx *context.Context) {
		var options []map[string]any
		var room *model2.Rooms
		var isUpdate bool

		// 1. Get params from query
		roomID, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"data": map[string]any{"room id": ctx.Query("id")},
			})
			return
		}
		semID, err := strconv.ParseInt(ctx.Query("value"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"data": map[string]any{"value": ctx.Query("value")},
			})
			return
		}
		schoolID, err := strconv.ParseInt(ctx.Query("school_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"data": map[string]any{"schoolid": ctx.Query("school_id")},
			})
			return
		}

		// 2. List all timeslots by school_id and semester_id
		slotService := services.NewTimeslotService(dbConn)
		slots, err := slotService.List(
			ctx.Request.Context(), schoolID, semID, 0,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"data": map[string]any{},
			})
			return
		}
		// 3. Check update
		if roomID != 0 {
			isUpdate = true
		}

		// 3. Get teacher
		roomService := services.NewRoomService(dbConn)
		if isUpdate {
			room, err = roomService.GetRoom(ctx.Request.Context(), roomID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code": http.StatusBadRequest,
					"msg":  err.Error(),
					"data": map[string]any{},
				})
				return
			}
		}
		// 4. create return data
		for _, s := range slots {
			opt := map[string]any{
				"text": fmt.Sprintf("%d: %s-%s",
					s.DayOfWeek,
					s.StartTime.Format(model2.TimeSlotLayout),
					s.EndTime.Format(model2.TimeSlotLayout),
				),
				"value": fmt.Sprint(s.ID),
			}

			if isUpdate {
				if exists := slices.ContainsFunc(room.Timeslots,
					func(t *model2.Timeslots) bool {
						return t.ID == s.ID
					}); exists {
					opt["selected"] = true
				}
			}
			options = append(options, opt)
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"msg":  "Success",
			"data": options,
		})
	}
}
