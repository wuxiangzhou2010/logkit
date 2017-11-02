// +build linux

package system

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/qiniu/logkit/metric"
	"github.com/qiniu/logkit/utils"
)

const (
	TypeMetricKernelVmstat   = "kernel_vmstat"
	MetricKernelVmstatUsages = "内核(kernel_vmstat)"
)

type KernelVmstat struct {
	statFile string
}

func (k *KernelVmstat) Name() string {
	return TypeMetricKernelVmstat
}

func (k *KernelVmstat) Usages() string {
	return MetricKernelVmstatUsages
}

func (k *KernelVmstat) Config() []utils.Option {
	return []utils.Option{}
}

func (k *KernelVmstat) Attributes() []utils.KeyValue {
	return []utils.KeyValue{}
}

func (k *KernelVmstat) Collect() (datas []map[string]interface{}, err error) {
	data, err := k.getProcVmstat()
	if err != nil {
		return nil, err
	}

	fields := make(map[string]interface{})

	dataFields := bytes.Fields(data)
	for i, field := range dataFields {

		// dataFields is an array of {"stat1_name", "stat1_value", "stat2_name",
		// "stat2_value", ...}
		// We only want the even number index as that contain the stat name.
		if i%2 == 0 {
			// Convert the stat value into an integer.
			m, err := strconv.Atoi(string(dataFields[i+1]))
			if err != nil {
				return nil, err
			}

			fields[string(field)] = int64(m)
		}
	}
	datas = append(datas, fields)
	return
}

func (k *KernelVmstat) getProcVmstat() ([]byte, error) {
	if _, err := os.Stat(k.statFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("kernel_vmstat: %s does not exist!", k.statFile)
	} else if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(k.statFile)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func init() {
	metric.Add(TypeMetricKernelVmstat, func() metric.Collector {
		return &KernelVmstat{
			statFile: "/proc/vmstat",
		}
	})
}
