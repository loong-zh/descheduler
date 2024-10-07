package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	dumpDir = "/var"
)

// PodEvictedReport defines the structure of the report.
type PodEvictedReport struct {
	NodeName     string            `json:"nodeName" yaml:"nodeName"`
	NameSpace    string            `json:"namespace" yaml:"namespace"`
	PodName      string            `json:"podName" yaml:"podName"`
	Reason       string            `json:"reason" yaml:"reason"`
	StrategyName string            `json:"strategyName" yaml:"strategyName"`
	Labels       map[string]string `json:"labels" yaml:"labels"`
	EvictedTime  time.Time         `json:"evictedTime" yaml:"evictedTime"`
}

func (p *PodEvictedReport) Dump(s string) error {
	now := time.Now()
	timeStamp := now.Format("2006-01-02")
	file := fmt.Sprintf("descheduler_%s.%s", timeStamp, s)

	filePath := filepath.Join(dumpDir, file)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	switch s {
	case "json":
		err = p.DumpJson(f)
	case "yaml":
		err = p.DumpYaml(f)
	case "txt":
		err = p.DumpText(f)
	default:
		fmt.Printf("Unsupported format: %s\n", s)
	}

	return err
}

func ensureOutput(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func (p *PodEvictedReport) DumpText(f *os.File) error {
	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "%s podName=%s nameSpace=%s node=%s reason=%s strategName=%s labels=%s\n", p.EvictedTime, p.PodName, p.NameSpace, p.NodeName, p.Reason, p.StrategyName, p.Labels)
	w.Flush()
	return nil
}

func (p *PodEvictedReport) DumpYaml(f *os.File) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%v\n", string(data))
	return nil
}

func (p *PodEvictedReport) DumpJson(f *os.File) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%v\n", string(data))
	return nil
}
