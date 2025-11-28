package enrollment_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"

	courseSdk "github.com/DanyJDuque/go_course_skd/course"
	mockCourseSdk "github.com/DanyJDuque/go_course_skd/course/mock"
	userSdk "github.com/DanyJDuque/go_course_skd/user"
	mockUserSdk "github.com/DanyJDuque/go_course_skd/user/mock"

	"github.com/DanyJDuque/go_lib_response/response"
	"github.com/DanyJDuque/gocourse_domain/domain"
	"github.com/DanyJDuque/gocourse_enrollment/internal/enrollment"
	"github.com/stretchr/testify/assert"
)

func TestCreateEndpoint(t *testing.T) {

	l := log.New(io.Discard, "", 0)

	t.Run("should return bad request error when user id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{})

		assert.NotNil(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrUserIdRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("should return bad request when course id is empty", func(t *testing.T) {
		endpoint := enrollment.MakeEndpoints(nil, enrollment.Config{})
		_, err := endpoint.Create(context.Background(), enrollment.CreateReq{UserID: "123"})

		assert.Error(t, err)

		resp := err.(response.Response)
		assert.EqualError(t, enrollment.ErrCourseIdRequired, resp.Error())
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	obj := []struct {
		tag            string
		repositoryMock enrollment.Repository
		userSdkMock    userSdk.Transport
		courseSdkMock  courseSdk.Transport
		wantErr        error
		wantCode       int
		wantResponse   *domain.Enrollment
	}{
		{
			tag: "should return and error if user sdk resturns an unexpected error",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return and error if user does not exist",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, userSdk.ErrNotFound{Message: "user not found"}
				},
			},
			wantErr:  userSdk.ErrNotFound{Message: "user not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return and error if course sdk resturns an unexpected error",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSDKMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return and error if course does not exist",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSDKMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, courseSdk.ErrNotFound{Message: "course not found"}
				},
			},
			wantErr:  courseSdk.ErrNotFound{Message: "course not found"},
			wantCode: http.StatusNotFound,
		},
		{
			tag: "should return an error if repository returns an unexpected error",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSDKMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRespository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) error {
					return errors.New("unexpected error")
				},
			},
			wantErr:  errors.New("unexpected error"),
			wantCode: http.StatusInternalServerError,
		},
		{
			tag: "should return the enrollment",
			userSdkMock: &mockUserSdk.UserSDKMock{
				GetMock: func(id string) (*domain.User, error) {
					return nil, nil
				},
			},
			courseSdkMock: &mockCourseSdk.CourseSDKMock{
				GetMock: func(id string) (*domain.Course, error) {
					return nil, nil
				},
			},
			repositoryMock: &mockRespository{
				CreateMock: func(ctx context.Context, enrollment *domain.Enrollment) error {
					enrollment.ID = "10010"
					return nil
				},
			},
			wantCode: http.StatusCreated,
			wantResponse: &domain.Enrollment{
				ID:       "10010",
				UserID:   "1",
				CourseID: "4",
				Status:   "P",
			},
		},
	}

	for _, obj := range obj {
		t.Run(obj.tag, func(t *testing.T) {
			service := enrollment.NewService(l, obj.userSdkMock, obj.courseSdkMock, obj.repositoryMock)
			endpoint := enrollment.MakeEndpoints(service, enrollment.Config{})
			resp, err := endpoint.Create(context.Background(), enrollment.CreateReq{
				UserID:   "1",
				CourseID: "4",
			})

			if obj.wantErr != nil {
				assert.NotNil(t, err)
				assert.Nil(t, resp)

				respErr := err.(response.Response)
				assert.EqualError(t, obj.wantErr, respErr.Error())
				assert.Equal(t, obj.wantCode, respErr.StatusCode())
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, err)

				r := resp.(response.Response)
				assert.Equal(t, obj.wantCode, r.StatusCode())
				assert.Empty(t, r.Error())

				enrollment := r.GetData().(*domain.Enrollment)
				assert.Equal(t, obj.wantResponse.ID, enrollment.ID)
				assert.Equal(t, obj.wantResponse.UserID, enrollment.UserID)
				assert.Equal(t, obj.wantResponse.CourseID, enrollment.CourseID)
				assert.Equal(t, obj.wantResponse.Status, enrollment.Status)
			}
		})

	}
}
