package httpserver

import (
	"fmt"

	"github.com/gookit/rux"
)

func Start() {
	r := rux.New()

	// Add Routes:
	r.GET("/", func(c *rux.Context) {
		c.Text(200, "hello")
	})
	r.GET("/hello/{name}", func(c *rux.Context) {
		c.Text(200, "hello "+c.Param("name"))
	})
	r.POST("/post", func(c *rux.Context) {
		c.Text(200, "hello")
	})
	// add multi method support for an route path
	r.Add("/post[/{id}]", func(c *rux.Context) {
		if c.Param("id") == "" {
			// do create post
			c.Text(200, "created")
			return
		}

		id := c.Params.Int("id")
		// do update post
		c.Text(200, "updated "+fmt.Sprint(id))
	}, rux.POST, rux.PUT)

	// Start server
	r.Listen(":8088")
	// can also
	// http.ListenAndServe(":8080", r)
}
