package main

import "tyrant"
import "fmt"

func main() {
        fmt.Printf("trying to use tyrant package to connect to tokyo tyrant\n");
        connection, ok := tyrant.Connect();
        if !ok {
                fmt.Printf("no connection...")
        } else {
                columns := make(map[string]string);
                columns["name"] = "falcon";
                columns["age"] = "31";
                columns["lang"] = "ja";
                if connection.Put("12345", columns) {
                        fmt.Printf("put ok\n")
                } else {
                        fmt.Printf("put failed\n")
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
        }

        // not sure why this doesn't work:
        //        tyrant.Close(connection);

        fmt.Printf("done\n");
}
