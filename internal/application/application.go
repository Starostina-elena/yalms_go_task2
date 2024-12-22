package application

import (
	//"bufio"
	"encoding/json"
	//"errors"
	"fmt"
	//"log"
	"net/http"
	"os"
	//"strings"

	"github.com/Starostina-elena/yalms_go_task2/pkg/rpn"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}


type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

func RPNHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := Response{}
	result_calc, err := rpn.Calc(request.Expression)
	if err != nil {
		w.WriteHeader(422)
		result.Error = "Expression is not valid"
	} else {
		result.Result = fmt.Sprintf("%f", result_calc)
	}
	jsonBytes, _ := json.Marshal(result)
    fmt.Fprintf(w, string(jsonBytes))
}

func Answer500(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
				result := Response{}
                result.Error = "Internal server error"
                jsonBytes, _ := json.Marshal(result)
				w.WriteHeader(500)
                fmt.Fprintf(w, string(jsonBytes))
            }
        }()
        next.ServeHTTP(w, r)
    })
}


func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", Answer500(RPNHandler))
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
