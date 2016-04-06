var sample;
(function (sample) {
    sample.send = function(url, imgurl) {
        var retries = 3;

        function get() {
            if(!url) {
                throw "url is required.";
            }

            var req = new XMLHttpRequest();
            req.onload = function() {
                if (req.readyState === 4) {
                    if (req.status === 200) {
                        document.getElementById("target").setAttribute("src", imgurl);

                        var output = document.getElementById("output");
                        output.innerHTML = "";
                        var node = JsonHuman.format(JSON.parse(req.responseText), {
                            // Show or hide Array-Indices in the output
                            showArrayIndex: true,

                            // Hyperlinks Option
                            // Enable <a> tag in the output html based on object keys
                            // Supports only strings and arrays
                            hyperlinks: {
                                enable: true,
                                keys: ['url'],          // Keys which will be output as links
                                target: '_blank'       // 'target' attribute of a
                            },

                            // Options for displaying bool
                            bool: {
                                // Show text? And what text for true & false?
                                showText: true,
                                text: {
                                    true: "Yes",
                                    false: "No"
                                },

                                // Show image? And which images (urls)?
                                showImage: true,
                                img: {
                                    true: 'css/true.png',
                                    false: 'css/false.png'
                                }
                            }
                        });

                        output.appendChild(node);
                    } else if (req.status === 500) {
                        retries--;
                        if(retries > 0) {
                            console.log("Retrying...");
                            setTimeout(function(){post()}, 100);
                        } else {
                            console.error("get error");
                        }
                    } else {
                        console.error(req.status);
                    }
                }
            };

            try {
                req.open("GET", url, true);
                req.send(null);
            } catch(e) {
                throw "Error retrieving data file. Some browsers only accept cross-domain request with HTTP.";
            }
        }
        get();
    };

    sample.submit = function() {
        document.getElementById("target").setAttribute("src", "");
        document.getElementById("output").innerHTML = "";

        var imgurl = document.getElementById("imgurl").value;

        sample.send("/api/1/vision?imgurl=" + imgurl, imgurl);
    };
})(sample || (sample = {}));
