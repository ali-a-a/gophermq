package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ali-a-a/gophermq/internal/app/gophermq/broker"
	"github.com/ali-a-a/gophermq/pkg/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	engine *echo.Echo
	broker *FakeBroker
}

type FakeBroker struct {
	publishErr   error
	fetchErr     error
	subscribeErr error

	data [][]byte
	id   string
}

func (fb *FakeBroker) Publish(_ string, _ []byte) error {
	return fb.publishErr
}

func (fb *FakeBroker) Subscribe(subject string) (*broker.Subscriber, error) {
	return &broker.Subscriber{
		ID:   fb.id,
		Subj: subject,
	}, fb.subscribeErr
}

func (fb *FakeBroker) Fetch(_ string, _ string) ([][]byte, error) {
	return fb.data, fb.fetchErr
}

func (suite *HandlerSuite) SetupSuite() {
	b := FakeBroker{}
	suite.broker = &b

	suite.engine = router.New()
	h := Handler{
		mq:         suite.broker,
		pubReqChan: make(chan *PublishReq),
	}

	suite.engine.POST("/api/publish", h.Publish)
	suite.engine.POST("/api/subscribe", h.Subscribe)
	suite.engine.POST("/api/fetch", h.Fetch)
	suite.engine.POST("/api/publish/async", h.PublishAsync)
}

func (suite *HandlerSuite) TestHandler_Publish() {
	cases := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectedResp   string
		err            error
	}{
		{
			name: "successful",
			request: PublishReq{
				Subject: "test.a",
				Data:    ":(",
			},
			expectedStatus: http.StatusOK,
			expectedResp:   "{\"status\":\"ok\"}\n",
		},
		{
			name: "bad subject",
			request: PublishReq{
				Subject: "test a",
				Data:    ":(",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "{\"message\":\"invalid subject\"}\n",
			err:            broker.ErrBadSubject,
		},
		{
			name: "bad subject",
			request: PublishReq{
				Subject: "test.a",
				Data:    ":(",
			},
			expectedStatus: http.StatusNotFound,
			expectedResp:   "{\"message\":\"subscriber not found\"}\n",
			err:            broker.ErrSubscriberNotFound,
		},
		{
			name: "overflow",
			request: PublishReq{
				Subject: "test.a",
				Data:    ":(",
			},
			expectedStatus: http.StatusTooManyRequests,
			expectedResp:   "{\"message\":\"broker overflow\"}\n",
			err:            broker.ErrMaxPending,
		},
		{
			name:           "bad request",
			request:        "bad request",
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "{\"message\":\"request's body is invalid\"}\n",
		},
		{
			name: "unknown error",
			request: PublishReq{
				Subject: "test.a",
				Data:    ":(",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "{\"message\":\"server error\"}\n",
			err:            errors.New("unknown error"),
		},
	}

	for _, tt := range cases {
		tt := tt

		suite.Run(tt.name, func() {
			data, err := json.Marshal(tt.request)
			suite.NoError(err)

			suite.broker.publishErr = tt.err

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/publish", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tt.expectedStatus, w.Code, tt.name)

			body, _ := ioutil.ReadAll(w.Body)
			suite.Equal(tt.expectedResp, string(body))
		})
	}
}

func (suite *HandlerSuite) TestHandler_Subscribe() {
	cases := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectedResp   string
		id             string
		err            error
	}{
		{
			name: "successful",
			request: SubscribeReq{
				Subject: "test.a",
			},
			id:             "sad id",
			expectedStatus: http.StatusOK,
			expectedResp:   "{\"id\":\"sad id\",\"subject\":\"test.a\"}\n",
		},
		{
			name: "bad subject",
			request: SubscribeReq{
				Subject: "test a",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "{\"message\":\"invalid subject\"}\n",
			err:            broker.ErrBadSubject,
		},
		{
			name:           "bad request",
			request:        "bad request",
			expectedStatus: http.StatusBadRequest,
			expectedResp:   "{\"message\":\"request's body is invalid\"}\n",
		},
		{
			name: "unknown error",
			request: SubscribeReq{
				Subject: "test.a",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "{\"message\":\"server error\"}\n",
			err:            errors.New("unknown error"),
		},
	}

	for _, tt := range cases {
		tt := tt

		suite.Run(tt.name, func() {
			data, err := json.Marshal(tt.request)
			suite.NoError(err)

			suite.broker.subscribeErr = tt.err
			suite.broker.id = tt.id

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/subscribe", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tt.expectedStatus, w.Code, tt.name)

			body, _ := ioutil.ReadAll(w.Body)
			suite.Equal(tt.expectedResp, string(body))
		})
	}
}

func (suite *HandlerSuite) TestHandler_Fetch() {
	cases := []struct {
		name           string
		request        interface{}
		expectedStatus int
		err            error
	}{
		{
			name: "successful",
			request: FetchReq{
				Subject: "test.a",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "bad subject",
			request: FetchReq{
				Subject: "test a",
			},
			expectedStatus: http.StatusNotFound,
			err:            broker.ErrBadID,
		},
		{
			name:           "bad request",
			request:        "bad request",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "unknown error",
			request: FetchReq{
				Subject: "test.a",
			},
			expectedStatus: http.StatusInternalServerError,
			err:            errors.New("unknown error"),
		},
	}

	for _, tt := range cases {
		tt := tt

		suite.Run(tt.name, func() {
			data, err := json.Marshal(tt.request)
			suite.NoError(err)

			suite.broker.fetchErr = tt.err

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/fetch", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tt.expectedStatus, w.Code, tt.name)
		})
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, &HandlerSuite{})
}
