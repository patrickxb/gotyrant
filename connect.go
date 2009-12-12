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

package main

import (
        "tyrant";
        "fmt";
        "os";
)

func main() {
        fmt.Printf("trying to use tyrant package to connect to tokyo tyrant\n");
        connection, err := tyrant.ConnectDefault();
        if err != nil {
                fmt.Printf("no connection: %s\n", err.String());
                os.Exit(1);
        }

        columns := make(map[string]string);
        columns["name"] = "falcon";
        columns["age"] = "31";
        columns["lang"] = "ja";
        err = connection.Put("12345", columns);
        if err != nil {
                fmt.Printf("put failed: %s\n", err.String())
        }

        query := connection.MakeQuery();
        query.AddCondition("name", tyrant.StringBeginsWith(), "f");
        result := connection.Execute(query);
        fmt.Printf("%d rows returned\n", len(result.Rows));
        for index, row := range result.Rows {
                fmt.Printf("ROW %d ------------------\n", index);
                for name, value := range row.Data {
                        fmt.Printf("%s => %s\n", name, value)
                }
                fmt.Printf("\n");
        }

        fmt.Printf("done\n");
}
