# GOJS
Manipulate dom with javascript via golang, useful when used in conjunction with sys ipc. Dev server included in package for testing.

## Example
Simple application using the dev server.

```html
<html>
    <style>
        .center {
            text-align: center;
        }
    </style>
    <body>        
        <div class="center">
            <h1 class="ui"> Hello World! </h1>  
        </div>
    </body>
</html>
```

```go
package main

var DEFAULT_PORT = 3001

func main() {
    ui := gojs.New(
		gojs.NewDevServer(DEFAULT_PORT),
		&gojs.UIConfig{
			HTMLDocPath: "main.html",
		},
	)

    ui.Show()
}

```