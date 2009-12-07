#ifndef _TT_WRAPPER_H_
#define _TT_WRAPPER_H_

void* xtcrdb_new(void); 
int xtcrdb_open(void* rdb, const char* host, int port);
int xtcrdb_close(void* rdb);
void xtcrdb_del(void* rdb);

int xtcrdb_tblput(void* rdb, const char* pkey, void* cols);
int xtcrdb_tblputkeep(void* rdb, const char* pkey, void* cols);

int xtcrdb_ecode(void* rdb);
const char* xtcrdb_errmsg(int ecode);

void* xtc_mapnew();
void xtc_mapput(void* cols, const char* name, const char* value);
void xtc_mapdel(void* cols);
const char* xtc_mapget2(void* cols, const char* name);
void xtc_mapiterinit(void* map);
const char* xtc_mapiternext2(void* map);


void* xtcrdb_qrynew(void* rdb);
void xtcrdb_qryaddcond(void* query, const char* column_name, int operation, const char* expression);
void xtcrdb_qrysetlimit(void* query, int limit, int offset);
void* xtcrdb_qrysearch(void* query);

int xtc_listnum(void* list);
const char* xtc_listval(void* list, int index);
void* xtcrdb_tblget(void* connection, const char* pkey);

int x_streq();
int x_strinc();
int x_strbw();
int x_numlt();

#endif
