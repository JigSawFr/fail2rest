package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Sean-Der/fail2go"
	"github.com/go-chi/chi"
)

func jailGetHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	currentlyFailed, totalFailed, fileList, currentlyBanned, totalBanned, IPList, err := fail2goConn.JailStatus(chi.URLParam(req, "jail"))
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	failRegexes, _ := fail2goConn.JailFailRegex(chi.URLParam(req, "jail"))
	findTime, _ := fail2goConn.JailFindTime(chi.URLParam(req, "jail"))
	useDNS, _ := fail2goConn.JailUseDNS(chi.URLParam(req, "jail"))
	maxRetry, _ := fail2goConn.JailMaxRetry(chi.URLParam(req, "jail"))
	actions, _ := fail2goConn.JailActions(chi.URLParam(req, "jail"))

	if IPList == nil {
		IPList = []string{}
	}
	if failRegexes == nil {
		failRegexes = []string{}
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{
		"currentlyFailed": currentlyFailed,
		"totalFailed":     totalFailed,
		"fileList":        fileList,
		"currentlyBanned": currentlyBanned,
		"totalBanned":     totalBanned,
		"IPList":          IPList,
		"failRegexes":     failRegexes,
		"findTime":        findTime,
		"useDNS":          useDNS,
		"maxRetry":        maxRetry,
		"actions":         actions})
	res.Write(encodedOutput)
}

type jailBanIPBody struct {
	IP string
}

func jailBanIPHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailBanIPBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailBanIP(chi.URLParam(req, "jail"), input.IP)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"bannedIP": output})
	res.Write(encodedOutput)
}

func jailUnbanIPHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailBanIPBody
	json.NewDecoder(req.Body).Decode(&input)
	output, err := fail2goConn.JailUnbanIP(chi.URLParam(req, "jail"), input.IP)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"unBannedIP": output})
	res.Write(encodedOutput)
}

type jailFailRegexBody struct {
	FailRegex string
}

func jailAddFailRegexHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailFailRegexBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailAddFailRegex(chi.URLParam(req, "jail"), input.FailRegex)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"FailRegex": output})
	res.Write(encodedOutput)
}

func jailDeleteFailRegexHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailFailRegexBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailDeleteFailRegex(chi.URLParam(req, "jail"), input.FailRegex)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"FailRegex": output})
	res.Write(encodedOutput)
}

type RegexResult struct {
	Line  string
	Match bool
}

func jailTestFailRegexHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailFailRegexBody
	json.NewDecoder(req.Body).Decode(&input)

	regexp, err := regexp.Compile(strings.Replace(input.FailRegex, "<HOST>", "(?:::f{4,6}:)?(?P<host>\\S+)", -1))

	if err != nil {
		writeHTTPError(res, err)
		return
	}

	_, _, fileList, _, _, _, err := fail2goConn.JailStatus(chi.URLParam(req, "jail"))
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	output := make(map[string][]RegexResult)
	for _, fileName := range fileList {
		file, err := os.Open(fileName)
		if err != nil {
			writeHTTPError(res, err)
			return
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			output[fileName] = append(output[fileName], RegexResult{Match: regexp.MatchString(scanner.Text()), Line: scanner.Text()})
		}
	}

	encodedOutput, _ := json.Marshal(output)
	res.Write(encodedOutput)
}

type jailFindTimeBody struct {
	FindTime int
}

func jailSetFindTimeHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailFindTimeBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailSetFindTime(chi.URLParam(req, "jail"), input.FindTime)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"FindTime": output})
	res.Write(encodedOutput)
}

type jailUseDNSBody struct {
	UseDNS string
}

func jailSetUseDNSHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailUseDNSBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailSetUseDNS(chi.URLParam(req, "jail"), input.UseDNS)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"useDNS": output})
	res.Write(encodedOutput)
}

type jailMaxRetryBody struct {
	MaxRetry int
}

func jailSetMaxRetryHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	var input jailMaxRetryBody
	json.NewDecoder(req.Body).Decode(&input)

	output, err := fail2goConn.JailSetMaxRetry(chi.URLParam(req, "jail"), input.MaxRetry)
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{"maxRetry": output})
	res.Write(encodedOutput)
}

func jailActionHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	port, err := fail2goConn.JailActionProperty(chi.URLParam(req, "jail"), chi.URLParam(req, "action"), "port")
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(map[string]interface{}{
		"port": port})
	res.Write(encodedOutput)
}

func jailHandler(r *chi.Mux, fail2goConn *fail2go.Conn) {

	r.Route("/jail", func(r chi.Router) {
		r.Post("/{jail}/bannedip", func(w http.ResponseWriter, r *http.Request) {
			jailBanIPHandler(w, r, fail2goConn)
		})
		r.Delete("/{jail}/bannedip", func(w http.ResponseWriter, r *http.Request) {
			jailUnbanIPHandler(w, r, fail2goConn)
		})

		r.Post("/{jail}/failregex", func(w http.ResponseWriter, r *http.Request) {
			jailAddFailRegexHandler(w, r, fail2goConn)
		})
		r.Delete("/{jail}/failregex", func(w http.ResponseWriter, r *http.Request) {
			jailDeleteFailRegexHandler(w, r, fail2goConn)
		})

		r.Post("/{jail}/testfailregex", func(w http.ResponseWriter, r *http.Request) {
			jailTestFailRegexHandler(w, r, fail2goConn)
		})

		r.Post("/{jail}/findtime", func(w http.ResponseWriter, r *http.Request) {
			jailSetFindTimeHandler(w, r, fail2goConn)
		})

		r.Post("/{jail}/usedns", func(w http.ResponseWriter, r *http.Request) {
			jailSetUseDNSHandler(w, r, fail2goConn)
		})

		r.Post("/{jail}/maxretry", func(w http.ResponseWriter, r *http.Request) {
			jailSetMaxRetryHandler(w, r, fail2goConn)
		})

		r.Get("/{jail}/action/{action}", func(w http.ResponseWriter, r *http.Request) {
			jailActionHandler(w, r, fail2goConn)
		})

		r.Get("/{jail}", func(w http.ResponseWriter, r *http.Request) {
			jailGetHandler(w, r, fail2goConn)
		})
	})
}
