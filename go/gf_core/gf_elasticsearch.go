/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_core

import (
	"fmt"
	"context"
	"syscall"
	"errors"
	"net/http"
	"time"
	"github.com/olivere/elastic"
)

//-------------------------------------------------
// ELASTICSEARCH_CONNECTION_RETRY
type Gf_elasticsearch_retrier struct {
	backoff elastic.Backoff
}

func new_gf_elasticsearch_retrier() *Gf_elasticsearch_retrier {
	return &Gf_elasticsearch_retrier{
		backoff: elastic.NewExponentialBackoff(10 * time.Millisecond, 8 * time.Second),
	}
}

func (p_retrier *Gf_elasticsearch_retrier) Retry(p_ctx context.Context,
	p_retry_int int,
	p_req       *http.Request,
	p_resp      *http.Response,
	p_err       error) (time.Duration, bool, error) {

	// dont attempt to retry if a connection is refused
	if p_err == syscall.ECONNREFUSED {
		return 0, false, errors.New("elasticsearch server or network is down")
	}

	// stop retries after a certain number of them has already happened
	if p_retry_int >= 5 {
		return 0, false, nil
	}

	// have retirer determine the next wait period
	wait, stop_bool := p_retrier.backoff.Next(p_retry_int)
	return wait, stop_bool, nil
}

//-------------------------------------------------
func Elastic__get_client(p_es_host_str string, p_runtime_sys *RuntimeSys) (*elastic.Client, *GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_elasticsearch.Elastic__get_client()")

	// es_host_str := "127.0.0.1:9200"
	p_runtime_sys.LogFun("INFO", fmt.Sprintf("es_host - %s", p_es_host_str))

	url_str := fmt.Sprintf("http://%s", p_es_host_str)

	elasticsearch_client, err := elastic.NewClient(elastic.SetURL(url_str),
		elastic.SetRetrier(new_gf_elasticsearch_retrier()))

	if err != nil {
		gf_err := ErrorCreate("failed to create an ElasticSearch client",
			"elasticsearch_get_client",
			map[string]interface{}{"es_host_str": p_es_host_str,},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err	
	}

	// ping elasticsearch server
	ctx                     := context.Background()
	ping_url_str            := fmt.Sprintf("http://%s", p_es_host_str)
	es_info, resp_code, err := elasticsearch_client.Ping(ping_url_str).Do(ctx)
	if err != nil {
		gf_err := ErrorCreate("failed to ping ElasticSearch server with a client",
			"elasticsearch_ping",
			map[string]interface{}{"ping_url_str": ping_url_str,},
			err, "gf_core", p_runtime_sys)
		return nil, gf_err	
	}

	fmt.Printf("elasticsearch - resp_code/version %d/%s\n", resp_code, es_info.Version.Number)

	return elasticsearch_client, nil
}