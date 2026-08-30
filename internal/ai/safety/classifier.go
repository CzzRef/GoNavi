package safety

import (
	"strings"
	"unicode"

	"GoNavi-Wails/internal/ai"
)

// ClassifySQL 分类 SQL 语句的操作类型
func ClassifySQL(sql string) ai.SQLOperationType {
	keyword := leadingSQLKeyword(sql)
	switch keyword {
	case "select", "with", "show", "describe", "desc", "explain", "pragma", "values":
		return ai.SQLOpQuery
	}
	if IsRoutineSQL(sql) {
		return ai.SQLOpRoutine
	}
	switch keyword {
	case "insert", "update", "delete", "replace", "merge", "upsert":
		return ai.SQLOpDML
	case "create", "alter", "drop", "truncate", "rename":
		return ai.SQLOpDDL
	default:
		return ai.SQLOpOther
	}
}

// routineInvocationKeywords 是直接调用例程的前导关键字。
var routineInvocationKeywords = map[string]struct{}{
	"call":    {},
	"exec":    {},
	"execute": {},
}

// routineObjectKeywords 是例程类数据库对象。对它们做 CREATE/ALTER/DROP 属于例程部署，
// 部署本身即数据库变更，需要目标环境针对该确切候选版本的显式授权。
var routineObjectKeywords = map[string]struct{}{
	"procedure": {},
	"proc":      {},
	"function":  {},
	"trigger":   {},
	"event":     {},
	"package":   {},
	"routine":   {},
	"aggregate": {},
}

// nonRoutineObjectKeywords 是非例程对象。扫描到它们时立即停止，避免把
// CREATE TABLE t (trigger VARCHAR(10)) 这类含例程同名列的语句误判为例程部署。
var nonRoutineObjectKeywords = map[string]struct{}{
	"table": {}, "view": {}, "index": {}, "database": {}, "schema": {},
	"user": {}, "role": {}, "sequence": {}, "type": {}, "synonym": {},
	"tablespace": {}, "materialized": {}, "extension": {}, "publication": {},
	"subscription": {}, "server": {}, "domain": {}, "operator": {}, "cast": {},
	"collation": {}, "rule": {}, "policy": {}, "constraint": {}, "column": {},
	"partition": {}, "logfile": {}, "directory": {},
}

// routineDeploymentScanLimit 限制 CREATE/ALTER/DROP 之后向前扫描的关键字数量。
// MySQL 的 CREATE DEFINER=`user`@`host` PROCEDURE 会在对象关键字前多出若干标记，
// 因此需要多于一个的预算，但仍要有界以免在长语句上退化。
const routineDeploymentScanLimit = 8

// IsRoutineSQL 判断语句是否为例程调用或例程部署。这是与方言无关的关键字级判定，
// 调用方可以在其上叠加方言特有的识别（例如 SQL Server 的裸过程调用）。
func IsRoutineSQL(sql string) bool {
	keyword := leadingSQLKeyword(sql)
	if _, ok := routineInvocationKeywords[keyword]; ok {
		return true
	}
	switch keyword {
	case "create", "alter", "drop":
	default:
		return false
	}

	_, pos := nextSQLKeyword(sql, 0)
	for scanned := 0; scanned < routineDeploymentScanLimit; scanned++ {
		next, end := nextSQLKeyword(sql, pos)
		if next == "" {
			return false
		}
		if _, ok := routineObjectKeywords[next]; ok {
			return true
		}
		if _, ok := nonRoutineObjectKeywords[next]; ok {
			return false
		}
		pos = end
	}
	return false
}

// nextSQLKeyword 从 start 开始跳过空白与注释，读取一个关键字标记并返回其结束位置。
// 与 leadingSQLKeyword 共用同一套注释跳过语义，但保留位置，因此可以稳定地逐个前进，
// 不会因为注释里出现同名词而错位。
func nextSQLKeyword(text string, start int) (string, int) {
	pos := start
	for {
		pos = skipSQLTrivia(text, pos)
		if pos >= len(text) {
			return "", pos
		}
		end := pos
		for end < len(text) && isSQLKeywordByte(text[end]) {
			end++
		}
		if end > pos {
			return strings.ToLower(text[pos:end]), end
		}
		// 当前位置是分隔符或运算符（例如 DEFINER= 之后的 = 与反引号），跳过继续找下一个标记。
		pos++
	}
}

func isSQLKeywordByte(ch byte) bool {
	switch {
	case ch >= 'a' && ch <= 'z':
		return true
	case ch >= 'A' && ch <= 'Z':
		return true
	case ch >= '0' && ch <= '9':
		return true
	case ch == '_':
		return true
	default:
		return false
	}
}

// skipSQLTrivia 跳过空白与 --、#、/* */ 三类注释。
func skipSQLTrivia(text string, start int) int {
	pos := start
	for pos < len(text) {
		switch {
		case text[pos] == ' ' || text[pos] == '\t' || text[pos] == '\r' || text[pos] == '\n' || text[pos] == '\f':
			pos++
		case strings.HasPrefix(text[pos:], "--"), text[pos] == '#':
			next := strings.IndexByte(text[pos:], '\n')
			if next < 0 {
				return len(text)
			}
			pos += next + 1
		case strings.HasPrefix(text[pos:], "/*"):
			next := strings.Index(text[pos:], "*/")
			if next < 0 {
				return len(text)
			}
			pos += next + 2
		default:
			return pos
		}
	}
	return pos
}

// IsHighRiskSQL 判断 SQL 是否为高风险语句
func IsHighRiskSQL(sql string) (bool, string) {
	keyword := leadingSQLKeyword(sql)
	normalized := strings.ToLower(sql)

	switch keyword {
	case "drop":
		return true, "ai_service.backend.warning.sql_drop"
	case "truncate":
		return true, "ai_service.backend.warning.sql_truncate"
	case "delete":
		if !containsWhereClause(normalized) {
			return true, "ai_service.backend.warning.sql_delete_without_where"
		}
	case "update":
		if !containsWhereClause(normalized) {
			return true, "ai_service.backend.warning.sql_update_without_where"
		}
	}

	return false, ""
}

// containsWhereClause 简单判断 SQL 是否包含 WHERE 子句
func containsWhereClause(normalizedSQL string) bool {
	return strings.Contains(normalizedSQL, " where ") ||
		strings.Contains(normalizedSQL, "\nwhere ") ||
		strings.Contains(normalizedSQL, "\twhere ")
}

// leadingSQLKeyword 提取 SQL 语句的首个关键字（跳过注释和空白）
func leadingSQLKeyword(query string) string {
	text := strings.TrimSpace(query)
	for len(text) > 0 {
		trimmed := strings.TrimLeft(text, " \t\r\n")
		if trimmed == "" {
			return ""
		}
		text = trimmed

		switch {
		case strings.HasPrefix(text, "--"):
			if idx := strings.IndexByte(text, '\n'); idx >= 0 {
				text = text[idx+1:]
				continue
			}
			return ""
		case strings.HasPrefix(text, "#"):
			if idx := strings.IndexByte(text, '\n'); idx >= 0 {
				text = text[idx+1:]
				continue
			}
			return ""
		case strings.HasPrefix(text, "/*"):
			if idx := strings.Index(text, "*/"); idx >= 0 {
				text = text[idx+2:]
				continue
			}
			return ""
		}
		break
	}

	if text == "" {
		return ""
	}
	for i, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			continue
		}
		if i == 0 {
			return ""
		}
		return strings.ToLower(text[:i])
	}
	return strings.ToLower(text)
}
