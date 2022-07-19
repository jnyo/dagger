package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dagger/cloak/engine"
	dagger "github.com/dagger/cloak/sdk/go"
)

func main() {
	/*
		if err := engine.RunGraphiQL(context.Background(), 8080); err != nil {
			panic(err)
		}
	*/

	err := engine.Start(context.Background(), func(ctx context.Context) error {
		var output string
		var err error

		_, err = dagger.Do(ctx, `mutation{import(ref:"alpine"){name}}`)
		if err != nil {
			return err
		}
		_, err = dagger.Do(ctx, `mutation{import(ref:"helloworld_ts"){name}}`)
		if err != nil {
			return err
		}

		/*
			output, err = dagger.Do(ctx, tools.IntrospectionQuery)
			if err != nil {
				return err
			}
			fmt.Printf("schema: %s\n", output)
		*/

		// output, err = dagger.Do(ctx, `{core{image(ref:"alpine:3.15"){fs}}}`)
		output, err = dagger.Do(ctx, `{alpine{build(pkgs:["gcc","python3"])}}`)
		// output, err = dagger.Do(ctx, `{helloworld_ts{echo(message:"hi"){fs}}}`)
		if err != nil {
			return err
		}
		fmt.Printf("output: %s\n", output)
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return err
		}

		/*
					output, err = dagger.Do(ctx, fmt.Sprintf(`
			{
			 core {
			  exec(input: {
			   args: ["true"],
			   mounts: [{path:"/", fs:%q}]
			  }) {
					 mount(path:"/")
			   }
			 }
			}`, result["core"].(map[string]interface{})["image"].(map[string]interface{})["fs"].(string)))
					if err != nil {
						return err
					}
					fmt.Printf("output: %s\n", output)
					if err := json.Unmarshal([]byte(output), &result); err != nil {
						return err
					}
		*/

		output, err = dagger.Do(ctx, fmt.Sprintf(`mutation{evaluate(fs:%q)}`,
			result["alpine"].(map[string]interface{})["build"].(string)))
		// result["helloworld_ts"].(map[string]interface{})["echo"].(map[string]interface{})["fs"].(string)))
		// result["core"].(map[string]interface{})["exec"].(map[string]interface{})["mount"].(string)))
		if err != nil {
			return err
		}
		fmt.Printf("output: %s\n", output)
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return err
		}

		if err := engine.Shell(ctx, result["evaluate"].(string)); err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
