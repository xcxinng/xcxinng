package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func genRegexpByRange(start, end int) string {
	switch {
	case 0 <= start && 9 >= start:
		switch {
		case 0 < end && 9 >= end: // 1-9
			return fmt.Sprintf("[%d-%d]", start, end)

		case end >= 10 && end <= 19: // 1-19
			return fmt.Sprintf("[%d-9]|1[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 20 && end <= 29: // 1-29
			return fmt.Sprintf("[%d-9]|1[0-9]|2[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 30 && end <= 39: // 1-39
			return fmt.Sprintf("[%d-9]|[1-2][0-9]|3[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 40 && end <= 49: // 1-49
			return fmt.Sprintf("[%d-9]|[1-3][0-9]|4[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 50 && end <= 59: // 1-59
			return fmt.Sprintf("[%d-9]|[1-4][0-9]|5[0-%d]", GetOnesPlace(start), GetOnesPlace(end))
		}

	case start >= 10 && start <= 19:
		switch {
		case end >= 10 && end <= 19: // 10-19
			return fmt.Sprintf("1[%d-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 20 && end <= 29: // 10-29
			return fmt.Sprintf("1[%d-9]|2[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 30 && end <= 39: // 10-39
			return fmt.Sprintf("1[%d-9]|2[0-9]|3[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 40 && end <= 49: // 10-49
			return fmt.Sprintf("1[%d-9]|[2-3][0-9]|4[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 50 && end <= 59: // 10-59
			return fmt.Sprintf("1[%d-9]|[2-4][0-9]|5[0-%d]", GetOnesPlace(start), GetOnesPlace(end))
		}

	case start >= 20 && start <= 29:
		switch {
		case end >= 20 && end <= 29: // 20-29
			return fmt.Sprintf("2[%d-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 30 && end <= 39: // 20-39
			return fmt.Sprintf("2[%d-9]|3[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 40 && end <= 49: // 20-49
			return fmt.Sprintf("2[%d-9]|3[0-9]|4[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 50 && end <= 59: // 20-59
			return fmt.Sprintf("2[%d-9]|[3-4][0-9]|5[0-%d]", GetOnesPlace(start), GetOnesPlace(end))
		}

	case start >= 30 && start <= 39:
		switch {
		case end >= 30 && end <= 39: // 30-39
			return fmt.Sprintf("3[%d-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 40 && end <= 49: // 30-49
			return fmt.Sprintf("3[%d-9]|4[0-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 50 && end <= 59: // 30-59
			return fmt.Sprintf("3[%d-9]|4[0-9]|5[0-%d]", GetOnesPlace(start), GetOnesPlace(end))
		}

	case start >= 40 && start <= 49:
		switch {
		case end >= 40 && end <= 49: // 30-49
			return fmt.Sprintf("4[%d-%d]", GetOnesPlace(start), GetOnesPlace(end))

		case end >= 50 && end <= 59: // 30-59
			return fmt.Sprintf("4[%d-9]|5[0-%d]", GetOnesPlace(start), GetOnesPlace(end))
		}

	case start >= 50 && start <= 59:
		return fmt.Sprintf("5[%d-%d]", GetOnesPlace(start), GetOnesPlace(end))
	}
	return "端口超出范围"
}

// GetOnesPlace 函数接受一个整数作为参数，并返回其个位数
func GetOnesPlace(number int) int {
	// 取余操作可以得到一个整数的最低位数字
	onesPlace := number % 10
	if onesPlace < 0 {
		// 处理负数的情况，将其转换为正数
		onesPlace = -onesPlace
	}
	return onesPlace
}

// ReplaceStringsWithSameNumber 函数接受一个包含类似[1-1],[2-2]的字符串，并替换为[1],[2]等
func ReplaceStringsWithSameNumber(input string) string {
	// 定义正则表达式，匹配 [number-number]
	re := regexp.MustCompile(`\[(\d+)-(\d+)\]`)

	// 使用正则表达式替换匹配的部分
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// 提取匹配的数字部分
		matches := re.FindStringSubmatch(match)
		if len(matches) == 3 && matches[1] == matches[2] {
			// 如果数字相同，则只保留一个数字
			return matches[1]
		}
		// 否则保留原始匹配
		return match
	})

	return result
}

func FinalGenerateRegexp(start, end int) (string, error) {
	if start >= end {
		return "", errors.New("开始不能大于等于结束")
	}
	if start <= 0 || end <= 0 {
		return "", errors.New("开始或结束不能小于或等于0")
	}
	if start > 58 || end > 58 {
		return "", errors.New("开始或结束必须处于[1,58]区间范围")
	}
	re := genRegexpByRange(start, end)
	return ReplaceStringsWithSameNumber(re), nil
}

func GeneratePortRegexp(ratePrefix string,
	stacking bool,
	portSplitting bool,
	slotNumber int,
	ranges string) (string, error) {

	slotString := "1"
	portSplittingSuffix := ":[1-4]"
	portRegexp := ""
	var portMultiRegexps []string

	if stacking {
		slotString = "([1-2])"
	}

	err := PortRanges(ranges).ForEachRange(func(start, end int) error {
		portRangeRegexp, err := FinalGenerateRegexp(int(start), int(end))
		if err != nil {
			return err
		}
		portMultiRegexps = append(portMultiRegexps, portRangeRegexp)
		return nil
	})
	if err != nil {
		return "", err
	}

	switch slotNumber {
	case 2:
		// 10GE1/rangeStr
		// stacking 10GE[1-2]/rangeStr
		portRegexp = fmt.Sprintf("%s%s/(%s)", ratePrefix, slotString, strings.Join(portMultiRegexps, "|"))

	case 3:
		// 10GE1/0/rangeStr
		// stacking 10GE[1-2]/0/rangeStr
		portRegexp = fmt.Sprintf("%s%s/0/(%s)", ratePrefix, slotString, strings.Join(portMultiRegexps, "|"))

	default:
		return "", fmt.Errorf("不支持的槽位数:%d", slotNumber)
	}

	portRegexp = "^" + portRegexp
	if portSplitting {
		portRegexp += portSplittingSuffix
	}
	portRegexp += "$"
	_, err = regexp.Compile(portRegexp)
	if err != nil {
		return "", fmt.Errorf("正则表达式:%s 编译报错: %v", portRegexp, err)
	}
	fmt.Println("regexp generated: ", portRegexp)
	return portRegexp, nil
}

type PortRanges string

func (p PortRanges) GetCount() (int, error) {
	total := 0
	err := p.ForEachRange(func(start, end int) error { total += (end - start + 1); return nil })
	return int(total), err
}

func (p PortRanges) ForEachRange(action func(start, end int) error) error {
	for _, pr := range strings.Split(string(p), ",") {
		// pr: 1-42
		prStr := strings.Split(pr, "-")
		if len(prStr) != 2 {
			return fmt.Errorf("端口范围格式有误: %s", pr)
		}
		var start, end int64
		var err error
		start, err = strconv.ParseInt(prStr[0], 10, 32)
		if err != nil {
			return err
		}

		end, err = strconv.ParseInt(prStr[1], 10, 32)
		if err != nil {
			return err
		}

		err = action(int(start), int(end))
		if err != nil {
			return err
		}
	}
	return nil
}

type PortQuotaInfo struct {
	RatePrefix string     `json:"rate_prefix"`
	NamingType int        `json:"slot_number"`
	Ranges     PortRanges `json:"ranges"`

	// These two options are only valid for DownlinkPort(ServerPort).
	Stacking            bool   `json:"stacking,omitempty"`
	StackingRanges      string `json:"stacking_ranges"`
	PortSplitting       bool   `json:"port_splitting,omitempty"`
	PortSplittingRanges string `json:"port_splitting_ranges"`
}

type ModelPortQuotaParameter struct {
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
	// TotalPort    int            `json:"total_port" binding:"required"`
	DownlinkPort *PortQuotaInfo `json:"downlink_port" binding:"required"`
	UplinkPort   *PortQuotaInfo `json:"uplink_port" binding:"required"`
	ReservedPort *PortQuotaInfo `json:"reserved_port" binding:"required"`
}

var ErrInValidParameter = errors.New("入参数有误")

func EmulatePostHandler(param ModelPortQuotaParameter) (resp PortQuotaDocument, err error) {
	if param.DownlinkPort == nil {
		err = ErrInValidParameter
		return
	}
	if param.ReservedPort == nil {
		err = ErrInValidParameter
		return
	}
	if param.UplinkPort == nil {
		err = ErrInValidParameter
		return
	}

	resp.UplinkPortMetadata = *param.UplinkPort
	resp.UplinkPortTotal, err = param.UplinkPort.Ranges.GetCount()
	if err != nil {
		return
	}
	resp.UplinkPortRegexp, err = GeneratePortRegexp(
		param.UplinkPort.RatePrefix,
		param.UplinkPort.Stacking,
		param.UplinkPort.PortSplitting,
		param.UplinkPort.NamingType, string(param.UplinkPort.Ranges),
	)
	if err != nil {
		return
	}

	resp.ServerPortMetadata = *param.DownlinkPort
	resp.ServerPortTotal, err = param.DownlinkPort.Ranges.GetCount()
	if err != nil {
		return
	}
	resp.ServerportRegexp, err = GeneratePortRegexp(
		param.DownlinkPort.RatePrefix,
		param.DownlinkPort.Stacking,
		param.DownlinkPort.PortSplitting,
		param.DownlinkPort.NamingType, string(param.DownlinkPort.Ranges),
	)
	if err != nil {
		return
	}

	resp.ReservedPortMetadata = *param.ReservedPort
	resp.ReservedPortTotal, err = param.ReservedPort.Ranges.GetCount()
	if err != nil {
		return
	}
	resp.ReservedPortRegexp, err = GeneratePortRegexp(
		param.ReservedPort.RatePrefix,
		param.ReservedPort.Stacking,
		param.ReservedPort.PortSplitting,
		param.ReservedPort.NamingType, string(param.ReservedPort.Ranges),
	)
	if err != nil {
		return
	}
	resp.PortTotal = resp.ServerPortTotal + resp.UplinkPortTotal + resp.ReservedPortTotal
	return
}

type PortQuotaDocument struct {
	ID        int
	Vendor    string
	Model     string
	PortTotal int

	ServerPortTotal    int           `json:"server_port_total"`
	ServerportRegexp   string        `json:"serverport_regexp"`
	ServerPortMetadata PortQuotaInfo `json:"server_port_metadata"`

	UplinkPortTotal    int
	UplinkPortRegexp   string
	UplinkPortMetadata PortQuotaInfo

	ReservedPortTotal    int
	ReservedPortRegexp   string
	ReservedPortMetadata PortQuotaInfo
}
