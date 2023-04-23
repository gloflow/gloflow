/*
MIT License

Copyright (c) 2023 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_rpc_lib

import (
	"bytes"
	"fmt"
	"strconv"
	"net"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type GFudpServerDef struct {
	ExternalIPtoUseStr string
	PortStr            string

	// unique package pattern for the defined server that allows clients
	// to register with this server
	RegistrationPacketLst []byte
	RegistrationConfirmationPacketLst []byte

	// object headers to sync across connected clients
	ObjectsToSyncLst []PacketHeaderObjUpdate
}

type GFudpClient struct {
	IDint   int
	IPstr   string
	PortInt int
	Address *net.UDPAddr
}

type PacketHeaderObjUpdate []byte

//-------------------------------------------------

func UDPstartListening(pServerDef GFudpServerDef,
	pRuntimeSys *gf_core.RuntimeSys) {

	go func() {

		clientsTableMap := map[string]GFudpClient{}

		connExt, gfErr := UDPlistenForPackets(pServerDef.ExternalIPtoUseStr,
			pServerDef.PortStr,
			pRuntimeSys)
		if gfErr != nil {
			panic("failed to open UDP socket")
		}
		
		buff := make([]byte, 65536)
		for {

			// read UDP socket
			bytesReadLengthInt, udpAddressOfSender, err := connExt.ReadFromUDP(buff)
			if err != nil {

				// dont stop processing packets if reading one failed
				continue
			}


			
			packetDataLst := buff[0:bytesReadLengthInt]

			pRuntimeSys.LogNewFun("INFO", "received UDP data", map[string]interface{}{
				"package_data_lst": packetDataLst,
				"address_str":      udpAddressOfSender,
			})
			

			// packet header - 6 bytes
			packetHeaderLst := packetDataLst[0:6]

			clientCompositeKeyStr := fmt.Sprintf("%s:%d", udpAddressOfSender.IP, udpAddressOfSender.Port)



			//------------------------------
			// CLIENT_REGISTRATION
			// check if this client sending the packet is registering with this server.
			if bytes.Equal(packetHeaderLst, pServerDef.RegistrationPacketLst) {

				// client is already registered, so ignore this registration and just send a confirmation back
				// to unblock the client and continue to next packet
				if _, ok := clientsTableMap[clientCompositeKeyStr]; ok {

					// send confirmation packet back to client
					gfErr := UDPsendRegistrationConfirmation(udpAddressOfSender,
						connExt,
						pServerDef,
						pRuntimeSys)
					if gfErr != nil {

					}
					
					continue
				}

				pRuntimeSys.LogNewFun("INFO", "new client registering with server", map[string]interface{}{
					"address_str": udpAddressOfSender,
				})

				
				playerIDbyteLst := packetDataLst[6:7]
				playerIDint  := int(playerIDbyteLst[0])

				client := GFudpClient{
					IDint:   playerIDint,
					IPstr:   fmt.Sprint(udpAddressOfSender.IP),
					PortInt: udpAddressOfSender.Port,
					Address: udpAddressOfSender,
				}

				// register client in the clients_table
				clientsTableMap[clientCompositeKeyStr] = client

				// send confirmation packet back to client
				gfErr := UDPsendRegistrationConfirmation(udpAddressOfSender,
					connExt,
					pServerDef,
					pRuntimeSys)
				if gfErr != nil {

				}

				// skip the rest and process the next packet
				continue
			}

			//------------------------------

			for _, headerLst := range pServerDef.ObjectsToSyncLst {
				if bytes.Equal(packetHeaderLst, []byte(headerLst)) {

					UDPsendObjUpdateToAllClients(packetDataLst,
						clientCompositeKeyStr,
						clientsTableMap,
						connExt,
						pServerDef,
						pRuntimeSys)
				}
			}

			// send data back to client right away
			connExt.WriteTo(packetDataLst, udpAddressOfSender)
		}
	}()
}

//-------------------------------------------------

func UDPsendObjUpdateToAllClients(pPacketDataLst []byte,
	pUpdateOriginatorClientCompositeKeyStr string,
	pClientsTableMap map[string]GFudpClient,
	pConnExt         *net.UDPConn,
	pServerDef       GFudpServerDef,
	pRuntimeSys      *gf_core.RuntimeSys) {
	
	// iterate through all connected clients, and send them all this object update
	for k, udpClient := range pClientsTableMap {

		// skip sending this update packet to the player that originated it
		if k == pUpdateOriginatorClientCompositeKeyStr {
			continue
		}

		// fmt.Println(" -- sending packet to", udpClient.Address, pPacketDataLst)

		// send this packet to all other clients connected to this server
		_, err := pConnExt.WriteTo(pPacketDataLst, udpClient.Address)
		if err != nil {
			msgStr  := "failed to write obj update packet to UDP socket of one of the clients"
			typeStr := "udp_write_packge_to_socket_error"

			// only printing error info, since packet sending frequency is too high
			// and handling them as an gf_error would overload the error handling systems
			// (DB, Sentry, etc.)
			pRuntimeSys.LogNewFun("ERROR", msgStr, map[string]interface{}{
				"error_type_str": typeStr,
				"address_str":    udpClient.Address,
			})
		}
	}
}

//-------------------------------------------------

func UDPsendRegistrationConfirmation(pUDPaddressOfSender *net.UDPAddr,
	pConnExt    *net.UDPConn,
	pServerDef  GFudpServerDef,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	_, err := pConnExt.WriteTo(pServerDef.RegistrationConfirmationPacketLst, pUDPaddressOfSender)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to write client registration confirmation packet to UDP socket",
			"udp_write_packge_to_socket_error",
			map[string]interface{}{
				"ip_str":   fmt.Sprint(pUDPaddressOfSender.IP),
				"port_str": fmt.Sprint(pUDPaddressOfSender.Port),
			},
			err, "gf_rpc_lib", pRuntimeSys)
		return gfErr
	}
	return nil
}

//-------------------------------------------------

func UDPlistenForPackets(pExternalIPtoUseStr string,
	pPortStr    string,
	pRuntimeSys *gf_core.RuntimeSys) (*net.UDPConn, *gf_core.GFerror) {

	portInt, _ := strconv.Atoi(pPortStr)

	localAddr := net.UDPAddr{
		Port: portInt,

		// IMPORTANT!! - the IP has to be declared here for this servers packets that are sent to clients to have the 
		//               IP of a (possibly) load balancer.
		//               this is critical for there to exist a static public IP address that is reliable and not changing.
		//               for users that are behind NAT's we want all subsequent packets after the first one to go to the
		//               same IP/port (which is the one of the load balancer).
		IP: net.ParseIP(pExternalIPtoUseStr),
	}

	conn, err := net.ListenUDP("udp", &localAddr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to open UDP socket listening on a particular port",
			"udp_open_socket_error",
			map[string]interface{}{
				"port_str":               pPortStr,
				"external_ip_to_use_str": pExternalIPtoUseStr,
			},
			err, "gf_rpc_lib", pRuntimeSys)
		return nil, gfErr
	}

	// set buffering of UDP socket to make sure data isnt lost if server is taking too long
	// to read data.
	socketBufferSizeInt := 2 * 1024 * 1024 // 2MB
	conn.SetReadBuffer(socketBufferSizeInt)

	fmt.Printf("listening on UDP IP/port - %s:%s\n", pExternalIPtoUseStr, pPortStr)
	return conn, nil
}