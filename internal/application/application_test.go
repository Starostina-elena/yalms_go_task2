package application_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"strconv"

    "github.com/Starostina-elena/yalms_go_task2/internal/application"
)

func TestRPNHandler(t *testing.T) {

	cases := []struct {
        name string
        requestBody application.Request
		requestMethod string
        expectedStatus int
        expectedResult application.Response
    }{
        {
            name:           "Valid Expression",
            requestBody:    application.Request{Expression: "1+1"},
			requestMethod: "POST",
            expectedStatus: http.StatusOK,
            expectedResult: application.Response{Result: "2"},
        },
        {
            name:           "Invalid Expression",
            requestBody:    application.Request{Expression: "1+"},
			requestMethod: "POST",
            expectedStatus: 422,
            expectedResult: application.Response{Error: "Expression is not valid"},
        },
        {
            name:           "Empty Expression",
            requestBody:    application.Request{Expression: ""},
			requestMethod: "POST",
            expectedStatus: 422,
            expectedResult: application.Response{Error: "Expression is not valid"},
        },
		{
            name:           "Expression with letters",
            requestBody:    application.Request{Expression: "a+b"},
			requestMethod: "POST",
            expectedStatus: 422,
            expectedResult: application.Response{Error: "Expression is not valid"},
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            requestBody, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest(tc.requestMethod, "/api/v1/calculate", bytes.NewBuffer(requestBody))

            rr := httptest.NewRecorder()
            handler := http.HandlerFunc(application.RPNHandler)
            handler.ServeHTTP(rr, req)

            if rr.Code != tc.expectedStatus {
                t.Errorf("Test %s. Expected status %v; got %v", tc.name, tc.expectedStatus, rr.Code)
            }

            var answer application.Response
            err := json.NewDecoder(rr.Body).Decode(&answer)
            if err != nil {
                t.Fatalf("Test %s. Decoding error: %v", tc.name, err)
            }

			result, _ := strconv.ParseFloat(answer.Result, 64)
			want, _ := strconv.ParseFloat(tc.expectedResult.Result, 64)
            if result != want {
                t.Errorf("Test %s. Expected result %v; got %v", tc.name, tc.expectedResult, result)
            }
        })
    }
}
