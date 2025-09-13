package handlers

import (
	"net/http"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {

	//start := r.URL.Query().Get("start")
	//end := r.URL.Query().Get("end")
	//
	//if start == "" || end == "" {
	//	err := fmt.Errorf("StatsHandler: parameters start or end are empty")
	//	logger.Error("RateHandler error", err)
	//	http.Error(w, "Missing parameters: from and to required", http.StatusBadRequest)
	//	return
	//}
	//
	//data, err := os.ReadFile("./app.log")
	//
	//data, err := json.Marshal(rate)
	//if err != nil {
	//	err = fmt.Errorf("handler: failed to marshal response: %w", err)
	//	logger.Error("RateHandler error", err)
	//	http.Error(w, "Internal server error", http.StatusInternalServerError)
	//	return
	//}
	//w.Write(data)

}
