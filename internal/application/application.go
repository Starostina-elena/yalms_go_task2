package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
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
		log.Printf("[ERROR] Error while decoding json: %s", err)
		return
	}

	result := Response{}
	result_calc, err := rpn.Calc(request.Expression)
	if err != nil {
		w.WriteHeader(422)
		result.Error = "Expression is not valid"
		log.Printf("Rejected %s: wrong expression", request.Expression)
	} else {
		result.Result = fmt.Sprintf("%f", result_calc)
	}
	jsonBytes, _ := json.Marshal(result)
    fmt.Fprintf(w, string(jsonBytes))
	log.Printf("Finished with request")
}

func Answer500(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
				log.Printf("[ERROR] Internal server error: %s on request", err)
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

func CheckMethodIsPost(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got new request")
        if r.Method != "POST" {
			log.Println("Rejected request: wrong method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
        next.ServeHTTP(w, r)
    })
}


func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CheckMethodIsPost(Answer500(RPNHandler)))
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
