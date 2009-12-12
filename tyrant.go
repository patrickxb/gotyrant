/**
 * Copyright 2009 Patrick Crosby, XB Labs LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tyrant

// this comment coming up isn't really a comment...tells C package what to include

// #include "ttwrapper.h"
import "C"

import (
        "fmt";
        "os";
        "unsafe";
)

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

func Connect(host string, port int) (connection *Connection, err os.Error) {
        connection = new(Connection);
        connection.Tyrant = C.xtcrdb_new();
        open := C.xtcrdb_open(connection.Tyrant, C.CString("localhost"), 1978);
        if open == 0 {
                fmt.Printf("couldn't open database\n");
                connection.ErrorDisplay();
                return nil, os.NewError(connection.ErrorMessage());
        }
        fmt.Printf("connected to database localhost:1978\n");

        return connection, nil;
}

func ConnectDefault() (Connection *Connection, err os.Error) {
        return Connect("localhost", 1978)
}

func (connection *Connection) Close() os.Error {
        ok := C.xtcrdb_close(connection.Tyrant);
        if ok == 0 {
                return os.NewError(connection.ErrorMessage())
        }

        C.xtcrdb_del(connection.Tyrant);

        return nil;
}

// If record exists with this key, it is overwritten
func (connection *Connection) Put(primary_key string, columns ColumnMap) (err os.Error) {
        fmt.Printf("storing %v => %v\n", primary_key, columns);
        cols := C.xtc_mapnew();
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblput(connection.Tyrant, C.CString(primary_key), cols) == 0 {
                return os.NewError(connection.ErrorMessage())
        }
        // XXX use 'defer' for this?
        C.xtc_mapdel(cols);
        return nil;
}

// If record exists with this key, nothing happens
func (connection *Connection) Create(primaryKey string, columns ColumnMap) (err os.Error) {
        fmt.Printf("Create[%s]\n", primaryKey);
        cols := C.xtc_mapnew();
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblputkeep(connection.Tyrant, C.CString(primaryKey), cols) == 0 {
                return os.NewError(connection.ErrorMessage())
        }
        // XXX use 'defer' for this?
        C.xtc_mapdel(cols);
        return nil;
}

func (connection *Connection) Get(primaryKey string) *ColumnMap {
        cols := C.xtcrdb_tblget(connection.Tyrant, C.CString(primaryKey));
        if cols == nil {
                return nil
        }

        result := make(ColumnMap);
        result["PKEY"] = primaryKey;
        C.xtc_mapiterinit(cols);
        name := C.xtc_mapiternext2(cols);
        for name != nil {
                result[C.GoString(name)] = C.GoString(C.xtc_mapget2(cols, name));
                name = C.xtc_mapiternext2(cols);
        }
        return &result;
}

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
