package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type rowK8sevents struct {
	K8seventTs      time.Time
	K8seventFirstTs time.Time
	K8seventName    string
	K8seventReason  string
	K8seventType    string
	K8seventMessage string
}

func (s *appServer) GetMinDate() (time.Time, error) {
	var dbegin time.Time
	rqt := "select min(firstEventTs) from k8sevents;"
	rows, err := s.db.Query(rqt)
	if err != nil {
		return dbegin, err
	} else {
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&dbegin)
			if err != nil {
				dbegin = time.Now()
			}
		}
	}
	return dbegin, err
}

// func (s *appServer) GetMaxDate() (time.Time, error) {
// 	var dend time.Time
// 	rqt := "select max(date) from k8sevents;"
// 	rows, err := s.db.Query(rqt)
// 	if err != nil {
// 		return dend, err
// 	} else {
// 		defer rows.Close()
// 		for rows.Next() {
// 			err = rows.Scan(&dend)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// 	return dend, err
// }

func (s *appServer) nbResult(rqt string) (int, error) {
	var cnt int
	rows, err := s.db.Query(rqt)
	if err != nil {
		return cnt, err
	} else {
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&cnt)
			if err != nil {
				panic(err)
			}
		}
	}
	return cnt, err
}

func (s *appServer) makeRqtEvents(minDateTime time.Time, maxDateTime time.Time, searchName string, typeEvent string, reason string, message string, page int) (string, string) {
	var limitClause string
	whereClause := "where firstEventTs between '" + minDateTime.Format("2006-01-02 15:04") + "' and '" + maxDateTime.Format("2006-01-02 15:04") + "'"
	if typeEvent != "" {
		whereClause = whereClause + "and type='" + typeEvent + "'"
	}
	if reason != "" {
		whereClause = whereClause + "and reason='" + reason + "'"
	}
	if searchName != "" {
		whereClause = whereClause + "and name like '%" + searchName + "%'"
	}

	limitClause = fmt.Sprintf("limit %d offset %d", 50, (page)*50)
	rqt := "select firstEventTs,eventTs,name,reason,type,message from k8sevents " + whereClause + " order by EventTs,firstEventTs desc " + limitClause
	rqtCnt := "select count(*) from k8sevents " + whereClause
	return rqt, rqtCnt
}

func (s *appServer) calcPages(nbResults int) []string {
	// fmt.Println("nbresults=", nbResults)
	// fmt.Println("nbresults/50=", nbResults/50)
	nbPages := nbResults / 50

	if nbResults%50 != 0 {
		nbPages++
	}

	var res []string
	for i := 0; i < nbPages-1; i++ {
		res = append(res, strconv.Itoa(i))
	}
	return res

}

func (s *appServer) IndexHandler(response http.ResponseWriter, request *http.Request) {
	type dataIndex struct {
		Rows           []rowK8sevents
		Reasons        []string
		Types          []string
		Dbegin         time.Time
		Dend           time.Time
		Dmin           time.Time
		Dmax           time.Time
		TypeSelected   string
		ReasonSelected string
		Search         string
		NbResults      int
		Pages          []string
		Page           int
	}
	var data dataIndex
	var err error

	// fmt.Println("REASON=", request.FormValue("reason"))
	// fmt.Println("TYPE=", request.FormValue("type"))
	// fmt.Println("DEND=", request.FormValue("dend"))
	// fmt.Println("DBEGIN=", request.FormValue("dbegin"))

	data.Dmin, _ = s.GetMinDate()
	data.Dmax = time.Now()

	if request.FormValue("type") != "" {
		data.TypeSelected = request.FormValue("type")
	}
	if request.FormValue("reason") != "" {
		data.ReasonSelected = request.FormValue("reason")
	}
	if request.FormValue("page") == "" {
		data.Page = 0
	} else {
		data.Page, _ = strconv.Atoi(request.FormValue("page"))
	}
	data.Search = request.FormValue("search")

	if request.FormValue("dbegin") == "" {
		data.Dbegin = data.Dmin
	} else {
		data.Dbegin, err = time.Parse("2006-01-02T15:04", request.FormValue("dbegin"))
		if err != nil {
			fmt.Println(err.Error())
			data.Dbegin = data.Dmin
		}
	}

	if request.FormValue("dend") == "" {
		data.Dend = data.Dmax
	} else {
		data.Dend, err = time.Parse("2006-01-02T15:04", request.FormValue("dend"))
		if err != nil {
			fmt.Println(err.Error())
			data.Dend = data.Dmax
		}
	}

	tmplt := template.New("index.html")
	tmplt, _ = tmplt.ParseFiles("./templates/index.html")

	// rqt := "select date,name,reason,type from k8sevents order by date desc limit 50;"
	rqt, rqtCnt := s.makeRqtEvents(data.Dbegin, data.Dend, request.FormValue("search"), request.FormValue("type"), request.FormValue("reason"), request.FormValue("message"), data.Page)
	// fmt.Println(rqt)

	nbResults, _ := s.nbResult(rqtCnt)
	data.NbResults = nbResults
	data.Pages = s.calcPages(nbResults)

	rows, err := s.db.Query(rqt)
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
	} else {
		defer rows.Close()
		for rows.Next() {
			var rowRes rowK8sevents
			err = rows.Scan(&rowRes.K8seventFirstTs, &rowRes.K8seventTs, &rowRes.K8seventName, &rowRes.K8seventReason, &rowRes.K8seventType, &rowRes.K8seventMessage)
			if err != nil {
				panic(err)
			}
			data.Rows = append(data.Rows, rowRes)
			// fmt.Println(rowRes)
		}
	}
	// fmt.Println(len(data.Rows))
	rqt = "select distinct(reason) from k8sevents order by 1"
	rows, err = s.db.Query(rqt)
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
	} else {
		defer rows.Close()
		for rows.Next() {
			var reason string
			err = rows.Scan(&reason)
			if err != nil {
				panic(err)
			}
			data.Reasons = append(data.Reasons, reason)
			// fmt.Println(rowRes)
		}
	}

	rqt = "select distinct(type) from k8sevents order by 1"
	rows, err = s.db.Query(rqt)
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
	} else {
		defer rows.Close()
		for rows.Next() {
			var k8stype string
			err = rows.Scan(&k8stype)
			if err != nil {
				panic(err)
			}
			data.Types = append(data.Types, k8stype)
			// fmt.Println(rowRes)
		}
	}

	err = tmplt.Execute(response, data)
	if err != nil {
		fmt.Printf("Error when generating template index: %s\n", err.Error())
	}
}
