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
        "strconv";
        "unsafe";
)

type Connection struct {
        Tyrant unsafe.Pointer;
}

type Query struct {
        Tyrant unsafe.Pointer;
}

func (connection *Connection) ErrorCode() int {
        return int(C.xtcrdb_ecode(connection.Tyrant));
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
        open := C.xtcrdb_open(connection.Tyrant, C.CString(host), C.int(port));
        if open == 0 {
                fmt.Printf("couldn't open database %s:%d\n", host, port);
                connection.ErrorDisplay();
                return nil, os.NewError(connection.ErrorMessage());
        }
        fmt.Printf("connected to database %s:%d\n", host, port);

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
        defer C.xtc_mapdel(cols);
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblput(connection.Tyrant, C.CString(primary_key), cols) == 0 {
                return os.NewError(connection.ErrorMessage())
        }
        return nil;
}

// If record exists with this key, nothing happens
func (connection *Connection) Create(primaryKey string, columns ColumnMap) (exists bool, err os.Error) {
        fmt.Printf("Create[%s]\n", primaryKey);
        cols := C.xtc_mapnew();
        defer C.xtc_mapdel(cols);
        for name, value := range columns {
                C.xtc_mapput(cols, C.CString(name), C.CString(value))
        }
        if C.xtcrdb_tblputkeep(connection.Tyrant, C.CString(primaryKey), cols) == 0 {
                if connection.ErrorCode() != ErrCodeKeep() {
                        return exists, os.NewError(connection.ErrorMessage())
                }
                exists = true
        }
        return exists, nil;
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
func IncludesToken() int { return int(C.x_stror()) }

func OrderStrAsc() int { return int(C.x_strasc()) }
func OrderStrDesc() int { return int(C.x_strdesc()) }
func OrderNumAsc() int { return int(C.x_numasc()) }
func OrderNumDesc() int { return int(C.x_numdesc()) }

func ErrCodeKeep() int { return int(C.x_errcode_keep()) }

func (query *Query) AddCondition(column_name string, op int, expression string) {
        fmt.Printf("adding condition on column '%s', operation = %d, expression = '%s'\n",
                column_name, op, expression);
        // C.xtcrdb_qryaddcond(query.Tyrant, C.CString(column_name), _C_int(op), C.CString(expression));
        C.xtcrdb_qryaddcond(query.Tyrant, C.CString(column_name), C.int(op), C.CString(expression));
}

func (query *Query) SetLimit(limit int) {
        C.xtcrdb_qrysetlimit(query.Tyrant, C.int(limit), 0)
}

func (query *Query) SetLimitOffset(limit int, offset int) {
        C.xtcrdb_qrysetlimit(query.Tyrant, C.int(limit), C.int(offset))
}

func (query *Query) SetOrder(columnName string, order int) {
        C.xtcrdb_qrysetorder(query.Tyrant, C.CString(columnName), C.int(order));
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
                pk := C.xtc_listval(list, C.int(i));
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
        result.Count = list_size;
        return result;
}

func (query *Query) Count() int {
        result := C.xtcrdb_qrysearchcount(query.Tyrant);
        return int(result);
}

// Removes all records that match the query
func (query *Query) Remove() os.Error {
        result := C.xtcrdb_qrysearchout(query.Tyrant);
        if result == 0 {
                return os.ErrorString("Error removing records")
        }
        return nil
}

func (cmap *ColumnMap) GetInt64(key string) (int64, os.Error) {
        frow := *cmap;
        s, ok := frow[key];
        if !ok {
                //return 0, os.NewError("'%s' not found in column map", key);
                return 0, os.NewError("key not found in column map");
        }
        result, err := strconv.Atoi64(s);
        if err != nil {
                return 0, os.NewError("error converting column to int");
        }
        return result, nil;
}
