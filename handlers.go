package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type rowK8sevents struct {
	ExportedTime time.Time
	EventTime    time.Time
	FirstTime    time.Time
	Name         string
	Reason       string
	Type         string
	Message      string
}

func (s *appServer) getMinDate() (time.Time, error) {
	var dbegin time.Time
	rqt := "select min(exportedTime) from k8sevents;"
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

func (s *appServer) nbResult(minDateTime time.Time, maxDateTime time.Time, searchName string, typeEvent string, reason string, message string, page int) (int, error) {
	var cnt int
	rows, err := s.makeRqtEvents(true, minDateTime, maxDateTime, searchName, typeEvent, reason, message, page)
	if err != nil {
		return cnt, err
	} else {
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&cnt)
			if err != nil {
				return 0, err
			}
		}
	}
	return cnt, err
}

func transformToCount(rqt string, limitClause string) string {
	rqt = strings.ReplaceAll(rqt, limitClause, "")
	rqt = strings.ReplaceAll(rqt, "exportedTime,firstTime,eventTime,name,reason,type,message", "count(*)")
	rqt = strings.ReplaceAll(rqt, "order by exportedTime desc", "")
	return rqt
}

func (s *appServer) makeRqtEvents(countRequest bool, minDateTime time.Time, maxDateTime time.Time, searchName string, typeEvent string, reason string, message string, page int) (*sql.Rows, error) {
	limitClause := fmt.Sprintf("limit %d offset %d", 50, (page)*50)

	switch {
	case typeEvent != "" && reason != "" && searchName != "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where type=$1 and reason=$2 and name like $3 and exportedTime between $4 and $5 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, typeEvent, reason, searchName, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent != "" && reason != "" && searchName == "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where type=$1 and reason=$2 and exportedTime between $3 and $4 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, typeEvent, reason, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent != "" && reason == "" && searchName != "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where type=$1 and name like $2 and exportedTime between $3 and $4 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, typeEvent, searchName, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent == "" && reason != "" && searchName != "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where reason=$1 and name like $2 and exportedTime between $3 and $4 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, reason, searchName, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent == "" && reason == "" && searchName != "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where name like $1 and exportedTime between $2 and $3 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, searchName, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent == "" && reason != "" && searchName == "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where reason=$1 and exportedTime between $2 and $3 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, reason, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent != "" && reason == "" && searchName == "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where type=$1 and exportedTime between $2 and $3 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		return s.db.Query(rqt, typeEvent, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	case typeEvent == "" && reason == "" && searchName == "":
		rqt := "select exportedTime,firstTime,eventTime,name,reason,type,message from k8sevents where exportedTime between $1 and $2 order by exportedTime desc " + limitClause
		if countRequest {
			rqt = transformToCount(rqt, limitClause)
		}
		// fmt.Println(rqt)
		return s.db.Query(rqt, minDateTime.Format("2006-01-02 15:04"), maxDateTime.Format("2006-01-02 15:04"))
	}

	return nil, errors.New("unknown case. Please create an issue on the github project")
}

func (s *appServer) calcPages(nbResults int) []string {
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

func (s *appServer) getDistinct(column string) (reasons []string, err error) {
	rqt := "select distinct(" + column + ") from k8sevents order by 1"
	rows, err := s.db.Query(rqt)
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var reason string
			err = rows.Scan(&reason)
			if err != nil {
				return nil, err
			}
			reasons = append(reasons, reason)
		}
	}
	return reasons, err
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
	data.Dmin, _ = s.getMinDate()
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

	// Note the call to ParseFS instead of Parse
	tmplt, err := template.ParseFS(htmlFiles, "templates/index.html")
	if err != nil {
		panic(err)
	}
	// tmplt := template.New("index.html")
	// tmplt, _ = tmplt.ParseFiles("./templates/index.html")

	// rqt, rqtCnt := s.makeRqtEvents(data.Dbegin, data.Dend, request.FormValue("search"), request.FormValue("type"), request.FormValue("reason"), request.FormValue("message"), data.Page)

	nbResults, err := s.nbResult(data.Dbegin, data.Dend, request.FormValue("search"), request.FormValue("type"), request.FormValue("reason"), request.FormValue("message"), data.Page)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	data.NbResults = nbResults
	data.Pages = s.calcPages(nbResults)

	rows, err := s.makeRqtEvents(false, data.Dbegin, data.Dend, request.FormValue("search"), request.FormValue("type"), request.FormValue("reason"), request.FormValue("message"), data.Page)
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var rowRes rowK8sevents
			err = rows.Scan(&rowRes.ExportedTime, &rowRes.FirstTime, &rowRes.EventTime, &rowRes.Name, &rowRes.Reason, &rowRes.Type, &rowRes.Message)
			if err != nil {
				panic(err)
			}
			data.Rows = append(data.Rows, rowRes)
		}
	}

	data.Reasons, err = s.getDistinct("reason")
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
	}
	data.Types, err = s.getDistinct("type")
	if err != nil {
		var d dataErr
		d.ErrorMsg = err.Error()
		s.HandlerError(response, d)
	}

	err = tmplt.Execute(response, data)
	if err != nil {
		fmt.Printf("Error when generating template index: %s\n", err.Error())
	}
}
