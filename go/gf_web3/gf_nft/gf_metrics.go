// SPDX-License-Identifier: GPL-2.0
/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_nft

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

//-------------------------------------------------

type GFmetrics struct {
	NftDBinsertsCount        prometheus.Counter
	NftAlchemyDBinsertsCount prometheus.Counter
}

//-------------------------------------------------

func MetricsCreate() *GFmetrics {

	// NFT_DB_INSERTS__COUNT
	nftDBinsertsCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_web3__nft_db_inserts__count"),
		Help: "how many new NFTs were inserted into the DB",
	})
	prometheus.MustRegister(nftDBinsertsCount)

	// NFT_ALCHEMY_DB_INSERTS__COUNT
	nftAlchemyDBinsertsCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_web3__nft_alchemy_db_inserts__count"),
		Help: "how many new NFTs from Alchemy were inserted into the DB",
	})
	prometheus.MustRegister(nftAlchemyDBinsertsCount)

	metrics := &GFmetrics{
		NftDBinsertsCount:        nftDBinsertsCount,
		NftAlchemyDBinsertsCount: nftAlchemyDBinsertsCount,
	}
	return metrics
}