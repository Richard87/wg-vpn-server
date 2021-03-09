import {useState} from "react";
import {MDBBtn, MDBCol, MDBIcon, MDBInput, MDBModal, MDBModalBody, MDBModalFooter, MDBModalHeader, MDBRow} from "mdbreact";
import qrImage from "./qrgen.png";
import styled from "styled-components";

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

export default function CreateClient ({onSubmit}) {
    const [name, setName] = useState("")
    const [ip, setIp] = useState("")
    const [publicKey, setPublicKey] = useState("")
    const [showNewClient, setShowNewClient] = useState(false)

    const onLocalSubmit = () => {
        setShowNewClient(false)
        onSubmit({name, ip, publicKey})
    }

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
                        <MDBBtn color="primary">Generate private key</MDBBtn>
                    </MDBCol>
                    <MDBCol>
                        <img src={qrImage} className="img-fluid" alt="Generated QR code"/>
                        <MDBBtn color="secondary">Show config</MDBBtn>
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