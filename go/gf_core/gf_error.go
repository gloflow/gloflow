/*
MIT License

Copyright (c) 2019 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_core

import (
	"fmt"
	"time"
	"strings"
	"errors"
	"context"
	"runtime"
	"runtime/debug"
	"github.com/fatih/color"
	"github.com/globalsign/mgo/bson"
	"github.com/getsentry/sentry-go"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFerror struct {
	Id                   bson.ObjectId          `bson:"_id,omitempty"`
	Id_str               string                 `bson:"id_str"` 
	T_str                string                 `bson:"t"`                    // "gf_error"
	Creation_unix_time_f float64                `bson:"creation_unix_time_f"`
	Type_str             string                 `bson:"type_str"`
	User_msg_str         string                 `bson:"user_msg_str"`
	Data_map             map[string]interface{} `bson:"data_map"`
	Descr_str            string                 `bson:"descr_str"`
	Error                error                  `bson:"error"`
	ServiceNameStr       string                 `bson:"service_name_str"`
	Subsystem_name_str   string                 `bson:"subsystem_name_str"`   // major portion of functionality, a particular package, or a logical group of functions
	Stack_trace_str      string                 `bson:"stack_trace_str"`
	Function_name_str    string                 `bson:"func_name_str"`
	File_str             string                 `bson:"file_str"`
	Line_num_int         int                    `bson:"line_num_int"`
}

//-------------------------------------------------

func PanicCheckAndHandle(pUserMsgStr string,
	pPanicDataMap     map[string]interface{},
	pOnCompleteFun    func(),
	pSubsystemNameStr string,
	pRuntimeSys       *RuntimeSys) {

	// call to recover stops the unwinding and returns the argument passed to panic
	// If the goroutine is not panicking, recover returns nil.
	if panicInfo := recover(); panicInfo != nil {

		err := errors.New(fmt.Sprint(panicInfo))

		fmt.Println("PANIC >>>>>")
		spew.Dump(panicInfo)
		fmt.Println(err)

		//--------------------
		// SENTRY
		if pRuntimeSys.ErrorsSendToSentryBool {
			
			/*sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetExtra("gf_error.service_name",   gf_error.ServiceNameStr)
				scope.SetExtra("gf_error.subsystem_name", gf_error.Subsystem_name_str)
				scope.SetExtra("gf_error.type",           gf_error.Type_str)
			})*/

			sentry.WithScope(func(scope *sentry.Scope) {

				scope.SetTag(fmt.Sprintf("%s_panic.service_name",   pRuntimeSys.NamesPrefixStr), pRuntimeSys.ServiceNameStr)
				scope.SetTag(fmt.Sprintf("%s_panic.subsystem_name", pRuntimeSys.NamesPrefixStr), pSubsystemNameStr)
				scope.SetTag(fmt.Sprintf("%s_panic.type",           pRuntimeSys.NamesPrefixStr), "panic_error")

				for k, v := range pPanicDataMap {
					scope.SetTag(fmt.Sprintf("%s_panic.%s", pRuntimeSys.NamesPrefixStr, k),
						fmt.Sprint(v))
				}

				// scope.SetLevel(sentry.LevelWarning);

				sentry.CaptureException(err)
			})

			// FLUSH
			defer sentry.Flush(2 * time.Second)
		}

		//--------------------

		if pOnCompleteFun != nil {
			pOnCompleteFun()
		}
	}
}

//-------------------------------------------------

func ErrorCreateWithHook(pUserMsgStr string,
	pErrorTypeStr          string,
	pErrorDataMap          map[string]interface{},
	pError                 error,
	pSubsystemNameStr      string,
	pSkipStackFramesNumInt int,
	p_hook_fun             func(*GFerror) map[string]interface{},
	pRuntimeSys            *RuntimeSys) *GFerror {

	gfErr := ErrorCreateWithStackSkip(pUserMsgStr,
		pErrorTypeStr,
		pErrorDataMap,
		pError,
		pSubsystemNameStr,
		3,
		pRuntimeSys)

	p_hook_fun(gfErr)
	return gfErr
}

//-------------------------------------------------

func ErrorCreateWithStackSkip(pUserMsgStr string,
	pErrorTypeStr          string,
	pErrorDataMap          map[string]interface{},
	pError                 error,
	pSubsystemNameStr      string,
	pSkipStackFramesNumInt int,
	pRuntimeSys            *RuntimeSys) *GFerror {
	
	error_defs_map := errorGetDefs()
	
	gf_err := ErrorCreateWithDefs(pUserMsgStr,
		pErrorTypeStr,
		pErrorDataMap,
		pError,
		pSubsystemNameStr,
		error_defs_map,
		pSkipStackFramesNumInt,
		pRuntimeSys)

	return gf_err
}

//-------------------------------------------------

func ErrorCreate(pUserMsgStr string,
	pErrorTypeStr     string,
	pErrorDataMap     map[string]interface{},
	pError            error,
	pSubsystemNameStr string,
	pRuntimeSys       *RuntimeSys) *GFerror {

	error_defs_map := errorGetDefs()
	
	gf_err := ErrorCreateWithDefs(pUserMsgStr,
		pErrorTypeStr,
		pErrorDataMap,
		pError,
		pSubsystemNameStr,
		error_defs_map,

		// IMPORTANT!! - ErrorCreateWithDefs() is 2 stack levels away from the caller
		//               of ErrorCreate() so its important to account for that to get 
		//               the proper caller information.
		2, // pSkipStackFramesNumInt
		pRuntimeSys)

	return gf_err
}

//-------------------------------------------------

func ErrorCreateWithDefs(pUserMsgStr string,
	pErrorTypeStr     string,
	pErrorDataMap     map[string]interface{},
	pError            error,
	pSubsystemNameStr string,
	pErrDefsMap       map[string]ErrorDef,

	// IMPORTANT!! - number of stack frames to skip before recording. without skipping 
	//               we would get info on this function, not its caller which is where
	//               the error occured.
	pSkipStackFramesNumInt int,

	pRuntimeSys *RuntimeSys) *GFerror {

	

	creationUNIXtimeF := float64(time.Now().UnixNano()) / 1000000000.0
	idStr             := fmt.Sprintf("%s:%f", pErrorTypeStr, creationUNIXtimeF)
	stackTraceStr     := string(debug.Stack())

	// // IMPORTANT!! - number of stack frames to skip before recording. without skipping 
	// //               we would get info on this function, not its caller which is where
	// //               the error occured.
	// skip_stack_frames_num_int := 1

	// https://golang.org/pkg/runtime/#Caller
	programCounter, fileStr, line_num_int,_ := runtime.Caller(pSkipStackFramesNumInt)

	// FuncForPC - returns a *Func describing the function that contains the given program counter address
	function        := runtime.FuncForPC(programCounter)
	functionNameStr := function.Name()

	//--------------------
	// ERROR_DEF

	errorDef, ok := pErrDefsMap[pErrorTypeStr]
	if !ok {
		panic(fmt.Sprintf("unknown gf_error type encountered - %s", pErrorTypeStr))
	}

	//--------------------

	gfErr := GFerror{
		Id_str:               idStr,
		T_str:                "gf_error",
		Creation_unix_time_f: creationUNIXtimeF,
		Type_str:             pErrorTypeStr,
		User_msg_str:         pUserMsgStr,
		Data_map:             pErrorDataMap,
		Descr_str:            errorDef.DescrStr,
		Error:                pError,
		ServiceNameStr:       pRuntimeSys.ServiceNameStr,
		Subsystem_name_str:   pSubsystemNameStr,
		Stack_trace_str:      stackTraceStr,
		Function_name_str:    functionNameStr,
		File_str:             fileStr,
		Line_num_int:         line_num_int,
	}

	red      := color.New(color.FgRed).SprintFunc()
	cyan     := color.New(color.FgCyan, color.BgWhite).SprintFunc()
	yellow   := color.New(color.FgYellow).SprintFunc()
	yellowBg := color.New(color.FgBlack, color.BgYellow).SprintFunc()
	green    := color.New(color.FgBlack, color.BgGreen).SprintFunc()

	//--------------------
	// VIEW

	fmt.Printf("GF_ERROR:\n")
	fmt.Printf("file           - %s\n", yellowBg(gfErr.File_str))
	fmt.Printf("line_num       - %s\n", yellowBg(gfErr.Line_num_int))
	fmt.Printf("user_msg       - %s\n", yellowBg(gfErr.User_msg_str))
	fmt.Printf("id             - %s\n", yellow(gfErr.Id_str))
	fmt.Printf("type           - %s\n", yellow(gfErr.Type_str))
	fmt.Printf("service_name   - %s\n", yellow(gfErr.ServiceNameStr))
	fmt.Printf("subsystem_name - %s\n", yellow(gfErr.Subsystem_name_str))
	fmt.Printf("function_name  - %s\n", yellow(gfErr.Function_name_str))
	fmt.Printf("data           - %s\n", yellow(gfErr.Data_map))
	fmt.Printf("error          - %s\n", red(pError))
	fmt.Printf("%s:\n%s\n", cyan("STACK TRACE"), green(gfErr.Stack_trace_str))

	fmt.Printf("\n\n")

	var namesPrefixStr string
	if pRuntimeSys.NamesPrefixStr != "" {
		namesPrefixStr = pRuntimeSys.NamesPrefixStr
	} else {
		namesPrefixStr = "gf"
	}

	//--------------------
	// DB_PERSIST
	if pRuntimeSys.Errors_send_to_mongodb_bool {
		
		ctx := context.Background()
		errsDBcollNameStr := fmt.Sprintf("%s_errors", namesPrefixStr)

		_, err := pRuntimeSys.Mongo_db.Collection(errsDBcollNameStr).InsertOne(ctx, gfErr)
		if err != nil {

		}
	}
	
	//--------------------
	// SENTRY
	if pRuntimeSys.ErrorsSendToSentryBool {
		
		/*
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetExtra("gf_error.service_name",   gf_error.ServiceNameStr)
			scope.SetExtra("gf_error.subsystem_name", gf_error.Subsystem_name_str)
			scope.SetExtra("gf_error.type",           gf_error.Type_str)
		})
		*/

		sentry.WithScope(func(pScope *sentry.Scope) {

			pScope.SetTag(fmt.Sprintf("%s_error.service_name",   namesPrefixStr), gfErr.ServiceNameStr)
			pScope.SetTag(fmt.Sprintf("%s_error.subsystem_name", namesPrefixStr), gfErr.Subsystem_name_str)
			pScope.SetTag(fmt.Sprintf("%s_error.type",           namesPrefixStr), gfErr.Type_str)

			for k, v := range gfErr.Data_map {
				pScope.SetTag(fmt.Sprintf("%s_error.%s", namesPrefixStr, k),
					fmt.Sprint(v))
			}

			// scope.SetLevel(sentry.LevelWarning);

			if pError != nil {
				sentry.CaptureException(pError)
			} else {

				// IMPORTANT!! - in case the GF_error doesnt have a correspoting
				//               golang error. this is for GF error conditions that are 
				//               not caused by a golang error.
				err := errors.New(fmt.Sprintf("%s error - %s", strings.ToUpper(namesPrefixStr), gfErr.Type_str))
				sentry.CaptureException(err)
			}
		})

		// FLUSH
		defer sentry.Flush(2 * time.Second)
	}
	
	//--------------------
	// METRICS - prometheus metrics
	if pRuntimeSys.Metrics != nil {
		pRuntimeSys.Metrics.ErrorsCounter.Inc()
	}

	//--------------------

	return &gfErr
}

//-------------------------------------------------

func ErrorInitSentry(pSentryEndpointStr string,
	pTransactionsTracingRateMap map[string]float64,
	pSampleRateDefaultF         float64) error {

	fmt.Println("INIT SENTRY")

	err := sentry.Init(sentry.ClientOptions{
		Dsn: pSentryEndpointStr,

		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,

		//--------------------
		// TRACING

		EnableTracing: true,
		TracesSampler: sentry.TracesSampler(func(pSamplingCtx sentry.SamplingContext) float64 {
			
			transactionNameStr := pSamplingCtx.Span.Name
			if sampleRateF, ok := pTransactionsTracingRateMap[transactionNameStr]; ok {
				return sampleRateF
			}

			// if no specific sampling rate for this transaction was defined, use the default rate
			return pSampleRateDefaultF
		}),
		// TracesSampleRate: pSampleRateDefaultF, // 1.0,

		//--------------------
	})
	if err != nil {
		return err
	}
	
	return nil
}

//-------------------------------------------------

func AttachMapToSentryScope(pMap map[string]interface{},
	pScope *sentry.Scope) {

	for key, value := range pMap {
        var stringValue string
        switch v := value.(type) {
        case string:
            stringValue = v
        case int:
            stringValue = fmt.Sprintf("%d", v)
        case float64:
            stringValue = fmt.Sprintf("%f", v)
        case bool:
            stringValue = fmt.Sprintf("%t", v)
        case []byte:
            stringValue = string(v)
        default:
            stringValue = fmt.Sprintf("%v", v)
        }
        pScope.SetTag(key, stringValue)
    }
}