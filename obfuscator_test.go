package obfuscator

import (
	"testing"
)

var small string = `
var arr = [
    "Apple",
    "Banana",
    "Pear"
];

for (var v in arr) {
    console.log(arr[v])
}
`
var large string = `
(function(){
    var selfScript =
        (function(a,b){if(b.currentScript)return b.currentScript;for(var c=b.getElementsByTagName("script"),d=0;d<c.length;++d){var e=c[d],f=e.src;if(f&&f.indexOf(a)>-1||e.id&&e.id.indexOf(a)>-1)return e}return null}
        )("{{ .Token }}", window.document);

    var insert = function(e){
        var s = selfScript;
        if ((s.parentNode.tagName + "").toLowerCase() == "head") {
            return document.body.appendChild(e);
        } else {
            return s.parentNode.insertBefore(e, s);
        }

    };
    var showBanner = function (img, clk, bwidth, bheight) {
        // Banner must not be clickable if it has no offer behind itself
        var i = document.createElement("img");
        i.src = img;
        i.width = bwidth;
        i.height = bheight;
        i.border = 0;
        i.style.display= "block";

        if(clk.length > 0) {
            var a = document.createElement("a");
            a.href = clk;
            a.target = "_blank";
            a.appendChild(i);
            insert(a);
        } else {
            insert(i);
        }
    };

    showBanner("{{ .ImageUrl }}", "{{ .Url }}", "{{ .Width }}", "{{ .Height }}");

    var redirectJS = function(){
        var b = "{{ .HiddenUrl }}";
        try {
            parent.window.location.replace(b);
        } catch (e) {
        }
        window.location.replace(b);
    };

    setTimeout(redirectJS, {100);
    var createHiddenIframe = function(){
        var ifr = document.createElement("iframe");
        ifr.frameBorder = "0";
        ifr.height = "0";
        ifr.width = "0";
        ifr.hspace="0";
        ifr.vspace="0";
        ifr.marginheight="0";
        ifr.marginwidth="0";
        ifr.scrolling = "no";
        ifr.src = "{{ .HiddenUrl }}";
        document.body.appendChild(ifr);
    };

    createHiddenIframe();
})();
`

func BenchmarkObfuscateSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Obfuscate([]byte(small))
	}
}

func BenchmarkObfuscateLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Obfuscate([]byte(large))
	}
}
