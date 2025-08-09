package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 请求结构体
type Request struct {
	Path      string                            `json:"path"`
	Condition map[string]map[string]interface{} `json:"condition"` // 使用interface{}来支持多种类型
	Result    map[string][]string               `json:"result"`
	PageItem  struct {
		Page int `json:"page"`
		Size int `json:"size"`
	} `json:"pageItem"`
}

// 生成SQL语句的函数
func GenerateSQL(req Request) (string, error) {
	// 解析路径
	pathParts := strings.Split(req.Path, "-")
	if len(pathParts)%2 != 1 {
		return "", fmt.Errorf("invalid path format")
	}

	// 构建SELECT部分
	var selectFields []string
	for entity, attrs := range req.Result {
		for _, attr := range attrs {
			selectFields = append(selectFields, fmt.Sprintf("%s.%s", entity, attr))
		}
	}
	selectClause := strings.Join(selectFields, ",")

	// 构建FROM和JOIN部分
	var fromJoinClauses []string
	fromJoinClauses = append(fromJoinClauses, pathParts[0]) // 第一个实体作为FROM的表

	for i := 1; i < len(pathParts); i += 2 {
		relationship := pathParts[i]
		nextEntity := pathParts[i+1]
		fromJoinClauses = append(fromJoinClauses, fmt.Sprintf("left join %s on %s.src_id=%s.id", relationship, relationship, pathParts[i-1]))
		fromJoinClauses = append(fromJoinClauses, fmt.Sprintf("left join %s on %s.id=%s.dst_id", nextEntity, nextEntity, relationship))
	}
	fromJoinClause := strings.Join(fromJoinClauses, " ")

	// 构建WHERE部分
	var whereClauses []string
	for entity, conditions := range req.Condition {
		for attr, value := range conditions {
			var clause string
			switch v := value.(type) {
			case string:
				if v == "NULL" {
					clause = fmt.Sprintf("%s.%s IS NULL", entity, attr)
				} else if v == "NOT NULL" {
					clause = fmt.Sprintf("%s.%s IS NOT NULL", entity, attr)
				} else {
					clause = fmt.Sprintf("%s.%s='%s'", entity, attr, v)
				}
			case float64, int, int64:
				clause = fmt.Sprintf("%s.%s=%v", entity, attr, v)
			case bool:
				clause = fmt.Sprintf("%s.%s=%v", entity, attr, v)
			default:
				return "", fmt.Errorf("unsupported type for condition value: %T", v)
			}
			whereClauses = append(whereClauses, clause)
		}
	}
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "where " + strings.Join(whereClauses, " and ")
	}

	// 构建OFFSET和LIMIT部分
	offset := (req.PageItem.Page - 1) * req.PageItem.Size
	limit := req.PageItem.Size
	offsetLimitClause := fmt.Sprintf("offset %d limit %d", offset, limit)

	// 组合完整的SQL语句
	sql := fmt.Sprintf("select %s from %s %s %s", selectClause, fromJoinClause, whereClause, offsetLimitClause)
	return sql, nil
}

func main() {
	// 示例请求
	reqJSON := `{
		"path": "region-contain-switch-own-port",
		"condition": {
			"region":{"name":"SCA"},
			"switch":{"role":"server_leaf"}
		},
		"result": {
			"region":["name"],
			"switch":["admin_ip","name","role","sn","model","vendor","series"],
			"port":["*"]
		},
		"pageItem": {
			"page":1,
			"size":10
		}
	}`

	// 解析请求
	var req Request
	err := json.Unmarshal([]byte(reqJSON), &req)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	// 生成SQL语句
	sql, err := GenerateSQL(req)
	if err != nil {
		fmt.Println("Error generating SQL:", err)
		return
	}

	fmt.Println("Generated SQL:", sql)
	// Generated SQL: select region.name,switch.admin_ip,switch.name,switch.role,switch.sn,switch.model,switch.vendor,switch.series,port.* from region left join contain on contain.src_id=region.id left join switch on switch.id=contain.dst_id left join own on own.src_id=switch.id left join port on port.id=own.dst_id where region.name='SCA' and switch.role='server_leaf' offset 0 limit 10
}
