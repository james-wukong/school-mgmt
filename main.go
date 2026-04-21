package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"                 // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres" // sql driver
	_ "github.com/GoAdminGroup/themes/adminlte"                      // ui theme

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/gin-gonic/gin"

	"github.com/james-wukong/online-school-mgmt/internal/api/ajax"
	"github.com/james-wukong/online-school-mgmt/internal/api/export"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/tables"
	"github.com/james-wukong/online-school-mgmt/pages"
)

func main() {
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()
	eng.AddConfigFromYAML("./config.yml")
	// Get the initialized connection
	conn := eng.PostgresqlConnection()

	// Check if connection is valid
	if conn == nil {
		panic(errors.New("database connection is nil"))
	}
	db := models.Init(conn)

	if err := eng.AddGenerators(tables.Generators).
		AddGenerators(tables.GetGenerators(db)).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	// AjAX api
	eng.Data("POST", "/admin/ajax/teacher/sem_timeslot",
		ajax.AjaxTeacherSemesterTSHandler(db),
		false,
	)
	eng.Data("POST", "/admin/ajax/room/sem_timeslot",
		ajax.AjaxRoomSemesterTSHandler(db),
		false,
	)
	eng.Data("POST", "/export/reports",
		export.ExportReportHandler(db),
		true,
	)

	eng.HTML("GET", "/admin", pages.GetDashBoard)
	eng.HTMLFile("GET", "/admin/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	// if err := r.Run(":8091"); err != nil {
	// 	panic(err)
	// }
	srv := &http.Server{
		Addr:         ":8091",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.PostgresqlConnection().Close()
}
