import {MDBBtn, MDBCard, MDBCardBody, MDBCardText, MDBCardTitle} from "mdbreact";
import {formatDistance} from "date-fns"

export default function ClientCard({client}) {

    const latestHandshake = client?.latestHandshake ? new Date(client.latestHandshake) : null
    const latestHandshakeDistance = latestHandshake ? formatDistance(latestHandshake, new Date(), { addSuffix: true }) : ""

    return <MDBCard>
        <MDBCardBody>
            <MDBCardTitle>{client?.name ?? "TEST"}</MDBCardTitle>
            <MDBCardText>
                <strong>peer: </strong><span
                style={{wordBreak: "break-all"}}>{client?.publicKey}</span><br/>
                <strong>endpoint: </strong>{client?.endpoint}<br/>
                <strong>allowed ips: </strong>{(client?.ip ?? []).join(", ")}<br/>
                <strong>latest handshake: </strong>{latestHandshakeDistance}<br/>
                <strong>transfer: </strong>{bytesToSize(client?.receivedBytes)} received, {bytesToSize(client?.sentBytes)} sent<br/>
            </MDBCardText>
            <MDBBtn href="#">MDBBtn</MDBBtn>
        </MDBCardBody>
    </MDBCard>
}

function bytesToSize(bytes) {
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 Byte';
    const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
    return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
}