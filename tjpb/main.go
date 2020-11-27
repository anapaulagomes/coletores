package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/coletores"
)

var gitCommit string

func main() {
	outputFolder := os.Getenv("OUTPUT_FOLDER")
	if outputFolder == "" {
		outputFolder = "./output"
	}
	month, err := strconv.Atoi(os.Getenv("MONTH"))
	if err != nil {
		logError("Invalid month (\"%s\"): %q", os.Getenv("MONTH"), err)
		os.Exit(1)
	}
	year, err := strconv.Atoi(os.Getenv("YEAR"))
	if err != nil {
		logError("Invalid year (\"%s\"): %q", os.Getenv("YEAR"), err)
		os.Exit(1)
	}

	if err := os.Mkdir(outputFolder, os.ModePerm); err != nil && !os.IsExist(err) {
		logError("Error creating output folder(%s): %q", outputFolder, err)
		os.Exit(1)
	}

	files, err := crawl(outputFolder, month, year)
	if err != nil {
		logError("Error crawling (%d,%d,%s) error: %q", month, year, outputFolder, err)
		os.Exit(1)
	}

	//teste := []string{"transparencia_202001_servidores_0.pdf", "remuneracoes-magistrados-tjpb-01-2020.pdf"}
	files, allEmployees, err := genEmployees(files, outputFolder, month, year)
	if err != nil {
		logError("Error generating employees, error: %v", err)
		os.Exit(1)
	}

	//teste2 := "transparencia_202004_servidores_0_0.pdf"
	//teste3 := "remuneracoes-magistrados-tjpb-01-2020.pdf"
	er := coletores.ExecutionResult{Cr: newCrawlingResult(allEmployees, files, month, year)}
	b, err := json.MarshalIndent(er, "", "  ")
	if err != nil {
		logError("JSON marshaling error: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func newCrawlingResult(emps []coletores.Employee, files []string, month, year int) coletores.CrawlingResult {
	crawlerInfo := coletores.Crawler{
		CrawlerID:      "tjpb",
		CrawlerVersion: gitCommit,
	}
	cr := coletores.CrawlingResult{
		AgencyID:  "tjpb",
		Month:     month,
		Year:      year,
		Files:     files,
		Employees: emps,
		Crawler:   crawlerInfo,
		Timestamp: time.Now(),
	}
	return cr
}

// genEmployees navigate
func genEmployees(files []string, outputFolder string, month, year int) ([]string, []coletores.Employee, error) {
	var allEmployees []coletores.Employee
	var pathFixed []string
	template, err := checkTemplate(month, year)
	if err != nil {
		return nil, nil, fmt.Errorf("error trying to check in which template this month/year belongs %v", err)
	}
	for i, f := range files {
		pathFixed = append(pathFixed, fmt.Sprintf("%v/%v", outputFolder, filepath.Base(f)))
		switch {
		case strings.Contains(f, "magistrados") && template == "E":
			emps, err := parserMagE(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template D: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "servidores") && template == "E":
			emps, err := parserServerE(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing servant of template D: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "magistrados") && template == "D":
			emps, err := parserMagD(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template E: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "servidores") && template == "D":
			emps, err := parserServD(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing servants of template E: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "magistrados") && template == "C":
			emps, err := parserMagC(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template C: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "servidores") && template == "C":
			emps, err := parserServC(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template C: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "magistrados") && template == "B":
			emps, err := parserMagB(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template B: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "servidores") && template == "B":
			emps, err := parserServB(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template B: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		case strings.Contains(f, "magistrados") && template == "A":
			emps, err := parserMagA(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing magistrate of template A: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		default:
			emps, err := parserServA(pathFixed[i])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing servant of template A: %v", err)
			}
			allEmployees = append(allEmployees, emps...)
		}
		files = append(files, strings.Replace(f, ".pdf", ".csv", 1))
	}
	return files, allEmployees, nil
}
