package tyrant

// #include "ttwrapper.h"
import "C"

import "fmt"
import "unsafe"

// this doesn't work...
//type Connection unsafe.Pointer;

type Connection struct {
        Tyrant unsafe.Pointer;
}

type Query struct {
        Tyrant unsafe.Pointer;
}

func (connection *Connection) ErrorMessage() string {
        return C.GoString(C.xtcrdb_errmsg(C.xtcrdb_ecode(connection.Tyrant)))
}

func (connection *Connection) ErrorDisplay() {
        fmt.Printf("TT Error:  " + connection.ErrorMessage() + "\n")
}

func Connect() (connection *Connection, ok bool) {
        connection = new(Connection);
        connection.Tyrant = C.xtcrdb_new();
        open := C.xtcrdb_open(connection.Tyrant, C.CString("localhost"), 1978);
        if open == 0 {
                fmt.Printf("couldn't open database\n");
                connection.ErrorDisplay();
                ok = false;
        } else {
                fmt.Printf("connected to database localhost:1978\n");
                ok = true;
        }
        return connection, ok;
}

func (connection *Connection) Close() {
        ok := C.xtcrdb_close(connection.Tyrant);
        if ok == 0 {
                connection.ErrorDisplay()
        }

        C.xtcrdb_del(connection.Tyrant);
}

// If record exists with this key, it is overwritten
func (connection *Connection) Put(primary_key string, columns ColumnMap) (ok bool) {
        fmt.Printf("storing %v => %v\n", primary_key, columns);
        ok = true;
        cols := C.xtc_mapnew();
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblput(connection.Tyrant, C.CString(primary_key), cols) == 0 {
                connection.ErrorDisplay();
                ok = false;
        }
        // XXX use 'defer' for this?
        C.xtc_mapdel(cols);
        return ok;
}

// If record exists with this key, nothing happens
func (connection *Connection) Create(primary_key string, columns ColumnMap) (ok bool) {
        ok = true;
        cols := C.xtc_mapnew();
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblputkeep(connection.Tyrant, C.CString(primary_key), cols) == 0 {
                connection.ErrorDisplay();
                ok = false;
        }
        // XXX use 'defer' for this?
        C.xtc_mapdel(cols);
        return ok;
}

/*
func MakeQuery(connection unsafe.Pointer) (query unsafe.Pointer) {
        return C.xtcrdb_qrynew(connection);
}
*/

func (connection *Connection) MakeQuery() (query *Query) {
        query = new(Query);
        query.Tyrant = C.xtcrdb_qrynew(connection.Tyrant);
        return query;
}

func StringEqual() int { return int(C.x_streq()) }

func StringIncluded() int { return int(C.x_strinc()) }

func StringBeginsWith() int { return int(C.x_strbw()) }

func NumLessThan() int { return int(C.x_numlt()) }

func (query *Query) AddCondition(column_name string, op int, expression string) {
        fmt.Printf("adding condition on column '%s', operation = %d, expression = '%s'\n",
                column_name, op, expression);
        C.xtcrdb_qryaddcond(query.Tyrant, C.CString(column_name), _C_int(op), C.CString(expression));
}

func (query *Query) SetLimit(limit int) {
        C.xtcrdb_qrysetlimit(query.Tyrant, _C_int(limit), 0)
}

type ColumnMap map[string]string

type Row struct {
        Data ColumnMap;
}

type SearchResult struct {
        Rows    []Row;
        Count   int;
}

// XXX: just return a Vector of map[string] string...
func (connection *Connection) Execute(query *Query) SearchResult {
        list := C.xtcrdb_qrysearch(query.Tyrant);
        list_size := int(C.xtc_listnum(list));
        fmt.Printf("list size: %d\n", list_size);
        rows := make([]Row, list_size);
        for i := 0; i < list_size; i++ {
                pk := C.xtc_listval(list, _C_int(i));
                cols := C.xtcrdb_tblget(connection.Tyrant, pk);
                if cols != nil {
                        var row Row;
                        row.Data = make(ColumnMap);
                        row.Data["PKEY"] = C.GoString(pk);
                        C.xtc_mapiterinit(cols);
                        name := C.xtc_mapiternext2(cols);
                        for name != nil {
                                row.Data[C.GoString(name)] = C.GoString(C.xtc_mapget2(cols, name));
                                name = C.xtc_mapiternext2(cols);
                        }
                        rows[i] = row;
                }
        }
        var result SearchResult;
        result.Rows = rows;
        return result;
}
