import ClientCard from "./ClientCard";
import styled from "styled-components";
import {useQuery} from "react-query";
import CreateClient from "./CreateClient";


const Grid = styled.div`
  margin-top: 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(270px, 1fr));
  grid-gap: 2rem;
`


export default function Dashboard() {
    const {data, isSuccess, refetch} = useQuery(`${process.env.REACT_APP_API_SERVER}/clients`)

    const onSubmit = (newClient) => {
        return fetch(`${process.env.REACT_APP_API_SERVER}/clients`, {
            method: "POST",
            body: JSON.stringify(newClient)
        }).then(() => refetch())
    }

    const onDelete = client => {
        return fetch(`${process.env.REACT_APP_API_SERVER}/clients/${client.id}`, {method: "DELETE"}).then(() => refetch())
    }

    return <>
        <Grid>
            {isSuccess && data.map(client => <ClientCard onDelete={() => onDelete(client)} key={client.id} client={client}/>)}
        </Grid>
        <CreateClient onSubmit={onSubmit}/>
    </>
}