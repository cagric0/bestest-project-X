package main

//func (a *App) home(w http.ResponseWriter, r *http.Request) {
//	ctx := context.Background()
//	tests, err := a.Hz.GetTestList(ctx)
//	if err != nil {
//		log.Print("HZ GetTestList: ", err) // log it
//		return
//	}
//	a.createPage(w, "homepage", tests)
//}
//
//func (a *App) testRunIDs(w http.ResponseWriter, r *http.Request) {
//	ctx := context.Background()
//	testName, _ := mux.Vars(r)["test-name"]
//	testRunIDs, err := a.Hz.GetTestRunIDs(ctx, testName)
//	if err != nil {
//		log.Print("HZ GetTestRunIDs: ", err) // log it
//		return
//	}
//	pageData := struct {
//		TestName string
//		RunIDs   interface{}
//	}{
//		TestName: testName,
//		RunIDs:   testRunIDs,
//	}
//	a.createPage(w, "testrun", pageData)
//}
//
//func (a *App) testLogs(w http.ResponseWriter, r *http.Request) {
//	ctx := context.Background()
//	logIdentifier, _ := mux.Vars(r)["log-identifier"]
//
//	logs, err := a.Hz.GetLogs(ctx, logIdentifier)
//	if err != nil {
//		log.Print("HZ GetLogs: ", err) // log it
//		return
//	}
//
//	logMap := logs.(map[string]string)
//	logNames := make([]string, 0, len(logMap))
//	for k := range logMap {
//		logNames = append(logNames, k)
//	}
//
//	pageData := struct {
//		LogIdentifier string
//		LogNames      interface{}
//	}{
//		LogIdentifier: logIdentifier,
//		LogNames:      logNames,
//	}
//	a.createPage(w, "logs", pageData)
//}
//
//func (a *App) testLogDetail(w http.ResponseWriter, r *http.Request) {
//	ctx := context.Background()
//	logIdentifier, _ := mux.Vars(r)["log-identifier"]
//	logName, _ := mux.Vars(r)["log-name"]
//
//	logs, err := a.Hz.GetLogs(ctx, logIdentifier)
//	if err != nil {
//		log.Print("HZ GetLogs: ", err) // log it
//		return
//	}
//
//	logMap := logs.(map[string]string)
//	logDetail := logMap[logName]
//
//	pageData := struct {
//		LogIdentifier string
//		LogName       string
//		LogDetail     string
//	}{
//		LogIdentifier: logIdentifier,
//		LogName:       logName,
//		LogDetail:     logDetail,
//	}
//	a.createPage(w, "logDetail", pageData)
//}
//
//func (a *App) createPage(w http.ResponseWriter, page string, data interface{}) {
//	t, err := template.ParseFiles("template/" + page + ".html") //parse the html file homepage.html
//	if err != nil {                                             // if there is an error
//		log.Print("template parsing error: ", err) // log it
//	}
//	err = t.Execute(w, data) //execute the template and pass it the HomePageVars struct to fill in the gaps
//	if err != nil {          // if there is an error
//		log.Print("template executing error: ", err) //log it
//	}
//}
