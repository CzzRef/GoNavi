package app

import (
	"testing"

	"GoNavi-Wails/internal/ai"
)

// InspectSQL 是 MCP 与 CLI 共用的判定来源。例程语句必须在这里就被标记，
// 否则下游两个执行面只能各自猜测。
func TestInspectSQLMarksRoutineStatements(t *testing.T) {
	cases := []struct {
		name   string
		dbType string
		sql    string
	}{
		{"mysql call", "mysql", "CALL sp_fix_orders(1)"},
		{"mysql create procedure", "mysql", "CREATE PROCEDURE p() BEGIN END"},
		{"mysql create definer procedure", "mysql", "CREATE DEFINER=`root`@`localhost` PROCEDURE p() BEGIN END"},
		{"mysql drop procedure", "mysql", "DROP PROCEDURE IF EXISTS p"},
		{"mysql create trigger", "mysql", "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN END"},
		{"postgres call", "postgres", "CALL refresh_stats()"},
		{"oracle create function", "oracle", "CREATE OR REPLACE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;"},
		{"sqlserver exec", "sqlserver", "EXEC dbo.sp_rebuild"},
		{"sqlserver bare sp_", "sqlserver", "sp_helpdb 'app'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inspection := InspectSQL(tc.dbType, tc.sql)
			if inspection.StatementCount != 1 {
				t.Fatalf("StatementCount = %d, want 1", inspection.StatementCount)
			}
			stmt := inspection.Statements[0]
			if !stmt.Routine {
				t.Fatalf("InspectSQL(%q, %q).Routine = false, want true", tc.dbType, tc.sql)
			}
			if stmt.ReadOnly {
				t.Fatalf("a routine statement must not be read-only: %q", tc.sql)
			}
		})
	}
}

// 只读判定优先于例程判定。SQL Server 的裸过程调用启发式对任意标识符开头的语句都成立，
// 若不让只读先胜出，普通查询会被误判为例程并阻断。
func TestInspectSQLDoesNotMarkNonRoutineStatements(t *testing.T) {
	cases := []struct {
		name   string
		dbType string
		sql    string
	}{
		{"select", "mysql", "SELECT * FROM users"},
		{"insert", "mysql", "INSERT INTO users VALUES (1)"},
		{"update", "mysql", "UPDATE users SET name = 'x' WHERE id = 1"},
		{"delete", "mysql", "DELETE FROM users WHERE id = 1"},
		{"create table", "mysql", "CREATE TABLE t (id INT)"},
		{"create table with routine-named column", "mysql", "CREATE TABLE t (id INT, trigger VARCHAR(10))"},
		{"drop table", "mysql", "DROP TABLE t"},
		{"alter table", "mysql", "ALTER TABLE t ADD col INT"},
		{"create index", "postgres", "CREATE INDEX idx ON t (id)"},
		{"create materialized view", "postgres", "CREATE MATERIALIZED VIEW v AS SELECT 1"},
		{"sqlserver select", "sqlserver", "SELECT TOP 10 * FROM users"},
		{"sqlserver update", "sqlserver", "UPDATE users SET name = 'x' WHERE id = 1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inspection := InspectSQL(tc.dbType, tc.sql)
			if inspection.StatementCount != 1 {
				t.Fatalf("StatementCount = %d, want 1", inspection.StatementCount)
			}
			if inspection.Statements[0].Routine {
				t.Fatalf("InspectSQL(%q, %q).Routine = true, want false", tc.dbType, tc.sql)
			}
		})
	}
}

// EXPLAIN ANALYZE 会真正执行语句体，因此不能成为包裹例程调用的绕行通道。
func TestInspectSQLMarksRoutineInsideExplainAnalyze(t *testing.T) {
	inspection := InspectSQL("postgres", "EXPLAIN ANALYZE EXECUTE prepared_stmt")
	if inspection.StatementCount != 1 {
		t.Fatalf("StatementCount = %d, want 1", inspection.StatementCount)
	}
	if !inspection.Statements[0].Routine {
		t.Fatal("EXPLAIN ANALYZE EXECUTE should be marked as a routine statement")
	}

	plain := InspectSQL("postgres", "EXPLAIN SELECT 1")
	if plain.Statements[0].Routine {
		t.Fatal("plain EXPLAIN must not be marked as a routine statement")
	}
}

func TestHeadlessSafetyDeniesRoutineAtEveryLevel(t *testing.T) {
	levels := []ai.SQLPermissionLevel{
		ai.PermissionReadOnly,
		ai.PermissionReadWrite,
		ai.PermissionFull,
		ai.SQLPermissionLevel("unknown-level"),
	}
	statements := []string{
		"CALL sp_fix_orders(1)",
		"CREATE PROCEDURE p() BEGIN END",
		"DROP PROCEDURE IF EXISTS p",
	}

	for _, level := range levels {
		for _, sql := range statements {
			decision := evaluateHeadlessSQLSafety(level, "mysql", sql)
			if len(decision.Disallowed) != 1 {
				t.Fatalf("evaluateHeadlessSQLSafety(%s, %q).Disallowed = %d, want 1", level, sql, len(decision.Disallowed))
			}
			if got := decision.Disallowed[0].Operation; got != ai.SQLOpRoutine {
				t.Fatalf("operation for %q = %s, want %s", sql, got, ai.SQLOpRoutine)
			}
			if len(decision.ConfirmRequired) != 0 {
				t.Fatalf("a denied routine must not appear as merely confirm-required: %q", sql)
			}
		}
	}
}

// 完全模式仍然放行普通 DDL/DML，例程收紧不得连带把它们一起挡掉。
func TestHeadlessSafetyStillAllowsNonRoutineAtFull(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (id INT)",
		"DROP TABLE t",
		"UPDATE t SET x = 1 WHERE id = 1",
		"SELECT * FROM t",
	} {
		decision := evaluateHeadlessSQLSafety(ai.PermissionFull, "mysql", sql)
		if len(decision.Disallowed) != 0 {
			t.Fatalf("evaluateHeadlessSQLSafety(full, %q).Disallowed = %d, want 0", sql, len(decision.Disallowed))
		}
	}
}
