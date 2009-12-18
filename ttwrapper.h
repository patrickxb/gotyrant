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

#ifndef _TT_WRAPPER_H_
#define _TT_WRAPPER_H_

/* Tokyo Cabinet Maps */
void* xtc_mapnew();
void xtc_mapput(void* cols, const char* name, const char* value);
void xtc_mapdel(void* cols);
const char* xtc_mapget2(void* cols, const char* name);
void xtc_mapiterinit(void* map);
const char* xtc_mapiternext2(void* map);

/* Tokyo Cabinet Lists */
int xtc_listnum(void* list);
const char* xtc_listval(void* list, int index);

/* Tokyo Tyrant new/open/close/del */
void* xtcrdb_new(void); 
int xtcrdb_open(void* rdb, const char* host, int port);
int xtcrdb_close(void* rdb);
void xtcrdb_del(void* rdb);

/* Tokyo Tyrant errors */
int xtcrdb_ecode(void* rdb);
const char* xtcrdb_errmsg(int ecode);

/* Tokyo Tyrant tables */
void* xtcrdb_tblget(void* connection, const char* pkey);
int xtcrdb_tblput(void* rdb, const char* pkey, void* cols);
int xtcrdb_tblputkeep(void* rdb, const char* pkey, void* cols);

/* Tokyo Tyrant queries */
void* xtcrdb_qrynew(void* rdb);
void xtcrdb_qryaddcond(void* query, const char* column_name, int operation, const char* expression);
void xtcrdb_qrysetlimit(void* query, int limit, int offset);
void xtcrdb_qrysetorder(void* query, const char* column_name, int order);
void* xtcrdb_qrysearch(void* query);
int xtcrdb_qrysearchout(void* query);
int xtcrdb_qrysearchcount(void* query);

/* Tokyo Tyrant query conditions */
int x_streq();
int x_strinc();
int x_strbw();
int x_numlt();

/* Tokyo Tyrant query orders */
int x_strasc();
int x_strdesc();
int x_numasc();
int x_numdesc();

#endif
