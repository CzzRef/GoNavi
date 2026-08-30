package safety

import (
	"GoNavi-Wails/internal/ai"
	"testing"
)

func TestClassifySQL(t *testing.T) {
	tests := []struct {
		sql  string
		want ai.SQLOperationType
	}{
		{"SELECT * FROM users", ai.SQLOpQuery},
		{"  select id from t", ai.SQLOpQuery},
		{"SHOW TABLES", ai.SQLOpQuery},
		{"DESCRIBE users", ai.SQLOpQuery},
		{"DESC users", ai.SQLOpQuery},
		{"EXPLAIN SELECT 1", ai.SQLOpQuery},
		{"WITH cte AS (SELECT 1) SELECT * FROM cte", ai.SQLOpQuery},
		{"PRAGMA table_info(t)", ai.SQLOpQuery},
		{"VALUES (1, 2)", ai.SQLOpQuery},
		{"INSERT INTO users VALUES (1)", ai.SQLOpDML},
		{"UPDATE users SET name='x'", ai.SQLOpDML},
		{"DELETE FROM users WHERE id=1", ai.SQLOpDML},
		{"REPLACE INTO users VALUES (1)", ai.SQLOpDML},
		{"MERGE INTO t USING s ON t.id=s.id", ai.SQLOpDML},
		{"CREATE TABLE t (id INT)", ai.SQLOpDDL},
		{"ALTER TABLE t ADD col INT", ai.SQLOpDDL},
		{"DROP TABLE t", ai.SQLOpDDL},
		{"TRUNCATE TABLE t", ai.SQLOpDDL},
		{"RENAME TABLE old TO new", ai.SQLOpDDL},
		{"/* comment */ SELECT 1", ai.SQLOpQuery},
		{"-- comment\nDELETE FROM t", ai.SQLOpDML},
		{"-- line1\n-- line2\nSELECT 1", ai.SQLOpQuery},
		{"/* block */ -- line\nUPDATE t SET x=1", ai.SQLOpDML},
		{"", ai.SQLOpOther},
		{"   ", ai.SQLOpOther},
		{"-- only comment", ai.SQLOpOther},
	}
	for _, tt := range tests {
		got := ClassifySQL(tt.sql)
		if got != tt.want {
			t.Errorf("ClassifySQL(%q) = %s, want %s", tt.sql, got, tt.want)
		}
	}
}

func TestIsHighRiskSQL(t *testing.T) {
	tests := []struct {
		sql         string
		highRisk    bool
		warningKey  string
	}{
		{"DROP TABLE users", true, "ai_service.backend.warning.sql_drop"},
		{"DROP DATABASE test", true, "ai_service.backend.warning.sql_drop"},
		{"TRUNCATE TABLE users", true, "ai_service.backend.warning.sql_truncate"},
		{"DELETE FROM users", true, "ai_service.backend.warning.sql_delete_without_where"},             // 无 WHERE
		{"DELETE FROM users WHERE id=1", false, ""},                                                    // 有 WHERE
		{"UPDATE users SET name='x'", true, "ai_service.backend.warning.sql_update_without_where"},    // 无 WHERE
		{"UPDATE users SET name='x' WHERE id=1", false, ""},                                            // 有 WHERE
		{"SELECT * FROM users", false, ""},
		{"INSERT INTO users VALUES (1)", false, ""},
	}
	for _, tt := range tests {
		highRisk, warningKey := IsHighRiskSQL(tt.sql)
		if highRisk != tt.highRisk {
			t.Errorf("IsHighRiskSQL(%q) = %v, want %v", tt.sql, highRisk, tt.highRisk)
		}
		if warningKey != tt.warningKey {
			t.Errorf("IsHighRiskSQL(%q) warning = %q, want %q", tt.sql, warningKey, tt.warningKey)
		}
	}
}

func TestGuard_ReadOnly(t *testing.T) {
	g := NewGuard(ai.PermissionReadOnly)
	tests := []struct {
		sql     string
		allowed bool
	}{
		{"SELECT * FROM t", true},
		{"INSERT INTO t VALUES (1)", false},
		{"UPDATE t SET x=1", false},
		{"DELETE FROM t", false},
		{"DROP TABLE t", false},
		{"CREATE TABLE t (id INT)", false},
	}
	for _, tt := range tests {
		result := g.Check(tt.sql)
		if result.Allowed != tt.allowed {
			t.Errorf("Guard[readonly].Check(%q).Allowed = %v, want %v", tt.sql, result.Allowed, tt.allowed)
		}
	}
}

func TestGuard_ReadWrite(t *testing.T) {
	g := NewGuard(ai.PermissionReadWrite)
	tests := []struct {
		sql     string
		allowed bool
		confirm bool
	}{
		{"SELECT * FROM t", true, false},
		{"INSERT INTO t VALUES (1)", true, true},
		{"UPDATE t SET x=1", true, true},       // 允许但需确认
		{"DELETE FROM t WHERE id=1", true, true}, // 允许但需确认
		{"DROP TABLE t", false, true},            // DDL 不允许
		{"CREATE TABLE t (id INT)", false, true},
	}
	for _, tt := range tests {
		result := g.Check(tt.sql)
		if result.Allowed != tt.allowed {
			t.Errorf("Guard[readwrite].Check(%q).Allowed = %v, want %v", tt.sql, result.Allowed, tt.allowed)
		}
		if result.RequiresConfirm != tt.confirm {
			t.Errorf("Guard[readwrite].Check(%q).RequiresConfirm = %v, want %v", tt.sql, result.RequiresConfirm, tt.confirm)
		}
	}
}

func TestGuard_Full(t *testing.T) {
	g := NewGuard(ai.PermissionFull)
	tests := []struct {
		sql     string
		allowed bool
	}{
		{"SELECT * FROM t", true},
		{"INSERT INTO t VALUES (1)", true},
		{"DROP TABLE t", true},
		{"CREATE TABLE t (id INT)", true},
		// 例程调用与例程部署在完全模式下也不放行：Agent 不执行例程，只交出候选 SQL。
		{"CALL bulk_insert_users(100000)", false},
		{"EXEC dbo.sp_rebuild", false},
		{"CREATE PROCEDURE p() BEGIN END", false},
		{"DROP FUNCTION IF EXISTS f", false},
	}
	for _, tt := range tests {
		result := g.Check(tt.sql)
		if result.Allowed != tt.allowed {
			t.Errorf("Guard[full].Check(%q).Allowed = %v, want %v", tt.sql, result.Allowed, tt.allowed)
		}
	}
}

func TestIsRoutineSQL(t *testing.T) {
	routine := []string{
		"CALL sp_fix_orders(1)",
		"call sp_fix_orders(1)",
		"EXEC dbo.sp_rebuild",
		"EXECUTE sp_rebuild",
		"CREATE PROCEDURE p() BEGIN END",
		"CREATE OR REPLACE PROCEDURE p AS BEGIN NULL; END;",
		"CREATE DEFINER=`root`@`localhost` PROCEDURE p() BEGIN END",
		"CREATE OR REPLACE EDITIONABLE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;",
		"ALTER FUNCTION f COST 100",
		"DROP PROCEDURE IF EXISTS p",
		"DROP TRIGGER tr",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN END",
		"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN END",
		"CREATE PACKAGE BODY pkg AS END;",
		"-- deploy\nCREATE PROCEDURE p() BEGIN END",
		"/* deploy */ DROP PROCEDURE p",
	}
	for _, sql := range routine {
		if !IsRoutineSQL(sql) {
			t.Errorf("IsRoutineSQL(%q) = false, want true", sql)
		}
	}

	notRoutine := []string{
		"SELECT * FROM t",
		"INSERT INTO t VALUES (1)",
		"UPDATE t SET x = 1",
		"DELETE FROM t WHERE id = 1",
		"CREATE TABLE t (id INT)",
		"DROP TABLE t",
		"ALTER TABLE t ADD col INT",
		"TRUNCATE TABLE t",
		"CREATE INDEX idx ON t (id)",
		"CREATE MATERIALIZED VIEW v AS SELECT 1",
		"CREATE DATABASE d",
		"CREATE USER u",
		// 列名与例程对象同名不得触发例程判定：扫描在第一个非例程对象关键字处停止。
		"CREATE TABLE t (trigger VARCHAR(10), function VARCHAR(10))",
		"CREATE VIEW v AS SELECT procedure_id FROM t",
		"",
		"-- only comment",
	}
	for _, sql := range notRoutine {
		if IsRoutineSQL(sql) {
			t.Errorf("IsRoutineSQL(%q) = true, want false", sql)
		}
	}
}

func TestClassifySQLRoutine(t *testing.T) {
	tests := []struct {
		sql  string
		want ai.SQLOperationType
	}{
		{"CALL p(1)", ai.SQLOpRoutine},
		{"EXEC p", ai.SQLOpRoutine},
		{"CREATE PROCEDURE p() BEGIN END", ai.SQLOpRoutine},
		{"ALTER PROCEDURE p COMMENT 'x'", ai.SQLOpRoutine},
		{"DROP TRIGGER tr", ai.SQLOpRoutine},
		// 例程分类不得抢走普通 DDL 与 DML。
		{"CREATE TABLE t (id INT)", ai.SQLOpDDL},
		{"DROP TABLE t", ai.SQLOpDDL},
		{"INSERT INTO t VALUES (1)", ai.SQLOpDML},
		{"SELECT * FROM t", ai.SQLOpQuery},
	}
	for _, tt := range tests {
		if got := ClassifySQL(tt.sql); got != tt.want {
			t.Errorf("ClassifySQL(%q) = %s, want %s", tt.sql, got, tt.want)
		}
	}
}

func TestGuardRoutineDeniedAtEveryLevel(t *testing.T) {
	levels := []ai.SQLPermissionLevel{
		ai.PermissionReadOnly,
		ai.PermissionReadWrite,
		ai.PermissionFull,
		ai.SQLPermissionLevel("unknown-level"),
	}
	statements := []string{
		"CALL sp_fix_orders(1)",
		"EXEC dbo.sp_rebuild",
		"CREATE PROCEDURE p() BEGIN END",
		"CREATE DEFINER=`root`@`localhost` PROCEDURE p() BEGIN END",
		"DROP FUNCTION f",
	}
	for _, level := range levels {
		g := NewGuard(level)
		for _, sql := range statements {
			result := g.Check(sql)
			if result.Allowed {
				t.Errorf("Guard[%s].Check(%q).Allowed = true, want false", level, sql)
			}
			if result.OperationType != ai.SQLOpRoutine {
				t.Errorf("Guard[%s].Check(%q).OperationType = %s, want %s", level, sql, result.OperationType, ai.SQLOpRoutine)
			}
		}
	}
}

func TestGuard_HighRiskWarning(t *testing.T) {
	g := NewGuard(ai.PermissionFull)
	result := g.Check("DELETE FROM users")
	if result.WarningMessage == "" {
		t.Error("expected high-risk warning for DELETE without WHERE")
	}
	if !result.RequiresConfirm {
		t.Error("expected RequiresConfirm for high-risk SQL")
	}
}
