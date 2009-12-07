#include "ttwrapper.h"
#include <tcrdb.h>

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

void* 
xtc_mapnew()
{
        return tcmapnew();
}

void
xtc_mapdel(void* cols)
{
        tcmapdel(cols);
}

void 
xtc_mapput(void* cols, const char* name, const char* value)
{
        tcmapput2(cols, name, value);
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

int 
xtcrdb_tblput(void* rdb, const char* pkey, void* cols)
{
        // XXX from sample code, but is this necessary??? strlen?
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblput(rdb, pkbuf, pksiz, cols);
}

int 
xtcrdb_tblputkeep(void* rdb, const char* pkey, void* cols)
{
        // XXX from sample code, but is this necessary??? strlen?
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblputkeep(rdb, pkbuf, pksiz, cols);
}

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

int 
x_streq() 
{
        return RDBQCSTREQ;
}

int
x_strinc()
{
        return RDBQCSTRINC;
}

int
x_strbw()
{
        return RDBQCSTRBW;
}

int 
x_numlt()
{
        return RDBQCNUMLT;
}

void*
xtcrdb_qrysearch(void* query)
{
        return tcrdbqrysearch(query);
}

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

void* 
xtcrdb_tblget(void* connection, const char* pkey)
{
        // XXX from sample code, but is this necessary??? strlen?
        char pkbuf[256];
        int pksiz = sprintf(pkbuf, "%s", pkey);
        return tcrdbtblget(connection, pkbuf, pksiz);
}
