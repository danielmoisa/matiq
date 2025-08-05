package oracle

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielmoisa/matiq/internal/actionruntime/common"
	parser_sql "github.com/danielmoisa/matiq/internal/utils/parser/sql"
	"github.com/danielmoisa/matiq/internal/utils/resourcelist"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type Connector struct {
	resourceOptions Resource
	actionOptions   Action
}

func (o *Connector) ValidateResourceOptions(resourceOptions map[string]interface{}) (common.ValidateResult, error) {
	// format resource options
	if err := mapstructure.Decode(resourceOptions, &o.resourceOptions); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	// validate oracle options
	validate := validator.New()
	if err := validate.Struct(o.resourceOptions); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	return common.ValidateResult{Valid: true}, nil
}

func (o *Connector) ValidateActionTemplate(actionOptions map[string]interface{}) (common.ValidateResult, error) {
	// format action options
	if err := mapstructure.Decode(actionOptions, &o.actionOptions); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	// validate oracle options
	validate := validator.New()
	if err := validate.Struct(o.actionOptions); err != nil {
		return common.ValidateResult{Valid: false}, err
	}

	return common.ValidateResult{Valid: true}, nil
}

func (o *Connector) TestConnection(resourceOptions map[string]interface{}) (common.ConnectionResult, error) {
	// get oracle connection
	db, err := o.getConnectionWithOptions(resourceOptions)
	if err != nil {
		return common.ConnectionResult{Success: false}, err
	}
	defer db.Close()

	// test oracle connection
	if err := db.Ping(); err != nil {
		return common.ConnectionResult{Success: false}, err
	}

	return common.ConnectionResult{Success: true}, nil
}

func (o *Connector) GetMetaInfo(resourceOptions map[string]interface{}) (common.MetaInfoResult, error) {
	// get oracle connection
	db, err := o.getConnectionWithOptions(resourceOptions)
	if err != nil {
		return common.MetaInfoResult{Success: false}, err
	}
	defer db.Close()

	// test oracle connection
	if err := db.Ping(); err != nil {
		return common.MetaInfoResult{Success: false}, err
	}

	columns := mapColumns(db)

	return common.MetaInfoResult{
		Success: true,
		Schema:  columns,
	}, nil
}

func (o *Connector) Run(resourceOptions map[string]interface{}, actionOptions map[string]interface{}, rawActionOptions map[string]interface{}) (common.RuntimeResult, error) {
	// get Oracle connection
	db, err := o.getConnectionWithOptions(resourceOptions)
	if err != nil {
		return common.RuntimeResult{Success: false}, errors.New("failed to get oracle connection")
	}
	defer db.Close()
	// format query
	if err := mapstructure.Decode(actionOptions, &o.actionOptions); err != nil {
		return common.RuntimeResult{Success: false}, err
	}
	// set context field
	errInSetRawQuery := o.actionOptions.SetRawQueryAndContext(rawActionOptions)
	if errInSetRawQuery != nil {
		return common.RuntimeResult{Success: false}, errInSetRawQuery
	}

	queryResult := common.RuntimeResult{Success: false}
	queryResult.Rows = make([]map[string]interface{}, 0, 0)
	queryResult.Extra = make(map[string]interface{})
	err = nil

	// start a default context
	ctx, cancel := context.WithTimeout(context.TODO(), common.DEFAULT_QUERY_AND_EXEC_TIMEOUT)
	defer cancel()

	// action mode switch
	switch o.actionOptions.Mode {
	case ACTION_SQL_MODE:
		fallthrough
	case ACTION_SQL_SAFE_MODE:
		sqlEscaper := parser_sql.NewSQLEscaper(resourcelist.TYPE_ORACLE_ID)
		escapedSQL, sqlArgs, errInEscapeSQL := sqlEscaper.EscapeSQLActionTemplate(o.actionOptions.RawQuery, o.actionOptions.Context, o.actionOptions.IsSafeMode())
		if errInEscapeSQL != nil {
			return queryResult, errInEscapeSQL
		}
		// check if o.actionOptions.Opts.Raw is select query
		isSelectQuery := false

		var query SQL
		if err := mapstructure.Decode(o.actionOptions.Opts, &query); err != nil {
			return queryResult, errors.New("type error of action content")
		}
		lexer := parser_sql.NewLexer(query.Raw)
		isSelectQuery, err = parser_sql.IsSelectSQL(lexer)
		if err != nil {
			return common.RuntimeResult{Success: false}, err
		}

		// fetch data
		if isSelectQuery && o.actionOptions.IsSafeMode() {
			fmt.Printf("[oracle] [RUN] isSelectQuery, IsSafeMode, escapedSQL: %s\n", escapedSQL)
			rows, err := db.QueryContext(ctx, escapedSQL, sqlArgs...)
			if err != nil {
				return queryResult, err
			}
			mapRes, err := common.RetrieveToMap(rows)
			if err != nil {
				return queryResult, err
			}
			defer rows.Close()
			queryResult.Success = true
			queryResult.Rows = mapRes
		} else if isSelectQuery && !o.actionOptions.IsSafeMode() {
			fmt.Printf("[oracle] [RUN] isSelectQuery, !IsSafeMode, query.Raw: %s\n", query.Raw)
			rows, err := db.QueryContext(ctx, escapedSQL)
			if err != nil {
				return queryResult, err
			}
			mapRes, err := common.RetrieveToMap(rows)
			if err != nil {
				return queryResult, err
			}
			defer rows.Close()
			queryResult.Success = true
			queryResult.Rows = mapRes
		} else if !isSelectQuery && o.actionOptions.IsSafeMode() {
			fmt.Printf("[oracle] [RUN] !isSelectQuery, IsSafeMode, escapedSQL: %s\n", escapedSQL)
			execResult, err := db.ExecContext(ctx, escapedSQL, sqlArgs...)
			if err != nil {
				return queryResult, err
			}
			affectedRows, err := execResult.RowsAffected()
			if err != nil {
				return queryResult, err
			}
			queryResult.Success = true
			queryResult.Extra["message"] = fmt.Sprintf("Affeted %d rows.", affectedRows)
		} else if !isSelectQuery && !o.actionOptions.IsSafeMode() {
			fmt.Printf("[oracle] [RUN] !isSelectQuery, !IsSafeMode, query.Raw: %s\n", query.Raw)
			execResult, err := db.ExecContext(ctx, escapedSQL)
			if err != nil {
				return queryResult, err
			}
			affectedRows, err := execResult.RowsAffected()
			if err != nil {
				return queryResult, err
			}
			queryResult.Success = true
			queryResult.Extra["message"] = fmt.Sprintf("Affeted %d rows.", affectedRows)
		}
	case ACTION_GUI_MODE:
		err = errors.New("unsupported action mode")
	default:
		err = errors.New("unsupported action mode")
	}
	return queryResult, err
}
