package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	courseSdk "github.com/DanyJDuque/go_course_skd/course/mock"
	userSdk "github.com/DanyJDuque/go_course_skd/user/mock"
	"github.com/DanyJDuque/gocourse_domain/domain"
	"github.com/DanyJDuque/gocourse_enrollment/internal/enrollment"
	"github.com/DanyJDuque/gocourse_enrollment/pkg/bootstrap"
	"github.com/DanyJDuque/gocourse_enrollment/pkg/handler"
	"github.com/joho/godotenv"
	"github.com/ncostamagna/go_http_client/client"
)

var cli client.Transport

func TestMain(m *testing.M) {

	_ = godotenv.Load("../.env")
	l := log.New(io.Discard, "", 0)

	db, err := bootstrap.DBConection()
	if err != nil {
		l.Fatal(err)
	}

	tx := db.Begin()

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	userSdk := &userSdk.UserSDKMock{
		GetMock: func(id string) (*domain.User, error) {
			return nil, nil

		},
	}

	courseSdk := &courseSdk.CourseSDKMock{
		GetMock: func(id string) (*domain.Course, error) {
			return nil, nil
		},
	}

	ctx := context.Background()
	enrollRepo := enrollment.NewRepo(l, tx)
	enrollSrv := enrollment.NewService(l, userSdk, courseSdk, enrollRepo)

	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollSrv, enrollment.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	// header := http.Header{}
	// header.Set("Content-Type", "application/json")

	cli = client.New(nil, "http://"+address, 0, false)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	r := m.Run()

	if err := srv.Shutdown(context.Background()); err != nil {
		l.Println(err)
	}

	tx.Rollback()
	os.Exit(r) // el cierre de los tests

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)

	})
}
