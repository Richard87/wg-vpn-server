import {useState} from "react";
import {MDBBtn, MDBCol, MDBIcon, MDBInput, MDBModal, MDBModalBody, MDBModalFooter, MDBModalHeader, MDBRow} from "mdbreact";
import styled from "styled-components";
import {useQuery} from "react-query";
import QRCode from 'qrcode.react';


const FloatingButton = styled(MDBBtn)`
  position: fixed !important;
  right: 2rem;
  bottom: 2rem;
  border-radius: 50%;
  width: 4rem;
  height: 4rem;
  padding: 1rem;

  &:hover {
    transform: scale(1.1, 1.1);
  }
`

const wgConfig = (name, address, privateKey,DNS, endpoint, publicKey) => {
    return `
[Interface]
# Name = ${name}
Address = ${address}
PrivateKey = ${privateKey}
DNS = ${DNS}

[Peer]
# Name = ${endpoint}
Endpoint = ${endpoint}
PublicKey = ${publicKey}
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25`
}

export default function CreateClient ({onSubmit}) {
    const [name, setName] = useState("")
    const [ip, setIp] = useState("")
    const [publicKey, setPublicKey] = useState("")
    const [showNewClient, setShowNewClient] = useState(false)
    const [showConfig, setShowConfig] = useState(false)
    const [privateKey, setPrivateKey] = useState(null)

    const {data: config, refetch} = useQuery("config", {
        queryFn: () => fetch(`${process.env.REACT_APP_API_SERVER}/config`).then(res => res.json()),
        onSuccess: data => setIp(data.nextAvailableIp4)
    })

    const onLocalSubmit = () => {
        setShowNewClient(false)
        onSubmit({name, allowedIps: [ip], publicKey}).then(refetch)
    }

    const onGenerate = () => {
        const {publicKey, privateKey} = window.wireguard.generateKeypair()
        setPrivateKey(privateKey)
        setPublicKey(publicKey)
    }

    const currentConfig = wgConfig(name,ip,privateKey,config?.recommendedDNS,config?.endpoint,config?.publicKey)

    return <>
        <FloatingButton onClick={() => setShowNewClient(true)} gradient="purple">
            <MDBIcon icon="plus"/>
        </FloatingButton>
        <MDBModal centered size={"lg"} isOpen={showNewClient} toggle={() => setShowNewClient(false)}>
            <MDBModalHeader toggle={() => setShowNewClient(false)}>New client</MDBModalHeader>
            <MDBModalBody>
                <MDBRow>
                    <MDBCol>
                        <MDBInput icon="signature"
                                  label="Name"
                                  value={name}
                                  onChange={e => setName(e.target.value)}
                                  name="name"/>
                        <MDBInput icon="globe-europe"
                                  label="IP Address"
                                  value={ip}
                                  onChange={e => setIp(e.target.value)}
                                  name="ip"/>
                        <MDBInput icon="key"
                                  label="Public key"
                                  value={publicKey}
                                  onChange={e => setPublicKey(e.target.value)}
                                  name="publicKey"/>
                        <MDBBtn onClick={onGenerate} color="primary">Generate private key</MDBBtn>
                    </MDBCol>
                    <MDBCol>
                        {showConfig
                            ? <pre><code>{currentConfig}</code></pre>
                            : <QRCode style={{width: "100%", height: "auto"}} renderAs={"svg"} value={currentConfig} />
                        }
                        <MDBBtn color="secondary" onClick={() => setShowConfig(show => !show)}>{showConfig ? 'Show QR Code' : 'Show config'}</MDBBtn>
                    </MDBCol>
                </MDBRow>
            </MDBModalBody>
            <MDBModalFooter>
                <MDBBtn color="secondary" onClick={() => setShowNewClient(false)}>Close</MDBBtn>
                <MDBBtn color="primary" onClick={onLocalSubmit}>Save changes</MDBBtn>
            </MDBModalFooter>
        </MDBModal>
    </>
}