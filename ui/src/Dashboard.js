import ClientCard from "./ClientCard";
import styled from "styled-components";
import {useQuery} from "react-query";
import CreateClient from "./CreateClient";
import {MDBBtn, MDBCol, MDBModal, MDBModalBody, MDBModalFooter, MDBModalHeader, MDBRow} from "mdbreact";
import {useState} from "react";
import {authFetch} from "./index";


const Grid = styled.div`
  margin-top: 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
  grid-gap: 2rem;
`


export default function Dashboard() {
    const {data, isSuccess, refetch} = useQuery(`clients`)
    const [remove, setRemove] = useState(null)

    const onSubmit = (newClient) => {
        return authFetch(`clients`, {method: "POST", json: newClient}).then(() => refetch())
    }

    const onDelete = client => {
        setRemove(null)
        return authFetch(`clients/${client.id}`, {method: "DELETE"}).then(() => refetch())
    }

    return <>
        <Grid>
            {isSuccess && data.map(client => <ClientCard onDelete={() => setRemove(client)} key={client.id} client={client}/>)}
        </Grid>
        <RemoveClient onConfirm={() => onDelete(remove)} onClose={() => setRemove(null)} name={remove?.name} />
        <CreateClient onSubmit={onSubmit}/>
    </>
}

const RemoveClient = ({onClose, onConfirm, name}) => (
    <MDBModal centered size={"lg"} isOpen={!!name} toggle={onClose}>
        <MDBModalHeader toggle={onClose}>Remove {name}</MDBModalHeader>
        <MDBModalBody>
            <MDBRow>
                <MDBCol>
                    <h1>Are you sure you want to remove {name}?</h1>
                </MDBCol>
            </MDBRow>
        </MDBModalBody>
        <MDBModalFooter>
            <MDBBtn color="secondary" onClick={onClose}>Cancel</MDBBtn>
            <MDBBtn color="primary" onClick={onConfirm}>Remove</MDBBtn>
        </MDBModalFooter>
    </MDBModal>
)