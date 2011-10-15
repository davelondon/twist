package twist

import (
	"http"
	"fmt"
)

type Writer struct {
	Output    http.ResponseWriter
	Buffer    string
	Templates []Template
	SendRoot  bool
	SendHtml  bool
}

func Root(w *Writer) *Item {

	return &Item{
		id:         "root",
		template:   nil,
		writer:     w,
		Attributes: make(map[string]string),
		Styles:     make(map[string]string),
	}

}

func NewWriter(o http.ResponseWriter, sendRoot bool, sendHtml bool) *Writer {
	return &Writer{
		Output:    o,
		Buffer:    "",
		Templates: make([]Template, 0),
		SendRoot:  sendRoot,
		SendHtml:  sendHtml,
	}
}

func (w *Writer) RegisterTemplate(t Template) {

	for i := 0; i < len(w.Templates); i++ {
		if w.Templates[i].name == t.name {
			return
		}
	}
	w.Templates = append(w.Templates, t)

}

func (c *Context) Send() {
	if c.Writer.SendRoot {
		c.Writer.sendPage(c.Root)
	} else {
		c.Writer.sendFragment()
	}

}

func (w *Writer) sendPage(item *Item) {

	root := `
<script src="/static/jquery.js"></script>
<script src="/static/json.js"></script>
<script src="/static/helpers.js"></script>
<script src="/static/native.history.js"></script>
<div id="head"></div>
<div id="root">`
	root += item.RenderHtml()
	root += `
</div>
<script>

	(function(window,undefined){

	    // Prepare
	    var History = window.History; // Note: We are using a capital H instead of a lower h
	    if ( !History.enabled ) {
	         // History.js is disabled for this browser.
	         // This is because we can optionally choose to support HTML4 browsers or not.
	        return false;
	    }

	    // Bind to StateChange Event
	    History.Adapter.bind(window,'statechange',function(){ // Note: We are using statechange instead of popstate
			var State = History.getState(); // Note: We are using History.getState() instead of event.state
			//History.log(State.data, State.title, State.url);
			$.post(State.url, null, function(data){$("#head").append($("<div>").html(data))}, "html");
	    });

	    // Change our States
//	    History.pushState({state:1}, "State 1", "?state=1"); // logs {state:1}, "State 1", "?state=1"
//	    History.pushState({state:2}, "State 2", "?state=2"); // logs {state:2}, "State 2", "?state=2"
//	    History.replaceState({state:3}, "State 3", "?state=3"); // logs {state:3}, "State 3", "?state=3"
//	    History.pushState(null, null, "?state=4"); // logs {}, '', "?state=4"
//	    History.back(); // logs {state:3}, "State 3", "?state=3"
//	    History.back(); // logs {state:1}, "State 1", "?state=1"
//	    History.back(); // logs {}, "Home Page", "?"
//	    History.go(2); // logs {state:3}, "State 3", "?state=3"

	})(window);

`

	templates := ""
	script := ""

	if len(w.Templates) > 0 {
		templates = `
var templatesToLoad = ` + fmt.Sprint(len(w.Templates)) + `;
var templatesLoaded = 0;`
		for i := 0; i < len(w.Templates); i++ {
			templates += `
$("#head").append($("<div>").load("/template_` + w.Templates[i].name + `", function() {templatesLoaded++;if(templatesLoaded == templatesToLoad){runScript();}}));`
		}

		script = `
function runScript()
{` + w.Buffer + `
}
</script>`

	} else {
		script = w.Buffer + `
</script>`
	}

	fmt.Fprint(w.Output, root+templates+script)
	w.Buffer = ""

}
func (w *Writer) sendFragment() {

	templates := ``
	script := ``
	if len(w.Templates) > 0 {
		templates = `
var templatesToLoad = ` + fmt.Sprint(len(w.Templates)) + `;
var templatesLoaded = 0;`
		for i := 0; i < len(w.Templates); i++ {
			templates += `
$("#head").append($("<div>").load("/template_` + w.Templates[i].name + `", function() {templatesLoaded++;if(templatesLoaded == templatesToLoad){runScript();}}));`
		}
		script = `
function runScript()
{` + w.Buffer + `
}`
	} else {
		script = w.Buffer
	}

	fmt.Fprint(w.Output,
		`<script>`+templates+script+`</script>`)
	w.Buffer = ""

}
