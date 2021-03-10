import {useQuery} from "react-query"
import styled from "styled-components"
import {MDBContainer, MDBNavbar, MDBNavbarBrand,} from "mdbreact";
import ClientCard from "./ClientCard";
import CreateClient from "./CreateClient";

const Grid = styled.div`
  margin-top: 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
  grid-gap: 2rem;
`

function App() {
    const {data, isSuccess, refetch} = useQuery(["clients"], {
        queryFn: () => fetch(`${process.env.REACT_APP_API_SERVER}/clients`).then(res => res.json())
    })

    const onSubmit = (newClient) => {
        return fetch(`${process.env.REACT_APP_API_SERVER}/clients`, {
            method: "POST",
            body: JSON.stringify(newClient)
        }).then(() => refetch())
    }

    return (
        <div className="App">
            <MDBNavbar color="indigo" dark expand="md">
                <MDBNavbarBrand>
                    <strong className="white-text">WG VPN Server</strong>
                </MDBNavbarBrand>
            </MDBNavbar>
            <MDBContainer>
                <Grid>
                    {isSuccess && data.map(client => <ClientCard key={client.id} client={client}/>)}
                </Grid>
            </MDBContainer>
            <CreateClient onSubmit={onSubmit}/>
        </div>
    );
}

export default App;


