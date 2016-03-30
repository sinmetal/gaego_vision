package gaego_vision

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"google.golang.org/api/vision/v1"
)

func init() {
	http.HandleFunc("/api/1/vision", handleVision)
}

func handleVision(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	imgurl := r.FormValue("imgurl")
	if imgurl == "" {
		http.Error(w, "required imgurl parameter", http.StatusBadRequest)
		return
	}

	client := urlfetch.Client(ctx)
	resp, err := client.Get(imgurl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		errmsg := fmt.Sprintf("%s fetch status = %s", imgurl, resp.Status)
		log.Infof(ctx, "%s", errmsg)
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warningf(ctx, "%s fetch body read error. err = %s", imgurl, err.Error())
		http.Error(w, fmt.Sprintf("%s fetch body read error", imgurl), http.StatusBadRequest)
		return
	}

	vresp, err := callVision(ctx, base64.StdEncoding.EncodeToString(blob))
	if err != nil {
		log.Warningf(ctx, "%s call vision api error. err = %s", imgurl, err.Error())
		http.Error(w, fmt.Sprintf("%s call vision api error", imgurl), http.StatusBadRequest)
		return
	}
	body, err := json.Marshal(vresp.Responses[0].LabelAnnotations)
	if err != nil {
		log.Warningf(ctx, "vision api response json marshal error. err = %s", err.Error())
		http.Error(w, fmt.Sprintf("%s vision api response json marshal error", imgurl), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func callVision(ctx context.Context, enc string) (*vision.BatchAnnotateImagesResponse, error) {
	img := &vision.Image{Content: enc}

	feature := &vision.Feature{
		Type:       "LABEL_DETECTION",
		MaxResults: 10,
	}
	req := &vision.AnnotateImageRequest{
		Image:    img,
		Features: []*vision.Feature{feature},
	}
	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, vision.CloudPlatformScope),
			Base:   &urlfetch.Transport{Context: ctx},
		},
	}
	svc, err := vision.New(client)
	if err != nil {
		return nil, err
	}

	return svc.Images.Annotate(batch).Do()
}
