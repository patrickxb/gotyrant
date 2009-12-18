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

#include "ttwrapper.h"
#include <tcrdb.h>

/* Tokyo Cabinet Maps */

void* 
xtc_mapnew()
{
        return tcmapnew();
}

void 
xtc_mapput(void* cols, const char* name, const char* value)
{
        tcmapput2(cols, name, value);
}

void
xtc_mapdel(void* cols)
{
        tcmapdel(cols);
}

const char* 
xtc_mapget2(void* cols, const char* name)
{
        return tcmapget2(cols, name);
}

void 
xtc_mapiterinit(void* map)
{
        tcmapiterinit(map);
}

const char*
xtc_mapiternext2(void* map)
{
        return tcmapiternext2(map);
}

/* Tokyo Cabinet Lists */

int
xtc_listnum(void* list)
{
        return tclistnum(list);
}

const char* 
xtc_listval(void* list, int index)
{
        return tclistval2(list, index);
}

/* Tokyo Tyrant new/open/close/del */

void*
xtcrdb_new()
{
        return tcrdbnew();
}

int
xtcrdb_open(void* rdb, const char* host, int port)
{
        return (int)tcrdbopen(rdb, host, port);
}

int 
xtcrdb_close(void* rdb)
{
        return (int)tcrdbclose(rdb);
}

void 
xtcrdb_del(void* rdb)
{
        tcrdbdel(rdb);
}

/* Tokyo Tyrant errors */

int
xtcrdb_ecode(void* rdb)
{
        return tcrdbecode(rdb);
}

const char* 
xtcrdb_errmsg(int ecode)
{
        return tcrdberrmsg(ecode);
}

/* Tokyo Tyrant tables */

void* 
xtcrdb_tblget(void* connection, const char* pkey)
{
        // XXX from sample code, but is this necessary??? strlen?
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblget(connection, pkbuf, pksiz);
}

int 
xtcrdb_tblput(void* rdb, const char* pkey, void* cols)
{
        // XXX from sample code, I assume for safety...
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblput(rdb, pkbuf, pksiz, cols);
}

int 
xtcrdb_tblputkeep(void* rdb, const char* pkey, void* cols)
{
        // XXX from sample code, I assume for safety...
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblputkeep(rdb, pkbuf, pksiz, cols);
}

/* Tokyo Tyrant queries */

void* 
xtcrdb_qrynew(void* rdb)
{
        return tcrdbqrynew(rdb);
}

void 
xtcrdb_qryaddcond(void* query, const char* column_name, int operation, const char* expression)
{
        tcrdbqryaddcond(query, column_name, operation, expression);
}

void 
xtcrdb_qrysetlimit(void* query, int limit, int offset)
{
        tcrdbqrysetlimit(query, limit, offset);
}

void 
xtcrdb_qrysetorder(void* query, const char* column_name, int order)
{
        tcrdbqrysetorder(query, column_name, order);
}

void*
xtcrdb_qrysearch(void* query)
{
        return tcrdbqrysearch(query);
}

int
xtcrdb_qrysearchout(void* query)
{
        return tcrdbqrysearchout(query);
}

int 
xtcrdb_qrysearchcount(void* query)
{
        return tcrdbqrysearchcount(query);
}

/* Tokyo Tyrant query conditions */

int x_streq() { return RDBQCSTREQ; }
int x_strinc() { return RDBQCSTRINC; }
int x_strbw() { return RDBQCSTRBW; }
int x_numlt() { return RDBQCNUMLT; }

/* Tokyo Tyrant query orders */

int x_strasc() { return RDBQOSTRASC; }
int x_strdesc() { return RDBQOSTRDESC; }
int x_numasc() { return RDBQONUMASC; }
int x_numdesc() { return RDBQONUMDESC; }


