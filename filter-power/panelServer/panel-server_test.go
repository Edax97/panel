package panelServer

import (
	"filter-power/wailonServer"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func GetMaps() (imeisFromIds map[string]string, variableFromSymbol map[string]string) {
	// (time)id_variable>value: should login at imeis[id]
	// should send value to var_names[variable]
	// should not send if time is not :00

	imeisFromIds = make(map[string]string)
	// When tests run from package directory, this relative path reaches repo map
	path := "../../map-imei/map-imei.csv"
	content, err := os.ReadFile(path)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, ln := range lines {
			if strings.TrimSpace(ln) == "" {
				continue
			}
			parts := strings.Split(ln, ",")
			if len(parts) < 3 {
				continue
			}
			id := strings.TrimSpace(parts[0])
			imei := strings.TrimSpace(parts[2])
			if imei != "." {
				if _, convErr := strconv.Atoi(imei); convErr == nil {
					if len(imei) < 15 {
						imei = strings.Repeat("0", 15-len(imei)) + imei
					}
				}
			}
			imeisFromIds[id] = imei
		}
	}

	variableFromSymbol = map[string]string{
		// energy (non-transformer and transformer share these)
		"WHr_I":   "watth",
		"VARHr_I": "varh",
		"VARHr_O": "varo",
		// transformer extra fields
		"IA":     "Ia",
		"IB":     "Ib",
		"IC":     "Ic",
		"VAB":    "Vab",
		"VBC":    "Vbc",
		"VCA":    "Vca",
		"VAN":    "Van",
		"VBN":    "Vbn",
		"VCN":    "Vcn",
		"VAHrIn": "varin",
		"PFTtl":  "pfttl",
	}
	return
	// transformer_id
	// should send aditional params
	// should send correct params
}

func ContainsRE(out string, want string) bool {
	reg, _ := regexp.Compile(want)
	return reg.MatchString(out)
}

type InputCase struct {
	file      string
	time      string
	id        string
	isTrans   bool
	variables []string
	values    []string
}

func (i InputCase) GenInput() [][]string {
	data := [][]string{
		{"sep="},
		{"Element Id"},
		{"Device Name"},
		{"Device Type"},
		{"Measurement Name"},
		{"Measurement Unit"},
		{i.time},
	}
	for j, v := range i.variables {
		data[1] = append(data[1], i.id+"_"+v)
		data[2] = append(data[2], "-")
		data[3] = append(data[3], "-")
		data[4] = append(data[4], "-")
		data[5] = append(data[5], "-")
		data[6] = append(data[6], i.values[j])
	}
	return data
}

func TestTransformers(t *testing.T) {
	inputList := []InputCase{
		{
			file:      "data_3.csv",
			time:      "2027/01/01 10:00:00",
			id:        "zigbee:47_zd",
			variables: []string{"WHr_I", "VARHr_I", "VARHr_O"},
			values:    []string{"100.0", "5.0", "1.0"},
		},
		{
			file:      "data_3.csv",
			time:      "2028/01/01 11:15:00",
			id:        "zigbee:14_zd",
			variables: []string{"WHr_I"},
			values:    []string{"200"},
		},
		{
			file:      "data_1.csv",
			time:      "2029/12/01 23:00:00",
			id:        "zigbee:17_zd",
			variables: []string{"WHr_I", "VARHr_I"},
			values:    []string{"100.0", "NaN"},
		},
		{
			file: "data_1.csv",
			time: "2030/01/01 10:00:00",
			id:   "modbus:1_mb",
			variables: []string{"WHr_I", "VARHr_I", "VARHr_O",
				"IA", "IB", "IC",
				"VAB", "VBC", "VCA",
				"VAN", "VBN", "VCN",
				"VAHrIn", "PFTtl",
			},
			values: []string{"1000120", "500021", "123",
				"10", "2", "12",
				"330", "300", "320",
				"220", "210", "240",
				"1500", "0.98",
			},
			isTrans: true,
		},
	}

	imeiMap, variableMap := GetMaps()

	ser := wailonServer.NewMockServer()
	panelServer, err := NewPanelServer()
	if err != nil {
		t.Fatal("Could not init Panel Server object")
	}

	for _, i := range inputList {
		data := i.GenInput()
		ser.FlushOut()
		err := panelServer.SendPanelServer(data, i.file, ser)
		if err != nil {
			t.Fatalf("Error when parsing data: %v", err)
		}
		out := ser.OutBuffer

		t.Logf("Test id: %s\n {%s}", i.id, ser.OutBuffer)
		// check if 00
		if !ContainsRE(i.time, "[0-9]:00:00") {
			if len(out) != 0 {
				t.Errorf("should be empty at %s\ngot:%s", i.time, out)
			}
			continue
		}

		// login
		imei := imeiMap[i.file+"_"+i.id+"_WHr_I"]
		if !strings.Contains(out, imei) {
			t.Errorf("want: %s,\ngot: %s", imei, out)
		}

		// var:3:value
		for j, v := range i.variables {
			field := fmt.Sprintf("%s:[0-9]:%s", variableMap[v], i.values[j])
			if !ContainsRE(out, field) {
				t.Errorf("want: %s\ngot: %s", field, out)
			}
		}

		if i.isTrans {
			wantFields := []string{
				"Iaprim",
				"Ibprim",
				"Icprim",
				"Vabprim",
				"Vbcprim",
				"Vcaprim",
				"Pprim",
				"Qprim",
				"Sprim",
			}
			for _, w := range wantFields {
				field := fmt.Sprintf("%s:[0-9]", w)
				notField := fmt.Sprintf("%s:[0-9]:NaN", w)
				if !ContainsRE(out, field) {
					t.Errorf("want: %s\ngot: %s", field, out)
				}
				if ContainsRE(out, notField) {
					t.Errorf("dont want: %s", notField)
				}
			}
		}

	}

}
