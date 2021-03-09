import {MDBBtn, MDBCard, MDBCardBody, MDBCardText, MDBCardTitle} from "mdbreact";

export default function ClientCard({client}) {
    return <MDBCard>
        <MDBCardBody>
            <MDBCardTitle>{client?.name ?? "TEST"}</MDBCardTitle>
            <MDBCardText>
                <strong>peer: </strong><span
                style={{wordBreak: "break-all"}}>neN+LnFpJ2FiuBhWVGr/VLl4ubu9cOKyI1K0VUZFSnk=</span><br/>
                <strong>endpoint: </strong>77.18.62.145:15427<br/>
                <strong>allowed ips: </strong>192.168.43.101/32<br/>
                <strong>latest handshake: </strong> 4 days, 1 hour, 22 minutes, 58 seconds ago<br/>
                <strong>transfer: </strong>301.63 MiB received, 361.54 MiB sent<br/>
            </MDBCardText>
            <MDBBtn href="#">MDBBtn</MDBBtn>
        </MDBCardBody>
    </MDBCard>
}
