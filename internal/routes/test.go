package routes

import (
	"net/http"

	"github.com/bojurgess/backend.beno.lol/internal/config"
	"github.com/bojurgess/backend.beno.lol/internal/database"
)

type Test struct {
	DB     *database.Database
	Config *config.Config
}

func (p *Test) Route(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
<html>
<head>
    <title>SSE Example</title>
</head>
<body>
    <div id="sse-data"></div>

    <script>
        const eventSource = new EventSource('http://localhost:3000/user/8fzywdklm84r9hupsurfxdoj2');
        eventSource.onmessage = function(event) {
            const dataElement = document.getElementById('sse-data');
			const json = JSON.parse(event.data);
            dataElement.innerHTML += json.item.name + '<br>';
        };
    </script>
</body>
</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
