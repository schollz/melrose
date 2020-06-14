package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antonmedv/expr/file"
	"github.com/emicklei/melrose"
	"github.com/emicklei/melrose/dsl"
	"github.com/emicklei/melrose/notify"
)

// LanguageServer can execute DSL statements received over HTTP
type LanguageServer struct {
	store     dsl.VariableStorage
	address   string
	evaluator *dsl.Evaluator
}

// NewLanguageServer returns a new LanguageService. It is not started.
func NewLanguageServer(store dsl.VariableStorage, loopControl melrose.LoopController, addr string) *LanguageServer {
	return &LanguageServer{store: store, address: addr, evaluator: dsl.NewEvaluator(store, loopControl)}
}

// Start will start a HTTP listener for serving DSL statements
// curl -v -d 'n = note("C")' http://localhost:8118/v1/statements
func (l *LanguageServer) Start() error {
	http.HandleFunc("/v1/statements", l.statementHandler)
	return http.ListenAndServe(l.address, nil)
}

func (l *LanguageServer) statementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	trace := query.Get("trace") == "true"
	if trace {
		log.Printf("[melrose.trace] %s\n", r.URL.String())
	}
	// get line
	line := 1
	lineString := query.Get("line")
	if len(lineString) > 0 {
		if i, err := strconv.Atoi(lineString); err == nil {
			line = i
		}
	}
	// get expression
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	returnValue, err := l.evaluator.EvaluateProgram(string(data))
	var response evaluationResult
	if err != nil {
		// evaluation failed.
		w.WriteHeader(http.StatusInternalServerError)
		response = resultFrom(line, err)
	} else {
		// evaluation was ok.

		if query.Get("action") == "inspect" {
			melrose.PrintValue(returnValue)
		}

		// check if play was requested and is playable
		if query.Get("action") == "play" {
			if s, ok := returnValue.(melrose.Sequenceable); ok {
				melrose.Context().AudioDevice.Play(
					s,
					melrose.Context().LoopControl.BPM(),
					time.Now())
			}
			if lp, ok := returnValue.(*melrose.Loop); ok {
				melrose.Context().LoopControl.Begin(lp)
			}
		}
		// loop operation
		if query.Get("action") == "begin" {
			if lp, ok := returnValue.(*melrose.Loop); ok {
				if !lp.IsRunning() {
					melrose.Context().LoopControl.Begin(lp)
				}
			}
			// ignore if not Loop
		}
		// loop operation
		if query.Get("action") == "end" {
			if lp, ok := returnValue.(*melrose.Loop); ok {
				if lp.IsRunning() {
					lp.Stop()
				}
			}
			// ignore if not Loop
		}
		if query.Get("action") == "kill" {
			// kill the play and any loop
			melrose.Context().LoopControl.Reset()
			melrose.Context().AudioDevice.Reset()
		}
		response = resultFrom(line, returnValue)
	}
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(response)
	if err != nil {
		log.Printf("[melrose.error] %#v\n", err)
	}
	if response.IsError {
		notify.Print(notify.Error(response.Object.(error)))
	} else {
		melrose.PrintValue(response.Object)
	}
	if trace {
		// doit again
		buf := bytes.Buffer{}
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		err = enc.Encode(response)
		log.Printf("[melrose.trace] %#v, error:%v\n", buf.String(), err)
	}
}

type evaluationResult struct {
	Type     string      `json:"type"`
	IsError  bool        `json:"is-error"`
	Message  string      `json:"message"`
	Filename string      `json:"file"`
	Line     int         `json:"line"`
	Column   int         `json:"column"`
	Object   interface{} `json:"object"`
}

func resultFrom(line int, val interface{}) evaluationResult {
	t := fmt.Sprintf("%T", val)
	if err, ok := val.(error); ok {
		// patch Location of error
		if fe, ok := err.(*file.Error); ok {
			fe.Location.Line = fe.Location.Line - 1 + line
		}
		return evaluationResult{
			Type:     t,
			IsError:  true,
			Filename: "yours.mel",
			Message:  err.Error(),
			Line:     line,
			Object:   val,
		}
	}
	// no error
	var msg string
	if stor, ok := val.(melrose.Storable); ok {
		msg = stor.Storex()
	} else {
		msg = fmt.Sprintf("%v", val)
	}
	// no Object if ok
	return evaluationResult{Type: t, IsError: false, Message: msg}
}
