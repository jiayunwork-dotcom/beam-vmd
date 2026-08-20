// Command beam-vmd is the entry point of the Euler-Bernoulli beam solver. It
// exposes the solver behind a small CLI and an HTTP server, and serves the
// interactive web console from the embedded web/ directory.
//
// CLI:
//
//	beam-vmd -solve example/simply.json        solve a case file and print
//	beam-vmd -solve example/simply.json -table print a text table instead
//	beam-vmd -http :8080                      serve the web console
//
// HTTP:
//
//	POST /api/solve   {L, EI, support, points, distrib} -> reactions, V/M/y samples, extrema, checks
//	GET  /api/example -> the bundled simply-supported example
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"beam-vmd/internal/beam"
	"beam-vmd/internal/model"
)

//go:embed web
var webFS embed.FS

//go:embed example/simply.json
var simplyJSON []byte

func main() {
	httpAddr := flag.String("http", "", "serve the web console on this address (e.g. :8080)")
	solveFile := flag.String("solve", "", "solve a case JSON file and print the result")
	tableOut := flag.Bool("table", false, "print a text table instead of JSON for -solve")
	flag.Parse()

	if *httpAddr != "" {
		if err := serveHTTP(*httpAddr); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *solveFile != "" {
		if err := runCLI(*solveFile, *tableOut); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		return
	}
	flag.Usage()
	os.Exit(2)
}

// runCLI loads a case file, solves it and prints the result as JSON, or as a
// human-readable table when table is set.
func runCLI(path string, table bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	b, err := model.ParseBeam(data)
	if err != nil {
		return err
	}
	res, err := beam.Solve(b)
	if err != nil {
		return err
	}
	if table {
		fmt.Print(beam.Table(res))
		return nil
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// apiError is the JSON body returned for illegal or non-converged models.
type apiError struct {
	Error string `json:"error"`
}

// apiSolveResponse is the JSON payload returned by POST /api/solve.
type apiSolveResponse struct {
	L         float64           `json:"L"`
	EI        float64           `json:"EI"`
	Support   string            `json:"support"`
	Reactions []model.Reaction  `json:"reactions"`
	Samples   []model.Sample    `json:"samples"`
	Extrema   model.Extrema     `json:"extrema"`
	Checks    model.CheckReport `json:"checks"`
	Error     string            `json:"error,omitempty"`
}

// handleSolve serves POST /api/solve: it parses a beam document, solves it and
// returns the reactions, internal-force/deflection samples, extrema and the
// cross-validation report. Illegal input produces a JSON error body.
func handleSolve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiError{Error: "method not allowed"})
		return
	}
	body, err := readBody(r, 1<<20)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}
	b, err := model.ParseBeam(body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}
	res, err := beam.Solve(b)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, apiSolveResponse{
		L:         res.L,
		EI:        res.EI,
		Support:   res.Support,
		Reactions: res.Reactions,
		Samples:   res.Samples,
		Extrema:   res.Extrema,
		Checks:    res.Checks,
	})
}

// handleExample serves GET /api/example returning the bundled simply-supported
// case so the page can load a working input with one click.
func handleExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, http.StatusOK, json.RawMessage(simplyJSON))
}

// serveHTTP wires the routes and serves the embedded web assets.
func serveHTTP(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/solve", handleSolve)
	mux.HandleFunc("/api/example", handleExample)
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))
	log.Printf("beam-vmd web console on http://localhost%s", addr)
	return http.ListenAndServe(addr, mux)
}

// readBody reads at most max bytes from the request body.
func readBody(r *http.Request, max int64) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("empty request body")
	}
	defer r.Body.Close()
	if r.ContentLength > max {
		return nil, fmt.Errorf("request body too large")
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, max+1))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if int64(len(data)) > max {
		return nil, fmt.Errorf("request body too large")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty request body")
	}
	return data, nil
}

// writeJSON encodes v as JSON into the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
